package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/info"
)

func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
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
func (s *Storage) initIndexSegment(ctx context.Context, path string, opt *pairStorageInitIndexSegment) (seg segment.Segment, err error) {
	rp := s.getAbsPath(path)

	input := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
	}

	output, err := s.service.CreateMultipartUpload(input)
	if err != nil {
		return
	}

	id := aws.StringValue(output.UploadId)

	seg = segment.NewIndexBasedSegment(path, id)
	return seg, nil
}
func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (err error) {
	marker := ""
	delimiter := "/"
	rp := s.getAbsPath(dir)

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
					ObjectMeta: info.NewObjectMeta(),
				}

				opt.DirFunc(o)
			}
		}

		if opt.HasFileFunc {
			for _, v := range output.Contents {
				o, err := s.formatFileObject(v)
				if err != nil {
					return err
				}

				opt.FileFunc(o)
			}
		}

		marker = aws.StringValue(output.StartAfter)
		if !aws.BoolValue(output.IsTruncated) {
			break
		}
	}
	return
}
func (s *Storage) listPrefix(ctx context.Context, prefix string, opt *pairStorageListPrefix) (err error) {
	marker := ""
	rp := s.getAbsPath(prefix)

	var output *s3.ListObjectsV2Output
	for {
		output, err = s.service.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:     aws.String(s.name),
			Prefix:     aws.String(rp),
			MaxKeys:    aws.Int64(1000),
			StartAfter: aws.String(marker),
		})
		if err != nil {
			return err
		}

		for _, v := range output.Contents {
			o, err := s.formatFileObject(v)
			if err != nil {
				return err
			}

			opt.ObjectFunc(o)
		}

		marker = aws.StringValue(output.StartAfter)
		if !aws.BoolValue(output.IsTruncated) {
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

	var output *s3.ListMultipartUploadsOutput
	for {
		output, err = s.service.ListMultipartUploads(&s3.ListMultipartUploadsInput{
			Bucket:         aws.String(s.name),
			KeyMarker:      aws.String(keyMarker),
			Prefix:         aws.String(rp),
			UploadIdMarker: aws.String(uploadIDMarker),
			MaxUploads:     aws.Int64(int64(limit)),
		})
		if err != nil {
			return
		}

		for _, v := range output.Uploads {
			seg := segment.NewIndexBasedSegment(*v.Key, *v.UploadId)

			opt.SegmentFunc(seg)
		}

		keyMarker = aws.StringValue(output.NextKeyMarker)
		uploadIDMarker = aws.StringValue(output.NextUploadIdMarker)
		if keyMarker == "" && uploadIDMarker == "" {
			break
		}
		if !aws.BoolValue(output.IsTruncated) {
			break
		}
	}
	return
}
func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta info.StorageMeta, err error) {
	meta = info.NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	return meta, nil
}
func (s *Storage) read(ctx context.Context, path string, opt *pairStorageRead) (rc io.ReadCloser, err error) {
	rp := s.getAbsPath(path)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
	}

	output, err := s.service.GetObject(input)
	if err != nil {
		return nil, err
	}

	rc = output.Body
	if opt.HasReadCallbackFunc {
		rc = iowrap.CallbackReadCloser(rc, opt.ReadCallbackFunc)
	}
	return rc, nil
}
func (s *Storage) stat(ctx context.Context, path string, opt *pairStorageStat) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
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
		ObjectMeta: info.NewObjectMeta(),
	}

	if output.ContentType != nil {
		o.SetContentType(*output.ContentType)
	}
	if output.ETag != nil {
		o.SetETag(*output.ETag)
	}
	if v := aws.StringValue(output.StorageClass); v != "" {
		setStorageClass(o.ObjectMeta, v)
	}
	return o, nil
}
func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (err error) {
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	rp := s.getAbsPath(path)

	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.name),
		Key:           aws.String(rp),
		ContentLength: &opt.Size,
		Body:          aws.ReadSeekCloser(r),
	}
	if opt.HasChecksum {
		input.ContentMD5 = &opt.Checksum
	}
	if opt.HasStorageClass {
		input.StorageClass = &opt.StorageClass
	}

	_, err = s.service.PutObject(input)
	if err != nil {
		return err
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

	_, err = s.service.UploadPart(&s3.UploadPartInput{
		Body:          aws.ReadSeekCloser(r),
		Bucket:        aws.String(s.name),
		ContentLength: aws.Int64(size),
		Key:           aws.String(rp),
		PartNumber:    aws.Int64(int64(p.Index)),
		UploadId:      aws.String(seg.ID()),
	})
	if err != nil {
		return
	}
	return
}

func (s *Storage) abortSegment(ctx context.Context, seg segment.Segment, opt *pairStorageAbortSegment) (err error) {
	parts := seg.(*segment.IndexBasedSegment).Parts()
	objectParts := make([]*s3.CompletedPart, 0, len(parts))
	for _, v := range parts {
		objectParts = append(objectParts, &s3.CompletedPart{
			PartNumber: aws.Int64(int64(v.Index)),
		})
	}

	rp := s.getAbsPath(seg.Path())

	_, err = s.service.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(s.name),
		Key:             aws.String(rp),
		MultipartUpload: &s3.CompletedMultipartUpload{Parts: objectParts},
		UploadId:        aws.String(seg.ID()),
	})
	if err != nil {
		return
	}
	return
}
func (s *Storage) completeSegment(ctx context.Context, seg segment.Segment, opt *pairStorageCompleteSegment) (err error) {
	rp := s.getAbsPath(seg.Path())

	_, err = s.service.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
		Bucket:   aws.String(s.name),
		Key:      aws.String(rp),
		UploadId: aws.String(seg.ID()),
	})
	if err != nil {
		return
	}
	return
}
