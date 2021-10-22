package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"

	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	. "go.beyondstorage.io/v5/types"
)

func (s *Storage) completeMultipart(ctx context.Context, o *Object, parts []*Part, opt pairStorageCompleteMultipart) (err error) {
	input := s.formatCompleteMultipartUploadInput(o, parts, opt)
	_, err = s.service.CompleteMultipartUpload(ctx, input)
	if err != nil {
		return
	}

	o.Mode.Del(ModePart)
	o.Mode.Add(ModeRead)
	return
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	rp := s.getAbsPath(path)

	// Handle create multipart object separately.
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

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *Object, err error) {
	if !s.features.VirtualDir {
		err = NewOperationNotImplementedError("create_dir")
		return
	}

	rp := s.getAbsPath(path)

	// Add `/` at the end of `path` to simulate a directory.
	//ref: https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-folders.html
	rp += "/"

	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.name),
		Key:           aws.String(rp),
		ContentLength: int64(0),
	}
	if opt.HasStorageClass {
		input.StorageClass = s3types.StorageClass(opt.StorageClass)
	}
	if opt.HasExceptedBucketOwner {
		input.ExpectedBucketOwner = &opt.ExceptedBucketOwner
	}
	output, err := s.service.PutObject(ctx, input)
	if err != nil {
		return
	}
	o = s.newObject(true)
	o.Mode = ModeDir
	o.ID = rp
	o.Path = path
	o.SetEtag(aws.ToString(output.ETag))
	var sm ObjectSystemMetadata
	//output.ServerSideEncryption's type is s3types.ServerSideEncryption, which is equivalent to string
	sm.ServerSideEncryption = string(output.ServerSideEncryption)
	if v := aws.ToString(output.SSEKMSKeyId); v != "" {
		sm.ServerSideEncryptionAwsKmsKeyID = v
	}
	if v := aws.ToString(output.SSEKMSEncryptionContext); v != "" {
		sm.ServerSideEncryptionContext = v
	}
	if v := aws.ToString(output.SSECustomerAlgorithm); v != "" {
		sm.ServerSideEncryptionCustomerAlgorithm = v
	}
	if v := aws.ToString(output.SSECustomerKeyMD5); v != "" {
		sm.ServerSideEncryptionCustomerKeyMd5 = v
	}
	sm.ServerSideEncryptionBucketKeyEnabled = output.BucketKeyEnabled
	o.SetSystemMetadata(sm)

	return o, nil
}

// metadataLinkTargetHeader is the name of the user-defined metadata name used to store the link target.
const metadataLinkTargetHeader = "x-amz-meta-bs-link-target"

func (s *Storage) createLink(ctx context.Context, path string, target string, opt pairStorageCreateLink) (o *Object, err error) {
	rt := s.getAbsPath(target)
	rp := s.getAbsPath(path)

	input := &s3.PutObjectInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
		// As s3 does not support symlink, we can only use user-defined metadata to simulate it.
		// ref: https://github.com/beyondstorage/go-service-s3/blob/master/rfcs/178-add-virtual-link-support.md
		Metadata: map[string]string{
			metadataLinkTargetHeader: rt,
		},
	}
	output, err := s.service.PutObject(ctx, input)
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
		// s3 does not have an absolute path, so when we call `getAbsPath`, it will remove the prefix `/`.
		// To ensure that the path matches the one the user gets, we should re-add `/` here.
		o.SetLinkTarget("/" + rt)
		o.Mode |= ModeLink
	}

	var sm ObjectSystemMetadata
	//output.ServerSideEncryption's type is s3types.ServerSideEncryption, which is equivalent to string
	sm.ServerSideEncryption = string(output.ServerSideEncryption)
	if v := aws.ToString(output.SSEKMSKeyId); v != "" {
		sm.ServerSideEncryptionAwsKmsKeyID = v
	}
	if v := aws.ToString(output.SSEKMSEncryptionContext); v != "" {
		sm.ServerSideEncryptionContext = v
	}
	if v := aws.ToString(output.SSECustomerAlgorithm); v != "" {
		sm.ServerSideEncryptionCustomerAlgorithm = v
	}
	if v := aws.ToString(output.SSECustomerKeyMD5); v != "" {
		sm.ServerSideEncryptionCustomerKeyMd5 = v
	}
	sm.ServerSideEncryptionBucketKeyEnabled = output.BucketKeyEnabled

	o.SetSystemMetadata(sm)

	return
}

func (s *Storage) createMultipart(ctx context.Context, path string, opt pairStorageCreateMultipart) (o *Object, err error) {
	rp := s.getAbsPath(path)

	input, err := s.formatCreateMultipartUploadInput(path, opt)
	if err != nil {
		return nil, err
	}
	output, err := s.service.CreateMultipartUpload(ctx, input)
	if err != nil {
		return
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.Mode |= ModePart
	o.SetMultipartID(aws.ToString(output.UploadId))
	var sm ObjectSystemMetadata
	//output.ServerSideEncryption's type is s3types.ServerSideEncryption, which is equivalent to string
	sm.ServerSideEncryption = string(output.ServerSideEncryption)
	if v := aws.ToString(output.SSEKMSKeyId); v != "" {
		sm.ServerSideEncryptionAwsKmsKeyID = v
	}
	if v := aws.ToString(output.SSEKMSEncryptionContext); v != "" {
		sm.ServerSideEncryptionContext = v
	}
	if v := aws.ToString(output.SSECustomerAlgorithm); v != "" {
		sm.ServerSideEncryptionCustomerAlgorithm = v
	}
	if v := aws.ToString(output.SSECustomerKeyMD5); v != "" {
		sm.ServerSideEncryptionCustomerKeyMd5 = v
	}
	sm.ServerSideEncryptionBucketKeyEnabled = output.BucketKeyEnabled

	o.SetSystemMetadata(sm)
	return o, nil
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	if opt.HasMultipartID {
		abortInput := s.formatAbortMultipartUploadInput(path, opt)

		_, err = s.service.AbortMultipartUpload(ctx, abortInput)

		if err != nil {
			// AbortMultipartUpload is idempotent in s3, but non-idempotent in minio, we need to omit `NoSuchUpload` error for compatibility.
			// ref: [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
			e := &s3types.NoSuchUpload{}
			if errors.As(err, &e) {
				err = nil
			}
		}
		if err != nil {
			return err
		}
		return nil
	}

	input, err := s.formatDeleteObjectInput(path, opt)
	if err != nil {
		return err
	}

	// S3 DeleteObject is idempotent, so we don't need to check NoSuchKey error.
	//
	// References
	// - [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
	// - https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html
	_, err = s.service.DeleteObject(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	input := &objectPageStatus{
		maxKeys: 200,
		prefix:  s.getAbsPath(path),
	}

	if opt.HasExceptedBucketOwner {
		input.expectedBucketOwner = opt.ExceptedBucketOwner
	}

	if !opt.HasListMode {
		// Support `ListModePrefix` as the default `ListMode`.
		// ref: [GSP-46](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/654-unify-list-behavior.md)
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
		maxParts: 200,
		key:      o.ID,
		uploadId: o.MustGetMultipartID(),
	}
	if opt.HasExceptedBucketOwner {
		input.expectedBucketOwner = opt.ExceptedBucketOwner
	}

	return NewPartIterator(ctx, s.nextPartPage, input), nil
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	meta = NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	// set write restriction
	meta.SetWriteSizeMaximum(writeSizeMaximum)
	// set multipart restrictions
	meta.SetMultipartNumberMaximum(multipartNumberMaximum)
	meta.SetMultipartSizeMaximum(multipartSizeMaximum)
	meta.SetMultipartSizeMinimum(multipartSizeMinimum)
	return meta
}

func (s *Storage) nextObjectPageByDir(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	listInput := &s3.ListObjectsV2Input{
		Bucket:            &s.name,
		Delimiter:         &input.delimiter,
		MaxKeys:           int32(input.maxKeys),
		ContinuationToken: input.getServiceContinuationToken(),
		Prefix:            &input.prefix,
	}
	if input.expectedBucketOwner != "" {
		listInput.ExpectedBucketOwner = &input.expectedBucketOwner
	}
	output, err := s.service.ListObjectsV2(ctx, listInput)
	if err != nil {
		return err
	}

	for _, v := range output.CommonPrefixes {
		o := s.newObject(true)
		o.ID = *v.Prefix
		o.Path = s.getRelPath(*v.Prefix)
		o.Mode |= ModeDir

		page.Data = append(page.Data, o)
	}

	for _, v := range output.Contents {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if !output.IsTruncated {
		return IterateDone
	}

	input.continuationToken = *output.NextContinuationToken
	return nil
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	listInput := &s3.ListObjectsV2Input{
		Bucket:            &s.name,
		MaxKeys:           int32(input.maxKeys),
		ContinuationToken: input.getServiceContinuationToken(),
		Prefix:            &input.prefix,
	}
	if input.expectedBucketOwner != "" {
		listInput.ExpectedBucketOwner = &input.expectedBucketOwner
	}
	output, err := s.service.ListObjectsV2(ctx, listInput)
	if err != nil {
		return err
	}

	for _, v := range output.Contents {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}
	if !output.IsTruncated {
		return IterateDone
	}
	input.continuationToken = aws.ToString(output.NextContinuationToken)
	return nil
}

func (s *Storage) nextPartObjectPageByPrefix(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)
	listInput := &s3.ListMultipartUploadsInput{
		Bucket:         &s.name,
		KeyMarker:      &input.keyMarker,
		MaxUploads:     int32(input.maxKeys),
		Prefix:         &input.prefix,
		UploadIdMarker: &input.uploadIdMarker,
	}
	if input.expectedBucketOwner != "" {
		listInput.ExpectedBucketOwner = &input.expectedBucketOwner
	}
	output, err := s.service.ListMultipartUploads(ctx, listInput)
	if err != nil {
		return err
	}

	for _, v := range output.Uploads {
		o := s.newObject(true)
		o.ID = *v.Key
		o.Path = s.getRelPath(*v.Key)
		o.Mode |= ModePart
		o.SetMultipartID(*v.UploadId)

		page.Data = append(page.Data, o)
	}
	if !output.IsTruncated {
		return IterateDone
	}
	input.keyMarker = aws.ToString(output.KeyMarker)
	input.uploadIdMarker = aws.ToString(output.UploadIdMarker)
	return nil
}

func (s *Storage) nextPartPage(ctx context.Context, page *PartPage) error {
	input := page.Status.(*partPageStatus)
	listInput := &s3.ListPartsInput{
		Bucket:           &s.name,
		Key:              &input.key,
		MaxParts:         int32(input.maxParts),
		PartNumberMarker: &input.partNumberMarker,
		UploadId:         &input.uploadId,
	}
	if input.expectedBucketOwner != "" {
		listInput.ExpectedBucketOwner = &input.expectedBucketOwner
	}
	output, err := s.service.ListParts(ctx, listInput)
	if err != nil {
		return err
	}

	for _, v := range output.Parts {
		p := &Part{
			// The returned `PartNumber` is [1, 10000].
			// Set Index=*v.PartNumber-1 here to make the `PartNumber` zero-based for user.
			Index: int(v.PartNumber) - 1,
			Size:  v.Size,
			ETag:  aws.ToString(v.ETag),
		}

		page.Data = append(page.Data, p)
	}
	if !output.IsTruncated {
		return IterateDone
	}
	input.partNumberMarker = aws.ToString(output.NextPartNumberMarker)
	return nil
}

func (s *Storage) querySignHTTPCompleteMultipart(ctx context.Context, o *Object, parts []*Part, expire time.Duration, opt pairStorageQuerySignHTTPCompleteMultipart) (req *http.Request, err error) {

	// Currently presign only support Get/Put/Head object & UploadPart
	// We don't support  querySignHTTPCompleteMultipart for now
	return nil, services.ErrCapabilityInsufficient
}

func (s *Storage) querySignHTTPCreateMultipart(ctx context.Context, path string, expire time.Duration, opt pairStorageQuerySignHTTPCreateMultipart) (req *http.Request, err error) {
	// Currently presign only support Get/Put/Head object & UploadPart
	// We don't support  querySignHTTPCreateMultipart for now
	return nil, services.ErrCapabilityInsufficient
}

func (s *Storage) querySignHTTPDelete(ctx context.Context, path string, expire time.Duration, opt pairStorageQuerySignHTTPDelete) (req *http.Request, err error) {
	// Currently presign only support Get/Put/Head object & UploadPart
	// We don't support  querySignHTTPDelete for now
	return nil, services.ErrCapabilityInsufficient
}

func (s *Storage) querySignHTTPListMultipart(ctx context.Context, o *Object, expire time.Duration, opt pairStorageQuerySignHTTPListMultipart) (req *http.Request, err error) {
	// Currently presign only support Get/Put/Head object & UploadPart
	// We don't support  querySignHTTPListMultipart for now
	return nil, services.ErrCapabilityInsufficient
}

func (s *Storage) querySignHTTPRead(ctx context.Context, path string, expire time.Duration, opt pairStorageQuerySignHTTPRead) (req *http.Request, err error) {
	pairs, err := s.parsePairStorageRead(opt.pairs)
	if err != nil {
		return
	}

	input, err := s.formatGetObjectInput(path, pairs)
	if err != nil {
		return
	}
	presignClient := s3.NewPresignClient(s.service, func(options *s3.PresignOptions) {
		options.Expires = expire
	})
	getReq, err := presignClient.PresignGetObject(ctx, input)
	if err != nil {
		return
	}
	req, err = http.NewRequest("GET", getReq.URL, nil)
	if err != nil {
		return
	}
	req.Header = getReq.SignedHeader
	return
}

func (s *Storage) querySignHTTPWrite(ctx context.Context, path string, size int64, expire time.Duration, opt pairStorageQuerySignHTTPWrite) (req *http.Request, err error) {
	pairs, err := s.parsePairStorageWrite(opt.pairs)
	if err != nil {
		return nil, err
	}

	input, err := s.formatPutObjectInput(path, size, pairs)
	if err != nil {
		return nil, err
	}
	presignClient := s3.NewPresignClient(s.service, func(options *s3.PresignOptions) {
		options.Expires = expire
	})
	putReq, err := presignClient.PresignPutObject(ctx, input)
	if err != nil {
		return nil, err
	}
	req, err = http.NewRequest("PUT", putReq.URL, nil)
	if err != nil {
		return
	}
	req.Header = putReq.SignedHeader
	req.ContentLength = size
	return
}

func (s *Storage) querySignHTTPWriteMultipart(ctx context.Context, o *Object, size int64, index int, expire time.Duration, opt pairStorageQuerySignHTTPWriteMultipart) (req *http.Request, err error) {
	pairs, err := s.parsePairStorageWriteMultipart(opt.pairs)
	if err != nil {
		return nil, err
	}

	input, err := s.formatUploadPartInput(o, size, index, pairs)
	if err != nil {
		return nil, err
	}

	presignClient := s3.NewPresignClient(s.service, func(options *s3.PresignOptions) {
		options.Expires = expire
	})
	putReq, err := presignClient.PresignUploadPart(ctx, input)
	if err != nil {
		return nil, err
	}
	req, err = http.NewRequest("PUT", putReq.URL, nil)
	if err != nil {
		return
	}
	req.Header = putReq.SignedHeader
	req.ContentLength = size
	return
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	input, err := s.formatGetObjectInput(path, opt)
	if err != nil {
		return
	}
	output, err := s.service.GetObject(ctx, input)
	if err != nil {
		return
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
		listInput := &s3.ListPartsInput{
			Bucket:   aws.String(s.name),
			Key:      aws.String(rp),
			UploadId: aws.String(opt.MultipartID),
		}
		if opt.HasExceptedBucketOwner {
			listInput.ExpectedBucketOwner = &opt.ExceptedBucketOwner
		}
		_, err = s.service.ListParts(ctx, listInput)
		if err != nil {
			return nil, err
		}

		o = s.newObject(true)
		o.ID = rp
		o.Path = path
		o.Mode.Add(ModePart)
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

	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
	}
	if opt.HasExceptedBucketOwner {
		input.ExpectedBucketOwner = &opt.ExceptedBucketOwner
	}
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		input.SSECustomerAlgorithm, input.SSECustomerKey, input.SSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return
		}
	}

	output, err := s.service.HeadObject(ctx, input)
	if err != nil {
		return nil, err
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path

	if output.Metadata != nil {
		metadata := output.Metadata
		if target, ok := metadata[metadataLinkTargetHeader]; ok {
			// The path is a symlink object.
			if !s.features.VirtualLink {
				// The virtual link is not enabled, so we set the object mode to `ModeRead`.
				o.Mode |= ModeRead
			} else {
				o.Mode |= ModeLink
				// s3 does not have an absolute path, so when we call `getAbsPath`, it will remove the prefix `/`.
				// To ensure that the path matches the one the user gets, we should re-add `/` here.
				o.SetLinkTarget("/" + target)
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
	o.SetContentLength(output.ContentLength)
	o.SetLastModified(aws.ToTime(output.LastModified))
	if output.ContentType != nil {
		o.SetContentType(*output.ContentType)
	}
	if output.ETag != nil {
		o.SetEtag(*output.ETag)
	}

	var sm ObjectSystemMetadata
	//output.StorageClass's type is s3types.StorageClass, which is equivalent to string
	sm.StorageClass = string(output.StorageClass)
	if output.ServerSideEncryption != "" {
		sm.ServerSideEncryption = "output.ServerSideEncryption"
	}
	if v := aws.ToString(output.SSEKMSKeyId); v != "" {
		sm.ServerSideEncryptionAwsKmsKeyID = v
	}
	if v := aws.ToString(output.SSECustomerAlgorithm); v != "" {
		sm.ServerSideEncryptionCustomerAlgorithm = v
	}
	if v := aws.ToString(output.SSECustomerKeyMD5); v != "" {
		sm.ServerSideEncryptionCustomerKeyMd5 = v
	}
	sm.ServerSideEncryptionBucketKeyEnabled = output.BucketKeyEnabled

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
	if (r == nil && size == 0) || (r != nil && size == 0) {
		r = bytes.NewReader([]byte{})
	} else if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	} else {
		r = io.LimitReader(r, size)
	}

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	input, err := s.formatPutObjectInput(path, size, opt)
	if err != nil {
		return
	}

	input.Body = r
	_, err = s.service.PutObject(ctx, input)
	if err != nil {
		return
	}
	return size, nil
}

func (s *Storage) writeMultipart(ctx context.Context, o *Object, r io.Reader, size int64, index int, opt pairStorageWriteMultipart) (n int64, part *Part, err error) {
	if size > multipartSizeMaximum {
		err = fmt.Errorf("size limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return
	}
	if index < 0 || index >= multipartNumberMaximum {
		err = fmt.Errorf("multipart number limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return
	}

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	input := &s3.UploadPartInput{
		Bucket: &s.name,
		// For S3, the `PartNumber` is [1, 10000]. But for users, the `PartNumber` is zero-based.
		// Set PartNumber=index+1 here to ensure pass in an effective `PartNumber` for `UploadPart`.
		// ref: https://docs.aws.amazon.com/AmazonS3/latest/userguide/mpuoverview.html
		PartNumber:    (int32(index + 1)),
		Key:           aws.String(o.ID),
		UploadId:      aws.String(o.MustGetMultipartID()),
		ContentLength: size,
		Body:          iowrap.SizedReadSeekCloser(r, size),
	}
	if opt.HasExceptedBucketOwner {
		input.ExpectedBucketOwner = &opt.ExceptedBucketOwner
	}
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		input.SSECustomerAlgorithm, input.SSECustomerKey, input.SSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return
		}
	}
	output, err := s.service.UploadPart(ctx, input)
	if err != nil {
		return
	}

	part = &Part{
		Index: index,
		Size:  size,
		ETag:  aws.ToString(output.ETag),
	}
	return size, part, nil
}
