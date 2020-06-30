package oss

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/cns-io/go-storage/v2/pkg/headers"
	"github.com/cns-io/go-storage/v2/pkg/iowrap"
	"github.com/cns-io/go-storage/v2/types"
	"github.com/cns-io/go-storage/v2/types/info"
)

func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	err = s.bucket.DeleteObject(rp)
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (err error) {
	marker := ""
	delimiter := "/"
	limit := 200

	rp := s.getAbsPath(dir)

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
					ObjectMeta: info.NewObjectMeta(),
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
		if !output.IsTruncated {
			break
		}
	}
	return
}
func (s *Storage) listPrefix(ctx context.Context, prefix string, opt *pairStorageListPrefix) (err error) {
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
		if !output.IsTruncated {
			break
		}
	}
	return
}
func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta info.StorageMeta, err error) {
	meta = info.NewStorageMeta()
	meta.Name = s.bucket.BucketName
	meta.WorkDir = s.workDir
	return
}
func (s *Storage) read(ctx context.Context, path string, opt *pairStorageRead) (rc io.ReadCloser, err error) {
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
func (s *Storage) stat(ctx context.Context, path string, opt *pairStorageStat) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	output, err := s.bucket.GetObjectMeta(rp)
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		ObjectMeta: info.NewObjectMeta(),
	}

	if v := output.Get(headers.ContentLength); v != "" {
		size, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		o.Size = size
	}

	if v := output.Get(headers.LastModified); v != "" {
		lastModified, err := time.Parse(time.RFC1123, v)
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
		setStorageClass(o.ObjectMeta, v)
	}

	return o, nil
}
func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (err error) {
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
