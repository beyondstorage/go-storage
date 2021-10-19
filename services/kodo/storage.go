package kodo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	qs "github.com/qiniu/go-sdk/v7/storage"

	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	. "go.beyondstorage.io/v5/types"
)

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

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *Object, err error) {
	if !s.features.VirtualDir {
		err = NewOperationNotImplementedError("create_dir")
		return
	}

	rp := s.getAbsPath(path)

	// Add `/` at the end of `path` to simulate a directory.
	// ref: https://developer.qiniu.com/kodo/kb/1705/how-to-create-the-folder-under-the-space
	rp += "/"

	// kodo `put` doesn't support `overwrite`, so we need to check whether the dir exists.
	// ref: [GSP-134](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/134-write-behavior-consistency.md)
	fi, err := s.bucket.Stat(s.name, rp)
	if err == nil {
		// The dir is exist.
		o = s.newObject(true)
		o.SetLastModified(convertUnixTimestampToTime(fi.PutTime))
		o.SetContentLength(fi.Fsize)

		var sm ObjectSystemMetadata
		sm.StorageClass = fi.Type
		o.SetSystemMetadata(sm)
	} else if !checkError(err, responseCodeResourceNotExist) {
		// Something error other then ResourceNotExist happened, return directly.
		return
	} else {
		// The dir is not exist, we should create the dir.
		uploader := qs.NewFormUploader(s.bucket.Cfg)
		ret := qs.PutRet{}
		err = uploader.Put(ctx,
			&ret, s.putPolicy.UploadToken(s.bucket.Mac), rp, io.LimitReader(nil, 0), 0, nil)
		if err != nil {
			return
		}

		o = s.newObject(false)
	}

	o.Path = path
	o.ID = rp
	o.Mode = ModeDir
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

	err = s.bucket.Delete(s.name, rp)
	if err != nil && checkError(err, responseCodeResourceNotExist) {
		// Omit `612`(resource to be deleted dose not exist) error code here
		//
		// References
		// - [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
		// - https://developer.qiniu.com/kodo/1257/delete
		err = nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	input := &objectPageStatus{
		limit:  1000,
		prefix: s.getAbsPath(path),
	}

	if !opt.HasListMode {
		// Support `ListModePrefix` as the default `ListMode`.
		// ref: [GSP-654](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/654-unify-list-behavior.md)
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
	return meta
}

func (s *Storage) nextObjectPageByDir(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	entries, commonPrefix, nextMarker, _, err := s.bucket.ListFiles(
		s.name,
		input.prefix,
		input.delimiter,
		input.marker,
		input.limit,
	)
	if err != nil {
		return err
	}

	for _, v := range commonPrefix {
		o := s.newObject(true)
		o.ID = v
		o.Path = s.getRelPath(v)
		o.Mode |= ModeDir

		page.Data = append(page.Data, o)
	}

	for _, v := range entries {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if nextMarker == "" {
		return IterateDone
	}

	input.marker = nextMarker
	return nil
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	entries, _, nextMarker, _, err := s.bucket.ListFiles(
		s.name,
		input.prefix,
		input.delimiter,
		input.marker,
		input.limit,
	)
	if err != nil {
		return err
	}

	for _, v := range entries {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if nextMarker == "" {
		return IterateDone
	}

	input.marker = nextMarker
	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)

	deadline := time.Now().Add(time.Hour).Unix()
	url := qs.MakePrivateURL(s.bucket.Mac, s.domain, rp, deadline)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	// ref: https://developer.qiniu.com/kodo/1232/download-process
	if opt.HasOffset && !opt.HasSize {
		rangeBytes := fmt.Sprintf("bytes=%d-", opt.Offset)
		req.Header.Add("Range", rangeBytes)
	} else if !opt.HasOffset && opt.HasSize {
		rangeBytes := fmt.Sprintf("bytes=0-%d", opt.Size-1)
		req.Header.Add("Range", rangeBytes)
	} else if opt.HasOffset && opt.HasSize {
		rangeBytes := fmt.Sprintf("bytes=%d-%d", opt.Offset, opt.Offset+opt.Size-1)
		req.Header.Add("Range", rangeBytes)
	}

	resp, err := s.bucket.Client.Do(ctx, req)
	if err != nil {
		return 0, err
	}

	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		err = qs.ResponseError(resp)
		return 0, err
	}

	rc := resp.Body

	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(resp.Body, opt.IoCallback)
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

	fi, err := s.bucket.Stat(s.name, rp)
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

	o.SetLastModified(convertUnixTimestampToTime(fi.PutTime))
	o.SetContentLength(fi.Fsize)

	if fi.Hash != "" {
		o.SetEtag(fi.Hash)
	}
	if fi.MimeType != "" {
		o.SetContentType(fi.MimeType)
	}

	var sm ObjectSystemMetadata
	sm.StorageClass = fi.Type
	o.SetSystemMetadata(sm)

	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	rp := s.getAbsPath(path)

	// According to GSP-751, we should allow the user to pass in a nil io.Reader.
	// ref: https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/751-write-empty-file-behavior.md
	if r == nil && size == 0 {
		r = bytes.NewReader([]byte{})
	} else if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	} else {
		r = io.LimitReader(r, size)
	}

	// kodo `put` doesn't support `overwrite`, so we need to check and delete the object if exists.
	// ref: [GSP-134](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/134-write-behavior-consistency.md)
	_, err = s.bucket.Stat(s.name, rp)
	if err == nil {
		err = s.bucket.Delete(s.name, rp)
		if err != nil {
			return
		}
	} else if !checkError(err, responseCodeResourceNotExist) {
		return
	}

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	uploader := qs.NewFormUploader(s.bucket.Cfg)
	ret := qs.PutRet{}
	err = uploader.Put(ctx,
		&ret, s.putPolicy.UploadToken(s.bucket.Mac), rp, r, size, nil)
	if err != nil {
		return
	}
	return size, nil
}
