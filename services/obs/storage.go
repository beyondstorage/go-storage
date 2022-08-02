package obs

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"

	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/pkg/iowrap"
	"github.com/beyondstorage/go-storage/v5/services"
	"github.com/beyondstorage/go-storage/v5/types"
)

func (s *Storage) create(path string, opt pairStorageCreate) (o *types.Object) {
	rp := s.getAbsPath(path)

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			return
		}

		rp += "/"
		o = s.newObject(true)
		o.Mode |= types.ModeDir
	} else {
		o = s.newObject(false)
		o.Mode |= types.ModeRead
	}

	o.ID = rp
	o.Path = path

	return o
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			err = services.PairUnsupportedError{Pair: ps.WithObjectMode(opt.ObjectMode)}
			return err
		}

		rp += "/"
	}

	input := &obs.DeleteObjectInput{
		Bucket: s.bucket,
		Key:    rp,
	}

	// obs DeleteObject is idempotent, so we don't need to check NoSuchKey error.
	//
	// References
	// - [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
	// - https://support.huaweicloud.com/api-obs/obs_04_0085.html#section0
	_, err = s.client.DeleteObject(input)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *types.ObjectIterator, err error) {
	if opt.ListMode.IsDir() {
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
	}

	input := &objectPageStatus{
		maxKeys: 200,
		prefix:  s.getAbsPath(path),
	}

	if !opt.HasListMode {
		// Support `ListModePrefix` as the default `ListMode`.
		// ref: [GSP-46](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/654-unify-list-behavior.md)
		opt.ListMode = types.ListModePrefix
	}

	var nextFn types.NextObjectFunc

	switch {
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

func (s *Storage) metadata(opt pairStorageMetadata) (meta *types.StorageMeta) {
	meta = types.NewStorageMeta()
	meta.Name = s.bucket
	meta.WorkDir = s.workDir
	return meta
}

func (s *Storage) nextObjectPageByDir(ctx context.Context, page *types.ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	listInput := &obs.ListObjectsInput{
		Bucket: s.bucket,
		Marker: input.marker,
	}
	listInput.Delimiter = input.delimiter
	listInput.Prefix = input.prefix
	listInput.MaxKeys = input.maxKeys

	output, err := s.client.ListObjects(listInput)
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

	if output.NextMarker == "" {
		return types.IterateDone
	}
	if !output.IsTruncated {
		return types.IterateDone
	}

	input.marker = output.NextMarker

	return err
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *types.ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	listInput := &obs.ListObjectsInput{
		Bucket: s.bucket,
		Marker: input.marker,
	}
	listInput.Delimiter = input.delimiter
	listInput.Prefix = input.prefix
	listInput.MaxKeys = input.maxKeys

	output, err := s.client.ListObjects(listInput)
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

	if output.NextMarker == "" {
		return types.IterateDone
	}
	if !output.IsTruncated {
		return types.IterateDone
	}

	input.marker = output.NextMarker

	return err
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)

	input := &obs.GetObjectInput{}
	input.Bucket = s.bucket
	input.Key = rp

	output, err := s.client.GetObject(input)
	if err != nil {
		return 0, err
	}

	rc := output.Body
	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}

	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			err = services.PairUnsupportedError{Pair: ps.WithObjectMode(opt.ObjectMode)}
			return nil, err
		}

		rp += "/"
	}

	input := &obs.GetObjectInput{}
	input.Bucket = s.bucket
	input.Key = rp

	output, err := s.client.GetObject(input)
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
	o.SetLastModified(output.LastModified)

	if output.ContentType != "" {
		o.SetContentType(output.ContentType)
	}
	if output.ETag != "" {
		o.SetEtag(output.ETag)
	}

	var sm ObjectSystemMetadata
	if v := output.StorageClass; v != "" {
		sm.StorageClass = string(v)
	}

	o.SetSystemMetadata(sm)

	return
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	if size > writeSizeMaximum {
		err = fmt.Errorf("size limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return 0, err
	}

	// According to GSP-751, we should allow the user to pass in a nil io.Reader.
	// Since obs supports reader passed in as nil, we do not need to determine the case where the reader is nil and the size is 0.
	// ref: https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/751-write-empty-file-behavior.md
	if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	}

	rp := s.getAbsPath(path)

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	input := &obs.PutObjectInput{
		Body: io.LimitReader(r, size),
	}
	input.Bucket = s.bucket
	input.Key = rp
	input.ContentLength = size

	if opt.HasContentMd5 {
		input.ContentMD5 = opt.ContentMd5
	}
	if opt.HasStorageClass {
		input.StorageClass = obs.StorageClassType(opt.StorageClass)
	}

	_, err = s.client.PutObject(input)
	if err != nil {
		return 0, err
	}

	return size, nil
}
