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
		"Storager qingstor {Name: %s, Location: %s, WorkDir: %s}",
		*s.properties.BucketName, *s.properties.Zone, "/"+s.workDir,
	)
}

// Metadata implements Storager.Metadata
func (s *Storage) Metadata() (m metadata.Storage, err error) {
	m = metadata.Storage{
		Name:     *s.properties.BucketName,
		WorkDir:  s.workDir,
		Metadata: make(metadata.Metadata),
	}
	m.SetLocation(*s.properties.Zone)
	return m, nil
}

// Statistical implements Storager.Statistical
func (s *Storage) Statistical() (m metadata.Metadata, err error) {
	const errorMessage = "%s Statistical: %w"

	output, err := s.bucket.GetStatistics()
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, s, err)
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
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	const errorMessage = "%s Stat [%s]: %w"

	input := &service.HeadObjectInput{}

	rp := s.getAbsPath(path)

	output, err := s.bucket.HeadObject(rp, input)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, s, path, err)
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
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	rp := s.getAbsPath(path)

	_, err = s.bucket.DeleteObject(rp)
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}

// Copy implements Storager.Copy
func (s *Storage) Copy(src, dst string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Copy from [%s] to [%s]: %w"

	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	_, err = s.bucket.PutObject(rd, &service.PutObjectInput{
		XQSCopySource: &rs,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, src, dst, err)
	}
	return nil
}

// Move implements Storager.Move
func (s *Storage) Move(src, dst string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Move from [%s] to [%s]: %w"

	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	_, err = s.bucket.PutObject(rd, &service.PutObjectInput{
		XQSMoveSource: &rs,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, src, dst, err)
	}
	return nil
}

// Reach implements Storager.Reach
func (s *Storage) Reach(path string, pairs ...*types.Pair) (url string, err error) {
	const errorMessage = "%s Reach [%s]: %w"

	opt, err := parseStoragePairReach(pairs...)
	if err != nil {
		return "", fmt.Errorf(errorMessage, s, path, err)
	}

	// FIXME: sdk should export GetObjectRequest as interface too?
	bucket := s.bucket.(*service.Bucket)

	rp := s.getAbsPath(path)

	r, _, err := bucket.GetObjectRequest(rp, nil)
	if err != nil {
		err = handleQingStorError(err)
		return "", fmt.Errorf(errorMessage, s, path, err)
	}
	if err = r.Build(); err != nil {
		err = handleQingStorError(err)
		return "", fmt.Errorf(errorMessage, s, path, err)
	}

	expire := 3600
	if opt.HasExpire {
		expire = opt.Expire
	}
	if err = r.SignQuery(expire); err != nil {
		err = handleQingStorError(err)
		return "", fmt.Errorf(errorMessage, s, path, err)
	}
	return r.HTTPRequest.URL.String(), nil
}

// List implements Storager.List
func (s *Storage) List(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s List [%s]: %w"

	opt, _ := parseStoragePairList(pairs...)

	marker := ""
	limit := 200

	rp := s.getAbsPath(path)

	var output *service.ListObjectsOutput
	for {
		output, err = s.bucket.ListObjects(&service.ListObjectsInput{
			Limit:  &limit,
			Marker: &marker,
			Prefix: &rp,
		})
		if err != nil {
			err = handleQingStorError(err)
			return fmt.Errorf(errorMessage, s, path, err)
		}

		for _, v := range output.CommonPrefixes {
			o := &types.Object{
				Name:     s.getRelPath(*v),
				Type:     types.ObjectTypeDir,
				Metadata: make(metadata.Metadata),
			}

			if opt.HasDirFunc {
				opt.DirFunc(o)
			}
		}

		for _, v := range output.Keys {
			o := &types.Object{
				Name:      s.getRelPath(*v.Key),
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
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	const errorMessage = "%s Read [%s]: %w"

	input := &service.GetObjectInput{}

	rp := s.getAbsPath(path)

	output, err := s.bucket.GetObject(rp, input)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}
	return output.Body, nil
}

// WriteFile implements Storager.WriteFile
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Write [%s]: %w"

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
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

	rp := s.getAbsPath(path)

	_, err = s.bucket.PutObject(rp, input)
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
}

// ListSegments implements Storager.ListSegments
func (s *Storage) ListSegments(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s ListSegments [%s]: %w"

	opt, err := parseStoragePairListSegments(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	keyMarker := ""
	uploadIDMarker := ""
	limit := 200

	rp := s.getAbsPath(path)

	var output *service.ListMultipartUploadsOutput
	for {
		output, err = s.bucket.ListMultipartUploads(&service.ListMultipartUploadsInput{
			KeyMarker:      &keyMarker,
			Limit:          &limit,
			Prefix:         &rp,
			UploadIDMarker: &uploadIDMarker,
		})
		if err != nil {
			err = handleQingStorError(err)
			return fmt.Errorf(errorMessage, s, path, err)
		}

		for _, v := range output.Uploads {
			seg := segment.NewSegment(*v.Key, *v.UploadID, 0)

			if opt.HasSegmentFunc {
				opt.SegmentFunc(seg)
			}

			s.segmentLock.Lock()
			// Update client's segments.
			s.segments[seg.ID] = seg
			s.segmentLock.Unlock()
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
func (s *Storage) InitSegment(path string, pairs ...*types.Pair) (id string, err error) {
	const errorMessage = "%s InitSegment [%s]: %w"

	opt, err := parseStoragePairInitSegment(pairs...)
	if err != nil {
		return "", fmt.Errorf(errorMessage, s, path, err)
	}

	input := &service.InitiateMultipartUploadInput{}

	rp := s.getAbsPath(path)

	output, err := s.bucket.InitiateMultipartUpload(rp, input)
	if err != nil {
		err = handleQingStorError(err)
		return "", fmt.Errorf(errorMessage, s, path, err)
	}

	id = *output.UploadID

	s.segmentLock.Lock()
	s.segments[id] = segment.NewSegment(path, id, opt.PartSize)
	s.segmentLock.Unlock()
	return
}

// WriteSegment implements Storager.WriteSegment
func (s *Storage) WriteSegment(id string, offset, size int64, r io.Reader, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s WriteSegment [%s]: %w"

	s.segmentLock.RLock()
	seg, ok := s.segments[id]
	if !ok {
		return fmt.Errorf(errorMessage, s, id, segment.ErrSegmentNotInitiated)
	}
	s.segmentLock.RUnlock()

	p, err := seg.InsertPart(offset, size)
	if err != nil {
		return fmt.Errorf(errorMessage, s, id, err)
	}

	rp := s.getAbsPath(seg.Path)

	_, err = s.bucket.UploadMultipart(rp, &service.UploadMultipartInput{
		PartNumber:    &p.Index,
		UploadID:      &seg.ID,
		ContentLength: &size,
		Body:          r,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, id, err)
	}
	return
}

// CompleteSegment implements Storager.CompleteSegment
func (s *Storage) CompleteSegment(id string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s CompleteSegment [%s]: %w"

	s.segmentLock.RLock()
	seg, ok := s.segments[id]
	if !ok {
		return fmt.Errorf(errorMessage, s, id, segment.ErrSegmentNotInitiated)
	}
	s.segmentLock.RUnlock()

	err = seg.ValidateParts()
	if err != nil {
		return fmt.Errorf(errorMessage, s, id, err)
	}

	parts := seg.SortedParts()
	objectParts := make([]*service.ObjectPartType, 0, len(parts))
	for k, v := range parts {
		k := k
		objectParts = append(objectParts, &service.ObjectPartType{
			PartNumber: &k,
			Size:       &v.Size,
		})
	}

	rp := s.getAbsPath(seg.Path)

	_, err = s.bucket.CompleteMultipartUpload(rp, &service.CompleteMultipartUploadInput{
		UploadID:    &seg.ID,
		ObjectParts: objectParts,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, id, err)
	}

	s.segmentLock.Lock()
	delete(s.segments, id)
	s.segmentLock.Unlock()
	return
}

// AbortSegment implements Storager.AbortSegment
func (s *Storage) AbortSegment(id string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s AbortSegment [%s]: %w"

	s.segmentLock.RLock()
	seg, ok := s.segments[id]
	if !ok {
		return fmt.Errorf(errorMessage, s, id, segment.ErrSegmentNotInitiated)
	}
	s.segmentLock.RUnlock()

	rp := s.getAbsPath(seg.Path)

	_, err = s.bucket.AbortMultipartUpload(rp, &service.AbortMultipartUploadInput{
		UploadID: &seg.ID,
	})
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, id, err)
	}

	s.segmentLock.Lock()
	delete(s.segments, id)
	s.segmentLock.Unlock()
	return
}
