package cos

import (
	"context"
	"io"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/aos-dev/go-storage/v2/pkg/headers"
	"github.com/aos-dev/go-storage/v2/pkg/iowrap"
	"github.com/aos-dev/go-storage/v2/types"
	"github.com/aos-dev/go-storage/v2/types/info"
)

func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	_, err = s.object.Delete(ctx, rp)
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

	for {
		req := &cos.BucketGetOptions{
			Prefix:    rp,
			MaxKeys:   limit,
			Marker:    marker,
			Delimiter: delimiter,
		}

		resp, _, err := s.bucket.Get(ctx, req)
		if err != nil {
			return err
		}

		if opt.HasDirFunc {
			for _, v := range resp.CommonPrefixes {
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
			for _, v := range resp.Contents {
				o, err := s.formatFileObject(v)
				if err != nil {
					return err
				}

				opt.FileFunc(o)
			}
		}

		marker = resp.NextMarker
		if !resp.IsTruncated {
			break
		}
	}

	return
}
func (s *Storage) listPrefix(ctx context.Context, prefix string, opt *pairStorageListPrefix) (err error) {
	marker := ""
	limit := 200

	rp := s.getAbsPath(prefix)

	for {
		req := &cos.BucketGetOptions{
			Prefix:  rp,
			MaxKeys: limit,
			Marker:  marker,
		}

		resp, _, err := s.bucket.Get(ctx, req)
		if err != nil {
			return err
		}

		for _, v := range resp.Contents {
			o, err := s.formatFileObject(v)
			if err != nil {
				return err
			}

			opt.ObjectFunc(o)
		}

		marker = resp.NextMarker
		if !resp.IsTruncated {
			break
		}
	}

	return
}
func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta info.StorageMeta, err error) {
	meta = info.NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	return
}
func (s *Storage) read(ctx context.Context, path string, opt *pairStorageRead) (rc io.ReadCloser, err error) {
	rp := s.getAbsPath(path)

	resp, err := s.object.Get(ctx, rp, nil)
	if err != nil {
		return nil, err
	}

	rc = resp.Body

	if opt.HasReadCallbackFunc {
		rc = iowrap.CallbackReadCloser(rc, opt.ReadCallbackFunc)
	}
	return
}
func (s *Storage) stat(ctx context.Context, path string, opt *pairStorageStat) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	output, err := s.object.Head(ctx, rp, nil)
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       output.ContentLength,
		ObjectMeta: info.NewObjectMeta(),
	}

	// COS uses RFC1123 format in HEAD
	//
	// > Last-Modified: Fri, 09 Aug 2019 10:20:56 GMT
	//
	// ref: https://cloud.tencent.com/document/product/436/7745
	if v := output.Header.Get(headers.LastModified); v != "" {
		lastModified, err := time.Parse(time.RFC1123, v)
		if err != nil {
			return nil, err
		}
		o.UpdatedAt = lastModified
	}

	if v := output.Header.Get(headers.ContentType); v != "" {
		o.SetContentType(v)
	}

	if v := output.Header.Get(headers.ETag); v != "" {
		o.SetETag(output.Header.Get(v))
	}

	if v := output.Header.Get(storageClassHeader); v != "" {
		setStorageClass(o.ObjectMeta, v)
	}

	return o, nil
}
func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (err error) {
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
		putOptions.XCosStorageClass = opt.StorageClass
	}
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	_, err = s.object.Put(ctx, rp, r, putOptions)
	if err != nil {
		return err
	}
	return
}
