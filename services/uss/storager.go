package uss

import (
	"errors"
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
	loose   bool
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager uss {Name: %s, WorkDir: %s}",
		s.name, "/"+s.workDir)
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

	maxListLevel := -1

	rp := s.getAbsPath(path)

	if !opt.HasObjectFunc {
		maxListLevel = 1
	}

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

			o := &types.Object{
				ID:         v.Name,
				Name:       s.getRelPath(v.Name),
				Type:       types.ObjectTypeFile,
				Size:       v.Size,
				UpdatedAt:  v.Time,
				ObjectMeta: metadata.NewObjectMeta(),
			}
			o.SetETag(v.ETag)
			o.SetContentType(v.ContentType)

			if opt.HasObjectFunc {
				opt.ObjectFunc(o)
			}
			if opt.HasFileFunc {
				opt.FileFunc(o)
			}
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

	opt, err := parseStoragePairRead(pairs...)
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

	opt, err := parseStoragePairWrite(pairs...)
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
	o.SetETag(output.ETag)
	o.SetContentType(output.ContentType)

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

	if s.loose && errors.Is(err, services.ErrCapabilityInsufficient) {
		return nil
	}

	return &services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}
