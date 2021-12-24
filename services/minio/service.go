package minio

import (
	"context"

	"github.com/minio/minio-go/v7"

	ps "go.beyondstorage.io/v5/pairs"
	. "go.beyondstorage.io/v5/types"
)

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store Storager, err error) {
	paris := append(opt.pairs, ps.WithName(name))
	st, err := s.newStorage(paris...)
	if err != nil {
		return nil, err
	}
	err = s.service.MakeBucket(ctx, name, minio.MakeBucketOptions{})
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) delete(ctx context.Context, name string, opt pairServiceDelete) (err error) {
	err = s.service.RemoveBucket(ctx, name)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) get(ctx context.Context, name string, opt pairServiceGet) (store Storager, err error) {
	st, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) list(ctx context.Context, opt pairServiceList) (sti *StoragerIterator, err error) {
	input := &storagePageStatus{}
	return NewStoragerIterator(ctx, s.nextStoragePage, input), nil
}

func (s *Service) nextStoragePage(ctx context.Context, page *StoragerPage) (err error) {
	input := page.Status.(*storagePageStatus)
	input.buckets, err = s.service.ListBuckets(ctx)
	if err != nil {
		return err
	}
	for _, v := range input.buckets {
		store, err := s.newStorage(ps.WithName(v.Name))
		if err != nil {
			return err
		}
		page.Data = append(page.Data, store)
	}
	return IterateDone
}
