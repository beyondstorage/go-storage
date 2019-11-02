package qingstor

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/pengsrc/go-shared/convert"
	iface "github.com/yunify/qingstor-sdk-go/v3/interface"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
)

// Client is the qingstor object storage client.
//
//go:generate go run ../../internal/cmd/meta_gen/main.go
//go:generate mockgen -package qingstor -destination mock_test.go github.com/yunify/qingstor-sdk-go/v3/interface Service,Bucket
type Client struct {
	bucket iface.Bucket

	segments    map[string]*segment.Segment
	segmentLock sync.RWMutex
}

// newClient will create a new client.
func newClient(bucket iface.Bucket) *Client {
	return &Client{
		bucket:   bucket,
		segments: make(map[string]*segment.Segment),
	}
}

// Metadata implements Storager.Metadata
func (c *Client) Metadata() (m types.Metadata, err error) {
	errorMessage := "qingstor Metadata: %w"

	output, err := c.bucket.GetStatistics()
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, err)
	}

	m = make(types.Metadata)
	if output.Name != nil {
		m.SetName(*output.Name)
	}
	if output.Location != nil {
		m.SetLocation(*output.Location)
	}
	if output.Size != nil {
		m.SetSize(*output.Size)
	}
	if output.Count != nil {
		m.SetCount(*output.Count)
	}
	return m, nil
}

// Stat implements Storager.Stat
func (c *Client) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	errorMessage := "qingstor Stat: %w"

	input := &service.HeadObjectInput{}

	output, err := c.bucket.HeadObject(path, input)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, err)
	}

	// TODO: Add dir support.

	o = &types.Object{
		Name:     path,
		Type:     types.ObjectTypeFile,
		Metadata: make(types.Metadata),
	}

	if output.ContentType != nil {
		o.SetType(service.StringValue(output.ContentType))
	}
	if output.ContentLength != nil {
		o.SetSize(*output.ContentLength)
	}
	if output.ETag != nil {
		o.SetChecksum(service.StringValue(output.ETag))
	}
	if output.XQSStorageClass != nil {
		o.SetStorageClass(service.StringValue(output.XQSStorageClass))
	}
	if output.LastModified != nil {
		o.SetUpdatedAt(service.TimeValue(output.LastModified))
	}
	return o, nil
}

// Delete implements Storager.Delete
func (c *Client) Delete(path string, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor Delete: %w"

	// TODO: support delete dir.

	_, err = c.bucket.DeleteObject(path)
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}

// Copy implements Storager.Copy
func (c *Client) Copy(src, dst string, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor Copy: %w"

	_, err = c.bucket.PutObject(dst, &service.PutObjectInput{
		XQSCopySource: &src,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}

// Move implements Storager.Move
func (c *Client) Move(src, dst string, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor Move: %w"

	_, err = c.bucket.PutObject(dst, &service.PutObjectInput{
		XQSMoveSource: &src,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}

// Reach implements Storager.Reach
func (c *Client) Reach(path string, pairs ...*types.Pair) (url string, err error) {
	errorMessage := "qingstor Reach: %w"

	opt, err := parseStoragePairReach(pairs...)
	if err != nil {
		return "", fmt.Errorf(errorMessage, err)
	}

	// FIXME: sdk should export GetObjectRequest as interface too?
	bucket := c.bucket.(*service.Bucket)

	r, _, err := bucket.GetObjectRequest(path, nil)
	if err != nil {
		err = handleQingStorError(err)
		return "", fmt.Errorf(errorMessage, err)
	}
	if err = r.Build(); err != nil {
		err = handleQingStorError(err)
		return "", fmt.Errorf(errorMessage, err)
	}

	expire := 3600
	if opt.HasExpire {
		expire = opt.Expire
	}
	if err = r.SignQuery(expire); err != nil {
		err = handleQingStorError(err)
		return "", fmt.Errorf(errorMessage, err)
	}
	return r.HTTPRequest.URL.String(), nil
}

// CreateDir implements Storager.CreateDir
func (c *Client) CreateDir(path string, option ...*types.Pair) (err error) {
	panic("not supported")
}

// ListDir implements Storager.ListDir
func (c *Client) ListDir(path string, pairs ...*types.Pair) (it iterator.ObjectIterator) {
	errorMessage := "qingstor ListDir: %w"

	opt, _ := parseStoragePairListDir(pairs...)

	marker := ""
	limit := 200
	delimiter := "/"
	if opt.HasRecursive && opt.Recursive {
		delimiter = ""
	}

	var output *service.ListObjectsOutput
	var err error

	fn := iterator.NextObjectFunc(func(objects *[]*types.Object) error {
		idx := 0
		buf := make([]*types.Object, limit)

		output, err = c.bucket.ListObjects(&service.ListObjectsInput{
			Limit:     &limit,
			Marker:    &marker,
			Prefix:    &path,
			Delimiter: &delimiter,
		})
		if err != nil {
			err = handleQingStorError(err)
			return fmt.Errorf(errorMessage, err)
		}

		for _, v := range output.CommonPrefixes {
			o := &types.Object{
				Name:     *v,
				Type:     types.ObjectTypeDir,
				Metadata: make(types.Metadata),
			}

			buf[idx] = o
			idx++
		}

		for _, v := range output.Keys {
			o := &types.Object{
				Name:     *v.Key,
				Metadata: make(types.Metadata),
			}

			// If Key end with delimiter or key's MimeType == DirectoryMIMEType,
			// we should treat this key as a Dir Object.
			if (delimiter != "" && strings.HasSuffix(*v.Key, delimiter)) ||
				service.StringValue(v.MimeType) == DirectoryMIMEType {
				o.Type = types.ObjectTypeDir
			} else {
				o.Type = types.ObjectTypeFile
			}

			if v.MimeType != nil {
				o.SetType(service.StringValue(v.MimeType))
			}
			if v.StorageClass != nil {
				o.SetStorageClass(service.StringValue(v.StorageClass))
			}
			if v.Etag != nil {
				o.SetChecksum(service.StringValue(v.Etag))
			}
			if v.Size != nil {
				o.SetSize(service.Int64Value(v.Size))
			}
			if v.Modified != nil {
				o.SetUpdatedAt(time.Unix(int64(service.IntValue(v.Modified)), 0))
			}

			buf[idx] = o
			idx++
		}

		// Set input objects
		*objects = buf[:idx]

		marker = convert.StringValue(output.NextMarker)
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

	it = iterator.NewObjectIterator(fn)
	return
}

// Read implements Storager.Read
func (c *Client) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	errorMessage := "qingstor ReadFile: %w"

	input := &service.GetObjectInput{}

	output, err := c.bucket.GetObject(path, input)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, err)
	}
	return output.Body, nil
}

// WriteFile implements Storager.WriteFile
func (c *Client) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor WriteFile for id %s: %w"

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, path, err)
	}

	input := &service.PutObjectInput{
		ContentLength: &opt.Size,
		Body:          r,
	}
	if opt.HasChecksum {
		input.ContentMD5 = &opt.Checksum
	}
	if opt.HasStorageClass {
		input.XQSStorageClass = &opt.StorageClass
	}

	_, err = c.bucket.PutObject(path, input)
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, path, err)
	}
	return nil
}

// ListSegments implements Storager.ListSegments
func (c *Client) ListSegments(path string, pairs ...*types.Pair) (it iterator.SegmentIterator) {
	errorMessage := "qingstor ListSegments: %w"

	keyMarker := ""
	uploadIDMarker := ""
	limit := 200

	var output *service.ListMultipartUploadsOutput
	var err error

	fn := iterator.NextSegmentFunc(func(segments *[]*segment.Segment) error {
		idx := 0
		buf := make([]*segment.Segment, limit)

		output, err = c.bucket.ListMultipartUploads(&service.ListMultipartUploadsInput{
			KeyMarker:      &keyMarker,
			Limit:          &limit,
			Prefix:         &path,
			UploadIDMarker: &uploadIDMarker,
		})
		if err != nil {
			err = handleQingStorError(err)
			return fmt.Errorf(errorMessage, err)
		}

		for _, v := range output.Uploads {
			s := segment.NewSegment(*v.Key, *v.UploadID)

			buf[idx] = s
			idx++

			c.segmentLock.Lock()
			// Update client's segments.
			c.segments[s.ID] = s
			c.segmentLock.Unlock()
		}

		// Set input objects
		*segments = buf[:idx]

		keyMarker = convert.StringValue(output.NextKeyMarker)
		uploadIDMarker = convert.StringValue(output.NextUploadIDMarker)
		if keyMarker == "" && uploadIDMarker == "" {
			return iterator.ErrDone
		}
		if output.HasMore != nil && !*output.HasMore {
			return iterator.ErrDone
		}
		return nil
	})

	it = iterator.NewSegmentIterator(fn)
	return
}

// InitSegment implements Storager.InitSegment
func (c *Client) InitSegment(path string, pairs ...*types.Pair) (id string, err error) {
	errorMessage := "qingstor InitSegment for id %s: %w"

	input := &service.InitiateMultipartUploadInput{}

	output, err := c.bucket.InitiateMultipartUpload(path, input)
	if err != nil {
		err = handleQingStorError(err)
		return "", fmt.Errorf(errorMessage, path, err)
	}

	id = *output.UploadID

	c.segmentLock.Lock()
	c.segments[id] = segment.NewSegment(path, id)
	c.segmentLock.Unlock()
	return
}

// WriteSegment implements Storager.WriteSegment
func (c *Client) WriteSegment(id string, offset, size int64, r io.Reader, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor WriteSegment for id %s: %w"

	c.segmentLock.RLock()
	s, ok := c.segments[id]
	if !ok {
		return fmt.Errorf(errorMessage, id, segment.ErrSegmentNotInitiated)
	}
	c.segmentLock.RUnlock()

	p := &segment.Part{
		Offset: offset,
		Size:   size,
	}

	partNumber, err := s.InsertPart(p)
	if err != nil {
		return fmt.Errorf(errorMessage, id, err)
	}

	_, err = c.bucket.UploadMultipart(s.Path, &service.UploadMultipartInput{
		PartNumber:    &partNumber,
		UploadID:      &s.ID,
		ContentLength: &size,
		Body:          r,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, id, err)
	}
	return
}

// CompleteSegment implements Storager.CompleteSegment
func (c *Client) CompleteSegment(id string, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor CompleteSegment for id %s: %w"

	c.segmentLock.RLock()
	s, ok := c.segments[id]
	if !ok {
		return fmt.Errorf(errorMessage, id, segment.ErrSegmentNotInitiated)
	}
	c.segmentLock.RUnlock()

	err = s.ValidateParts()
	if err != nil {
		return fmt.Errorf(errorMessage, id, err)
	}

	parts := s.SortedParts()
	objectParts := make([]*service.ObjectPartType, len(parts))
	for k, v := range parts {
		objectParts[k] = &service.ObjectPartType{
			PartNumber: &v.Index,
			Size:       &v.Size,
		}
	}

	_, err = c.bucket.CompleteMultipartUpload(s.Path, &service.CompleteMultipartUploadInput{
		UploadID:    &s.ID,
		ObjectParts: objectParts,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, id, err)
	}

	c.segmentLock.Lock()
	delete(c.segments, id)
	c.segmentLock.Unlock()
	return
}

// AbortSegment implements Storager.AbortSegment
func (c *Client) AbortSegment(id string, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor AbortSegment for id %s: %w"

	c.segmentLock.RLock()
	s, ok := c.segments[id]
	if !ok {
		return fmt.Errorf(errorMessage, id, segment.ErrSegmentNotInitiated)
	}
	c.segmentLock.RUnlock()

	_, err = c.bucket.AbortMultipartUpload(s.Path, &service.AbortMultipartUploadInput{
		UploadID: &s.ID,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, id, err)
	}

	c.segmentLock.Lock()
	delete(c.segments, id)
	c.segmentLock.Unlock()
	return
}
