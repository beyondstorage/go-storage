package minio

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"

	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

const defaultListObjectBufferSize = 100

func (s *Storage) copy(ctx context.Context, src string, dst string, opt pairStorageCopy) (err error) {
	srcOpts := minio.CopySrcOptions{
		Bucket: s.bucket,
		Object: s.getAbsPath(src),
	}
	dstOpts := minio.CopyDestOptions{
		Bucket: s.bucket,
		Object: s.getAbsPath(dst),
	}
	_, err = s.client.CopyObject(ctx, dstOpts, srcOpts)
	return err
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *types.Object) {
	rp := s.getAbsPath(path)
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			return
		}
		rp += "/"
		o = s.newObject(true)
		o.Mode = types.ModeDir
	} else {
		o = s.newObject(false)
		o.Mode = types.ModeRead
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
	err = s.client.RemoveObject(ctx, s.bucket, rp, minio.RemoveObjectOptions{})
	return err
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *types.ObjectIterator, err error) {
	rp := s.getAbsPath(path)
	options := minio.ListObjectsOptions{}
	if !opt.HasListMode || opt.ListMode.IsPrefix() {
		options.Recursive = true
	} else if opt.ListMode.IsDir() {
		if !strings.HasSuffix(rp, "/") {
			rp += "/"
		}
	} else {
		return nil, services.ListModeInvalidError{Actual: opt.ListMode}
	}
	options.Prefix = rp
	input := &objectPageStatus{
		bufferSize: defaultListObjectBufferSize,
		options:    options,
	}
	return types.NewObjectIterator(ctx, s.nextObjectPage, input), nil
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *types.StorageMeta) {
	meta = types.NewStorageMeta()
	meta.Name = s.bucket
	meta.WorkDir = s.workDir
	return meta
}

func (s *Storage) nextObjectPage(ctx context.Context, page *types.ObjectPage) error {
	input := page.Status.(*objectPageStatus)
	if input.objChan == nil {
		input.objChan = s.client.ListObjects(ctx, s.bucket, input.options)
	}
	for i := 0; i < input.bufferSize; i++ {
		v, ok := <-input.objChan
		if !ok {
			return types.IterateDone
		}
		if v.Err != nil {
			return v.Err
		}
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}
		page.Data = append(page.Data, o)
		input.counter++
	}
	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)
	output, err := s.client.GetObject(ctx, s.bucket, rp, minio.GetObjectOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		cerr := output.Close()
		if cerr != nil {
			err = cerr
		}
	}()
	if opt.HasOffset {
		_, err = output.Seek(opt.Offset, 0)
		if err != nil {
			return 0, err
		}
	}
	var rc io.ReadCloser = output
	if opt.HasSize {
		rc = iowrap.LimitReadCloser(rc, opt.Size)
	}
	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}
	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *types.Object, err error) {
	rp := s.getAbsPath(path)
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			err = services.PairUnsupportedError{Pair: ps.WithObjectMode(opt.ObjectMode)}
			return
		}
		rp += "/"
	}
	output, err := s.client.StatObject(ctx, s.bucket, rp, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	o, err = s.formatFileObject(output)
	if err != nil {
		return nil, err
	}
	// Object.Path is either the absolute path or the relative path based on the working directory depends on the input.
	o.Path = path
	return
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	// According to GSP-751, we should allow the user to pass in a nil io.Reader.
	// ref: https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/751-write-empty-file-behavior.md
	if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	}

	rp := s.getAbsPath(path)
	r = io.LimitReader(r, size)
	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}
	options := minio.PutObjectOptions{}
	if opt.HasContentType {
		options.ContentType = opt.ContentType
	}
	if opt.HasStorageClass {
		options.StorageClass = opt.StorageClass
	}
	_, err = s.client.PutObject(ctx, s.bucket, rp, r, size, options)
	if err != nil {
		return 0, err
	}
	return size, nil
}
