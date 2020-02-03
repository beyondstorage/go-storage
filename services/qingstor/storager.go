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
type Storage struct {
	bucket     iface.Bucket
	config     *qsconfig.Config
	properties *service.Properties

	// options for this storager.
	workDir string // workDir dir for all operation.

	segments    map[string]*segment.Segment
	segmentLock sync.RWMutex
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager qingstor {Name: %s, Location: %s, WorkDir: %s}",
		*s.properties.BucketName, *s.properties.Zone, "/"+s.workDir,
	)
}

// Metadata implements Storager.Metadata
func (s *Storage) Metadata(pairs ...*types.Pair) (m metadata.StorageMeta, err error) {
	m = metadata.NewStorageMeta()
	m.Name = *s.properties.BucketName
	m.WorkDir = s.workDir
	m.SetLocation(*s.properties.Zone)
	return m, nil
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
				ID:         *v,
				Name:       s.getRelPath(*v),
				Type:       types.ObjectTypeDir,
				ObjectMeta: metadata.NewObjectMeta(),
			}

			if opt.HasDirFunc {
				opt.DirFunc(o)
			}
		}

		for _, v := range output.Keys {
			o := &types.Object{
				ID:         *v.Key,
				Name:       s.getRelPath(*v.Key),
				Size:       service.Int64Value(v.Size),
				UpdatedAt:  convertUnixTimestampToTime(service.IntValue(v.Modified)),
				ObjectMeta: metadata.NewObjectMeta(),
			}

			if v.MimeType != nil {
				o.SetContentType(service.StringValue(v.MimeType))
			}
			if v.StorageClass != nil {
				storageClass, err := formatStorageClass(service.StringValue(v.StorageClass))
				if err != nil {
					return fmt.Errorf(errorMessage, s, path, err)
				}
				o.SetStorageClass(storageClass)
			}
			if v.Etag != nil {
				o.SetETag(service.StringValue(v.Etag))
			}

			// If key's content type == DirectoryContentType,
			// we should treat this key as a Dir ObjectMeta.
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
		storageClass, err := parseStorageClass(opt.StorageClass)
		if err != nil {
			return fmt.Errorf(errorMessage, s, path, err)
		}
		input.XQSStorageClass = service.String(storageClass)
	}

	rp := s.getAbsPath(path)

	_, err = s.bucket.PutObject(rp, input)
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return nil
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
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       service.Int64Value(output.ContentLength),
		UpdatedAt:  service.TimeValue(output.LastModified),
		ObjectMeta: metadata.NewObjectMeta(),
	}

	if output.ContentType != nil {
		o.SetContentType(service.StringValue(output.ContentType))
	}
	if output.ETag != nil {
		o.SetETag(service.StringValue(output.ETag))
	}

	storageClass, err := formatStorageClass(service.StringValue(output.XQSStorageClass))
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}
	o.SetStorageClass(storageClass)

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

	expire := opt.Expire
	if err = r.SignQuery(expire); err != nil {
		err = handleQingStorError(err)
		return "", fmt.Errorf(errorMessage, s, path, err)
	}
	return r.HTTPRequest.URL.String(), nil
}

// Statistical implements Storager.Statistical
func (s *Storage) Statistical(pairs ...*types.Pair) (m metadata.StorageStatistic, err error) {
	const errorMessage = "%s Statistical: %w"

	m = metadata.NewStorageStatistic()

	output, err := s.bucket.GetStatistics()
	if err != nil {
		err = handleQingStorError(err)
		return m, fmt.Errorf(errorMessage, s, err)
	}

	if output.Size != nil {
		m.SetSize(*output.Size)
	}
	if output.Count != nil {
		m.SetCount(*output.Count)
	}
	return m, nil
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

func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
}

func (s *Storage) getRelPath(path string) string {
	return strings.TrimPrefix(path, s.workDir+"/")
}
