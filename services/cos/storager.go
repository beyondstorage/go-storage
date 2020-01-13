package cos

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// Storage is the cos object storage service.
//
//go:generate ../../internal/bin/service
type Storage struct {
	bucket *cos.BucketService
	object *cos.ObjectService

	name     string
	location string
	workDir  string
}

// newStorage will create a new client.
func newStorage(bucketName, region string, client *http.Client) *Storage {
	s := &Storage{}

	url := cos.NewBucketURL(bucketName, region, true)
	c := cos.NewClient(&cos.BaseURL{BucketURL: url}, client)
	s.bucket = c.Bucket
	s.object = c.Object
	s.name = bucketName
	s.location = region
	return s
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager cos {Name: %s, WorkDir: %s}",
		s.name, "/"+s.workDir,
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
	limit := 200

	rp := s.getAbsPath(path)

	for {
		req := &cos.BucketGetOptions{
			Prefix:  rp,
			MaxKeys: limit,
			Marker:  marker,
		}

		resp, _, err := s.bucket.Get(opt.Context, req)
		if err != nil {
			return fmt.Errorf(errorMessage, s, path, err)
		}

		for _, v := range resp.Contents {
			// COS use ISO8601 format: 2019-05-27T11:26:14.000Z
			t, err := time.Parse("2006-01-02T15:04:05.999Z", v.LastModified)
			if err != nil {
				return fmt.Errorf(errorMessage, s, path, err)
			}

			o := &types.Object{
				ID:         v.Key,
				Name:       s.getRelPath(v.Key),
				Type:       types.ObjectTypeFile,
				Size:       int64(v.Size),
				UpdatedAt:  t,
				ObjectMeta: metadata.NewObjectMeta(),
			}
			o.SetETag(v.ETag)

			storageClass, err := formatStorageClass(v.StorageClass)
			if err != nil {
				return fmt.Errorf(errorMessage, s, path, err)
			}
			o.SetStorageClass(storageClass)

			opt.FileFunc(o)
		}

		marker = resp.NextMarker
		if !resp.IsTruncated {
			break
		}
	}

	return
}

// Read implements Storager.Read
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	const errorMessage = "%s Read [%s]: %w"

	opt, err := parseStoragePairRead(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	resp, err := s.object.Get(opt.Context, rp, nil)
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

	putOptions := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{},
	}
	if opt.HasChecksum {
		putOptions.ContentMD5 = opt.Checksum
	}
	if opt.HasSize {
		putOptions.ContentLength = int(opt.Size)
	}
	if opt.HasStorageClass {
		storageClass, err := parseStorageClass(opt.StorageClass)
		if err != nil {
			return fmt.Errorf(errorMessage, s, path, err)
		}
		putOptions.XCosStorageClass = storageClass
	}

	_, err = s.object.Put(opt.Context, rp, r, putOptions)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	const errorMessage = "%s Stat [%s]: %w"

	opt, err := parseStoragePairStat(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	output, err := s.object.Head(opt.Context, rp, nil)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	lastModified, err := time.Parse(time.RFC822, output.Header.Get("Last-Modified"))
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       output.ContentLength,
		UpdatedAt:  lastModified,
		ObjectMeta: metadata.NewObjectMeta(),
	}

	storageClass, err := formatStorageClass(output.Header.Get("x-cos-storage-class"))
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}
	o.SetStorageClass(storageClass)

	return o, nil
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	opt, err := parseStoragePairDelete(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	_, err = s.object.Delete(opt.Context, rp)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}
