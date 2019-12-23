package qingstor

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/pengsrc/go-shared/convert"
	qsconfig "github.com/yunify/qingstor-sdk-go/v3/config"
	iface "github.com/yunify/qingstor-sdk-go/v3/interface"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the qingstor object storage client.
//
//go:generate ../../internal/bin/meta
//go:generate mockgen -package qingstor -destination mock_test.go github.com/yunify/qingstor-sdk-go/v3/interface Service,Bucket
type Storage struct {
	bucket     iface.Bucket
	config     *qsconfig.Config
	properties *service.Properties

	// options for this storager.
	workDir string // workDir dir for all operation.

	segments    map[string]*segment.Segment
	segmentLock sync.RWMutex
}

// newStorage will create a new client.
func newStorage(bucket *service.Bucket) (*Storage, error) {
	c := &Storage{
		bucket:     bucket,
		config:     bucket.Config,
		properties: bucket.Properties,
		segments:   make(map[string]*segment.Segment),
	}
	return c, nil
}

// Init implements Storager.Init
func (c *Storage) Init(pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor Init: %w"

	opt, err := parseStoragePairInit(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	if opt.HasWorkDir {
		// TODO: we should validate workDir
		c.workDir = strings.TrimLeft(opt.WorkDir, "/")
	}

	return nil
}

// String implements Storager.String
func (c *Storage) String() string {
	return fmt.Sprintf(
		"qingstor Storager {Bucket: %s, Location: %s, WorkDir: %s}",
		*c.properties.BucketName, *c.properties.Zone, "/"+c.workDir,
	)
}

// Metadata implements Storager.Metadata
func (c *Storage) Metadata() (m metadata.Storage, err error) {
	m = metadata.Storage{
		Name:     *c.properties.BucketName,
		WorkDir:  c.workDir,
		Metadata: make(metadata.Metadata),
	}
	m.SetLocation(*c.properties.Zone)
	return m, nil
}

func (c *Storage) Statistical() (m metadata.Metadata, err error) {
	errorMessage := "qingstor Statistical: %w"

	output, err := c.bucket.GetStatistics()
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, err)
	}

	m = make(metadata.Metadata)
	if output.Size != nil {
		m.SetSize(*output.Size)
	}
	if output.Count != nil {
		m.SetCount(*output.Count)
	}
	return m, nil
}

// Stat implements Storager.Stat
func (c *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	errorMessage := "qingstor Stat: %w"

	input := &service.HeadObjectInput{}

	rp := c.getAbsPath(path)

	output, err := c.bucket.HeadObject(rp, input)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, err)
	}

	// TODO: Add dir support.

	o = &types.Object{
		Name:      path,
		Type:      types.ObjectTypeFile,
		Size:      service.Int64Value(output.ContentLength),
		UpdatedAt: service.TimeValue(output.LastModified),
		Metadata:  make(metadata.Metadata),
	}

	if output.ContentType != nil {
		o.SetType(service.StringValue(output.ContentType))
	}
	if output.ETag != nil {
		o.SetChecksum(service.StringValue(output.ETag))
	}
	if output.XQSStorageClass != nil {
		o.SetClass(service.StringValue(output.XQSStorageClass))
	}
	return o, nil
}

// Delete implements Storager.Delete
func (c *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor Delete: %w"

	rp := c.getAbsPath(path)

	_, err = c.bucket.DeleteObject(rp)
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}

// Copy implements Storager.Copy
func (c *Storage) Copy(src, dst string, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor Copy: %w"

	rs := c.getAbsPath(src)
	rd := c.getAbsPath(dst)

	_, err = c.bucket.PutObject(rd, &service.PutObjectInput{
		XQSCopySource: &rs,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}

// Move implements Storager.Move
func (c *Storage) Move(src, dst string, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor Move: %w"

	rs := c.getAbsPath(src)
	rd := c.getAbsPath(dst)

	_, err = c.bucket.PutObject(rd, &service.PutObjectInput{
		XQSMoveSource: &rs,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}

// Reach implements Storager.Reach
func (c *Storage) Reach(path string, pairs ...*types.Pair) (url string, err error) {
	errorMessage := "qingstor Reach: %w"

	opt, err := parseStoragePairReach(pairs...)
	if err != nil {
		return "", fmt.Errorf(errorMessage, err)
	}

	// FIXME: sdk should export GetObjectRequest as interface too?
	bucket := c.bucket.(*service.Bucket)

	rp := c.getAbsPath(path)

	r, _, err := bucket.GetObjectRequest(rp, nil)
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

// ListDir implements Storager.ListDir
func (c *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor ListDir: %w"

	opt, _ := parseStoragePairListDir(pairs...)

	marker := ""
	limit := 200

	rp := c.getAbsPath(path)

	var output *service.ListObjectsOutput
	for {
		output, err = c.bucket.ListObjects(&service.ListObjectsInput{
			Limit:  &limit,
			Marker: &marker,
			Prefix: &rp,
		})
		if err != nil {
			err = handleQingStorError(err)
			return fmt.Errorf(errorMessage, err)
		}

		for _, v := range output.CommonPrefixes {
			o := &types.Object{
				Name:     c.getRelPath(*v),
				Type:     types.ObjectTypeDir,
				Metadata: make(metadata.Metadata),
			}

			if opt.HasDirFunc {
				opt.DirFunc(o)
			}
		}

		for _, v := range output.Keys {
			o := &types.Object{
				Name:      c.getRelPath(*v.Key),
				Size:      service.Int64Value(v.Size),
				UpdatedAt: convertUnixTimestampToTime(service.IntValue(v.Modified)),
				Metadata:  make(metadata.Metadata),
			}

			if v.MimeType != nil {
				o.SetType(service.StringValue(v.MimeType))
			}
			if v.StorageClass != nil {
				o.SetClass(service.StringValue(v.StorageClass))
			}
			if v.Etag != nil {
				o.SetChecksum(service.StringValue(v.Etag))
			}

			// If key's content type == DirectoryContentType,
			// we should treat this key as a Dir Object.
			if service.StringValue(v.MimeType) == DirectoryContentType {
				o.Type = types.ObjectTypeDir

				if opt.HasDirFunc {
					opt.DirFunc(o)
				}
				continue
			}

			o.Type = types.ObjectTypeFile
			if opt.HasFileFunc {
				opt.FileFunc(o)
			}
		}

		marker = convert.StringValue(output.NextMarker)
		if marker == "" {
			break
		}
		if output.HasMore != nil && !*output.HasMore {
			break
		}
		if len(output.Keys) == 0 {
			break
		}
	}
	return
}

// Read implements Storager.Read
func (c *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	errorMessage := "qingstor ReadFile: %w"

	input := &service.GetObjectInput{}

	rp := c.getAbsPath(path)

	output, err := c.bucket.GetObject(rp, input)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, err)
	}
	return output.Body, nil
}

// WriteFile implements Storager.WriteFile
func (c *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
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

	rp := c.getAbsPath(path)

	_, err = c.bucket.PutObject(rp, input)
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, path, err)
	}
	return nil
}

// ListSegments implements Storager.ListSegments
func (c *Storage) ListSegments(path string, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor ListSegments: %w"

	opt, err := parseStoragePairListSegments(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	keyMarker := ""
	uploadIDMarker := ""
	limit := 200

	rp := c.getAbsPath(path)

	var output *service.ListMultipartUploadsOutput
	for {
		output, err = c.bucket.ListMultipartUploads(&service.ListMultipartUploadsInput{
			KeyMarker:      &keyMarker,
			Limit:          &limit,
			Prefix:         &rp,
			UploadIDMarker: &uploadIDMarker,
		})
		if err != nil {
			err = handleQingStorError(err)
			return fmt.Errorf(errorMessage, err)
		}

		for _, v := range output.Uploads {
			s := segment.NewSegment(*v.Key, *v.UploadID, 0)

			if opt.HasSegmentFunc {
				opt.SegmentFunc(s)
			}

			c.segmentLock.Lock()
			// Update client's segments.
			c.segments[s.ID] = s
			c.segmentLock.Unlock()
		}

		keyMarker = convert.StringValue(output.NextKeyMarker)
		uploadIDMarker = convert.StringValue(output.NextUploadIDMarker)
		if keyMarker == "" && uploadIDMarker == "" {
			break
		}
		if output.HasMore != nil && !*output.HasMore {
			break
		}
	}
	return
}

// InitSegment implements Storager.InitSegment
func (c *Storage) InitSegment(path string, pairs ...*types.Pair) (id string, err error) {
	errorMessage := "qingstor InitSegment for id %s: %w"

	opt, err := parseStoragePairInitSegment(pairs...)
	if err != nil {
		return "", fmt.Errorf(errorMessage, path, err)
	}

	input := &service.InitiateMultipartUploadInput{}

	rp := c.getAbsPath(path)

	output, err := c.bucket.InitiateMultipartUpload(rp, input)
	if err != nil {
		err = handleQingStorError(err)
		return "", fmt.Errorf(errorMessage, path, err)
	}

	id = *output.UploadID

	c.segmentLock.Lock()
	c.segments[id] = segment.NewSegment(path, id, opt.PartSize)
	c.segmentLock.Unlock()
	return
}

// WriteSegment implements Storager.WriteSegment
func (c *Storage) WriteSegment(id string, offset, size int64, r io.Reader, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor WriteSegment for id %s: %w"

	c.segmentLock.RLock()
	s, ok := c.segments[id]
	if !ok {
		return fmt.Errorf(errorMessage, id, segment.ErrSegmentNotInitiated)
	}
	c.segmentLock.RUnlock()

	p, err := s.InsertPart(offset, size)
	if err != nil {
		return fmt.Errorf(errorMessage, id, err)
	}

	rp := c.getAbsPath(s.Path)

	_, err = c.bucket.UploadMultipart(rp, &service.UploadMultipartInput{
		PartNumber:    &p.Index,
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
func (c *Storage) CompleteSegment(id string, pairs ...*types.Pair) (err error) {
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
	objectParts := make([]*service.ObjectPartType, 0, len(parts))
	for k, v := range parts {
		k := k
		objectParts = append(objectParts, &service.ObjectPartType{
			PartNumber: &k,
			Size:       &v.Size,
		})
	}

	rp := c.getAbsPath(s.Path)

	_, err = c.bucket.CompleteMultipartUpload(rp, &service.CompleteMultipartUploadInput{
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
func (c *Storage) AbortSegment(id string, pairs ...*types.Pair) (err error) {
	errorMessage := "qingstor AbortSegment for id %s: %w"

	c.segmentLock.RLock()
	s, ok := c.segments[id]
	if !ok {
		return fmt.Errorf(errorMessage, id, segment.ErrSegmentNotInitiated)
	}
	c.segmentLock.RUnlock()

	rp := c.getAbsPath(s.Path)

	_, err = c.bucket.AbortMultipartUpload(rp, &service.AbortMultipartUploadInput{
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
