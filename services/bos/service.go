package bos

import (
	"context"

	"go.beyondstorage.io/v5/types"
)

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store types.Storager, err error) {
	f := s.f
	f.Name = name
	st, err := f.newStorage()
	if err != nil {
		return nil, err
	}

	_, err = s.service.PutBucket(name)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (s *Service) delete(ctx context.Context, name string, opt pairServiceDelete) (err error) {
	err = s.service.DeleteBucket(name)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) get(ctx context.Context, name string, opt pairServiceGet) (store types.Storager, err error) {
	f := s.f
	f.Name = name
	st, err := f.newStorage()
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (s *Service) list(ctx context.Context, opt pairServiceList) (sti *types.StoragerIterator, err error) {
	input := &storagePageStatus{}

	return types.NewStoragerIterator(ctx, s.nextStoragePage, input), nil
}

func (s *Service) nextStoragePage(ctx context.Context, page *types.StoragerPage) error {
	output, err := s.service.ListBuckets()
	if err != nil {
		return err
	}

	for _, v := range output.Buckets {
		f := s.f
		f.Name = v.Name
		store, err := f.newStorage()
		if err != nil {
			return err
		}
		page.Data = append(page.Data, store)
	}

	return types.IterateDone
}
