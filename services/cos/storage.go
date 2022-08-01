package cos

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"

	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/pkg/headers"
	"github.com/beyondstorage/go-storage/v5/pkg/iowrap"
	"github.com/beyondstorage/go-storage/v5/services"
	"github.com/beyondstorage/go-storage/v5/types"
)

func (s *Storage) completeMultipart(ctx context.Context, o *types.Object, parts []*types.Part, opt pairStorageCompleteMultipart) (err error) {
	// Users should make sure the numbers of the uploaded parts are continuous and sorted in ascending order.
	// ref: https://intl.cloud.tencent.com/document/product/436/7742
	upload := &cos.CompleteMultipartUploadOptions{}
	for _, v := range parts {
		upload.Parts = append(upload.Parts, cos.Object{
			ETag: v.ETag,
			// For users the `PartNumber` is zero-based. But for COS, the effective `PartNumber` is [1, 10000].
			// Set PartNumber=v.Index+1 here to ensure pass in an effective `PartNumber` for `cos.Object`.
			PartNumber: v.Index + 1,
		})
	}

	_, _, err = s.object.CompleteMultipartUpload(ctx, o.ID, o.MustGetMultipartID(), upload)
	if err != nil {
		return
	}

	o.Mode.Del(types.ModePart)
	o.Mode.Add(types.ModeRead)
	return
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *types.Object) {
	rp := s.getAbsPath(path)

	// Handle create multipart object separately.
	if opt.HasMultipartID {
		o = s.newObject(true)
		o.Mode = types.ModePart
		o.SetMultipartID(opt.MultipartID)
	} else {
		if opt.HasObjectMode && opt.ObjectMode.IsDir() {
			if !s.features.VirtualDir {
				return
			}

			rp += "/"
			o = s.newObject(true)
			o.Mode = types.ModeDir
		} else {
			o = s.newObject(false)
			o.Mode = types.ModeRead
		}
	}

	o.ID = rp
	o.Path = path
	return o
}

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *types.Object, err error) {
	if !s.features.VirtualDir {
		err = types.NewOperationNotImplementedError("create_dir")
		return
	}

	rp := s.getAbsPath(path)

	// Add `/` at the end of `path` to simulate a directory.
	// ref: https://cloud.tencent.com/document/product/436/13324
	rp += "/"

	putOptions := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: 0,
		},
	}
	if opt.HasStorageClass {
		putOptions.XCosStorageClass = opt.StorageClass
	}

	_, err = s.object.Put(ctx, rp, io.LimitReader(nil, 0), putOptions)
	if err != nil {
		return
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.Mode |= types.ModeDir
	return
}

func (s *Storage) createMultipart(ctx context.Context, path string, opt pairStorageCreateMultipart) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	input := &cos.InitiateMultipartUploadOptions{}
	if opt.HasStorageClass {
		input.XCosStorageClass = opt.StorageClass
	}
	if opt.HasContentType {
		input.ContentType = opt.ContentType
	}
	// SSE-C
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		input.XCosSSECustomerAglo, input.XCosSSECustomerKey, input.XCosSSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return
		}
	}
	// SSE-COS or SSE-KMS
	if opt.HasServerSideEncryption {
		input.XCosServerSideEncryption = opt.ServerSideEncryption
		if opt.ServerSideEncryption == ServerSideEncryptionCosKms {
			// FIXME: we can remove the usage of `XOptionHeader` when cos' SDK supports SSE-KMS
			input.XOptionHeader = &http.Header{}
			if opt.HasServerSideEncryptionCosKmsKeyID {
				input.XOptionHeader.Set(serverSideEncryptionCosKmsKeyIdHeader, opt.ServerSideEncryptionCosKmsKeyID)
			}
			if opt.HasServerSideEncryptionContext {
				input.XOptionHeader.Set(serverSideEncryptionContextHeader, opt.ServerSideEncryptionContext)
			}
		}
	}

	output, _, err := s.object.InitiateMultipartUpload(ctx, rp, input)
	if err != nil {
		return
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.Mode |= types.ModePart
	o.SetMultipartID(output.UploadID)
	return o, nil
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	if opt.HasMultipartID {
		_, err = s.object.AbortMultipartUpload(ctx, rp, opt.MultipartID)
		if err != nil && checkError(err, responseCodeNoSuchUpload) {
			// Omit `NoSuchUpload` error here.
			// ref: [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
			err = nil
		}
		if err != nil {
			return err
		}
		return nil
	}

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			err = services.PairUnsupportedError{Pair: ps.WithObjectMode(opt.ObjectMode)}
			return
		}

		rp += "/"
	}

	_, err = s.object.Delete(ctx, rp)
	if err != nil && checkError(err, responseCodeNoSuchKey) {
		// Omit `NoSuchKey` error here.
		// ref: [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
		err = nil
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *types.ObjectIterator, err error) {
	if !opt.HasListMode {
		// Support `ListModePrefix` as the default `ListMode`.
		// ref: [GSP-46](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/654-unify-list-behavior.md)
		opt.ListMode = types.ListModePrefix
	}

	if opt.ListMode.IsDir() {
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
	}

	input := &objectPageStatus{
		maxKeys: 200,
		prefix:  s.getAbsPath(path),
	}

	var nextFn types.NextObjectFunc

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

	return types.NewObjectIterator(ctx, nextFn, input), nil
}

func (s *Storage) listMultipart(ctx context.Context, o *types.Object, opt pairStorageListMultipart) (pi *types.PartIterator, err error) {
	input := &partPageStatus{
		maxParts: "200",
		key:      o.ID,
		uploadId: o.MustGetMultipartID(),
	}

	return types.NewPartIterator(ctx, s.nextPartPage, input), nil
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *types.StorageMeta) {
	meta = types.NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	// set write restriction
	meta.SetWriteSizeMaximum(writeSizeMaximum)
	// set multipart restrictions
	meta.SetMultipartNumberMaximum(multipartNumberMaximum)
	meta.SetMultipartSizeMaximum(multipartSizeMaximum)
	meta.SetMultipartSizeMinimum(multipartSizeMinimum)
	return
}

func (s *Storage) nextObjectPageByDir(ctx context.Context, page *types.ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	output, _, err := s.bucket.Get(ctx, &cos.BucketGetOptions{
		Prefix:    input.prefix,
		Delimiter: input.delimiter,
		Marker:    input.keyMarker,
		MaxKeys:   input.maxKeys,
	})
	if err != nil {
		return err
	}

	for _, v := range output.CommonPrefixes {
		o := s.newObject(true)
		o.ID = v
		o.Path = s.getRelPath(v)
		o.Mode |= types.ModeDir

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
		return types.IterateDone
	}

	input.keyMarker = output.NextMarker
	return nil
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *types.ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	output, _, err := s.bucket.Get(ctx, &cos.BucketGetOptions{
		Prefix:  input.prefix,
		Marker:  input.keyMarker,
		MaxKeys: input.maxKeys,
	})
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
		return types.IterateDone
	}

	input.keyMarker = output.NextMarker
	return nil
}

func (s *Storage) nextPartObjectPageByPrefix(ctx context.Context, page *types.ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	listInput := &cos.ListMultipartUploadsOptions{
		Prefix:         input.prefix,
		MaxUploads:     input.maxKeys,
		KeyMarker:      input.keyMarker,
		UploadIDMarker: input.uploadIdMarker,
	}

	output, _, err := s.bucket.ListMultipartUploads(ctx, listInput)
	if err != nil {
		return err
	}

	for _, v := range output.Uploads {
		o := s.newObject(true)
		o.ID = v.Key
		o.Path = s.getRelPath(v.Key)
		o.Mode |= types.ModePart
		o.SetMultipartID(v.UploadID)

		page.Data = append(page.Data, o)
	}

	if !output.IsTruncated {
		return types.IterateDone
	}

	input.uploadIdMarker = output.NextUploadIDMarker
	input.keyMarker = output.NextKeyMarker
	return nil
}

func (s *Storage) nextPartPage(ctx context.Context, page *types.PartPage) error {
	input := page.Status.(*partPageStatus)

	output, _, err := s.object.ListParts(ctx, input.key, input.uploadId, &cos.ObjectListPartsOptions{
		MaxParts:         input.maxParts,
		PartNumberMarker: input.partNumberMarker,
	})
	if err != nil {
		return err
	}

	for _, v := range output.Parts {
		p := &types.Part{
			// The returned `PartNumber` is [1, 10000].
			// Set Index=*v.PartNumber-1 here to make the `PartNumber` zero-based for user.
			Index: v.PartNumber - 1,
			Size:  v.Size,
			ETag:  v.ETag,
		}

		page.Data = append(page.Data, p)
	}

	if !output.IsTruncated {
		return types.IterateDone
	}

	input.partNumberMarker = output.NextPartNumberMarker
	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)

	getOptions := &cos.ObjectGetOptions{}
	rangeOptions := &cos.RangeOptions{}
	if opt.HasOffset {
		rangeOptions.HasStart = true
		rangeOptions.Start = opt.Offset
	}
	if opt.HasSize {
		rangeOptions.HasEnd = true
		if opt.HasOffset {
			rangeOptions.End = rangeOptions.Start + opt.Size - 1
		} else {
			rangeOptions.HasStart = true
			rangeOptions.Start = 0
			rangeOptions.End = opt.Size - 1
		}
	}
	getOptions.Range = cos.FormatRangeOptions(rangeOptions)
	// SSE-C
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		getOptions.XCosSSECustomerAglo, getOptions.XCosSSECustomerKey, getOptions.XCosSSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return 0, err
		}
	}
	resp, err := s.object.Get(ctx, rp, getOptions)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	rc := resp.Body
	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}

	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	if opt.HasMultipartID {
		_, _, err = s.object.ListParts(ctx, rp, opt.MultipartID, nil)
		if err != nil {
			return nil, err
		}

		o = s.newObject(true)
		o.ID = rp
		o.Path = path
		o.Mode |= types.ModePart
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

	headOptions := &cos.ObjectHeadOptions{}
	// SSE-C
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		headOptions.XCosSSECustomerAglo, headOptions.XCosSSECustomerKey, headOptions.XCosSSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return
		}
	}
	output, err := s.object.Head(ctx, rp, headOptions)
	if err != nil {
		return nil, err
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o.Mode |= types.ModeDir
	} else {
		o.Mode |= types.ModeRead
	}

	o.SetContentLength(output.ContentLength)

	// COS uses RFC1123 format in HEAD
	//
	// > Last-Modified: Fri, 09 Aug 2019 10:20:56 GMT
	//
	// ref: https://cloud.tencent.com/document/product/436/7745
	if v := output.Header.Get(headers.LastModified); v != "" {
		lastModified, err := time.Parse(time.RFC1123, v)
		if err != nil {
			return nil, err
		}
		o.SetLastModified(lastModified)
	}

	if v := output.Header.Get(headers.ContentType); v != "" {
		o.SetContentType(v)
	}

	if v := output.Header.Get(headers.ETag); v != "" {
		o.SetEtag(v)
	}

	var sm ObjectSystemMetadata
	if v := output.Header.Get(storageClassHeader); v != "" {
		sm.StorageClass = v
	}
	if v := output.Header.Get(serverSideEncryptionHeader); v != "" {
		sm.ServerSideEncryption = v
	}
	if v := output.Header.Get(serverSideEncryptionCosKmsKeyIdHeader); v != "" {
		sm.ServerSideEncryptionCosKmsKeyID = v
	}
	if v := output.Header.Get(serverSideEncryptionCustomerAlgorithmHeader); v != "" {
		sm.ServerSideEncryptionCustomerAlgorithm = v
	}
	if v := output.Header.Get(serverSideEncryptionCustomerKeyMd5Header); v != "" {
		sm.ServerSideEncryptionCustomerKeyMd5 = v
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
	if r == nil && size == 0 {
		r = bytes.NewReader([]byte{})
	} else if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	} else {
		r = io.LimitReader(r, size)
	}

	rp := s.getAbsPath(path)

	putOptions := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: size,
		},
	}
	if opt.HasContentMd5 {
		putOptions.ContentMD5 = opt.ContentMd5
	}
	if opt.HasStorageClass {
		putOptions.XCosStorageClass = opt.StorageClass
	}
	// SSE-C
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		putOptions.XCosSSECustomerAglo, putOptions.XCosSSECustomerKey, putOptions.XCosSSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return
		}
	}
	// SSE-COS or SSE-KMS
	if opt.HasServerSideEncryption {
		putOptions.XCosServerSideEncryption = opt.ServerSideEncryption
		if opt.ServerSideEncryption == ServerSideEncryptionCosKms {
			// FIXME: we can remove the usage of `XOptionHeader` when cos' SDK supports SSE-KMS
			putOptions.XOptionHeader = &http.Header{}
			if opt.HasServerSideEncryptionCosKmsKeyID {
				putOptions.XOptionHeader.Set(serverSideEncryptionCosKmsKeyIdHeader, opt.ServerSideEncryptionCosKmsKeyID)
			}
			if opt.HasServerSideEncryptionContext {
				putOptions.XOptionHeader.Set(serverSideEncryptionContextHeader, opt.ServerSideEncryptionContext)
			}
		}
	}
	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	_, err = s.object.Put(ctx, rp, r, putOptions)
	if err != nil {
		return 0, err
	}

	return size, nil
}

func (s *Storage) writeMultipart(ctx context.Context, o *types.Object, r io.Reader, size int64, index int, opt pairStorageWriteMultipart) (n int64, part *types.Part, err error) {
	if size > multipartSizeMaximum {
		err = fmt.Errorf("size limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return
	}
	if index < 0 || index >= multipartNumberMaximum {
		err = fmt.Errorf("multipart number limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return
	}

	input := &cos.ObjectUploadPartOptions{
		ContentLength: size,
	}
	if opt.HasContentMd5 {
		input.ContentMD5 = opt.ContentMd5
	}

	// For COS, the `PartNumber` is [1, 10000]. But for users, the `PartNumber` is zero-based.
	// Set PartNumber=index+1 here to ensure pass in an effective `PartNumber` for `UploadPart`.
	// ref: https://cloud.tencent.com/document/product/436/7750
	output, err := s.object.UploadPart(ctx, o.ID, o.MustGetMultipartID(), index+1, r, input)
	if err != nil {
		return
	}

	part = &types.Part{
		Index: index,
		Size:  size,
		ETag:  output.Header.Get("ETag"),
	}
	return size, part, nil
}
