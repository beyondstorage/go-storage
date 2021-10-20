package oss

import (
	"context"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	ps "go.beyondstorage.io/v5/pairs"
	typ "go.beyondstorage.io/v5/types"
)

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store typ.Storager, err error) {
	st, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	err = s.service.CreateBucket(name)
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

func (s *Service) get(ctx context.Context, name string, opt pairServiceGet) (store typ.Storager, err error) {
	st, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) list(ctx context.Context, opt pairServiceList) (it *typ.StoragerIterator, err error) {
	input := &storagePageStatus{
		maxKeys: 200,
	}

	return typ.NewStoragerIterator(ctx, s.nextStoragePage, input), nil
}

func (s *Service) nextStoragePage(ctx context.Context, page *typ.StoragerPage) error {
	input := page.Status.(*storagePageStatus)

	output, err := s.service.ListBuckets(
		oss.Marker(input.marker),
		oss.MaxKeys(input.maxKeys),
	)
	if err != nil {
		return err
	}

	for _, v := range output.Buckets {
		store, err := s.newStorage(ps.WithName(v.Name))
		if err != nil {
			return err
		}

		page.Data = append(page.Data, store)
	}

	if !output.IsTruncated {
		return typ.IterateDone
	}

	input.marker = output.NextMarker
	return nil
}
