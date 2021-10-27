package tests

import (
	"context"
	"io"
	"net/http"
	"time"

	. "go.beyondstorage.io/v5/types"
)

type Storage struct {
	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	objects []*Object

	Pairs []Pair
	UnimplementedStorager
}

func (s *Storage) combineBlock(ctx context.Context, o *Object, bids []string, opt pairStorageCombineBlock) (err error) {
	panic("not implemented")
}

func (s *Storage) commitAppend(ctx context.Context, o *Object, opt pairStorageCommitAppend) (err error) {
	panic("not implemented")
}

func (s *Storage) completeMultipart(ctx context.Context, o *Object, parts []*Part, opt pairStorageCompleteMultipart) (err error) {
	panic("not implemented")
}

func (s *Storage) copy(ctx context.Context, src string, dst string, opt pairStorageCopy) (err error) {
	panic("not implemented")
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	panic("not implemented")
}

func (s *Storage) createAppend(ctx context.Context, path string, opt pairStorageCreateAppend) (o *Object, err error) {
	panic("not implemented")
}

func (s *Storage) createBlock(ctx context.Context, path string, opt pairStorageCreateBlock) (o *Object, err error) {
	panic("not implemented")
}

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *Object, err error) {
	panic("not implemented")
}

func (s *Storage) createLink(ctx context.Context, path string, target string, opt pairStorageCreateLink) (o *Object, err error) {
	panic("not implemented")
}

func (s *Storage) createMultipart(ctx context.Context, path string, opt pairStorageCreateMultipart) (o *Object, err error) {
	panic("not implemented")
}

func (s *Storage) createPage(ctx context.Context, path string, opt pairStorageCreatePage) (o *Object, err error) {
	panic("not implemented")
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	panic("not implemented")
}

func (s *Storage) fetch(ctx context.Context, path string, url string, opt pairStorageFetch) (err error) {
	panic("not implemented")
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	fn := NextObjectFunc(func(ctx context.Context, page *ObjectPage) error {
		page.Data = s.objects
		return nil
	})
	return NewObjectIterator(ctx, fn, nil), nil
}

func (s *Storage) listBlock(ctx context.Context, o *Object, opt pairStorageListBlock) (bi *BlockIterator, err error) {
	panic("not implemented")
}

func (s *Storage) listMultipart(ctx context.Context, o *Object, opt pairStorageListMultipart) (pi *PartIterator, err error) {
	panic("not implemented")
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	panic("not implemented")
}

func (s *Storage) move(ctx context.Context, src string, dst string, opt pairStorageMove) (err error) {
	panic("not implemented")
}

func (s *Storage) querySignHTTPCompleteMultipart(ctx context.Context, o *Object, parts []*Part, expire time.Duration, opt pairStorageQuerySignHTTPCompleteMultipart) (req *http.Request, err error) {
	panic("not implemented")
}

func (s *Storage) querySignHTTPCreateMultipart(ctx context.Context, path string, expire time.Duration, opt pairStorageQuerySignHTTPCreateMultipart) (req *http.Request, err error) {
	panic("not implemented")
}

func (s *Storage) querySignHTTPDelete(ctx context.Context, path string, expire time.Duration, opt pairStorageQuerySignHTTPDelete) (req *http.Request, err error) {
	panic("not implemented")
}

func (s *Storage) querySignHTTPListMultipart(ctx context.Context, o *Object, expire time.Duration, opt pairStorageQuerySignHTTPListMultipart) (req *http.Request, err error) {
	panic("not implemented")
}

func (s *Storage) querySignHTTPRead(ctx context.Context, path string, expire time.Duration, opt pairStorageQuerySignHTTPRead) (req *http.Request, err error) {
	panic("not implemented")
}

func (s *Storage) querySignHTTPWrite(ctx context.Context, path string, size int64, expire time.Duration, opt pairStorageQuerySignHTTPWrite) (req *http.Request, err error) {
	panic("not implemented")
}

func (s *Storage) querySignHTTPWriteMultipart(ctx context.Context, o *Object, size int64, index int, expire time.Duration, opt pairStorageQuerySignHTTPWriteMultipart) (req *http.Request, err error) {
	panic("not implemented")
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	panic("not implemented")
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	return nil, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	panic("not implemented")
}

func (s *Storage) writeAppend(ctx context.Context, o *Object, r io.Reader, size int64, opt pairStorageWriteAppend) (n int64, err error) {
	panic("not implemented")
}

func (s *Storage) writeBlock(ctx context.Context, o *Object, r io.Reader, size int64, bid string, opt pairStorageWriteBlock) (n int64, err error) {
	panic("not implemented")
}

func (s *Storage) writeMultipart(ctx context.Context, o *Object, r io.Reader, size int64, index int, opt pairStorageWriteMultipart) (n int64, part *Part, err error) {
	panic("not implemented")
}

func (s *Storage) writePage(ctx context.Context, o *Object, r io.Reader, size int64, offset int64, opt pairStorageWritePage) (n int64, err error) {
	panic("not implemented")
}
