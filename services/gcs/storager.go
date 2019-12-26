package gcs

import (
	"context"
	"fmt"
	"io"
	"strings"

	gs "cloud.google.com/go/storage"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
	"google.golang.org/api/iterator"
)

// Storage is the gcs service client.
//
//go:generate ../../internal/bin/meta
type Storage struct {
	bucket *gs.BucketHandle

	name    string
	workDir string
}

// newStorage will create a new client.
func newStorage(bucket *gs.BucketHandle, name string) *Storage {
	c := &Storage{
		bucket: bucket,
		name:   name,
	}
	return c
}

// String implements Storager.String
func (s Storage) String() string {
	return fmt.Sprintf(
		"Storager gcs {Name: %s, WorkDir: %s}",
		s.name, "/"+s.workDir,
	)
}

// Init implements Storager.Init
func (s Storage) Init(pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Init: %w"

	opt, err := parseStoragePairInit(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, err)
	}

	if opt.HasWorkDir {
		// TODO: we should validate workDir
		s.workDir = strings.TrimLeft(opt.WorkDir, "/")
	}

	return nil
}

// Metadata implements Storager.Metadata
func (s Storage) Metadata() (m metadata.Storage, err error) {
	m = metadata.Storage{
		Name:     s.name,
		WorkDir:  s.workDir,
		Metadata: make(metadata.Metadata),
	}
	return m, nil
}

// ListDir implements Storager.ListDir
func (s Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s ListDir [%s]: %w"

	opt, err := parseStoragePairListDir(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	for {
		it := s.bucket.Objects(context.TODO(), &gs.Query{
			Prefix: rp,
		})
		object, err := it.Next()
		if err != nil && err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf(errorMessage, s, path, err)
		}

		o := &types.Object{
			Name:      s.getRelPath(object.Name),
			Type:      types.ObjectTypeDir,
			Size:      object.Size,
			UpdatedAt: object.Updated,
			Metadata:  make(metadata.Metadata),
		}
		o.SetType(object.ContentType)
		o.SetClass(object.StorageClass)
		o.SetChecksum(string(object.MD5))

		opt.FileFunc(o)
	}
}

// Read implements Storager.Read
func (s Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	const errorMessage = "%s Read [%s]: %w"

	rp := s.getAbsPath(path)

	object := s.bucket.Object(rp)
	r, err = object.NewReader(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}
	return
}

// Write implements Storager.Write
func (s Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Write [%s]: %w"

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	object := s.bucket.Object(rp)
	w := object.NewWriter(context.TODO())
	defer w.Close()

	if opt.HasChecksum {
		w.MD5 = []byte(opt.Checksum)
	}
	if opt.HasSize {
		w.Size = opt.Size
	}
	if opt.HasStorageClass {
		w.StorageClass = opt.StorageClass
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}

// Stat implements Storager.Stat
func (s Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	const errorMessage = "%s Stat [%s]: %w"

	rp := s.getAbsPath(path)

	attr, err := s.bucket.Object(rp).Attrs(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	o = &types.Object{
		Name:      path,
		Type:      types.ObjectTypeFile,
		Size:      attr.Size,
		UpdatedAt: attr.Updated,
		Metadata:  make(metadata.Metadata),
	}
	return o, nil
}

// Delete implements Storager.Delete
func (s Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	rp := s.getAbsPath(path)

	err = s.bucket.Object(rp).Delete(context.TODO())
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}
