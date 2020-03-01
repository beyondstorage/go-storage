package gcs

import (
	"fmt"
	"io"
	"strings"

	gs "cloud.google.com/go/storage"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the gcs service client.
type Storage struct {
	bucket *gs.BucketHandle

	name    string
	workDir string
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager gcs {Name: %s, WorkDir: %s}",
		s.name, "/"+s.workDir,
	)
}

// Metadata implements Storager.Metadata
func (s *Storage) Metadata(pairs ...*types.Pair) (m metadata.StorageMeta, err error) {
	m = metadata.NewStorageMeta()
	m.Name = s.name
	m.WorkDir = s.workDir
	return m, nil
}

// List implements Storager.List
func (s *Storage) List(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list", err, path)
	}()

	opt, err := parseStoragePairList(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(path)

	for {
		it := s.bucket.Objects(opt.Context, &gs.Query{
			Prefix: rp,
		})
		object, err := it.Next()
		if err != nil && err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}

		o := &types.Object{
			ID:         object.Name,
			Name:       s.getRelPath(object.Name),
			Type:       types.ObjectTypeDir,
			Size:       object.Size,
			UpdatedAt:  object.Updated,
			ObjectMeta: metadata.NewObjectMeta(),
		}
		o.SetContentType(object.ContentType)
		o.SetContentMD5(string(object.MD5))

		storageClass, err := formatStorageClass(object.StorageClass)
		if err != nil {
			return err
		}
		o.SetStorageClass(storageClass)

		opt.FileFunc(o)
	}
}

// Read implements Storager.Read
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	defer func() {
		err = s.formatError("read", err, path)
	}()

	opt, err := parseStoragePairRead(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	object := s.bucket.Object(rp)
	r, err = object.NewReader(opt.Context)
	if err != nil {
		return nil, err
	}

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReadCloser(r, opt.ReadCallbackFunc)
	}
	return
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("write", err, path)
	}()

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(path)

	object := s.bucket.Object(rp)
	w := object.NewWriter(opt.Context)
	defer w.Close()

	w.Size = opt.Size
	if opt.HasChecksum {
		w.MD5 = []byte(opt.Checksum)
	}
	if opt.HasStorageClass {
		storageClass, err := parseStorageClass(opt.StorageClass)
		if err != nil {
			return err
		}
		w.StorageClass = storageClass
	}
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	defer func() {
		err = s.formatError("stat", err, path)
	}()

	opt, err := parseStoragePairStat(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	attr, err := s.bucket.Object(rp).Attrs(opt.Context)
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         attr.Name,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       attr.Size,
		UpdatedAt:  attr.Updated,
		ObjectMeta: metadata.NewObjectMeta(),
	}

	storageClass, err := formatStorageClass(attr.StorageClass)
	if err != nil {
		return nil, err
	}
	o.SetStorageClass(storageClass)

	return o, nil
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("delete", err, path)
	}()

	opt, err := parseStoragePairStat(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(path)

	err = s.bucket.Object(rp).Delete(opt.Context)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
}

func (s *Storage) getRelPath(path string) string {
	return strings.TrimPrefix(path, s.workDir+"/")
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	if err == nil {
		return nil
	}

	// Handle errors returned by gcs.
	e, ok := err.(*googleapi.Error)
	if ok {
		err = formatGcsError(e)
	}

	return &services.StorageError{
		Op:       op,
		Err:      err,
		Storager: s,
		Path:     path,
	}
}
