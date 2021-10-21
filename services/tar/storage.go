package tar

import (
	"context"
	"io"

	"go.beyondstorage.io/v5/services"
	. "go.beyondstorage.io/v5/types"
)

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	panic("not implemented")
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	panic("not implemented")
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	return NewObjectIterator(ctx, s.nextObjectPageByPrefix, nil), nil
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	panic("not implemented")
}

func (s *Storage) nextObjectPageByPrefix(ctx context.Context, page *ObjectPage) error {
	page.Data = s.objects
	return IterateDone
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	idx, ok := s.objectsIndex[path]
	if !ok {
		return 0, services.ErrObjectNotExist
	}
	offset, _ := s.objectsOffset[path]

	_, err = s.f.Seek(offset, io.SeekStart)
	if err != nil {
		panic(err)
	}

	return io.CopyN(w, s.f, s.objects[idx].MustGetContentLength())
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	idx, ok := s.objectsIndex[path]
	if !ok {
		return nil, services.ErrObjectNotExist
	}

	return s.objects[idx], nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	panic("not implemented")
}
