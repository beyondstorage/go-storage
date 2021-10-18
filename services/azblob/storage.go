package azblob

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"

	"github.com/Azure/azure-storage-blob-go/azblob"

	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	. "go.beyondstorage.io/v5/types"
)

func (s *Storage) commitAppend(ctx context.Context, o *Object, opt pairStorageCommitAppend) (err error) {
	return
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	rp := s.getAbsPath(path)

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

	o.ID = rp
	o.Path = path
	return o
}

func (s *Storage) createAppend(ctx context.Context, path string, opt pairStorageCreateAppend) (o *Object, err error) {
	rp := s.getAbsPath(path)

	headers := azblob.BlobHTTPHeaders{}
	if opt.HasContentType {
		headers.ContentType = opt.ContentType
	}

	var cpk azblob.ClientProvidedKeyOptions
	if opt.HasEncryptionKey {
		cpk, err = calculateEncryptionHeaders(opt.EncryptionKey, opt.EncryptionScope)
		if err != nil {
			return
		}
	}

	_, err = s.bucket.NewAppendBlobURL(rp).Create(ctx, headers, nil,
		azblob.BlobAccessConditions{}, nil, cpk)
	if err != nil {
		return
	}

	o = s.newObject(true)
	o.Mode = ModeRead | ModeAppend
	o.ID = rp
	o.Path = path
	o.SetAppendOffset(0)
	return o, nil
}

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *Object, err error) {
	if !s.features.VirtualDir {
		err = NewOperationNotImplementedError("create_dir")
		return
	}

	rp := s.getAbsPath(path)

	// Specify a character or string delimiter within a blob name to create a virtual hierarchy.
	// ref: https://docs.microsoft.com/en-us/rest/api/storageservices/naming-and-referencing-containers--blobs--and-metadata#resource-names
	rp += "/"

	accessTier := azblob.AccessTierNone
	if opt.HasAccessTier {
		accessTier = azblob.AccessTierType(opt.AccessTier)
	}

	_, err = s.bucket.NewBlockBlobURL(rp).Upload(
		ctx, iowrap.SizedReadSeekCloser(nil, 0), azblob.BlobHTTPHeaders{},
		azblob.Metadata{}, azblob.BlobAccessConditions{},
		accessTier, azblob.BlobTagsMap{}, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.Mode |= ModeDir
	return
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			err = services.PairUnsupportedError{Pair: ps.WithObjectMode(opt.ObjectMode)}
			return
		}

		rp += "/"
	}

	_, err = s.bucket.NewBlockBlobURL(rp).Delete(ctx,
		azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil && checkError(err, azblob.ServiceCodeBlobNotFound) {
		// Omit `BlobNotFound` error here
		// ref: [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
		err = nil
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	input := &objectPageStatus{
		maxResults: 200,
		prefix:     s.getAbsPath(path),
	}

	if !opt.HasListMode {
		// Support `ListModePrefix` as the default `ListMode`.
		opt.ListMode = ListModePrefix
	}

	var nextFn NextObjectFunc

	switch {
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

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	meta = NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	meta.SetWriteSizeMaximum(WriteSizeMaximum)
	meta.SetAppendSizeMaximum(AppendSizeMaximum)
	meta.SetAppendNumberMaximum(AppendNumberMaximum)
	meta.SetAppendTotalSizeMaximum(AppendBlobIfMaxSizeLessThanOrEqual)
	return meta
}

func (s *Storage) nextObjectPageByDir(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	output, err := s.bucket.ListBlobsHierarchySegment(ctx, input.marker, input.delimiter, azblob.ListBlobsSegmentOptions{
		Prefix:     input.prefix,
		MaxResults: input.maxResults,
	})
	if err != nil {
		return err
	}

	for _, v := range output.Segment.BlobPrefixes {
		o := s.newObject(true)
		o.ID = v.Name
		o.Path = s.getRelPath(v.Name)
		o.Mode |= ModeDir

		page.Data = append(page.Data, o)
	}

	for _, v := range output.Segment.BlobItems {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if !output.NextMarker.NotDone() {
		return IterateDone
	}

	input.marker = output.NextMarker
	return nil
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	output, err := s.bucket.ListBlobsFlatSegment(ctx, input.marker, azblob.ListBlobsSegmentOptions{
		Prefix:     input.prefix,
		MaxResults: input.maxResults,
	})
	if err != nil {
		return err
	}

	for _, v := range output.Segment.BlobItems {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if !output.NextMarker.NotDone() {
		return IterateDone
	}

	input.marker = output.NextMarker
	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)

	offset := int64(0)
	if opt.HasOffset {
		offset = opt.Offset
	}

	count := int64(azblob.CountToEnd)
	if opt.HasSize {
		count = opt.Size
	}

	var cpk azblob.ClientProvidedKeyOptions
	if opt.HasEncryptionKey {
		cpk, err = calculateEncryptionHeaders(opt.EncryptionKey, opt.EncryptionScope)
		if err != nil {
			return 0, err
		}
	}
	output, err := s.bucket.NewBlockBlobURL(rp).Download(
		ctx, offset, count,
		azblob.BlobAccessConditions{}, false, cpk)
	if err != nil {
		return 0, err
	}
	defer func() {
		cerr := output.Response().Body.Close()
		if cerr != nil {
			err = cerr
		}
	}()

	rc := output.Response().Body
	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}

	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	rp := s.getAbsPath(path)

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			err = services.PairUnsupportedError{Pair: ps.WithObjectMode(opt.ObjectMode)}
			return
		}

		rp += "/"
	}

	var cpk azblob.ClientProvidedKeyOptions
	if opt.HasEncryptionKey {
		cpk, err = calculateEncryptionHeaders(opt.EncryptionKey, opt.EncryptionScope)
		if err != nil {
			return
		}
	}

	output, err := s.bucket.NewBlockBlobURL(rp).GetProperties(ctx, azblob.BlobAccessConditions{}, cpk)
	if err != nil {
		return nil, err
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o.Mode |= ModeDir
	} else {
		o.Mode |= ModeRead
	}

	o.SetContentLength(output.ContentLength())
	o.SetLastModified(output.LastModified())

	if v := string(output.ETag()); v != "" {
		o.SetEtag(v)
	}
	if v := output.ContentType(); v != "" {
		o.SetContentType(v)
	}
	if v := output.ContentMD5(); len(v) > 0 {
		o.SetContentMd5(base64.StdEncoding.EncodeToString(v))
	}

	var sm ObjectSystemMetadata
	if v := output.AccessTier(); v != "" {
		sm.AccessTier = v
	}
	if v := output.EncryptionKeySha256(); v != "" {
		sm.EncryptionKeySha256 = v
	}
	if v := output.EncryptionScope(); v != "" {
		sm.EncryptionScope = v
	}
	if v, err := strconv.ParseBool(output.IsServerEncrypted()); err == nil {
		sm.ServerEncrypted = v
	}
	o.SetSystemMetadata(sm)

	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	if size > WriteSizeMaximum {
		err = fmt.Errorf("size limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return
	}

	// According to GSP-751, we should allow the user to pass in a nil io.Reader.
	// ref: https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/751-write-empty-file-behavior.md
	if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	} else {
		r = io.LimitReader(r, size)
	}

	rp := s.getAbsPath(path)

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	accessTier := azblob.AccessTierNone
	if opt.HasAccessTier {
		accessTier = azblob.AccessTierType(opt.AccessTier)
	}

	headers := azblob.BlobHTTPHeaders{}
	if opt.HasContentMd5 {
		headers.ContentMD5, err = base64.StdEncoding.DecodeString(opt.ContentMd5)
		if err != nil {
			return 0, err
		}
	}
	if opt.HasContentType {
		headers.ContentType = opt.ContentType
	}

	var cpk azblob.ClientProvidedKeyOptions
	if opt.HasEncryptionKey {
		cpk, err = calculateEncryptionHeaders(opt.EncryptionKey, opt.EncryptionScope)
		if err != nil {
			return 0, err
		}
	}
	_, err = s.bucket.NewBlockBlobURL(rp).Upload(
		ctx, iowrap.SizedReadSeekCloser(r, size),
		headers, azblob.Metadata{}, azblob.BlobAccessConditions{},
		accessTier, azblob.BlobTagsMap{}, cpk)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (s *Storage) writeAppend(ctx context.Context, o *Object, r io.Reader, size int64, opt pairStorageWriteAppend) (n int64, err error) {
	if size > AppendSizeMaximum {
		err = fmt.Errorf("size limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return
	}
	rp := o.GetID()

	offset, _ := o.GetAppendOffset()

	var cpk azblob.ClientProvidedKeyOptions
	if opt.HasEncryptionKey {
		cpk, err = calculateEncryptionHeaders(opt.EncryptionKey, opt.EncryptionScope)
		if err != nil {
			return
		}
	}

	var ac azblob.AppendBlobAccessConditions
	ac.AppendPositionAccessConditions.IfMaxSizeLessThanOrEqual = AppendBlobIfMaxSizeLessThanOrEqual
	if 0 == offset {
		ac.AppendPositionAccessConditions.IfAppendPositionEqual = -1
	} else {
		ac.AppendPositionAccessConditions.IfAppendPositionEqual = offset
	}

	appendResp, err := s.bucket.NewAppendBlobURL(rp).AppendBlock(
		ctx, iowrap.SizedReadSeekCloser(r, size),
		ac, nil, cpk)
	if err != nil {
		return
	}

	// BlobAppendOffset() returns the offset at which the block was committed, in bytes, but seems not the next append position.
	// ref: https://github.com/Azure/azure-storage-blob-go/blob/master/azblob/zt_url_append_blob_test.go
	offset, err = strconv.ParseInt(appendResp.BlobAppendOffset(), 10, 64)
	if err != nil {
		return
	}
	offset += size
	o.SetAppendOffset(offset)

	return size, nil
}
