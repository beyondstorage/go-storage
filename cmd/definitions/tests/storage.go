package tests

import (
	"context"
	"io"

	. "github.com/aos-dev/go-storage/v3/types"
)

type Storage struct {
	pairPolicy   PairPolicy
	defaultPairs DefaultStoragePairs
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

func (s *Storage) createMultipart(ctx context.Context, path string, opt pairStorageCreateMultipart) (o *Object, err error) {
	panic("not implemented")
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	panic("not implemented")
}

func (s *Storage) fetch(ctx context.Context, path string, url string, opt pairStorageFetch) (err error) {
	panic("not implemented")
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	panic("not implemented")
}

func (s *Storage) listMultipart(ctx context.Context, o *Object, opt pairStorageListMultipart) (pi *PartIterator, err error) {
	panic("not implemented")
}

func (s *Storage) metadata(ctx context.Context, opt pairStorageMetadata) (meta *StorageMeta, err error) {
	panic("not implemented")
}

func (s *Storage) move(ctx context.Context, src string, dst string, opt pairStorageMove) (err error) {
	panic("not implemented")
}

func (s *Storage) reach(ctx context.Context, path string, opt pairStorageReach) (url string, err error) {
	panic("not implemented")
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	panic("not implemented")
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	panic("not implemented")
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	panic("not implemented")
}

func (s *Storage) writeMultipart(ctx context.Context, o *Object, r io.Reader, size int64, index int, opt pairStorageWriteMultipart) (n int64, part *Part, err error) {
	panic("not implemented")
}
