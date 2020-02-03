package kodo

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	qs "github.com/qiniu/api.v7/v7/storage"

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

// List implements Storager.List
func (s *Storage) List(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s List [%s]: %w"

	opt, err := parseStoragePairList(pairs...)
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
				Name:       s.getRelPath(v.Key),
				Type:       types.ObjectTypeDir,
				Size:       v.Fsize,
				UpdatedAt:  convertUnixTimestampToTime(v.PutTime),
				ObjectMeta: metadata.NewObjectMeta(),
			}
			o.SetContentType(v.MimeType)
			o.SetETag(v.Hash)

			storageClass, err := formatStorageClass(v.Type)
			if err != nil {
				return fmt.Errorf(errorMessage, s, path, err)
			}
			o.SetStorageClass(storageClass)

			opt.FileFunc(o)
		}

		marker = nextMarker
		if marker == "" {
			return nil
		}
	}
}

// Read implements Storager.Read
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	const errorMessage = "%s Read [%s]: %w"

	rp := s.getAbsPath(path)

	url := qs.MakePrivateURL(s.bucket.Mac, s.domain, rp, 3600)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	r = resp.Body
	return
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Write [%s]: %w"

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	uploader := qs.NewFormUploader(s.bucket.Cfg)
	ret := qs.PutRet{}
	err = uploader.Put(opt.Context,
		&ret, s.putPolicy.UploadToken(s.bucket.Mac), rp, r, opt.Size, nil)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	const errorMessage = "%s Stat [%s]: %w"

	rp := s.getAbsPath(path)

	fi, err := s.bucket.Stat(s.name, rp)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       fi.Fsize,
		UpdatedAt:  convertUnixTimestampToTime(fi.PutTime),
		ObjectMeta: metadata.NewObjectMeta(),
	}
	o.SetETag(fi.Hash)

	storageClass, err := formatStorageClass(fi.Type)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}
	o.SetStorageClass(storageClass)

	return o, nil
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	rp := s.getAbsPath(path)

	err = s.bucket.Delete(s.name, rp)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}

func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
}

func (s *Storage) getRelPath(path string) string {
	return strings.TrimPrefix(path, s.workDir+"/")
}
