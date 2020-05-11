package qingstor

import (
	"fmt"
	"io"
	"strings"

	"github.com/pengsrc/go-shared/convert"
	qsconfig "github.com/qingstor/qingstor-sdk-go/v4/config"
	iface "github.com/qingstor/qingstor-sdk-go/v4/interface"
	"github.com/qingstor/qingstor-sdk-go/v4/service"

	"github.com/Xuanwo/storage/pkg/headers"
	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/info"
)

// Storage is the qingstor object storage client.
type Storage struct {
	bucket     iface.Bucket
	config     *qsconfig.Config
	properties *service.Properties

	// options for this storager.
	workDir string // workDir dir for all operation.
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
func (s *Storage) Metadata(pairs ...*types.Pair) (m info.StorageMeta, err error) {
	m = info.NewStorageMeta()
	m.Name = *s.properties.BucketName
	m.WorkDir = s.workDir
	m.SetLocation(*s.properties.Zone)
	return m, nil
}

// ListDir implements DirLister.ListDir
func (s *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpListDir, err, path)
	}()

	opt, _ := s.parsePairListDir(pairs...)

	marker := ""
	delimiter := "/"
	limit := 200

	rp := s.getAbsPath(path)

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
					ObjectMeta: info.NewObjectMeta(),
				}

				opt.DirFunc(o)
			}
		}

		if opt.HasFileFunc {
			for _, v := range output.Keys {
				o, err := s.formatFileObject(v)
				if err != nil {
					return err
				}

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

// ListPrefix implements PrefixLister.ListPrefix
func (s *Storage) ListPrefix(prefix string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpListPrefix, err, prefix)
	}()

	opt, _ := s.parsePairListPrefix(pairs...)

	marker := ""
	limit := 200

	rp := s.getAbsPath(prefix)

	var output *service.ListObjectsOutput
	for {
		output, err = s.bucket.ListObjects(&service.ListObjectsInput{
			Limit:  &limit,
			Marker: &marker,
			Prefix: &rp,
		})
		if err != nil {
			return
		}

		for _, v := range output.Keys {
			o, err := s.formatFileObject(v)
			if err != nil {
				return err
			}

			opt.ObjectFunc(o)
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
		err = s.formatError(services.OpRead, err, path)
	}()

	opt, err := s.parsePairRead(pairs...)
	if err != nil {
		return
	}

	input := &service.GetObjectInput{}

	if opt.HasOffset || opt.HasSize {
		rs := headers.FormatRange(opt.Offset, opt.Size)
		input.Range = &rs
	}

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
		err = s.formatError(services.OpWrite, err, path)
	}()

	opt, err := s.parsePairWrite(pairs...)
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
		input.XQSStorageClass = service.String(opt.StorageClass)
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
		err = s.formatError(services.OpStat, err, path)
	}()

	input := &service.HeadObjectInput{}

	rp := s.getAbsPath(path)

	output, err := s.bucket.HeadObject(rp, input)
	if err != nil {
		return
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       service.Int64Value(output.ContentLength),
		UpdatedAt:  service.TimeValue(output.LastModified),
		ObjectMeta: info.NewObjectMeta(),
	}

	if output.ContentType != nil {
		o.SetContentType(service.StringValue(output.ContentType))
	}
	if output.ETag != nil {
		o.SetETag(service.StringValue(output.ETag))
	}

	if v := service.StringValue(output.XQSStorageClass); v != "" {
		setStorageClass(o.ObjectMeta, v)
	}

	return o, nil
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpDelete, err, path)
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
		err = s.formatError(services.OpCopy, err, src, dst)
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
		err = s.formatError(services.OpMove, err, src, dst)
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
		err = s.formatError(services.OpReach, err, path)
	}()

	opt, err := s.parsePairReach(pairs...)
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
func (s *Storage) Statistical(pairs ...*types.Pair) (m info.StorageStatistic, err error) {
	defer func() {
		err = s.formatError(services.OpStatistical, err)
	}()

	m = info.NewStorageStatistic()

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

// ListPrefixSegments implements Storager.ListPrefixSegments
func (s *Storage) ListPrefixSegments(prefix string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpListPrefixSegments, err, prefix)
	}()

	opt, err := s.parsePairListPrefixSegments(pairs...)
	if err != nil {
		return
	}

	keyMarker := ""
	uploadIDMarker := ""
	limit := 200

	rp := s.getAbsPath(prefix)

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
			// TODO: we should handle rel prefix here.
			seg := segment.NewIndexBasedSegment(*v.Key, *v.UploadID)

			opt.SegmentFunc(seg)
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

// InitIndexSegment implements Storager.InitIndexSegment
func (s *Storage) InitIndexSegment(path string, pairs ...*types.Pair) (seg segment.Segment, err error) {
	defer func() {
		err = s.formatError(services.OpInitIndexSegment, err, path)
	}()

	_, err = s.parsePairInitIndexSegment(pairs...)
	if err != nil {
		return
	}

	input := &service.InitiateMultipartUploadInput{}

	rp := s.getAbsPath(path)

	output, err := s.bucket.InitiateMultipartUpload(rp, input)
	if err != nil {
		return
	}

	id := *output.UploadID

	seg = segment.NewIndexBasedSegment(path, id)
	return seg, nil
}

// WriteIndexSegment implements Storager.WriteIndexSegment
func (s *Storage) WriteIndexSegment(seg segment.Segment, r io.Reader, index int, size int64, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpWriteIndexSegment, err, seg.Path(), seg.ID())
	}()

	opt, err := s.parsePairWriteIndexSegment(pairs...)
	if err != nil {
		return
	}

	p, err := seg.(*segment.IndexBasedSegment).InsertPart(index, size)
	if err != nil {
		return
	}

	rp := s.getAbsPath(seg.Path())

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	_, err = s.bucket.UploadMultipart(rp, &service.UploadMultipartInput{
		PartNumber:    service.Int(p.Index),
		UploadID:      service.String(seg.ID()),
		ContentLength: &size,
		Body:          r,
	})
	if err != nil {
		return
	}
	return
}

// CompleteSegment implements Storager.CompleteSegment
func (s *Storage) CompleteSegment(seg segment.Segment, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpCompleteSegment, err, seg.Path(), seg.ID())
	}()

	parts := seg.(*segment.IndexBasedSegment).Parts()
	objectParts := make([]*service.ObjectPartType, 0, len(parts))
	for _, v := range parts {
		objectParts = append(objectParts, &service.ObjectPartType{
			PartNumber: service.Int(v.Index),
			Size:       service.Int64(v.Size),
		})
	}

	rp := s.getAbsPath(seg.Path())

	_, err = s.bucket.CompleteMultipartUpload(rp, &service.CompleteMultipartUploadInput{
		UploadID:    service.String(seg.ID()),
		ObjectParts: objectParts,
	})
	if err != nil {
		return
	}
	return
}

// AbortSegment implements Storager.AbortSegment
func (s *Storage) AbortSegment(seg segment.Segment, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpAbortSegment, err, seg.Path(), seg.ID())
	}()

	rp := s.getAbsPath(seg.Path())

	_, err = s.bucket.AbortMultipartUpload(rp, &service.AbortMultipartUploadInput{
		UploadID: service.String(seg.ID()),
	})
	if err != nil {
		return
	}
	return
}

// getAbsPath will calculate object storage's abs path
func (s *Storage) getAbsPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return prefix + path
}

// getRelPath will get object storage's rel path.
func (s *Storage) getRelPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return strings.TrimPrefix(path, prefix)
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

func (s *Storage) formatFileObject(v *service.KeyType) (o *types.Object, err error) {
	o = &types.Object{
		ID:         *v.Key,
		Name:       s.getRelPath(*v.Key),
		Type:       types.ObjectTypeFile,
		Size:       service.Int64Value(v.Size),
		UpdatedAt:  convertUnixTimestampToTime(service.IntValue(v.Modified)),
		ObjectMeta: info.NewObjectMeta(),
	}

	if v.MimeType != nil {
		o.SetContentType(service.StringValue(v.MimeType))
	}
	if value := service.StringValue(v.StorageClass); value != "" {
		setStorageClass(o.ObjectMeta, value)
	}
	if v.Etag != nil {
		o.SetETag(service.StringValue(v.Etag))
	}
	return o, nil
}
