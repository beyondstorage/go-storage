package azblob

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the azblob service client.
type Storage struct {
	bucket azblob.ContainerURL

	name    string
	workDir string
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager azblob {Name: %s, WorkDir: %s}",
		s.name, s.workDir,
	)
}

// Metadata implements Storager.Metadata
func (s *Storage) Metadata(pairs ...*types.Pair) (m metadata.StorageMeta, err error) {
	m = metadata.NewStorageMeta()
	m.Name = s.name
	m.WorkDir = s.workDir
	return m, nil
}

// ListDir implements Storager.ListDir
func (s *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpListDir, err, path)
	}()

	opt, err := s.parsePairListDir(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(path)

	marker := azblob.Marker{}

	var output *azblob.ListBlobsHierarchySegmentResponse
	for {
		output, err = s.bucket.ListBlobsHierarchySegment(opt.Context, marker, "/", azblob.ListBlobsSegmentOptions{
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

// ListPrefix implements Storager.ListPrefix
func (s *Storage) ListPrefix(prefix string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpListPrefix, err, prefix)
	}()

	opt, err := s.parsePairListPrefix(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(prefix)

	marker := azblob.Marker{}

	var output *azblob.ListBlobsFlatSegmentResponse
	for {
		output, err = s.bucket.ListBlobsFlatSegment(opt.Context, marker, azblob.ListBlobsSegmentOptions{
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

// Read implements Storager.Read
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	defer func() {
		err = s.formatError(services.OpRead, err, path)
	}()

	opt, err := s.parsePairRead(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	offset := int64(0)
	if opt.HasOffset {
		offset = opt.Offset
	}

	count := int64(azblob.CountToEnd)
	if opt.HasSize {
		count = opt.Size
	}

	output, err := s.bucket.NewBlockBlobURL(rp).Download(opt.Context, offset, count, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return nil, err
	}

	r = output.Response().Body
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReadCloser(r, opt.ReadCallbackFunc)
	}
	return r, nil
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpWrite, err, path)
	}()

	opt, err := s.parsePairWrite(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(path)

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	// TODO: add checksum and storage class support.
	_, err = s.bucket.NewBlockBlobURL(rp).Upload(opt.Context, iowrap.ReadSeekCloser(r),
		azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		return err
	}
	return nil
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	defer func() {
		err = s.formatError(services.OpStat, err, path)
	}()

	opt, err := s.parsePairStat(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	output, err := s.bucket.NewBlockBlobURL(rp).GetProperties(opt.Context, azblob.BlobAccessConditions{})
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       output.ContentLength(),
		UpdatedAt:  output.LastModified(),
		ObjectMeta: metadata.NewObjectMeta(),
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
	if v := formatStorageClass(azblob.AccessTierType(output.AccessTier())); v != "" {
		o.SetStorageClass(v)
	}

	return o, nil
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpDelete, err, path)
	}()

	opt, err := s.parsePairStat(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(path)

	_, err = s.bucket.NewBlockBlobURL(rp).Delete(opt.Context,
		azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		return err
	}
	return nil
}

// getAbsPath will calculate object storage's abs path
func (s *Storage) getAbsPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return prefix + path
}

// getRelPath will get object storage's rel path.
func (s *Storage) getRelPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return strings.TrimPrefix(path, prefix)
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

func (s *Storage) formatFileObject(v azblob.BlobItem) (o *types.Object, err error) {
	o = &types.Object{
		ID:         v.Name,
		Name:       s.getRelPath(v.Name),
		Type:       types.ObjectTypeFile,
		UpdatedAt:  v.Properties.LastModified,
		ObjectMeta: metadata.NewObjectMeta(),
	}

	o.SetETag(string(v.Properties.Etag))

	if v.Properties.ContentLength != nil {
		o.Size = *v.Properties.ContentLength
	}
	if v.Properties.ContentType != nil && *v.Properties.ContentType != "" {
		o.SetContentType(*v.Properties.ContentType)
	}
	if len(v.Properties.ContentMD5) > 0 {
		o.SetContentMD5(base64.StdEncoding.EncodeToString(v.Properties.ContentMD5))
	}
	if value := formatStorageClass(v.Properties.AccessTier); value != "" {
		o.SetStorageClass(value)
	}

	return o, nil
}

func (s *Storage) formatDirObject(v azblob.BlobPrefix) (o *types.Object) {
	o = &types.Object{
		ID:         v.Name,
		Name:       s.getRelPath(v.Name),
		Type:       types.ObjectTypeDir,
		Size:       0,
		ObjectMeta: metadata.NewObjectMeta(),
	}

	return o
}
