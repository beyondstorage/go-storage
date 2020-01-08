package oss

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the aliyun object storage service.
//
//go:generate ../../internal/bin/meta
//go:generate ../../internal/bin/context
type Storage struct {
	bucket *oss.Bucket

	name    string
	workDir string
}

// newStorage will create a new client.
func newStorage(bucket *oss.Bucket) *Storage {
	c := &Storage{
		bucket: bucket,
	}
	return c
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager oss {Name: %s, WorkDir: %s}",
		s.bucket.BucketName, "/"+s.workDir,
	)
}

// Init implements Storager.Init
func (s *Storage) Init(pairs ...*types.Pair) (err error) {
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
func (s *Storage) Metadata(pairs ...*types.Pair) (m metadata.StorageMeta, err error) {
	m = metadata.NewStorageMeta()
	m.Name = s.bucket.BucketName
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
	limit := 200

	rp := s.getAbsPath(path)

	var output oss.ListObjectsResult
	for {
		output, err = s.bucket.ListObjects(
			oss.Marker(marker),
			oss.MaxKeys(limit),
			oss.Prefix(rp),
		)
		if err != nil {
			return fmt.Errorf(errorMessage, s, path, err)
		}

		for _, v := range output.CommonPrefixes {
			o := &types.Object{
				ID:         v,
				Name:       s.getRelPath(v),
				Type:       types.ObjectTypeDir,
				ObjectMeta: metadata.NewObjectMeta(),
			}

			if opt.HasDirFunc {
				opt.DirFunc(o)
			}
		}

		for _, v := range output.Objects {
			o := &types.Object{
				ID:         v.Key,
				Name:       s.getRelPath(v.Key),
				Type:       types.ObjectTypeDir,
				Size:       v.Size,
				UpdatedAt:  v.LastModified,
				ObjectMeta: metadata.NewObjectMeta(),
			}

			o.SetContentType(v.Type)
			o.SetStorageClass(v.StorageClass)
			o.SetETag(v.ETag)

			if opt.HasFileFunc {
				opt.FileFunc(o)
			}
		}

		marker = output.NextMarker
		if output.IsTruncated {
			break
		}
	}
	return
}

// Read implements Storager.Read
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	const errorMessage = "%s Read [%s]: %w"

	rp := s.getAbsPath(path)

	output, err := s.bucket.GetObject(rp)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}
	return output, nil
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Write [%s]: %w"

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	options := make([]oss.Option, 0)
	if opt.HasChecksum {
		options = append(options, oss.ContentMD5(opt.Checksum))
	}
	if opt.HasSize {
		options = append(options, oss.ContentLength(opt.Size))
	}
	if opt.HasStorageClass {
		// TODO: we need to handle different storage class name between services.
		options = append(options, oss.StorageClass(oss.StorageClassType(opt.StorageClass)))
	}

	rp := s.getAbsPath(path)

	err = s.bucket.PutObject(rp, r, options...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	const errorMessage = "%s Stat [%s]: %w"

	rp := s.getAbsPath(path)

	output, err := s.bucket.GetObjectMeta(rp)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	// Parse content length.
	size, err := strconv.ParseInt(output.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}
	// Parse last modified.
	lastModified, err := time.Parse(time.RFC822, output.Get("Last-Modified"))
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	// TODO: get object's checksum and storage class.
	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       size,
		UpdatedAt:  lastModified,
		ObjectMeta: metadata.NewObjectMeta(),
	}
	return o, nil
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	rp := s.getAbsPath(path)

	err = s.bucket.DeleteObject(rp)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}
