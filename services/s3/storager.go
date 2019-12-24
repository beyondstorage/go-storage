package s3

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// Storage is the s3 object storage service.
//
//go:generate ../../internal/bin/meta
type Storage struct {
	service s3iface.S3API

	name    string
	workDir string

	segments    map[string]*segment.Segment
	segmentLock sync.RWMutex
}

// newStorage will create a new client.
func newStorage(service s3iface.S3API, bucketName string) (*Storage, error) {
	c := &Storage{
		service:  service,
		name:     bucketName,
		segments: make(map[string]*segment.Segment),
	}
	return c, nil
}

func (s *Storage) Init(pairs ...*types.Pair) (err error) {
	errorMessage := "s3 Init: %w"

	opt, err := parseStoragePairInit(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	if opt.HasWorkDir {
		// TODO: we should validate workDir
		s.workDir = strings.TrimLeft(opt.WorkDir, "/")
	}

	return nil
}

func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager s3 {Name: %s, WorkDir: %s}",
		s.name, "/"+s.workDir,
	)
}

func (s *Storage) Metadata() (m metadata.Storage, err error) {
	m = metadata.Storage{
		Name:     s.name,
		WorkDir:  s.workDir,
		Metadata: make(metadata.Metadata),
	}
	return m, nil
}

func (s *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	errorMessage := "s3 ListDir: %w"

	opt, err := parseStoragePairListDir(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	marker := ""
	rp := s.getAbsPath(path)

	var output *s3.ListObjectsV2Output
	for {
		output, err = s.service.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:     aws.String(s.name),
			Prefix:     aws.String(rp),
			MaxKeys:    aws.Int64(1000),
			StartAfter: aws.String(marker),
		})
		if err != nil {
			err = handleS3Error(err)
			return fmt.Errorf(errorMessage, err)
		}

		for _, v := range output.CommonPrefixes {
			o := &types.Object{
				Name:     s.getRelPath(*v.Prefix),
				Type:     types.ObjectTypeDir,
				Metadata: make(metadata.Metadata),
			}

			if opt.HasDirFunc {
				opt.DirFunc(o)
			}
		}

		for _, v := range output.Contents {
			o := &types.Object{
				Type:      types.ObjectTypeFile,
				Name:      s.getRelPath(*v.Key),
				Size:      aws.Int64Value(v.Size),
				UpdatedAt: aws.TimeValue(v.LastModified),
				Metadata:  make(metadata.Metadata),
			}

			if v.StorageClass != nil {
				o.SetClass(*v.StorageClass)
			}
			if v.ETag != nil {
				o.SetChecksum(*v.ETag)
			}

			if opt.HasFileFunc {
				opt.FileFunc(o)
			}
		}

		marker = aws.StringValue(output.StartAfter)
		if aws.BoolValue(output.IsTruncated) {
			break
		}
	}
	return
}

func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	panic("implement me")
}

func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	panic("implement me")
}

func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	panic("implement me")
}

func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}
