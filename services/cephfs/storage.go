package cephfs

import (
	"context"
	"io"

	"github.com/beyondstorage/go-storage/v5/types"
)

func (s *Storage) create(path string, opt pairStorageCreate) (o *types.Object) {
	panic("not implemented")
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	panic("not implemented")
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *types.ObjectIterator, err error) {
	panic("not implemented")
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *types.StorageMeta) {
	panic("not implemented")
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	panic("not implemented")
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *types.Object, err error) {
	panic("not implemented")
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	panic("not implemented")
}
