package onedrive

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	. "go.beyondstorage.io/v5/types"
)

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	o = s.newObject(false)
	o.Mode = ModeRead
	o.ID = s.getAbsPath(path)
	o.Path = path

	return o
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	absPath := s.getAbsPath(path)
	err = s.client.DeleteItem(ctx, absPath)
	if err != nil {
		if isNotFoundError(err) {
			err = nil
		}

		return err
	}

	return
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	input := &objectPageStatus{
		// need to passed by pairs?
		limit: 200,
		rp:    s.getAbsPath(path),
		dir:   filepath.ToSlash(path),

		continuationToken: opt.ContinuationToken,
	}

	if !opt.HasListMode || opt.ListMode.IsDir() {
		return NewObjectIterator(ctx, s.nextObjectPage, input), nil
	} else {
		return nil, services.ListModeInvalidError{Actual: opt.ListMode}
	}
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	meta = NewStorageMeta()
	meta.WorkDir = s.workDir

	return
}

func (s *Storage) nextObjectPage(ctx context.Context, page *ObjectPage) error {
	iteratorObj := page.Status.(*objectPageStatus)

	if iteratorObj.done {
		return IterateDone
	}

	items, continuationToken, err := s.client.List(ctx, iteratorObj.rp, iteratorObj.continuationToken, iteratorObj.limit)
	if err != nil {
		// empty dir
		if isNotFoundError(err) {
			err = nil
			iteratorObj.done = true
		}
		return err
	}

	// set skipToken(continuationToken)
	iteratorObj.continuationToken = continuationToken

	// make a new buffer at one time
	// avoid extra spending(append)
	page.Data = make([]*Object, len(items))
	for k, v := range items {
		page.Data[k] = s.formatObject(v, iteratorObj.dir, iteratorObj.rp)
	}

	if uint32(len(items)) < iteratorObj.limit {
		iteratorObj.done = true
	}

	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	var rc io.ReadCloser

	objectPath := s.getAbsPath(path)

	object, size, err := s.client.DownloadItem(ctx, objectPath)
	if err != nil {
		return n, err
	}

	if opt.HasOffset {
		if size < opt.Offset {
			return
		}

		// imitate seek
		tempbuffer := make([]byte, opt.Offset)
		_, err = object.Read(tempbuffer)
		if err != nil {
			return n, err
		}
	}

	rc = object

	if opt.HasSize {
		rc = iowrap.LimitReadCloser(rc, opt.Size)
	}

	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}

	return io.Copy(w, object)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	uniquePath := s.getAbsPath(path)
	objInfo, err := s.client.GetItem(ctx, uniquePath)
	if err != nil {
		return o, err
	}

	o = s.newObject(true)
	o.ID = uniquePath
	o.Path = path
	o.Mode |= ModeRead
	o.SetEtag(objInfo.Etag)
	o.SetLastModified(objInfo.LastModifiedDateTime)
	o.SetContentLength(objInfo.Size)

	return
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not nil")
	}

	n, err = s.client.Upload(ctx, s.getAbsPath(path), size, r, opt.Description)

	return
}
