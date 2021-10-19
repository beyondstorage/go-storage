package gcs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	gs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	. "go.beyondstorage.io/v5/types"
)

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	rp := s.getAbsPath(path)
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			return
		}
		// Add `/` at the end of path to simulate a directory.
		rp += "/"
		o = s.newObject(true)
		o.Mode = ModeDir
	} else {
		o = s.newObject(false)
		o.Mode = ModeRead
	}
	o.ID = rp
	o.Path = path
	return o
}

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *Object, err error) {
	if !s.features.VirtualDir {
		err = NewOperationNotImplementedError("create_dir")
		return
	}
	rp := s.getAbsPath(path)
	// Add `/` at the end of `path` to simulate a directory.
	// ref: https://cloud.google.com/storage/docs/naming-objects
	rp += "/"
	object := s.bucket.Object(rp)
	w := object.NewWriter(ctx)
	w.Size = 0
	if opt.HasStorageClass {
		w.StorageClass = opt.StorageClass
	}
	cerr := w.Close()
	if cerr != nil {
		err = cerr
	}
	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.Mode |= ModeDir
	return
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			err = services.PairUnsupportedError{Pair: ps.WithObjectMode(opt.ObjectMode)}
			return
		}
		rp += "/"
	}
	err = s.bucket.Object(rp).Delete(ctx)
	if err != nil && errors.Is(err, gs.ErrObjectNotExist) {
		// Omit `ErrObjectNotExist` error here.
		// ref: [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
		err = nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	input := &objectPageStatus{
		prefix: s.getAbsPath(path),
	}
	if !opt.HasListMode {
		// Support `ListModePrefix` as the default `ListMode`.
		// ref: [GSP-654](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/654-unify-list-behavior.md)
		opt.ListMode = ListModePrefix
	}
	var nextFn NextObjectFunc
	switch {
	case opt.ListMode.IsDir():
		input.delimiter = "/"
		nextFn = s.nextObjectPageByDir
	case opt.ListMode.IsPrefix():
		nextFn = s.nextObjectPageByPrefix
	default:
		return nil, services.ListModeInvalidError{Actual: opt.ListMode}
	}
	return NewObjectIterator(ctx, nextFn, input), nil
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	meta = NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	return
}

func (s *Storage) nextObjectPageByDir(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)
	it := s.bucket.Objects(ctx, &gs.Query{
		Prefix:    input.prefix,
		Delimiter: input.delimiter,
	})
	remaining := 200
	for remaining > 0 {
		object, err := it.Next()
		if err == iterator.Done {
			return IterateDone
		}
		if err != nil {
			return err
		}
		// Prefix is set only for ObjectAttrs which represent synthetic "directory
		// entries" when iterating over buckets using Query.Delimiter. See
		// ObjectIterator.Next. When set, no other fields in ObjectAttrs will be
		// populated.
		if object.Prefix != "" {
			o := s.newObject(true)
			o.ID = object.Prefix
			o.Path = s.getRelPath(object.Prefix)
			o.Mode |= ModeDir
			page.Data = append(page.Data, o)
			remaining -= 1
			continue
		}
		o, err := s.formatFileObject(object)
		if err != nil {
			return err
		}
		page.Data = append(page.Data, o)
		remaining -= 1
	}
	return nil
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)
	it := s.bucket.Objects(ctx, &gs.Query{
		Prefix: input.prefix,
	})
	remaining := 200
	for remaining > 0 {
		object, err := it.Next()
		if err == iterator.Done {
			return IterateDone
		}
		if err != nil {
			return err
		}
		o, err := s.formatFileObject(object)
		if err != nil {
			return err
		}
		page.Data = append(page.Data, o)
		remaining -= 1
	}
	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)
	var rc io.ReadCloser
	object := s.bucket.Object(rp)
	if opt.HasEncryptionKey {
		object = object.Key(opt.EncryptionKey)
	}
	if opt.HasOffset && !opt.HasSize {
		rc, err = object.NewRangeReader(ctx, opt.Offset, -1)
	} else if !opt.HasOffset && opt.HasSize {
		rc, err = object.NewRangeReader(ctx, 0, opt.Size)
	} else if opt.HasOffset && opt.HasSize {
		rc, err = object.NewRangeReader(ctx, opt.Offset, opt.Size)
	} else {
		rc, err = object.NewReader(ctx)
	}
	if err != nil {
		return 0, err
	}
	defer func() {
		cerr := rc.Close()
		if cerr != nil {
			err = cerr
		}
	}()
	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}
	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	rp := s.getAbsPath(path)
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			err = services.PairUnsupportedError{Pair: ps.WithObjectMode(opt.ObjectMode)}
			return
		}
		rp += "/"
	}
	attr, err := s.bucket.Object(rp).Attrs(ctx)
	if err != nil {
		return nil, err
	}
	o, err = s.formatFileObject(attr)
	if err != nil {
		return nil, err
	}
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o.Path = path
		o.Mode.Add(ModeDir)
	}
	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	// According to GSP-751, we should allow the user to pass in a nil io.Reader.
	// ref: https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/751-write-empty-file-behavior.md
	if r == nil && size == 0 {
		r = bytes.NewReader([]byte{})
	} else if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	} else {
		r = io.LimitReader(r, size)
	}

	rp := s.getAbsPath(path)
	object := s.bucket.Object(rp)
	if opt.HasEncryptionKey {
		object = object.Key(opt.EncryptionKey)
	}
	w := object.NewWriter(ctx)
	defer func() {
		cerr := w.Close()
		if cerr != nil {
			err = cerr
		}
	}()
	w.Size = size
	if opt.HasContentMd5 {
		// FIXME: we need to check value's encoding type.
		w.MD5 = []byte(opt.ContentMd5)
	}
	if opt.HasStorageClass {
		w.StorageClass = opt.StorageClass
	}
	if opt.HasKmsKeyName {
		w.KMSKeyName = opt.KmsKeyName
	}
	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}
	return io.Copy(w, r)
}
