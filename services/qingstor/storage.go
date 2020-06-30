package qingstor

import (
	"context"
	"io"

	"github.com/pengsrc/go-shared/convert"
	"github.com/qingstor/qingstor-sdk-go/v4/service"

	"github.com/aos-dev/go-storage/v2/pkg/headers"
	"github.com/aos-dev/go-storage/v2/pkg/iowrap"
	"github.com/aos-dev/go-storage/v2/pkg/segment"
	"github.com/aos-dev/go-storage/v2/types"
	"github.com/aos-dev/go-storage/v2/types/info"
)

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
				// add filter to exclude dir-key itself, which would exist if created in console, see issue #365
				if convert.StringValue(v.Key) == rp {
					continue
				}
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

func (s *Storage) copy(ctx context.Context, src string, dst string, opt *pairStorageCopy) (err error) {
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

func (s *Storage) move(ctx context.Context, src string, dst string, opt *pairStorageMove) (err error) {
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

func (s *Storage) abortSegment(ctx context.Context, seg segment.Segment, opt *pairStorageAbortSegment) (err error) {
	rp := s.getAbsPath(seg.Path())

	_, err = s.bucket.AbortMultipartUpload(rp, &service.AbortMultipartUploadInput{
		UploadID: service.String(seg.ID()),
	})
	if err != nil {
		return
	}
	return
}
func (s *Storage) completeSegment(ctx context.Context, seg segment.Segment, opt *pairStorageCompleteSegment) (err error) {
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
func (s *Storage) statistical(ctx context.Context, opt *pairStorageStatistical) (statistic info.StorageStatistic, err error) {
	statistic = info.NewStorageStatistic()

	output, err := s.bucket.GetStatistics()
	if err != nil {
		return
	}

	if output.Size != nil {
		statistic.SetSize(*output.Size)
	}
	if output.Count != nil {
		statistic.SetCount(*output.Count)
	}
	return statistic, nil
}
