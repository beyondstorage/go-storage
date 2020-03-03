package s3

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the s3 object storage service.
type Storage struct {
	service s3iface.S3API

	name    string
	workDir string
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager s3 {Name: %s, WorkDir: %s}",
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
	rp := s.getAbsPath(path)

	if !opt.HasObjectFunc {
		delimiter = "/"
	}

	var output *s3.ListObjectsV2Output
	for {
		output, err = s.service.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:     aws.String(s.name),
			Prefix:     aws.String(rp),
			MaxKeys:    aws.Int64(1000),
			StartAfter: aws.String(marker),
			Delimiter:  aws.String(delimiter),
		})
		if err != nil {
			return err
		}

		if opt.HasDirFunc {
			for _, v := range output.CommonPrefixes {
				o := &types.Object{
					ID:         *v.Prefix,
					Name:       s.getRelPath(*v.Prefix),
					Type:       types.ObjectTypeDir,
					ObjectMeta: metadata.NewObjectMeta(),
				}

				opt.DirFunc(o)
			}
		}

		if opt.HasObjectFunc || opt.HasFileFunc {
			for _, v := range output.Contents {
				o := &types.Object{
					ID:         *v.Key,
					Name:       s.getRelPath(*v.Key),
					Type:       types.ObjectTypeFile,
					Size:       aws.Int64Value(v.Size),
					UpdatedAt:  aws.TimeValue(v.LastModified),
					ObjectMeta: metadata.NewObjectMeta(),
				}

				if v.StorageClass != nil {
					storageClass, err := formatStorageClass(*v.StorageClass)
					if err != nil {
						return err
					}
					o.SetStorageClass(storageClass)
				}
				if v.ETag != nil {
					o.SetETag(*v.ETag)
				}

				if opt.HasObjectFunc {
					opt.ObjectFunc(o)
				}
				if opt.HasFileFunc {
					opt.FileFunc(o)
				}
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
	defer func() {
		err = s.formatError("read", err, path)
	}()

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
	}

	output, err := s.service.GetObject(input)
	if err != nil {
		return nil, err
	}

	r = output.Body
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReadCloser(r, opt.ReadCallbackFunc)
	}
	return r, nil
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

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
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
		storageClass, err := parseStorageClass(opt.StorageClass)
		if err != nil {
			return err
		}
		input.StorageClass = &storageClass
	}

	_, err = s.service.PutObject(input)
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

	input := &s3.HeadObjectInput{
		Key: aws.String(rp),
	}

	output, err := s.service.HeadObject(input)
	if err != nil {
		return nil, err
	}

	// TODO: Add dir support.

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       aws.Int64Value(output.ContentLength),
		UpdatedAt:  aws.TimeValue(output.LastModified),
		ObjectMeta: metadata.NewObjectMeta(),
	}

	if output.ContentType != nil {
		o.SetContentType(*output.ContentType)
	}
	if output.ETag != nil {
		o.SetETag(*output.ETag)
	}
	if output.StorageClass != nil {
		storageClass, err := formatStorageClass(*output.StorageClass)
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

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
	}

	_, err = s.service.DeleteObject(input)
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
