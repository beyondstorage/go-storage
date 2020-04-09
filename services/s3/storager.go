package s3

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the s3 object storage service.
type Storage struct {
	service s3iface.S3API

	name    string
	workDir string
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager s3 {Name: %s, WorkDir: %s}",
		s.name, "/"+s.workDir,
	)
}

// Metadata implements Storager.Metadata
func (s *Storage) Metadata(pairs ...*types.Pair) (m metadata.StorageMeta, err error) {
	m = metadata.NewStorageMeta()
	m.Name = s.name
	m.WorkDir = s.workDir
	return m, nil
}

// ListDir implements Storager.ListDir
func (s *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list_dir", err, path)
	}()

	opt, err := s.parsePairListDir(pairs...)
	if err != nil {
		return err
	}

	marker := ""
	delimiter := "/"
	rp := s.getAbsPath(path)

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
					ObjectMeta: metadata.NewObjectMeta(),
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
		if aws.BoolValue(output.IsTruncated) {
			break
		}
	}
	return
}

// ListPrefix implements Storager.ListPrefix
func (s *Storage) ListPrefix(prefix string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list_prefix", err, prefix)
	}()

	opt, err := s.parsePairListPrefix(pairs...)
	if err != nil {
		return err
	}

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
		if aws.BoolValue(output.IsTruncated) {
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

	opt, err := s.parsePairWrite(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
	}

	output, err := s.service.GetObject(input)
	if err != nil {
		return nil, err
	}

	r = output.Body
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReadCloser(r, opt.ReadCallbackFunc)
	}
	return r, nil
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("write", err, path)
	}()

	opt, err := s.parsePairWrite(pairs...)
	if err != nil {
		return err
	}

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	rp := s.getAbsPath(path)

	input := &s3.PutObjectInput{
		Key:           aws.String(rp),
		ContentLength: &opt.Size,
		Body:          aws.ReadSeekCloser(r),
	}
	if opt.HasChecksum {
		input.ContentMD5 = &opt.Checksum
	}
	if opt.HasStorageClass {
		storageClass, err := parseStorageClass(opt.StorageClass)
		if err != nil {
			return err
		}
		input.StorageClass = &storageClass
	}

	_, err = s.service.PutObject(input)
	if err != nil {
		return err
	}
	return nil
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	defer func() {
		err = s.formatError("stat", err, path)
	}()

	rp := s.getAbsPath(path)

	input := &s3.HeadObjectInput{
		Key: aws.String(rp),
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
		ObjectMeta: metadata.NewObjectMeta(),
	}

	if output.ContentType != nil {
		o.SetContentType(*output.ContentType)
	}
	if output.ETag != nil {
		o.SetETag(*output.ETag)
	}
	if v := formatStorageClass(aws.StringValue(output.StorageClass)); v != "" {
		o.SetStorageClass(v)
	}
	return o, nil
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("delete", err, path)
	}()

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

// ListPrefixSegments implements Storager.ListPrefixSegments
func (s *Storage) ListPrefixSegments(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list_prefix_segments", err, path)
	}()

	opt, err := s.parsePairListPrefixSegments(pairs...)
	if err != nil {
		return
	}

	keyMarker := ""
	uploadIDMarker := ""
	limit := 200

	rp := s.getAbsPath(path)

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

// InitSegment implements Storager.InitSegment
func (s *Storage) InitSegment(path string, pairs ...*types.Pair) (seg segment.Segment, err error) {
	defer func() {
		err = s.formatError("init_segment", err, path)
	}()

	_, err = s.parsePairInitSegment(pairs...)
	if err != nil {
		return
	}

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

// WriteSegment implements Storager.WriteSegment
func (s *Storage) WriteSegment(seg segment.Segment, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("write_segment", err, seg.Path(), seg.ID())
	}()

	opt, err := s.parsePairWriteSegment(pairs...)
	if err != nil {
		return
	}

	p, err := seg.(*segment.IndexBasedSegment).InsertPart(opt.Index, opt.Size)
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
		ContentLength: aws.Int64(opt.Size),
		Key:           aws.String(rp),
		PartNumber:    aws.Int64(int64(p.Index)),
		UploadId:      aws.String(seg.ID()),
	})
	if err != nil {
		return
	}
	return
}

// CompleteSegment implements Storager.CompleteSegment
func (s *Storage) CompleteSegment(seg segment.Segment, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("complete_segment", err, seg.Path(), seg.ID())
	}()

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

// AbortSegment implements Storager.AbortSegment
func (s *Storage) AbortSegment(seg segment.Segment, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("abort_segment", err, seg.Path(), seg.ID())
	}()

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

func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
}

func (s *Storage) getRelPath(path string) string {
	return strings.TrimPrefix(path, s.workDir+"/")
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

func (s *Storage) formatFileObject(v *s3.Object) (o *types.Object, err error) {
	o = &types.Object{
		ID:         *v.Key,
		Name:       s.getRelPath(*v.Key),
		Type:       types.ObjectTypeFile,
		Size:       aws.Int64Value(v.Size),
		UpdatedAt:  aws.TimeValue(v.LastModified),
		ObjectMeta: metadata.NewObjectMeta(),
	}

	if value := formatStorageClass(aws.StringValue(v.StorageClass)); value != "" {
		o.SetStorageClass(value)
	}
	if v.ETag != nil {
		o.SetETag(*v.ETag)
	}

	return
}
