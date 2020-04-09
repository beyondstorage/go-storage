package oss

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/Xuanwo/storage/pkg/headers"
	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the aliyun object storage service.
//
//go:generate ../../internal/bin/service
type Storage struct {
	bucket *oss.Bucket

	name    string
	workDir string
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager oss {Name: %s, WorkDir: %s}",
		s.bucket.BucketName, "/"+s.workDir,
	)
}

// Metadata implements Storager.Metadata
func (s *Storage) Metadata(pairs ...*types.Pair) (m metadata.StorageMeta, err error) {
	m = metadata.NewStorageMeta()
	m.Name = s.bucket.BucketName
	m.WorkDir = s.workDir
	return m, nil
}

// ListDir implements Storager.ListDir
func (s *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list_dir", err, path)
	}()

	opt, err := s.parsePairListDir(pairs...)
	if err != nil {
		return err
	}

	marker := ""
	delimiter := "/"
	limit := 200

	rp := s.getAbsPath(path)

	var output oss.ListObjectsResult
	for {
		output, err = s.bucket.ListObjects(
			oss.Marker(marker),
			oss.MaxKeys(limit),
			oss.Prefix(rp),
			oss.Delimiter(delimiter),
		)
		if err != nil {
			return err
		}

		if opt.HasDirFunc {
			for _, v := range output.CommonPrefixes {
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
			for _, v := range output.Objects {
				o, err := s.formatFileObject(v)
				if err != nil {
					return err
				}

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

// ListPrefix implements Storager.ListPrefix
func (s *Storage) ListPrefix(prefix string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list_prefix", err, prefix)
	}()

	opt, err := s.parsePairListPrefix(pairs...)
	if err != nil {
		return err
	}

	marker := ""
	limit := 200

	rp := s.getAbsPath(prefix)

	var output oss.ListObjectsResult
	for {
		output, err = s.bucket.ListObjects(
			oss.Marker(marker),
			oss.MaxKeys(limit),
			oss.Prefix(rp),
		)
		if err != nil {
			return err
		}

		for _, v := range output.Objects {
			o, err := s.formatFileObject(v)
			if err != nil {
				return err
			}

			opt.ObjectFunc(o)
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
	defer func() {
		err = s.formatError("read", err, path)
	}()

	opt, err := s.parsePairWrite(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	output, err := s.bucket.GetObject(rp)
	if err != nil {
		return nil, err
	}

	if opt.HasReadCallbackFunc {
		output = iowrap.CallbackReadCloser(output, opt.ReadCallbackFunc)
	}

	return output, nil
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("write", err, path)
	}()

	opt, err := s.parsePairWrite(pairs...)
	if err != nil {
		return err
	}

	options := make([]oss.Option, 0)
	options = append(options, oss.ContentLength(opt.Size))
	if opt.HasChecksum {
		options = append(options, oss.ContentMD5(opt.Checksum))
	}
	if opt.HasStorageClass {
		// TODO: we need to handle different storage class name between services.
		options = append(options, oss.StorageClass(oss.StorageClassType(opt.StorageClass)))
	}
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	rp := s.getAbsPath(path)

	err = s.bucket.PutObject(rp, r, options...)
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

	output, err := s.bucket.GetObjectMeta(rp)
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		ObjectMeta: metadata.NewObjectMeta(),
	}

	if v := output.Get(headers.ContentLength); v != "" {
		size, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		o.Size = size
	}

	if v := output.Get(headers.LastModified); v != "" {
		lastModified, err := time.Parse(time.RFC822, v)
		if err != nil {
			return nil, err
		}
		o.UpdatedAt = lastModified
	}

	// OSS advise us don't use Etag as Content-MD5.
	//
	// ref: https://help.aliyun.com/document_detail/31965.html
	if v := output.Get(headers.ETag); v != "" {
		o.SetETag(v)
	}

	if v := output.Get(headers.ContentType); v != "" {
		o.SetContentType(v)
	}

	if v := output.Get(storageClassHeader); v != "" {
		storageClass, err := formatStorageClass(v)
		if err != nil {
			return nil, err
		}
		o.SetStorageClass(storageClass)
	}

	return o, nil
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("delete", err, path)
	}()

	rp := s.getAbsPath(path)

	err = s.bucket.DeleteObject(rp)
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

	return &services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}

func (s *Storage) formatFileObject(v oss.ObjectProperties) (o *types.Object, err error) {
	o = &types.Object{
		ID:         v.Key,
		Name:       s.getRelPath(v.Key),
		Type:       types.ObjectTypeFile,
		Size:       v.Size,
		UpdatedAt:  v.LastModified,
		ObjectMeta: metadata.NewObjectMeta(),
	}

	if v.Type != "" {
		o.SetContentType(v.Type)
	}

	// OSS advise us don't use Etag as Content-MD5.
	//
	// ref: https://help.aliyun.com/document_detail/31965.html
	if v.ETag != "" {
		o.SetETag(v.ETag)
	}

	if v.Type != "" {
		storageClass, err := formatStorageClass(v.Type)
		if err != nil {
			return nil, err
		}
		o.SetStorageClass(storageClass)
	}

	return
}
