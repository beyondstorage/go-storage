package s3

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the s3 object storage service.
//
//go:generate ../../internal/bin/meta
type Storage struct {
	service s3iface.S3API

	name    string
	workDir string
}

// newStorage will create a new client.
func newStorage(service s3iface.S3API, bucketName string) (*Storage, error) {
	c := &Storage{
		service: service,
		name:    bucketName,
	}
	return c, nil
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

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager s3 {Name: %s, WorkDir: %s}",
		s.name, "/"+s.workDir,
	)
}

// Metadata implements Storager.Metadata
func (s *Storage) Metadata() (m metadata.Storage, err error) {
	m = metadata.Storage{
		Name:     s.name,
		WorkDir:  s.workDir,
		Metadata: make(metadata.Metadata),
	}
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
			return fmt.Errorf(errorMessage, s, path, err)
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

// Read implements Storager.Read
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	const errorMessage = "%s Read [%s]: %w"

	rp := s.getAbsPath(path)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
	}

	output, err := s.service.GetObject(input)
	if err != nil {
		err = handleS3Error(err)
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}
	return output.Body, nil
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Write [%s]: %w"

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	input := &s3.PutObjectInput{
		Key:           aws.String(rp),
		ContentLength: &opt.Size,
		Body:          aws.ReadSeekCloser(r),
	}
	if opt.HasChecksum {
		input.ContentMD5 = &opt.Checksum
	}
	if opt.HasStorageClass {
		input.StorageClass = &opt.StorageClass
	}

	_, err = s.service.PutObject(input)
	if err != nil {
		err = handleS3Error(err)
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	const errorMessage = "%s Stat [%s]: %w"

	rp := s.getAbsPath(path)

	input := &s3.HeadObjectInput{
		Key: aws.String(rp),
	}

	output, err := s.service.HeadObject(input)
	if err != nil {
		err = handleS3Error(err)
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	// TODO: Add dir support.

	o = &types.Object{
		Name:      path,
		Type:      types.ObjectTypeFile,
		Size:      aws.Int64Value(output.ContentLength),
		UpdatedAt: aws.TimeValue(output.LastModified),
		Metadata:  make(metadata.Metadata),
	}

	if output.ContentType != nil {
		o.SetType(*output.ContentType)
	}
	if output.ETag != nil {
		o.SetChecksum(*output.ETag)
	}
	if output.StorageClass != nil {
		o.SetClass(*output.StorageClass)
	}
	return o, nil
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	rp := s.getAbsPath(path)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
	}

	_, err = s.service.DeleteObject(input)
	if err != nil {
		err = handleS3Error(err)
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}
