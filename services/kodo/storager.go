package kodo

import (
	"fmt"
	"io"
	"strings"

	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
	qs "github.com/qiniu/api.v7/storage"
)

// Storage is the gcs service client.
//
//go:generate ../../internal/bin/meta
type Storage struct {
	bucket *qs.BucketManager

	name    string
	workDir string
}

// newStorage will create a new client.
func newStorage(bucket *qs.BucketManager, name string) *Storage {
	c := &Storage{
		bucket: bucket,
		name:   name,
	}

	return c
}

// String implements Storager.String
func (s Storage) String() string {
	return fmt.Sprintf(
		"Storager kodo {Name: %s, WorkDir: %s}",
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

	marker := ""
	rp := s.getAbsPath(path)

	for {
		entries, _, nextMarker, _, err := s.bucket.ListFiles(s.name, rp, "", marker, 1000)
		if err != nil {
			return fmt.Errorf(errorMessage, s, path, err)
		}

		for _, v := range entries {
			o := &types.Object{
				Name:      s.getRelPath(v.Key),
				Type:      types.ObjectTypeDir,
				Size:      v.Fsize,
				UpdatedAt: convertUnixTimestampToTime(v.PutTime),
				Metadata:  make(metadata.Metadata),
			}
			o.SetType(v.MimeType)
			// kodo's object checksum is not md5, let's keep empty.

			opt.FileFunc(o)
		}

		marker = nextMarker
		if marker == "" {
			return nil
		}
	}
}

// Read implements Storager.Read
func (s Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	panic("implement me")
}

// Write implements Storager.Write
func (s Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// Stat implements Storager.Stat
func (s Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	panic("implement me")
}

// Delete implements Storager.Delete
func (s Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}
