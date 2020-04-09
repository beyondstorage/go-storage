package kodo

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	qs "github.com/qiniu/api.v7/v7/storage"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the gcs service client.
type Storage struct {
	bucket    *qs.BucketManager
	domain    string
	putPolicy qs.PutPolicy // kodo need PutPolicy to generate upload token.

	name    string
	workDir string
	loose   bool
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager kodo {Name: %s, WorkDir: %s}",
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

// ListDir implements Storager.ListDir
func (s *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list_dir", err, path)
	}()

	opt, err := parseStoragePairListDir(pairs...)
	if err != nil {
		return err
	}

	marker := ""
	delimiter := "/"
	rp := s.getAbsPath(path)

	for {
		entries, commonPrefix, nextMarker, _, err := s.bucket.ListFiles(s.name, rp, delimiter, marker, 1000)
		if err != nil {
			return err
		}

		if opt.HasDirFunc {
			for _, v := range commonPrefix {
				o := &types.Object{
					ID:         v,
					Name:       s.getRelPath(v),
					Type:       types.ObjectTypeDir,
					ObjectMeta: metadata.NewObjectMeta(),
				}

				opt.DirFunc(o)
			}
		}

		if opt.HasFileFunc {
			for _, v := range entries {
				o, err := s.formatFileObject(v)
				if err != nil {
					return err
				}

				opt.FileFunc(o)
			}
		}

		marker = nextMarker
		if marker == "" {
			return nil
		}
	}
}

// ListPrefix implements Storager.ListPrefix
func (s *Storage) ListPrefix(prefix string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list_prefix", err, prefix)
	}()

	opt, err := parseStoragePairListPrefix(pairs...)
	if err != nil {
		return err
	}

	marker := ""
	rp := s.getAbsPath(prefix)

	for {
		entries, _, nextMarker, _, err := s.bucket.ListFiles(s.name, rp, "", marker, 1000)
		if err != nil {
			return err
		}

		for _, v := range entries {
			o, err := s.formatFileObject(v)
			if err != nil {
				return err
			}

			opt.ObjectFunc(o)
		}

		marker = nextMarker
		if marker == "" {
			return nil
		}
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

	url := qs.MakePrivateURL(s.bucket.Mac, s.domain, rp, 3600)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	r = resp.Body

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

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	rp := s.getAbsPath(path)

	uploader := qs.NewFormUploader(s.bucket.Cfg)
	ret := qs.PutRet{}
	err = uploader.Put(opt.Context,
		&ret, s.putPolicy.UploadToken(s.bucket.Mac), rp, r, opt.Size, nil)
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

	rp := s.getAbsPath(path)

	fi, err := s.bucket.Stat(s.name, rp)
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       fi.Fsize,
		UpdatedAt:  convertUnixTimestampToTime(fi.PutTime),
		ObjectMeta: metadata.NewObjectMeta(),
	}

	if fi.Hash != "" {
		o.SetETag(fi.Hash)
	}
	if fi.MimeType != "" {
		o.SetContentType(fi.MimeType)
	}

	storageClass, err := formatStorageClass(fi.Type)
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

	rp := s.getAbsPath(path)

	err = s.bucket.Delete(s.name, rp)
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

func (s *Storage) formatFileObject(v qs.ListItem) (o *types.Object, err error) {
	o = &types.Object{
		ID:         v.Key,
		Name:       s.getRelPath(v.Key),
		Type:       types.ObjectTypeFile,
		Size:       v.Fsize,
		UpdatedAt:  convertUnixTimestampToTime(v.PutTime),
		ObjectMeta: metadata.NewObjectMeta(),
	}

	if v.MimeType != "" {
		o.SetContentType(v.MimeType)
	}
	if v.Hash != "" {
		o.SetETag(v.Hash)
	}

	storageClass, err := formatStorageClass(v.Type)
	if err != nil {
		return nil, err
	}
	o.SetStorageClass(storageClass)

	return
}
