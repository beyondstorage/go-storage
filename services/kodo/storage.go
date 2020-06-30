package kodo

import (
	"context"
	"io"
	"net/http"
	"time"

	qs "github.com/qiniu/api.v7/v7/storage"

	"github.com/cns-io/go-storage/v2/pkg/iowrap"
	"github.com/cns-io/go-storage/v2/types"
	"github.com/cns-io/go-storage/v2/types/info"
)

func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	err = s.bucket.Delete(s.name, rp)
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (err error) {
	marker := ""
	delimiter := "/"
	rp := s.getAbsPath(dir)

	for {
		entries, commonPrefix, nextMarker, _, err := s.bucket.ListFiles(s.name, rp, delimiter, marker, 1000)
		if err != nil {
			return err
		}

		if opt.HasDirFunc {
			for _, v := range commonPrefix {
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
			for _, v := range entries {
				o, err := s.formatFileObject(v)
				if err != nil {
					return err
				}

				opt.FileFunc(o)
			}
		}

		marker = nextMarker
		if marker == "" {
			return nil
		}
	}
}
func (s *Storage) listPrefix(ctx context.Context, prefix string, opt *pairStorageListPrefix) (err error) {
	marker := ""
	rp := s.getAbsPath(prefix)

	for {
		entries, _, nextMarker, _, err := s.bucket.ListFiles(s.name, rp, "", marker, 1000)
		if err != nil {
			return err
		}

		for _, v := range entries {
			o, err := s.formatFileObject(v)
			if err != nil {
				return err
			}

			opt.ObjectFunc(o)
		}

		marker = nextMarker
		if marker == "" {
			return nil
		}
	}
}
func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta info.StorageMeta, err error) {
	meta = info.NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	return meta, nil
}
func (s *Storage) read(ctx context.Context, path string, opt *pairStorageRead) (rc io.ReadCloser, err error) {
	rp := s.getAbsPath(path)

	deadline := time.Now().Add(time.Hour).Unix()
	url := qs.MakePrivateURL(s.bucket.Mac, s.domain, rp, deadline)
	resp, err := s.bucket.Client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = qs.ResponseError(resp)
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

	fi, err := s.bucket.Stat(s.name, rp)
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       fi.Fsize,
		UpdatedAt:  convertUnixTimestampToTime(fi.PutTime),
		ObjectMeta: info.NewObjectMeta(),
	}

	if fi.Hash != "" {
		o.SetETag(fi.Hash)
	}
	if fi.MimeType != "" {
		o.SetContentType(fi.MimeType)
	}

	setStorageClass(o.ObjectMeta, fi.Type)

	return o, nil
}
func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (err error) {
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	rp := s.getAbsPath(path)

	uploader := qs.NewFormUploader(s.bucket.Cfg)
	ret := qs.PutRet{}
	err = uploader.Put(ctx,
		&ret, s.putPolicy.UploadToken(s.bucket.Mac), rp, r, opt.Size, nil)
	if err != nil {
		return err
	}
	return nil
}
