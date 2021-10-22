package qingstor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pengsrc/go-shared/convert"
	"github.com/qingstor/qingstor-sdk-go/v4/service"

	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	. "go.beyondstorage.io/v5/types"
)

func (s *Storage) commitAppend(ctx context.Context, o *Object, opt pairStorageCommitAppend) (err error) {
	return
}

func (s *Storage) completeMultipart(ctx context.Context, o *Object, parts []*Part, opt pairStorageCompleteMultipart) (err error) {
	objectParts := make([]*service.ObjectPartType, 0, len(parts))
	for _, v := range parts {
		objectParts = append(objectParts, &service.ObjectPartType{
			Etag:       service.String(v.ETag),
			PartNumber: service.Int(v.Index),
			Size:       service.Int64(v.Size),
		})
	}

	_, err = s.bucket.CompleteMultipartUploadWithContext(ctx, o.ID, &service.CompleteMultipartUploadInput{
		UploadID:    service.String(o.MustGetMultipartID()),
		ObjectParts: objectParts,
	})
	if err != nil {
		return
	}
	o.Mode.Del(ModePart)
	o.Mode.Add(ModeRead)
	return
}

func (s *Storage) copy(ctx context.Context, src string, dst string, opt pairStorageCopy) (err error) {
	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	srcPath := "/" + service.StringValue(s.properties.BucketName) + "/" + url.QueryEscape(rs)
	input := &service.PutObjectInput{
		XQSCopySource: &srcPath,
	}
	if opt.HasEncryptionCustomerAlgorithm {
		input.XQSEncryptionCustomerAlgorithm, input.XQSEncryptionCustomerKey, input.XQSEncryptionCustomerKeyMD5, err = calculateEncryptionHeaders(opt.EncryptionCustomerAlgorithm, opt.EncryptionCustomerKey)
		if err != nil {
			return
		}
	}
	if opt.HasCopySourceEncryptionCustomerAlgorithm {
		input.XQSCopySourceEncryptionCustomerAlgorithm, input.XQSCopySourceEncryptionCustomerKey, input.XQSCopySourceEncryptionCustomerKeyMD5, err = calculateEncryptionHeaders(opt.CopySourceEncryptionCustomerAlgorithm, opt.CopySourceEncryptionCustomerKey)
		if err != nil {
			return
		}
	}

	_, err = s.bucket.PutObjectWithContext(ctx, rd, input)
	if err != nil {
		return
	}
	return nil
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	rp := s.getAbsPath(path)

	// handle create multipart object separately
	// if opt has multipartID, set object done, because we can't stat multipart object in QingStor
	if opt.HasMultipartID {
		o = s.newObject(true)
		o.Mode = ModePart
		o.SetMultipartID(opt.MultipartID)
	} else {
		if opt.HasObjectMode && opt.ObjectMode.IsDir() {
			if !s.features.VirtualDir {
				return
			}

			rp += "/"
			o = s.newObject(true)
			o.Mode = ModeDir
		} else {
			o = s.newObject(false)
			o.Mode = ModeRead
		}
	}
	o.ID = rp
	o.Path = path
	return o
}

func (s *Storage) createAppend(ctx context.Context, path string, opt pairStorageCreateAppend) (o *Object, err error) {
	rp := s.getAbsPath(path)

	// We should set offset to 0 whether the object exists or not.
	var offset int64 = 0

	// If the object exists, we should delete it.
	headInput := &service.HeadObjectInput{}
	_, err = s.bucket.HeadObjectWithContext(ctx, rp, headInput)
	if err == nil {
		_, err = s.bucket.DeleteObject(rp)
		if err != nil {
			return nil, err
		}
	}

	input := &service.AppendObjectInput{
		Position: &offset,
	}
	if opt.HasContentType {
		input.ContentType = &opt.ContentType
	}
	if opt.HasStorageClass {
		input.XQSStorageClass = &opt.StorageClass
	}

	output, err := s.bucket.AppendObjectWithContext(ctx, rp, input)
	if err != nil {
		return
	}

	if output == nil || output.XQSNextAppendPosition == nil {
		err = ErrAppendNextPositionEmpty
		return
	} else {
		offset = *output.XQSNextAppendPosition
	}

	o = s.newObject(true)
	o.Mode = ModeRead | ModeAppend
	o.ID = rp
	o.Path = path
	o.SetAppendOffset(offset)
	// set metadata
	if opt.HasContentType {
		o.SetContentType(opt.ContentType)
	}
	var sm ObjectSystemMetadata
	if opt.HasStorageClass {
		sm.StorageClass = opt.StorageClass
	}
	return o, nil
}

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *Object, err error) {
	if !s.features.VirtualDir {
		err = NewOperationNotImplementedError("create_dir")
		return
	}

	rp := s.getAbsPath(path)

	// Add `/` at the end of path to simulate a directory.
	// ref: https://docs.qingcloud.com/qingstor/api/object/put.html
	rp += "/"

	input := &service.PutObjectInput{
		ContentLength: service.Int64(0),
	}
	if opt.HasStorageClass {
		input.XQSStorageClass = service.String(opt.StorageClass)
	}

	_, err = s.bucket.PutObjectWithContext(ctx, rp, input)
	if err != nil {
		return
	}

	o = s.newObject(true)
	o.Path = path
	o.ID = rp
	o.Mode = ModeDir
	return
}

// metadataLinkTargetHeader is the name of the user-defined metadata name used to store the target.
const metadataLinkTargetHeader = "x-qs-meta-bs-link-target"

func (s *Storage) createLink(ctx context.Context, path string, target string, opt pairStorageCreateLink) (o *Object, err error) {
	rt := s.getAbsPath(target)
	rp := s.getAbsPath(path)

	input := &service.PutObjectInput{
		// As qingstor does not support symlink, we can only use user-defined metadata to simulate it.
		// ref: https://github.com/beyondstorage/go-service-qingstor/blob/master/rfcs/79-add-virtual-link-support.md
		XQSMetaData: &map[string]string{
			metadataLinkTargetHeader: rt,
		},
	}

	_, err = s.bucket.PutObjectWithContext(ctx, rp, input)
	if err != nil {
		return nil, err
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path

	if !s.features.VirtualLink {
		// The virtual link is not enabled, so we set the object mode to `ModeRead`.
		o.Mode |= ModeRead
	} else {
		// qingstor does not have an absolute path, so when we call `getAbsPath`, it will remove the prefix `/`.
		// To ensure that the path matches the one the user gets, we should re-add `/` here.
		o.SetLinkTarget("/" + rt)
		o.Mode |= ModeLink
	}
	return
}

func (s *Storage) createMultipart(ctx context.Context, path string, opt pairStorageCreateMultipart) (o *Object, err error) {
	input := &service.InitiateMultipartUploadInput{}
	if opt.HasEncryptionCustomerAlgorithm {
		input.XQSEncryptionCustomerAlgorithm, input.XQSEncryptionCustomerKey, input.XQSEncryptionCustomerKeyMD5, err = calculateEncryptionHeaders(opt.EncryptionCustomerAlgorithm, opt.EncryptionCustomerKey)
		if err != nil {
			return
		}
	}

	rp := s.getAbsPath(path)

	output, err := s.bucket.InitiateMultipartUploadWithContext(ctx, rp, input)
	if err != nil {
		return
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.Mode |= ModePart
	o.SetMultipartID(*output.UploadID)

	return o, nil
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	if opt.HasMultipartID {
		// QingStor AbortMultipartUpload is idempotent, so we don't need to check upload_not_exists error.
		//
		// References
		// - [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
		// - https://docs.qingcloud.com/qingstor/api/object/multipart/abort_multipart_upload.html
		_, err = s.bucket.AbortMultipartUploadWithContext(ctx, rp, &service.AbortMultipartUploadInput{
			UploadID: service.String(opt.MultipartID),
		})
		if err != nil {
			return
		}
		return
	}

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			err = services.PairUnsupportedError{Pair: ps.WithObjectMode(opt.ObjectMode)}
			return
		}

		rp += "/"
	}

	// QingStor DeleteObject is idempotent, so we don't need to check object_not_exists error.
	//
	// - [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
	// - https://docs.qingcloud.com/qingstor/api/object/delete
	_, err = s.bucket.DeleteObjectWithContext(ctx, rp)
	if err != nil {
		return
	}
	return nil
}

func (s *Storage) fetch(ctx context.Context, path string, url string, opt pairStorageFetch) (err error) {
	_, err = s.bucket.PutObjectWithContext(ctx, path, &service.PutObjectInput{
		XQSFetchSource: service.String(url),
	})
	return err
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	input := &objectPageStatus{
		limit:  200,
		prefix: s.getAbsPath(path),
	}

	if !opt.HasListMode {
		// Support `ListModePrefix` as the default `ListMode`.
		// ref: [GSP-654](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/654-unify-list-behavior.md)
		opt.ListMode = ListModePrefix
	}

	var nextFn NextObjectFunc

	switch {
	case opt.ListMode.IsPart():
		nextFn = s.nextPartObjectPageByPrefix
	case opt.ListMode.IsDir():
		input.delimiter = "/"
		nextFn = s.nextObjectPageByDir
	case opt.ListMode.IsPrefix():
		nextFn = s.nextObjectPageByPrefix
	default:
		return nil, services.ListModeInvalidError{Actual: opt.ListMode}
	}

	return NewObjectIterator(ctx, nextFn, input), nil
}

func (s *Storage) listMultipart(ctx context.Context, o *Object, opt pairStorageListMultipart) (pi *PartIterator, err error) {
	input := &partPageStatus{
		limit:    200,
		prefix:   o.ID,
		uploadID: o.MustGetMultipartID(),
	}

	return NewPartIterator(ctx, s.nextPartPage, input), nil
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	meta = NewStorageMeta()
	meta.Name = *s.properties.BucketName
	meta.WorkDir = s.workDir
	meta.SetLocation(*s.properties.Zone)
	// set write restriction
	meta.SetWriteSizeMaximum(writeSizeMaximum)
	// set copy restriction
	meta.SetCopySizeMaximum(copySizeMaximum)
	// set append restrictions
	meta.SetAppendSizeMaximum(appendSizeMaximum)
	meta.SetAppendTotalSizeMaximum(appendTotalSizeMaximum)
	// set multipart restrictions
	meta.SetMultipartNumberMaximum(multipartNumberMaximum)
	meta.SetMultipartSizeMaximum(multipartSizeMaximum)
	meta.SetMultipartSizeMinimum(multipartSizeMinimum)
	return meta
}

func (s *Storage) move(ctx context.Context, src string, dst string, opt pairStorageMove) (err error) {
	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	srcPath := "/" + service.StringValue(s.properties.BucketName) + "/" + url.QueryEscape(rs)
	_, err = s.bucket.PutObjectWithContext(ctx, rd, &service.PutObjectInput{
		XQSMoveSource: &srcPath,
	})
	if err != nil {
		return
	}
	return nil
}

func (s *Storage) nextObjectPageByDir(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	output, err := s.bucket.ListObjectsWithContext(ctx, &service.ListObjectsInput{
		Delimiter: &input.delimiter,
		Limit:     &input.limit,
		Marker:    &input.marker,
		Prefix:    &input.prefix,
	})
	if err != nil {
		return err
	}

	for _, v := range output.CommonPrefixes {
		o := s.newObject(true)
		o.ID = *v
		o.Path = s.getRelPath(*v)
		o.Mode |= ModeDir

		page.Data = append(page.Data, o)
	}

	for _, v := range output.Keys {
		// add filter to exclude dir-key itself, which would exist if created in console, see issue #365
		if convert.StringValue(v.Key) == input.prefix {
			continue
		}
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if service.StringValue(output.NextMarker) == "" {
		return IterateDone
	}
	if !service.BoolValue(output.HasMore) {
		return IterateDone
	}
	if len(output.Keys) == 0 {
		return IterateDone
	}

	input.marker = *output.NextMarker
	return nil
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	output, err := s.bucket.ListObjectsWithContext(ctx, &service.ListObjectsInput{
		Limit:  &input.limit,
		Marker: &input.marker,
		Prefix: &input.prefix,
	})
	if err != nil {
		return err
	}

	for _, v := range output.Keys {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if service.StringValue(output.NextMarker) == "" {
		return IterateDone
	}
	if !service.BoolValue(output.HasMore) {
		return IterateDone
	}
	if len(output.Keys) == 0 {
		return IterateDone
	}

	input.marker = *output.NextMarker
	return nil
}

func (s *Storage) nextPartObjectPageByPrefix(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	output, err := s.bucket.ListMultipartUploadsWithContext(ctx, &service.ListMultipartUploadsInput{
		KeyMarker:      &input.marker,
		Limit:          &input.limit,
		Prefix:         &input.prefix,
		UploadIDMarker: &input.partIdMarker,
	})
	if err != nil {
		return err
	}

	for _, v := range output.Uploads {
		o := s.newObject(true)
		o.ID = *v.Key
		o.Path = s.getRelPath(*v.Key)
		o.Mode |= ModePart
		o.SetMultipartID(*v.UploadID)

		page.Data = append(page.Data, o)
	}

	nextKeyMarker := service.StringValue(output.NextKeyMarker)
	nextUploadIDMarker := service.StringValue(output.NextUploadIDMarker)

	if nextKeyMarker == "" && nextUploadIDMarker == "" {
		return IterateDone
	}
	if !service.BoolValue(output.HasMore) {
		return IterateDone
	}

	input.marker = nextKeyMarker
	input.partIdMarker = nextUploadIDMarker
	return nil
}

func (s *Storage) nextPartPage(ctx context.Context, page *PartPage) error {
	input := page.Status.(*partPageStatus)

	output, err := s.bucket.ListMultipartWithContext(ctx, input.prefix, &service.ListMultipartInput{
		Limit:            &input.limit,
		PartNumberMarker: &input.partNumberMarker,
		UploadID:         &input.uploadID,
	})
	if err != nil {
		return err
	}

	for _, v := range output.ObjectParts {
		p := &Part{
			Index: *v.PartNumber,
			Size:  *v.Size,
			ETag:  service.StringValue(v.Etag),
		}

		page.Data = append(page.Data, p)
	}

	// FIXME: QingStor ListMulitpart API looks like buggy.
	offset := input.partNumberMarker + len(output.ObjectParts)
	if offset >= service.IntValue(output.Count) {
		return IterateDone
	}

	input.partNumberMarker = offset
	return nil
}

func (s *Storage) querySignHTTPDelete(ctx context.Context, path string, expire time.Duration, opt pairStorageQuerySignHTTPDelete) (req *http.Request, err error) {
	panic("not implemented")
}

func (s *Storage) querySignHTTPRead(ctx context.Context, path string, expire time.Duration, opt pairStorageQuerySignHTTPRead) (req *http.Request, err error) {
	pairs, err := s.parsePairStorageRead(opt.pairs)
	if err != nil {
		return
	}

	input, err := s.formatGetObjectInput(pairs)
	if err != nil {
		return
	}

	bucket := s.bucket.(*service.Bucket)

	rp := s.getAbsPath(path)

	r, _, err := bucket.GetObjectRequest(rp, input)
	if err != nil {
		return
	}
	if err = r.BuildWithContext(ctx); err != nil {
		return
	}

	if err = r.SignQuery(int(expire.Seconds())); err != nil {
		return
	}

	return r.HTTPRequest, nil
}

func (s *Storage) querySignHTTPWrite(ctx context.Context, path string, size int64, expire time.Duration, opt pairStorageQuerySignHTTPWrite) (req *http.Request, err error) {
	pairs, err := s.parsePairStorageWrite(opt.pairs)
	if err != nil {
		return
	}

	input, err := s.formatPutObjectInput(size, pairs)
	if err != nil {
		return
	}

	bucket := s.bucket.(*service.Bucket)

	rp := s.getAbsPath(path)

	r, _, err := bucket.PutObjectRequest(rp, input)
	if err != nil {
		return
	}
	if err = r.BuildWithContext(ctx); err != nil {
		return
	}

	if err = r.SignQuery(int(expire.Seconds())); err != nil {
		return
	}

	return r.HTTPRequest, nil
}

func (s *Storage) reach(ctx context.Context, path string, opt pairStorageReach) (url string, err error) {
	// FIXME: sdk should export GetObjectRequest as interface too?
	bucket := s.bucket.(*service.Bucket)

	rp := s.getAbsPath(path)

	r, _, err := bucket.GetObjectRequest(rp, nil)
	if err != nil {
		return
	}
	if err = r.BuildWithContext(ctx); err != nil {
		return
	}

	expire := opt.Expire
	if err = r.SignQuery(int(expire.Seconds())); err != nil {
		return
	}
	return r.HTTPRequest.URL.String(), nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	input, err := s.formatGetObjectInput(opt)
	if err != nil {
		return
	}

	rp := s.getAbsPath(path)

	output, err := s.bucket.GetObjectWithContext(ctx, rp, input)
	if err != nil {
		return n, err
	}
	defer output.Body.Close()

	rc := output.Body
	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}

	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {

	rp := s.getAbsPath(path)

	if opt.HasMultipartID {
		input := &service.ListMultipartInput{
			UploadID: service.String(opt.MultipartID),
			Limit:    service.Int(0),
		}
		_, err := s.bucket.ListMultipartWithContext(ctx, rp, input)
		if err != nil {
			return nil, err
		}

		o = s.newObject(true)
		o.ID = rp
		o.Path = path
		o.Mode |= ModePart
		o.SetMultipartID(opt.MultipartID)
		return o, nil
	}

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			err = services.PairUnsupportedError{Pair: ps.WithObjectMode(opt.ObjectMode)}
			return
		}

		rp += "/"
	}

	input := &service.HeadObjectInput{}
	output, err := s.bucket.HeadObjectWithContext(ctx, rp, input)
	if err != nil {
		return
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path

	if output.XQSMetaData != nil {
		metadata := *output.XQSMetaData
		// By calling `HeadObject`, the first letter of the `key` of the object metadata will be capitalized.
		if v, ok := metadata[metadataLinkTargetHeader]; ok {
			// The path is a symlink object.
			if !s.features.VirtualLink {
				// The virtual link is not enabled, so we set the object mode to `ModeRead`.
				o.Mode |= ModeRead
			} else {
				// qingstor does not have an absolute path, so when we call `getAbsPath`, it will remove the prefix `/`.
				// To ensure that the path matches the one the user gets, we should re-add `/` here.
				o.SetLinkTarget("/" + v)
				o.Mode |= ModeLink
			}
		}
	}

	if o.Mode&ModeLink == 0 && o.Mode&ModeRead == 0 {
		if opt.HasObjectMode && opt.ObjectMode.IsDir() {
			o.Mode |= ModeDir
		} else {
			o.Mode |= ModeRead
		}
	}

	o.SetContentLength(service.Int64Value(output.ContentLength))
	o.SetLastModified(service.TimeValue(output.LastModified))

	if output.ContentType != nil {
		o.SetContentType(service.StringValue(output.ContentType))
	}
	if output.ETag != nil {
		o.SetEtag(service.StringValue(output.ETag))
	}

	var sm ObjectSystemMetadata
	if v := service.StringValue(output.XQSStorageClass); v != "" {
		sm.StorageClass = v
	}
	if v := service.StringValue(output.XQSEncryptionCustomerAlgorithm); v != "" {
		sm.EncryptionCustomerAlgorithm = v
	}
	o.SetSystemMetadata(sm)

	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	if size > writeSizeMaximum {
		err = fmt.Errorf("size limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return
	}

	// According to GSP-751, we should allow the user to pass in a nil io.Reader.
	// ref: https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/751-write-empty-file-behavior.md
	if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	}

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	input, err := s.formatPutObjectInput(size, opt)
	if err != nil {
		return
	}
	input.Body = io.LimitReader(r, size)

	rp := s.getAbsPath(path)

	_, err = s.bucket.PutObjectWithContext(ctx, rp, input)
	if err != nil {
		return
	}
	return size, nil
}

func (s *Storage) writeAppend(ctx context.Context, o *Object, r io.Reader, size int64, opt pairStorageWriteAppend) (n int64, err error) {
	if size > appendSizeMaximum {
		err = fmt.Errorf("size limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return
	}

	rp := o.GetID()

	offset, _ := o.GetAppendOffset()

	input := &service.AppendObjectInput{
		Position:      &offset,
		ContentLength: &size,
		Body:          io.LimitReader(r, size),
	}
	if opt.HasContentMd5 {
		input.ContentMD5 = &opt.ContentMd5
	}

	output, err := s.bucket.AppendObjectWithContext(ctx, rp, input)
	if err != nil {
		return
	}

	if output == nil || output.XQSNextAppendPosition == nil {
		err = ErrAppendNextPositionEmpty
		return
	} else {
		offset = *output.XQSNextAppendPosition
	}

	// We should reset the offset after calling `AppendObject` to prevent the offset being changed.
	o.SetAppendOffset(offset)

	return size, nil
}

func (s *Storage) writeMultipart(ctx context.Context, o *Object, r io.Reader, size int64, index int, opt pairStorageWriteMultipart) (n int64, part *Part, err error) {
	if index < multipartNumberMinimum || index > multipartNumberMaximum {
		err = ErrPartNumberInvalid
		return
	}
	if size > multipartSizeMaximum {
		err = fmt.Errorf("size limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return
	}

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	input := &service.UploadMultipartInput{
		PartNumber:    service.Int(index),
		UploadID:      service.String(o.MustGetMultipartID()),
		ContentLength: &size,
		Body:          io.LimitReader(r, size),
	}
	if opt.HasEncryptionCustomerAlgorithm {
		input.XQSEncryptionCustomerAlgorithm, input.XQSEncryptionCustomerKey, input.XQSEncryptionCustomerKeyMD5, err = calculateEncryptionHeaders(opt.EncryptionCustomerAlgorithm, opt.EncryptionCustomerKey)
		if err != nil {
			return
		}
	}

	output, err := s.bucket.UploadMultipartWithContext(ctx, o.ID, input)
	if err != nil {
		return
	}

	part = &Part{
		Index: index,
		Size:  size,
		ETag:  service.StringValue(output.ETag),
	}
	return size, part, nil
}
