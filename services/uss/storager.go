package uss

import (
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

// Storage is the uss service.
type Storage struct {
	bucket *upyun.UpYun

	name    string
	workDir string
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager uss {Name: %s, WorkDir: %s}",
		s.name, s.workDir)
}

// Metadata implements Storager.Metadata
func (s *Storage) Metadata(pairs ...*types.Pair) (m info.StorageMeta, err error) {
	m = info.NewStorageMeta()
	m.Name = s.name
	m.WorkDir = s.workDir
	return m, nil
}

// ListDir implements Storager.ListDir
func (s *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	// err could be updated in multiple goroutines, add explict lock to protect it.
	var errlock sync.Mutex

	defer func() {
		errlock.Lock()
		defer errlock.Unlock()

		err = s.formatError(services.OpListDir, err, path)
	}()

	opt, err := s.parsePairListDir(pairs...)
	if err != nil {
		return err
	}

	maxListLevel := 1

	rp := s.getAbsPath(path)

	// USS SDK will close this channel in List
	ch := make(chan *upyun.FileInfo, 200)

	go func() {
		errlock.Lock()
		defer errlock.Unlock()

		err = s.bucket.List(&upyun.GetObjectsConfig{
			Path:         rp,
			ObjectsChan:  ch,
			MaxListLevel: maxListLevel,
		})
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

// ListPrefix implements Storager.ListPrefix
func (s *Storage) ListPrefix(prefix string, pairs ...*types.Pair) (err error) {
	// err could be updated in multiple goroutines, add explict lock to protect it.
	var errlock sync.Mutex

	defer func() {
		errlock.Lock()
		defer errlock.Unlock()

		err = s.formatError(services.OpListPrefix, err, prefix)
	}()

	opt, err := s.parsePairListPrefix(pairs...)
	if err != nil {
		return err
	}

	maxListLevel := -1

	rp := s.getAbsPath(prefix)

	// USS SDK will close this channel in List
	ch := make(chan *upyun.FileInfo, 200)

	go func() {
		errlock.Lock()
		defer errlock.Unlock()

		err = s.bucket.List(&upyun.GetObjectsConfig{
			Path:         rp,
			ObjectsChan:  ch,
			MaxListLevel: maxListLevel,
		})
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

// Read implements Storager.Read
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	defer func() {
		err = s.formatError(services.OpRead, err, path)
	}()

	opt, err := s.parsePairRead(pairs...)
	if err != nil {
		return nil, err
	}

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
		r = iowrap.CallbackReadCloser(r, opt.ReadCallbackFunc)
	}
	return r, nil
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpWrite, err, path)
	}()

	opt, err := s.parsePairWrite(pairs...)
	if err != nil {
		return err
	}

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

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	defer func() {
		err = s.formatError(services.OpStat, err, path)
	}()

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

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpDelete, err, path)
	}()

	rp := s.getAbsPath(path)

	err = s.bucket.Delete(&upyun.DeleteObjectConfig{
		Path: rp,
	})
	if err != nil {
		return err
	}
	return
}

// getAbsPath will calculate object storage's abs path
func (s *Storage) getAbsPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return prefix + path
}

// getRelPath will get object storage's rel path.
func (s *Storage) getRelPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return strings.TrimPrefix(path, prefix)
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	if err == nil {
		return nil
	}

	return &services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}

func (s *Storage) formatFileObject(v *upyun.FileInfo) (o *types.Object, err error) {
	o = &types.Object{
		ID:         v.Name,
		Name:       s.getRelPath(v.Name),
		Type:       types.ObjectTypeFile,
		Size:       v.Size,
		UpdatedAt:  v.Time,
		ObjectMeta: info.NewObjectMeta(),
	}

	if v.ETag != "" {
		o.SetETag(v.ETag)
	}
	if v.ContentType != "" {
		o.SetContentType(v.ContentType)
	}

	return o, nil
}
