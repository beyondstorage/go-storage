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
func (s *Storage) Metadata() (m metadata.Storage, err error) {
	m = metadata.Storage{
		Name:     s.bucket.BucketName,
		WorkDir:  s.workDir,
		Metadata: make(metadata.Metadata),
	}
	return m, nil
}

// ListDir implements Storager.ListDir
func (s *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s ListDir [%s]: %w"

	opt, err := parseStoragePairListDir(pairs...)
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
				Name:     v,
				Type:     types.ObjectTypeDir,
				Metadata: make(metadata.Metadata),
			}

			if opt.HasDirFunc {
				opt.DirFunc(o)
			}
		}

		for _, v := range output.Objects {
			o := &types.Object{
				Name:      s.getRelPath(v.Key),
				Type:      types.ObjectTypeDir,
				Size:      v.Size,
				UpdatedAt: v.LastModified,
				Metadata:  make(metadata.Metadata),
			}

			o.SetType(v.Type)
			o.SetClass(v.StorageClass)
			o.SetChecksum(v.ETag)

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
		Name:      path,
		Type:      types.ObjectTypeFile,
		Size:      size,
		UpdatedAt: lastModified,
		Metadata:  make(metadata.Metadata),
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
