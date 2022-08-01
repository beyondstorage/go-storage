package tar

import (
	"context"
	"io"

	"github.com/beyondstorage/go-storage/v5/services"
	"github.com/beyondstorage/go-storage/v5/types"
)

func (s *Storage) create(path string, opt pairStorageCreate) (o *types.Object) {
	panic("not implemented")
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	panic("not implemented")
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *types.ObjectIterator, err error) {
	return types.NewObjectIterator(ctx, s.nextObjectPageByPrefix, nil), nil
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *types.StorageMeta) {
	panic("not implemented")
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *types.ObjectPage) error {
	page.Data = s.objects
	return types.IterateDone
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	idx, ok := s.objectsIndex[path]
	if !ok {
		return 0, services.ErrObjectNotExist
	}
	offset, _ := s.objectsOffset[path]

	_, err = s.file.Seek(offset, io.SeekStart)
	if err != nil {
		panic(err)
	}

	return io.CopyN(w, s.file, s.objects[idx].MustGetContentLength())
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *types.Object, err error) {
	idx, ok := s.objectsIndex[path]
	if !ok {
		return nil, services.ErrObjectNotExist
	}

	return s.objects[idx], nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	panic("not implemented")
}
