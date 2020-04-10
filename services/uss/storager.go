package uss

import (
	"fmt"
	"io"
	"strings"

	"github.com/upyun/go-sdk/upyun"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
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
func (s *Storage) Metadata(pairs ...*types.Pair) (m metadata.StorageMeta, err error) {
	m = metadata.NewStorageMeta()
	m.Name = s.name
	m.WorkDir = s.workDir
	return m, nil
}

// ListDir implements Storager.ListDir
func (s *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list_dir", err, path)
	}()

	opt, err := s.parsePairListDir(pairs...)
	if err != nil {
		return err
	}

	maxListLevel := 1

	rp := s.getAbsPath(path)

	ch := make(chan *upyun.FileInfo, 200)
	defer close(ch)

	go func() {
		for v := range ch {
			if v.IsDir {
				if !opt.HasDirFunc {
					continue
				}

				o := &types.Object{
					ID:         v.Name,
					Name:       s.getRelPath(v.Name),
					Type:       types.ObjectTypeDir,
					ObjectMeta: metadata.NewObjectMeta(),
				}

				opt.DirFunc(o)
				continue
			}

			o, err := s.formatFileObject(v)
			if err != nil {
				return
			}

			opt.FileFunc(o)
		}
	}()

	err = s.bucket.List(&upyun.GetObjectsConfig{
		Path:         rp,
		ObjectsChan:  ch,
		MaxListLevel: maxListLevel,
	})
	if err != nil {
		return err
	}
	return
}

// ListPrefix implements Storager.ListPrefix
func (s *Storage) ListPrefix(prefix string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list_prefix", err, prefix)
	}()

	opt, err := s.parsePairListPrefix(pairs...)
	if err != nil {
		return err
	}

	maxListLevel := -1

	rp := s.getAbsPath(prefix)

	ch := make(chan *upyun.FileInfo, 200)
	defer close(ch)

	go func() {
		for v := range ch {
			if v.IsDir {
				continue
			}

			o, err := s.formatFileObject(v)
			if err != nil {
				return
			}

			opt.ObjectFunc(o)
		}
	}()

	err = s.bucket.List(&upyun.GetObjectsConfig{
		Path:         rp,
		ObjectsChan:  ch,
		MaxListLevel: maxListLevel,
	})
	if err != nil {
		return err
	}
	return
}

// Read implements Storager.Read
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	defer func() {
		err = s.formatError("read", err, path)
	}()

	opt, err := s.parsePairRead(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	r, w := io.Pipe()

	_, err = s.bucket.Get(&upyun.GetObjectConfig{
		Path:   rp,
		Writer: w,
	})
	if err != nil {
		return nil, err
	}

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReadCloser(r, opt.ReadCallbackFunc)
	}
	return r, nil
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("write", err, path)
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
		err = s.formatError("stat", err, path)
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
		ObjectMeta: metadata.NewObjectMeta(),
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
		err = s.formatError("delete", err, path)
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
		ObjectMeta: metadata.NewObjectMeta(),
	}

	if v.ETag != "" {
		o.SetETag(v.ETag)
	}
	if v.ContentType != "" {
		o.SetContentType(v.ContentType)
	}

	return o, nil
}
