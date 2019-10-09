package qingstor

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
)

// DirectoryContentType is the content type for qingstor object directory.
const DirectoryContentType = "application/x-directory"

// Client is the qingstor object storage client.
//
//go:generate go run ../../internal/cmd/meta_gen/main.go
type Client struct {
	config  *Config
	service *service.Service
	bucket  *service.Bucket

	segments map[string]*segment.Segment
}

// setupBucket will setup bucket for client.
func (c *Client) setupBucket(bucketName, zoneName string) (err error) {
	errorMessage := "setup qingstor bucket failed: %w"

	if zoneName != "" {
		bucket, err := c.service.Bucket(bucketName, zoneName)
		if err != nil {
			return fmt.Errorf(errorMessage, err)
		}
		c.bucket = bucket
		return nil
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url := fmt.Sprintf("%s://%s.%s:%d", c.config.Protocol, bucketName, c.config.Host, c.config.Port)

	r, err := client.Head(url)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	if r.StatusCode != http.StatusTemporaryRedirect {
		err = fmt.Errorf("head status is %d instead of %d", r.StatusCode, http.StatusTemporaryRedirect)
		return fmt.Errorf(errorMessage, err)
	}

	// Example URL: https://bucket.zone.qingstor.com
	zoneName = strings.Split(r.Header.Get("Location"), ".")[1]
	bucket, err := c.service.Bucket(bucketName, zoneName)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	c.bucket = bucket
	return
}

// Stat implements Storager.Stat
func (c *Client) Stat(path string, option ...types.Option) (i types.Informer, err error) {
	errorMessage := "qingstor Stat failed: %w"

	input := &service.HeadObjectInput{}

	output, err := c.bucket.HeadObject(path, input)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	if *output.ContentType == DirectoryContentType {
		return &types.Dir{Name: path}, nil
	}
	return &types.File{
		Name: path,
		Size: *output.ContentLength,
		Type: *output.ContentType,
	}, nil
}

// Delete implements Storager.Delete
func (c *Client) Delete(path string, option ...types.Option) (err error) {
	errorMessage := "qingstor Delete failed: %w"

	_, err = c.bucket.DeleteObject(path)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}

// Copy implements Storager.Copy
func (c *Client) Copy(src, dst string, option ...types.Option) (err error) {
	errorMessage := "qingstor Copy failed: %w"

	_, err = c.bucket.PutObject(dst, &service.PutObjectInput{
		XQSCopySource: &src,
	})
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}

// Move implements Storager.Move
func (c *Client) Move(src, dst string, option ...types.Option) (err error) {
	errorMessage := "qingstor Move failed: %w"

	_, err = c.bucket.PutObject(dst, &service.PutObjectInput{
		XQSMoveSource: &src,
	})
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}

// ListDir implements Storager.ListDir
func (c *Client) ListDir(path string, option ...types.Option) (it iterator.Iterator) {
	errorMessage := "qingstor ListDir failed: %w"

	marker := ""
	limit := 200

	var output *service.ListObjectsOutput
	var err error

	fn := iterator.NextFunc(func(informer *[]types.Informer) error {
		idx := 0
		buf := make([]types.Informer, limit)

		output, err = c.bucket.ListObjects(&service.ListObjectsInput{
			Limit:  &limit,
			Marker: &marker,
			Prefix: &path,
		})
		if err != nil {
			return fmt.Errorf(errorMessage, err)
		}

		for _, v := range output.Keys {
			if *v.MimeType == DirectoryContentType {
				buf[idx] = &types.Dir{Name: *v.Key}
				idx++
				continue
			}
			buf[idx] = &types.File{
				Name: *v.Key,
				Size: *v.Size,
				Type: *v.MimeType,
			}
			idx++
		}

		marker = *output.NextMarker
		if marker == "" {
			return iterator.ErrDone
		}
		if output.HasMore != nil && !*output.HasMore {
			return iterator.ErrDone
		}
		if len(output.Keys) == 0 {
			return iterator.ErrDone
		}
		return nil
	})

	it = iterator.NewPrefixBasedIterator(fn)
	return
}

// ReadFile implements Storager.ReadFile
func (c *Client) ReadFile(path string, option ...types.Option) (r io.ReadCloser, err error) {
	errorMessage := "qingstor ReadFile failed: %w"

	_ = parseOptionReadFile(option...)
	input := &service.GetObjectInput{}

	output, err := c.bucket.GetObject(path, input)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}
	return output.Body, nil
}

// WriteFile implements Storager.WriteFile
func (c *Client) WriteFile(path string, size int64, r io.ReadCloser, option ...types.Option) (err error) {
	errorMessage := "qingstor WriteFile failed: %w"

	defer r.Close()

	opts := parseOptionWriteFile(option...)
	input := &service.PutObjectInput{
		ContentLength: &size,
		Body:          r,
	}
	if opts.HasMd5 {
		input.ContentMD5 = &opts.Md5
	}
	if opts.HasStorageClass {
		input.XQSStorageClass = &opts.StorageClass
	}

	_, err = c.bucket.PutObject(path, input)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}

// ReadStream implements Storager.ReadStream
func (c *Client) ReadStream(path string, option ...types.Option) (r io.ReadCloser, err error) {
	panic("not supported")
}

// WriteStream implements Storager.WriteStream
func (c *Client) WriteStream(path string, r io.ReadCloser, option ...types.Option) (err error) {
	panic("not supported")
}

// InitSegment implements Storager.InitSegment
func (c *Client) InitSegment(path string, size int64, option ...types.Option) (err error) {
	errorMessage := "qingstor InitSegment failed: %w"

	if _, ok := c.segments[path]; ok {
		return fmt.Errorf("Segment %s has been initiated", path)
	}

	_ = parseOptionInitSegment(option...)
	input := &service.InitiateMultipartUploadInput{}

	output, err := c.bucket.InitiateMultipartUpload(path, input)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	c.segments[path] = &segment.Segment{
		TotalSize: size,
		ID:        *output.UploadID,
		Parts:     make([]*segment.Part, 0),
	}
	return
}

// ReadSegment implements Storager.ReadSegment
func (c *Client) ReadSegment(path string, offset, size int64, option ...types.Option) (r io.ReadCloser, err error) {
	panic("implement me")
}

// WriteSegment implements Storager.WriteSegment
func (c *Client) WriteSegment(path string, offset, size int64, r io.ReadCloser, option ...types.Option) (err error) {
	errorMessage := "qingstor WriteSegment failed: %w"

	s, ok := c.segments[path]
	if !ok {
		return fmt.Errorf(errorMessage, fmt.Errorf("segment %s is not initiated", path))
	}

	p := &segment.Part{
		Offset: offset,
		Size:   size,
	}

	partNumber, err := s.GetPartIndex(p)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	_, err = c.bucket.UploadMultipart(path, &service.UploadMultipartInput{
		PartNumber:    &partNumber,
		UploadID:      &s.ID,
		ContentLength: &size,
		Body:          r,
	})
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	err = s.InsertPart(p)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	return
}

// CompleteSegment implements Storager.CompleteSegment
func (c *Client) CompleteSegment(path string, option ...types.Option) (err error) {
	errorMessage := "qingstor CompleteSegment failed: %w"

	s, ok := c.segments[path]
	if !ok {
		return fmt.Errorf(errorMessage, fmt.Errorf("segment %s is not initiated", path))
	}

	err = s.ValidateParts()
	if err != nil {
		return
	}

	objectParts := make([]*service.ObjectPartType, len(s.Parts))
	for k, v := range s.Parts {
		partNumber := k
		objectParts[k] = &service.ObjectPartType{
			PartNumber: &partNumber,
			Size:       &v.Size,
		}
	}

	_, err = c.bucket.CompleteMultipartUpload(path, &service.CompleteMultipartUploadInput{
		UploadID:    &s.ID,
		ObjectParts: objectParts,
	})
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	delete(c.segments, path)
	return
}

// AbortSegment implements Storager.AbortSegment
func (c *Client) AbortSegment(path string, option ...types.Option) (err error) {
	errorMessage := "qingstor AbortSegment failed: %w"

	s, ok := c.segments[path]
	if !ok {
		return fmt.Errorf(errorMessage, fmt.Errorf("segment %s is not initiated", path))
	}

	_, err = c.bucket.AbortMultipartUpload(path, &service.AbortMultipartUploadInput{
		UploadID: &s.ID,
	})
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	delete(c.segments, path)
	return
}
