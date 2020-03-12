package qingstor

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/pengsrc/go-shared/convert"
	qsconfig "github.com/yunify/qingstor-sdk-go/v3/config"
	iface "github.com/yunify/qingstor-sdk-go/v3/interface"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/services"
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
	loose   bool

	segments    map[string]*segment.Segment
	segmentLock sync.RWMutex
}

// String implements Storager.String
func (s *Storage) String() string {
	// qingstor work dir should start and end with "/"
	return fmt.Sprintf(
		"Storager qingstor {Name: %s, Location: %s, WorkDir: %s}",
		*s.properties.BucketName, *s.properties.Zone, s.workDir,
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
	defer func() {
		err = s.formatError("list", err, path)
	}()

	opt, _ := parseStoragePairList(pairs...)

	marker := ""
	delimiter := ""
	limit := 200

	rp := s.getAbsPath(path)

	if !opt.HasObjectFunc {
		delimiter = "/"
	}

	var output *service.ListObjectsOutput
	for {
		output, err = s.bucket.ListObjects(&service.ListObjectsInput{
			Limit:     &limit,
			Marker:    &marker,
			Prefix:    &rp,
			Delimiter: &delimiter,
		})
		if err != nil {
			return
		}

		if opt.HasDirFunc {
			for _, v := range output.CommonPrefixes {
				o := &types.Object{
					ID:         *v,
					Name:       s.getRelPath(*v),
					Type:       types.ObjectTypeDir,
					ObjectMeta: metadata.NewObjectMeta(),
				}

				opt.DirFunc(o)
			}
		}

		if opt.HasObjectFunc || opt.HasFileFunc {
			for _, v := range output.Keys {
				o := &types.Object{
					ID:         *v.Key,
					Name:       s.getRelPath(*v.Key),
					Type:       types.ObjectTypeFile,
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
						return err
					}
					o.SetStorageClass(storageClass)
				}
				if v.Etag != nil {
					o.SetETag(service.StringValue(v.Etag))
				}

				if opt.HasObjectFunc {
					opt.ObjectFunc(o)
				}
				if opt.HasFileFunc {
					opt.FileFunc(o)
				}
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
	defer func() {
		err = s.formatError("read", err, path)
	}()

	opt, err := parseStoragePairRead(pairs...)
	if err != nil {
		return
	}

	input := &service.GetObjectInput{}

	rp := s.getAbsPath(path)

	output, err := s.bucket.GetObject(rp, input)
	if err != nil {
		return
	}

	r = output.Body
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReadCloser(r, opt.ReadCallbackFunc)
	}
	return r, nil
}

// WriteFile implements Storager.WriteFile
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("write", err, path)
	}()

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return
	}

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
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
			return err
		}
		input.XQSStorageClass = service.String(storageClass)
	}

	rp := s.getAbsPath(path)

	_, err = s.bucket.PutObject(rp, input)
	if err != nil {
		return
	}
	return nil
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	defer func() {
		err = s.formatError("stat", err, path)
	}()

	input := &service.HeadObjectInput{}

	rp := s.getAbsPath(path)

	output, err := s.bucket.HeadObject(rp, input)
	if err != nil {
		return
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
		return
	}
	o.SetStorageClass(storageClass)

	return o, nil
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("delete", err, path)
	}()

	rp := s.getAbsPath(path)

	_, err = s.bucket.DeleteObject(rp)
	if err != nil {
		return
	}
	return nil
}

// Copy implements Storager.Copy
func (s *Storage) Copy(src, dst string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("copy", err, src, dst)
	}()

	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	_, err = s.bucket.PutObject(rd, &service.PutObjectInput{
		XQSCopySource: &rs,
	})
	if err != nil {
		return
	}
	return nil
}

// Move implements Storager.Move
func (s *Storage) Move(src, dst string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("move", err, src, dst)
	}()

	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	_, err = s.bucket.PutObject(rd, &service.PutObjectInput{
		XQSMoveSource: &rs,
	})
	if err != nil {
		return
	}
	return nil
}

// Reach implements Storager.Reach
func (s *Storage) Reach(path string, pairs ...*types.Pair) (url string, err error) {
	defer func() {
		err = s.formatError("reach", err, path)
	}()

	opt, err := parseStoragePairReach(pairs...)
	if err != nil {
		return
	}

	// FIXME: sdk should export GetObjectRequest as interface too?
	bucket := s.bucket.(*service.Bucket)

	rp := s.getAbsPath(path)

	r, _, err := bucket.GetObjectRequest(rp, nil)
	if err != nil {
		return
	}
	if err = r.Build(); err != nil {
		return
	}

	expire := opt.Expire
	if err = r.SignQuery(expire); err != nil {
		return
	}
	return r.HTTPRequest.URL.String(), nil
}

// Statistical implements Storager.Statistical
func (s *Storage) Statistical(pairs ...*types.Pair) (m metadata.StorageStatistic, err error) {
	defer func() {
		err = s.formatError("statistical", err)
	}()

	m = metadata.NewStorageStatistic()

	output, err := s.bucket.GetStatistics()
	if err != nil {
		return
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
	defer func() {
		err = s.formatError("list segments", err, path)
	}()

	opt, err := parseStoragePairListSegments(pairs...)
	if err != nil {
		return
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
			return
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
	defer func() {
		err = s.formatError("init segments", err, path)
	}()

	opt, err := parseStoragePairInitSegment(pairs...)
	if err != nil {
		return
	}

	input := &service.InitiateMultipartUploadInput{}

	rp := s.getAbsPath(path)

	output, err := s.bucket.InitiateMultipartUpload(rp, input)
	if err != nil {
		return
	}

	id = *output.UploadID

	s.segmentLock.Lock()
	s.segments[id] = segment.NewSegment(path, id, opt.PartSize)
	s.segmentLock.Unlock()
	return
}

// WriteSegment implements Storager.WriteSegment
func (s *Storage) WriteSegment(id string, offset, size int64, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("write segment", err, id)
	}()

	opt, err := parseStoragePairWriteSegment(pairs...)
	if err != nil {
		return
	}

	s.segmentLock.RLock()
	seg, ok := s.segments[id]
	if !ok {
		err = segment.ErrSegmentNotInitiated
		return
	}
	s.segmentLock.RUnlock()

	p, err := seg.InsertPart(offset, size)
	if err != nil {
		return
	}

	rp := s.getAbsPath(seg.Path)

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	_, err = s.bucket.UploadMultipart(rp, &service.UploadMultipartInput{
		PartNumber:    &p.Index,
		UploadID:      &seg.ID,
		ContentLength: &size,
		Body:          r,
	})
	if err != nil {
		return
	}
	return
}

// CompleteSegment implements Storager.CompleteSegment
func (s *Storage) CompleteSegment(id string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("complete segment", err, id)
	}()

	s.segmentLock.RLock()
	seg, ok := s.segments[id]
	if !ok {
		err = segment.ErrSegmentNotInitiated
		return
	}
	s.segmentLock.RUnlock()

	err = seg.ValidateParts()
	if err != nil {
		return
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
		return
	}

	s.segmentLock.Lock()
	delete(s.segments, id)
	s.segmentLock.Unlock()
	return
}

// AbortSegment implements Storager.AbortSegment
func (s *Storage) AbortSegment(id string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("abort segment", err, id)
	}()

	s.segmentLock.RLock()
	seg, ok := s.segments[id]
	if !ok {
		err = segment.ErrSegmentNotInitiated
		return
	}
	s.segmentLock.RUnlock()

	rp := s.getAbsPath(seg.Path)

	_, err = s.bucket.AbortMultipartUpload(rp, &service.AbortMultipartUploadInput{
		UploadID: &seg.ID,
	})
	if err != nil {
		return
	}

	s.segmentLock.Lock()
	delete(s.segments, id)
	s.segmentLock.Unlock()
	return
}

func (s *Storage) getAbsPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/") // qsPrefix should not start with "/"
	return prefix + path                         // qs abs path is the qsPrefix add path (path is not start with "/")
}

func (s *Storage) getRelPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/") // qsPrefix should not start with "/"
	return strings.TrimPrefix(path, prefix)      // qs rel path is the path trimmed qsPrefix
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	if err == nil {
		return nil
	}

	if s.loose && errors.Is(err, services.ErrCapabilityInsufficient) {
		return nil
	}

	return &services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}
