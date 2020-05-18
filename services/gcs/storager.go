package gcs

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	gs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/info"
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
		s.name, s.workDir,
	)
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
	defer func() {
		err = s.formatError(services.OpListDir, err, path)
	}()

	opt, err := s.parsePairListDir(pairs...)
	if err != nil {
		return err
	}

	delimiter := "/"

	rp := s.getAbsPath(path)

	it := s.bucket.Objects(opt.Context, &gs.Query{
		Prefix:    rp,
		Delimiter: delimiter,
	})

	for {
		object, err := it.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}

		// Prefix is set only for ObjectAttrs which represent synthetic "directory
		// entries" when iterating over buckets using Query.Delimiter. See
		// ObjectIterator.Next. When set, no other fields in ObjectAttrs will be
		// populated.
		if object.Prefix != "" {
			if !opt.HasDirFunc {
				continue
			}

			o := &types.Object{
				ID:         object.Prefix,
				Name:       s.getRelPath(object.Prefix),
				Type:       types.ObjectTypeDir,
				ObjectMeta: info.NewObjectMeta(),
			}

			opt.DirFunc(o)
			continue
		}

		o, err := s.formatFileObject(object)
		if err != nil {
			return err
		}

		if opt.HasFileFunc {
			opt.FileFunc(o)
		}
	}
}

// ListPrefix implements Storager.ListPrefix
func (s *Storage) ListPrefix(prefix string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpListPrefix, err, prefix)
	}()

	opt, err := s.parsePairListPrefix(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(prefix)

	it := s.bucket.Objects(opt.Context, &gs.Query{Prefix: rp})
	for {
		object, err := it.Next()
		if err != nil && err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}

		o, err := s.formatFileObject(object)
		if err != nil {
			return err
		}

		opt.ObjectFunc(o)
	}
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
		err = s.formatError(services.OpWrite, err, path)
	}()

	opt, err := s.parsePairWrite(pairs...)
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
		w.StorageClass = opt.StorageClass
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
		err = s.formatError(services.OpStat, err, path)
	}()

	opt, err := s.parsePairStat(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	attr, err := s.bucket.Object(rp).Attrs(opt.Context)
	if err != nil {
		return nil, err
	}

	return s.formatFileObject(attr)
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpDelete, err, path)
	}()

	opt, err := s.parsePairStat(pairs...)
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

func (s *Storage) formatFileObject(v *gs.ObjectAttrs) (o *types.Object, err error) {
	o = &types.Object{
		ID:         v.Name,
		Name:       s.getRelPath(v.Name),
		Type:       types.ObjectTypeFile,
		Size:       v.Size,
		UpdatedAt:  v.Updated,
		ObjectMeta: info.NewObjectMeta(),
	}

	if v.ContentType != "" {
		o.SetContentType(v.ContentType)
	}
	if v.Etag != "" {
		o.SetETag(v.Etag)
	}
	if len(v.MD5) > 0 {
		o.SetContentMD5(base64.StdEncoding.EncodeToString(v.MD5))
	}
	if value := v.StorageClass; value != "" {
		setStorageClass(o.ObjectMeta, value)
	}

	return
}
