package uss

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/upyun/go-sdk/upyun"

	"github.com/Xuanwo/storage/pkg/headers"
	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/info"
)

// delete implements Storager.Delete
//
// USS requires a short time between PUT and DELETE, or we will get this error:
// DELETE 429 {"msg":"concurrent put or delete","code":42900007,"id":"xxx"}
//
// Due to this problem, uss can't pass the storager integration tests.
func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	err = s.bucket.Delete(&upyun.DeleteObjectConfig{
		Path: rp,
	})
	if err != nil {
		return err
	}
	return
}
func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (err error) {
	// err could be updated in multiple goroutines, add explict lock to protect it.
	var errlock sync.Mutex

	rp := s.getAbsPath(dir)

	// USS SDK will close this channel in List
	ch := make(chan *upyun.FileInfo, 200)

	go func() {
		xerr := s.bucket.List(&upyun.GetObjectsConfig{
			Path:         rp,
			ObjectsChan:  ch,
			MaxListLevel: 1,
		})

		errlock.Lock()
		defer errlock.Unlock()
		err = xerr
	}()

	for v := range ch {
		if v.IsDir {
			if !opt.HasDirFunc {
				continue
			}

			o := &types.Object{
				ID:         v.Name,
				Name:       s.getRelPath(v.Name),
				Type:       types.ObjectTypeDir,
				ObjectMeta: info.NewObjectMeta(),
			}

			opt.DirFunc(o)
			continue
		}

		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		opt.FileFunc(o)
	}
	return
}
func (s *Storage) listPrefix(ctx context.Context, prefix string, opt *pairStorageListPrefix) (err error) {
	// err could be updated in multiple goroutines, add explict lock to protect it.
	var errlock sync.Mutex

	rp := s.getAbsPath(prefix)

	// USS SDK will close this channel in List
	ch := make(chan *upyun.FileInfo, 200)

	go func() {
		xerr := s.bucket.List(&upyun.GetObjectsConfig{
			Path:         rp,
			ObjectsChan:  ch,
			MaxListLevel: -1,
		})

		errlock.Lock()
		defer errlock.Unlock()
		err = xerr
	}()

	for v := range ch {
		if v.IsDir {
			continue
		}

		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		opt.ObjectFunc(o)
	}
	return
}
func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta info.StorageMeta, err error) {
	meta = info.NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	return meta, nil
}
func (s *Storage) read(ctx context.Context, path string, opt *pairStorageRead) (rc io.ReadCloser, err error) {
	rp := s.getAbsPath(path)

	var w *io.PipeWriter
	r, w = io.Pipe()

	go func() {
		defer w.Close()

		_, err = s.bucket.Get(&upyun.GetObjectConfig{
			Path:   rp,
			Writer: w,
		})
	}()

	if opt.HasReadCallbackFunc {
		rc = iowrap.CallbackReadCloser(rc, opt.ReadCallbackFunc)
	}
	return rc, nil
}
func (s *Storage) stat(ctx context.Context, path string, opt *pairStorageStat) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	output, err := s.bucket.GetInfo(rp)
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       output.Size,
		UpdatedAt:  output.Time,
		ObjectMeta: info.NewObjectMeta(),
	}
	if output.ETag != "" {
		o.SetETag(output.ETag)
	}
	if output.ContentType != "" {
		o.SetContentType(output.ContentType)
	}

	return o, nil
}
func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (err error) {
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	rp := s.getAbsPath(path)

	cfg := &upyun.PutObjectConfig{
		Path:   rp,
		Reader: r,
		Headers: map[string]string{
			headers.ContentLength: strconv.FormatInt(opt.Size, 10),
		},
	}

	err = s.bucket.Put(cfg)
	if err != nil {
		return err
	}
	return
}
