package bos

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bos/api"

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

		rp += "/"
		o = s.newObject(true)
		o.Mode |= ModeDir
	} else {
		o = s.newObject(false)
		o.Mode |= ModeRead
	}
	o.ID = rp
	o.Path = path

	return o
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

	err = s.client.DeleteObject(s.bucket, rp)
	if err != nil {
		if e, ok := err.(*bce.BceServiceError); ok && e.Code == "NoSuchKey" {
			// bos DeleteObject is not idempotent, so we need to check object_not_exists error.
			//
			// - [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
			// - https://cloud.baidu.com/doc/BOS/s/bkc5tsslq
			err = nil
		} else {
			return err
		}
	}

	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	input := &objectPageStatus{
		maxKeys: 200,
		prefix:  s.getAbsPath(path),
	}

	if !opt.HasListMode {
		// Support `ListModePrefix` as the default `ListMode`.
		// ref: [GSP-46](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/654-unify-list-behavior.md)
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
	meta.Name = s.bucket
	meta.WorkDir = s.workDir
	return meta
}

func (s *Storage) nextObjectPageByDir(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	listArgs := &api.ListObjectsArgs{
		Delimiter: input.delimiter,
		MaxKeys:   input.maxKeys,
		Prefix:    input.prefix,
		Marker:    input.marker,
	}

	output, err := s.client.ListObjects(s.bucket, listArgs)
	if err != nil {
		return err
	}

	for _, v := range output.CommonPrefixes {
		o := s.newObject(true)
		o.ID = v.Prefix
		o.Path = s.getRelPath(v.Prefix)
		o.Mode |= ModeDir

		page.Data = append(page.Data, o)
	}

	for _, v := range output.Contents {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if output.NextMarker == "" {
		return IterateDone
	}
	if !output.IsTruncated {
		return IterateDone
	}

	input.marker = output.NextMarker

	return nil
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	listArgs := &api.ListObjectsArgs{
		Delimiter: input.delimiter,
		MaxKeys:   input.maxKeys,
		Prefix:    input.prefix,
		Marker:    input.marker,
	}

	output, err := s.client.ListObjects(s.bucket, listArgs)
	if err != nil {
		return err
	}

	for _, v := range output.Contents {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if output.NextMarker == "" {
		return IterateDone
	}
	if !output.IsTruncated {
		return IterateDone
	}

	input.marker = output.NextMarker

	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)

	output := &api.GetObjectResult{}
	if opt.HasOffset && !opt.HasSize {
		output, err = s.client.GetObject(s.bucket, rp, nil, opt.Offset)
	} else if !opt.HasOffset && opt.HasSize {
		output, err = s.client.GetObject(s.bucket, rp, nil, 0, opt.Size-1)
	} else if opt.HasSize && opt.HasOffset {
		output, err = s.client.GetObject(s.bucket, rp, nil, opt.Offset, opt.Offset+opt.Size-1)
	} else {
		output, err = s.client.GetObject(s.bucket, rp, nil)
	}

	if err != nil {
		return 0, err
	}

	rc := output.Body
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
			return nil, err
		}

		rp += "/"
	}

	output, err := s.client.GetObject(s.bucket, rp, nil)
	if err != nil {
		return nil, err
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o.Mode |= ModeDir
	} else {
		o.Mode |= ModeRead
	}

	o.SetContentLength(output.ContentLength)
	// Last-Modified returns a format of :
	// Fri, 28 Jan 2011 20:10:32 GMT
	// ref:https://cloud.baidu.com/doc/BOS/s/xkc5pcmcj#%E7%A4%BA%E4%BE%8B
	lastModified, err := time.Parse(time.RFC1123, output.LastModified)
	if err != nil {
		return nil, err
	}
	o.SetLastModified(lastModified)

	if output.ContentType != "" {
		o.SetContentType(output.ContentType)
	}
	if output.ETag != "" {
		o.SetEtag(output.ETag)
	}

	var sm ObjectSystemMetadata
	if v := output.StorageClass; v != "" {
		sm.StorageClass = v
	}

	o.SetSystemMetadata(sm)

	return
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	if size > writeSizeMaximum {
		err = fmt.Errorf("size limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return 0, err
	}

	// According to GSP-751, we should allow the user to pass in a nil io.Reader.
	// ref: https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/751-write-empty-file-behavior.md
	if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	}

	rp := s.getAbsPath(path)

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	body, err := bce.NewBodyFromSizedReader(r, size)
	if err != nil {
		return 0, err
	}
	putArgs := &api.PutObjectArgs{
		ContentLength: size,
	}

	if opt.HasContentMd5 {
		putArgs.ContentMD5 = opt.ContentMd5
	}
	if opt.HasStorageClass {
		putArgs.StorageClass = opt.StorageClass
	}

	_, err = s.client.PutObject(s.bucket, rp, body, putArgs)
	if err != nil {
		return 0, err
	}

	return size, nil
}
