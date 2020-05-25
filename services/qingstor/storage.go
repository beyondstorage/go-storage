package qingstor

import (
	"context"
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

func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	_, err = s.bucket.DeleteObject(rp)
	if err != nil {
		return
	}
	return nil
}
func (s *Storage) initIndexSegment(ctx context.Context, path string, opt *pairStorageInitIndexSegment) (seg segment.Segment, err error) {
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
func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (err error) {
	marker := ""
	delimiter := "/"
	limit := 200

	rp := s.getAbsPath(dir)

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
func (s *Storage) listPrefix(ctx context.Context, prefix string, opt *pairStorageListPrefix) (err error) {
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
func (s *Storage) listPrefixSegments(ctx context.Context, prefix string, opt *pairStorageListPrefixSegments) (err error) {
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
func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta info.StorageMeta, err error) {
	meta = info.NewStorageMeta()
	meta.Name = *s.properties.BucketName
	meta.WorkDir = s.workDir
	meta.SetLocation(*s.properties.Zone)
	return meta, nil
}
func (s *Storage) reach(ctx context.Context, path string, opt *pairStorageReach) (url string, err error) {
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
func (s *Storage) read(ctx context.Context, path string, opt *pairStorageRead) (rc io.ReadCloser, err error) {
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

	rc = output.Body
	if opt.HasReadCallbackFunc {
		rc = iowrap.CallbackReadCloser(rc, opt.ReadCallbackFunc)
	}
	return rc, nil
}
func (s *Storage) stat(ctx context.Context, path string, opt *pairStorageStat) (o *types.Object, err error) {
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
func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (err error) {
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
func (s *Storage) writeIndexSegment(ctx context.Context, seg segment.Segment, r io.Reader, index int, size int64, opt *pairStorageWriteIndexSegment) (err error) {
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
