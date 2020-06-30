package azblob

import (
	"context"
	"encoding/base64"
	"io"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/cns-io/go-storage/v2/pkg/iowrap"
	"github.com/cns-io/go-storage/v2/types"
	"github.com/cns-io/go-storage/v2/types/info"
)

func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	_, err = s.bucket.NewBlockBlobURL(rp).Delete(ctx,
		azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (err error) {
	rp := s.getAbsPath(dir)

	marker := azblob.Marker{}

	var output *azblob.ListBlobsHierarchySegmentResponse
	for {
		output, err = s.bucket.ListBlobsHierarchySegment(ctx, marker, "/", azblob.ListBlobsSegmentOptions{
			Prefix: rp,
		})
		if err != nil {
			return err
		}

		if opt.HasDirFunc {
			for _, v := range output.Segment.BlobPrefixes {
				o := s.formatDirObject(v)

				opt.DirFunc(o)
			}
		}

		if opt.HasFileFunc {
			for _, v := range output.Segment.BlobItems {
				o, err := s.formatFileObject(v)
				if err != nil {
					return err
				}

				opt.FileFunc(o)
			}
		}

		marker = output.NextMarker
		if !marker.NotDone() {
			break
		}
	}
	return
}
func (s *Storage) listPrefix(ctx context.Context, prefix string, opt *pairStorageListPrefix) (err error) {
	rp := s.getAbsPath(prefix)

	marker := azblob.Marker{}

	var output *azblob.ListBlobsFlatSegmentResponse
	for {
		output, err = s.bucket.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{
			Prefix: rp,
		})
		if err != nil {
			return err
		}

		for _, v := range output.Segment.BlobItems {
			o, err := s.formatFileObject(v)
			if err != nil {
				return err
			}

			opt.ObjectFunc(o)
		}

		marker = output.NextMarker
		if !marker.NotDone() {
			break
		}
	}
	return
}
func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta info.StorageMeta, err error) {
	meta = info.NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	return meta, nil
}
func (s *Storage) read(ctx context.Context, path string, opt *pairStorageRead) (rc io.ReadCloser, err error) {
	rp := s.getAbsPath(path)

	offset := int64(0)
	if opt.HasOffset {
		offset = opt.Offset
	}

	count := int64(azblob.CountToEnd)
	if opt.HasSize {
		count = opt.Size
	}

	output, err := s.bucket.NewBlockBlobURL(rp).Download(ctx, offset, count, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return nil, err
	}

	rc = output.Response().Body
	if opt.HasReadCallbackFunc {
		rc = iowrap.CallbackReadCloser(rc, opt.ReadCallbackFunc)
	}
	return rc, nil
}
func (s *Storage) stat(ctx context.Context, path string, opt *pairStorageStat) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	output, err := s.bucket.NewBlockBlobURL(rp).GetProperties(ctx, azblob.BlobAccessConditions{})
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       output.ContentLength(),
		UpdatedAt:  output.LastModified(),
		ObjectMeta: info.NewObjectMeta(),
	}

	if v := string(output.ETag()); v != "" {
		o.SetETag(v)
	}
	if v := output.ContentType(); v != "" {
		o.SetContentType(v)
	}
	if v := output.ContentMD5(); len(v) > 0 {
		o.SetContentMD5(base64.StdEncoding.EncodeToString(v))
	}
	if v := StorageClass(output.AccessTier()); v != "" {
		setStorageClass(o.ObjectMeta, v)
	}

	return o, nil
}
func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (err error) {
	rp := s.getAbsPath(path)

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	// TODO: add checksum and storage class support.
	_, err = s.bucket.NewBlockBlobURL(rp).Upload(ctx, iowrap.SizedReadSeekCloser(r, opt.Size),
		azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		return err
	}
	return nil
}
