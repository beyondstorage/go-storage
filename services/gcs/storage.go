package gcs

import (
	"context"
	"io"

	gs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/aos-dev/go-storage/v2/pkg/iowrap"
	"github.com/aos-dev/go-storage/v2/types"
	"github.com/aos-dev/go-storage/v2/types/info"
)

func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	err = s.bucket.Object(rp).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (err error) {
	delimiter := "/"

	rp := s.getAbsPath(dir)

	it := s.bucket.Objects(ctx, &gs.Query{
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
func (s *Storage) listPrefix(ctx context.Context, prefix string, opt *pairStorageListPrefix) (err error) {
	rp := s.getAbsPath(prefix)

	it := s.bucket.Objects(ctx, &gs.Query{Prefix: rp})
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
func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta info.StorageMeta, err error) {
	meta = info.NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	return
}
func (s *Storage) read(ctx context.Context, path string, opt *pairStorageRead) (rc io.ReadCloser, err error) {
	rp := s.getAbsPath(path)

	object := s.bucket.Object(rp)
	rc, err = object.NewReader(ctx)
	if err != nil {
		return nil, err
	}

	if opt.HasReadCallbackFunc {
		rc = iowrap.CallbackReadCloser(rc, opt.ReadCallbackFunc)
	}
	return
}
func (s *Storage) stat(ctx context.Context, path string, opt *pairStorageStat) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	attr, err := s.bucket.Object(rp).Attrs(ctx)
	if err != nil {
		return nil, err
	}

	return s.formatFileObject(attr)
}
func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (err error) {
	rp := s.getAbsPath(path)

	object := s.bucket.Object(rp)
	w := object.NewWriter(ctx)
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
