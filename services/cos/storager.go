package cos

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the cos object storage service.
type Storage struct {
	bucket *cos.BucketService
	object *cos.ObjectService

	name     string
	location string
	workDir  string
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager cos {Name: %s, WorkDir: %s}",
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
	defer func() {
		err = s.formatError("list", err, path)
	}()

	opt, err := parseStoragePairList(pairs...)
	if err != nil {
		return err
	}

	marker := ""
	delimiter := ""
	limit := 200

	rp := s.getAbsPath(path)

	if !opt.HasObjectFunc {
		delimiter = "/"
	}

	for {
		req := &cos.BucketGetOptions{
			Prefix:    rp,
			MaxKeys:   limit,
			Marker:    marker,
			Delimiter: delimiter,
		}

		resp, _, err := s.bucket.Get(opt.Context, req)
		if err != nil {
			return err
		}

		if opt.HasObjectFunc || opt.HasFileFunc {
			for _, v := range resp.Contents {

				o := &types.Object{
					ID:         v.Key,
					Name:       s.getRelPath(v.Key),
					Type:       types.ObjectTypeFile,
					Size:       int64(v.Size),
					ObjectMeta: metadata.NewObjectMeta(),
				}

				// COS returns different value depends on object upload method or
				// encryption method, so we can't treat this value as content-md5
				//
				// ref: https://cloud.tencent.com/document/product/436/7729
				o.SetETag(v.ETag)

				// COS uses ISO8601 format: "2019-05-27T11:26:14.000Z" in List
				//
				// ref: https://cloud.tencent.com/document/product/436/7729
				t, err := time.Parse("2006-01-02T15:04:05.999Z", v.LastModified)
				if err != nil {
					return err
				}
				o.UpdatedAt = t

				storageClass, err := formatStorageClass(v.StorageClass)
				if err != nil {
					return err
				}
				o.SetStorageClass(storageClass)

				if opt.HasObjectFunc {
					opt.ObjectFunc(o)
				}
				if opt.HasFileFunc {
					opt.FileFunc(o)
				}
			}
		}

		if opt.HasDirFunc {
			for _, v := range resp.CommonPrefixes {
				o := &types.Object{
					ID:         v,
					Name:       s.getRelPath(v),
					Type:       types.ObjectTypeDir,
					ObjectMeta: metadata.NewObjectMeta(),
				}

				opt.DirFunc(o)
			}
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
	defer func() {
		err = s.formatError("read", err, path)
	}()

	opt, err := parseStoragePairRead(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	resp, err := s.object.Get(opt.Context, rp, nil)
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

	rp := s.getAbsPath(path)

	putOptions := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: int(opt.Size),
		},
	}
	if opt.HasChecksum {
		putOptions.ContentMD5 = opt.Checksum
	}
	if opt.HasStorageClass {
		storageClass, err := parseStorageClass(opt.StorageClass)
		if err != nil {
			return err
		}
		putOptions.XCosStorageClass = storageClass
	}
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	_, err = s.object.Put(opt.Context, rp, r, putOptions)
	if err != nil {
		return err
	}
	return
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	defer func() {
		err = s.formatError("stat", err, path)
	}()

	opt, err := parseStoragePairStat(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	output, err := s.object.Head(opt.Context, rp, nil)
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       output.ContentLength,
		ObjectMeta: metadata.NewObjectMeta(),
	}

	// COS uses RFC1123 format in HEAD
	//
	// > Last-Modified: Fri, 09 Aug 2019 10:20:56 GMT
	//
	// ref: https://cloud.tencent.com/document/product/436/7745
	lastModified, err := time.Parse(time.RFC1123, output.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	o.UpdatedAt = lastModified

	o.SetETag(output.Header.Get("ETag"))

	storageClass, err := formatStorageClass(output.Header.Get(storageClassHeader))
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

	opt, err := parseStoragePairDelete(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(path)

	_, err = s.object.Delete(opt.Context, rp)
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
