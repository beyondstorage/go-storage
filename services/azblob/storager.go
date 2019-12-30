package azblob

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the azblob service client.
//
//go:generate ../../internal/bin/meta
type Storage struct {
	bucket azblob.ContainerURL

	name    string
	workDir string
}

// newStorage will create a new client.
func newStorage(bucket azblob.ContainerURL, name string) *Storage {
	c := &Storage{
		bucket: bucket,
		name:   name,
	}
	return c
}

// String implements Storager.String
func (s Storage) String() string {
	return fmt.Sprintf(
		"Storager azblob {Name: %s, WorkDir: %s}",
		s.name, "/"+s.workDir,
	)
}

// Init implements Storager.Init
func (s Storage) Init(pairs ...*types.Pair) (err error) {
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
func (s Storage) Metadata() (m metadata.Storage, err error) {
	m = metadata.Storage{
		Name:     s.name,
		WorkDir:  s.workDir,
		Metadata: make(metadata.Metadata),
	}
	return m, nil
}

// ListDir implements Storager.ListDir
func (s Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s ListDir [%s]: %w"

	opt, err := parseStoragePairListDir(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	marker := azblob.Marker{}
	var output *azblob.ListBlobsFlatSegmentResponse
	for {
		output, err = s.bucket.ListBlobsFlatSegment(context.TODO(), marker, azblob.ListBlobsSegmentOptions{
			Prefix: rp,
		})
		if err != nil {
			return fmt.Errorf(errorMessage, s, path, err)
		}

		for _, v := range output.Segment.BlobItems {
			o := &types.Object{
				Name:      s.getRelPath(v.Name),
				Type:      types.ObjectTypeDir,
				Size:      *v.Properties.ContentLength,
				UpdatedAt: v.Properties.LastModified,
				Metadata:  make(metadata.Metadata),
			}
			o.SetType(*v.Properties.ContentType)
			o.SetClass(string(v.Properties.AccessTier))
			o.SetChecksum(string(v.Properties.ContentMD5))

			opt.FileFunc(o)
		}

		marker = output.NextMarker
		if !marker.NotDone() {
			break
		}
	}

	return
}

// Read implements Storager.Read
func (s Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	const errorMessage = "%s Read [%s]: %w"

	rp := s.getAbsPath(path)

	output, err := s.bucket.NewBlockBlobURL(rp).Download(context.TODO(), 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}
	return output.Body(azblob.RetryReaderOptions{}), nil
}

// Write implements Storager.Write
func (s Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Write [%s]: %w"

	_, err = parseStoragePairWrite(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	// TODO: add checksum and storage class support.
	_, err = s.bucket.NewBlockBlobURL(rp).Upload(context.TODO(), iowrap.NewReadSeekCloser(r),
		azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}

// Stat implements Storager.Stat
func (s Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	const errorMessage = "%s Stat [%s]: %w"

	rp := s.getAbsPath(path)

	output, err := s.bucket.NewBlockBlobURL(rp).GetProperties(context.TODO(), azblob.BlobAccessConditions{})
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	o = &types.Object{
		Name:      path,
		Type:      types.ObjectTypeFile,
		Size:      output.ContentLength(),
		UpdatedAt: output.LastModified(),
		Metadata:  make(metadata.Metadata),
	}
	return o, nil
}

// Delete implements Storager.Delete
func (s Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	rp := s.getAbsPath(path)

	_, err = s.bucket.NewBlockBlobURL(rp).Delete(context.TODO(),
		azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}
