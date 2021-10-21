package uss

import (
	"context"
	"io"
	"strconv"
	"sync"

	"github.com/upyun/go-sdk/v3/upyun"

	"go.beyondstorage.io/v5/pkg/headers"
	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	. "go.beyondstorage.io/v5/types"
)

const (
	// headerListIter indicates the start marker of page in header
	// headerListLimit indicates the size of page in response
	headerListIter  = "X-List-Iter"
	headerListLimit = "X-List-Limit"
	// iterEnd is Base64 code which indicates the last page of list
	// more detail at: https://docs.upyun.com/api/rest_api/#_13
	iterEnd = "g2gCZAAEbmV4dGQAA2VvZg"
)

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o = s.newObject(true)
		o.Mode = ModeDir
	} else {
		o = s.newObject(false)
		o.Mode = ModeRead
	}

	o.ID = s.getAbsPath(path)
	o.Path = path
	return o
}

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *Object, err error) {
	rp := s.getAbsPath(path)

	err = s.bucket.Mkdir(rp)
	if err != nil && checkErrorCode(err, responseCodeFolderAlreadyExist) {
		// Omit `folder already exists` error here.
		err = nil
	}
	if err != nil {
		return
	}

	o = s.newObject(false)
	o.Mode = ModeDir
	o.ID = rp
	o.Path = path
	return o, nil
}

// delete implements Storager.Delete
//
// USS requires a short time between PUT and DELETE, or we will get this error:
// DELETE 429 {"msg":"concurrent put or delete","code":42900007,"id":"xxx"}
//
// Due to this problem, uss can't pass the storager integration tests.
func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	config := &upyun.DeleteObjectConfig{
		Path: rp,
	}

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		config.Folder = true
	}

	err = s.bucket.Delete(config)
	if err != nil && checkErrorCode(err, responseCodeFileOrDirectoryNotFound) {
		// Omit `file or directory not found` error here.
		// ref: [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
		err = nil
	}
	if err != nil {
		return err
	}
	return
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	input := &objectPageStatus{
		// 50 is the recommended value in SDK
		// see more details at: https://github.com/upyun/go-sdk/blob/master/upyun/rest.go#L560
		limit:  "50",
		prefix: s.getAbsPath(path),
	}
	if opt.HasContinuationToken {
		input.iter = opt.ContinuationToken
	}

	if !opt.HasListMode {
		// Support `ListModePrefix` as the default `ListMode`.
		// ref: [GSP-654](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/654-unify-list-behavior.md)
		opt.ListMode = ListModePrefix
	}

	var nextFn NextObjectFunc

	switch {
	case opt.ListMode.IsDir():
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

func (s *Storage) nextObjectPageByDir(ctx context.Context, page *ObjectPage) (err error) {
	input := page.Status.(*objectPageStatus)

	header := make(map[string]string)
	header[headerListLimit] = input.limit
	header[headerListIter] = input.iter

	// err could be updated in multiple goroutines, add explict lock to protect it.
	var errlock sync.Mutex

	// USS SDK will close this channel in List
	ch := make(chan *upyun.FileInfo, 1)

	go func() {
		xerr := s.bucket.List(&upyun.GetObjectsConfig{
			Path:         input.prefix,
			ObjectsChan:  ch,
			MaxListLevel: 1, // 1 means not recursive
			Headers:      header,
		})

		errlock.Lock()
		defer errlock.Unlock()
		err = xerr
	}()

	for v := range ch {
		if v.IsDir {
			o := s.newObject(true)
			o.ID = v.Name
			o.Path = s.getRelPath(v.Name)
			o.Mode |= ModeDir
			// v.Meta means all the k-v in header with key which has prefix `x-upyun-meta-`
			// so we consider it as user's metadata
			// see more details at: https://github.com/upyun/go-sdk/blob/master/upyun/fileinfo.go#L39
			o.SetUserMetadata(v.Meta)

			page.Data = append(page.Data, o)
			continue
		}

		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if header[headerListIter] == iterEnd {
		return IterateDone
	}

	input.iter = header[headerListIter]
	return nil
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *ObjectPage) (err error) {
	input := page.Status.(*objectPageStatus)

	header := make(map[string]string)
	header[headerListLimit] = input.limit
	header[headerListIter] = input.iter

	// err could be updated in multiple goroutines, add explict lock to protect it.
	var errlock sync.Mutex

	// USS SDK will close this channel in List
	ch := make(chan *upyun.FileInfo, 1)

	go func() {
		xerr := s.bucket.List(&upyun.GetObjectsConfig{
			Path:         input.prefix,
			ObjectsChan:  ch,
			MaxListLevel: -1, // -1 means recursive
			Headers:      header,
		})

		errlock.Lock()
		defer errlock.Unlock()
		err = xerr
	}()

	for v := range ch {
		if v.IsDir {
			continue
		}

		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if header[headerListIter] == iterEnd {
		return IterateDone
	}

	input.iter = header[headerListIter]
	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)

	config := &upyun.GetObjectConfig{
		Path:   rp,
		Writer: w,
	}

	if opt.HasIoCallback {
		config.Writer = iowrap.CallbackWriter(w, opt.IoCallback)
	}

	f, err := s.bucket.Get(config)

	if err != nil {
		return 0, err
	}

	return f.Size, nil
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	rp := s.getAbsPath(path)

	output, err := s.bucket.GetInfo(rp)
	if err != nil {
		return nil, err
	}

	return s.formatFileObject(output)
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	rp := s.getAbsPath(path)

	cfg := &upyun.PutObjectConfig{
		Path:   rp,
		Reader: r,
		Headers: map[string]string{
			headers.ContentLength: strconv.FormatInt(size, 10),
		},
	}

	err = s.bucket.Put(cfg)
	if err != nil {
		return
	}
	return size, nil
}
