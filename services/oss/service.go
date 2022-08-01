package oss

import (
	"context"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	typ "github.com/beyondstorage/go-storage/v5/types"
)

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store typ.Storager, err error) {
	f := s.f
	f.Name = name
	st, err := f.newStorage()
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
	f := s.f
	f.Name = name
	st, err := f.newStorage()
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
		f := s.f
		f.Name = v.Name
		st, err := f.newStorage()
		if err != nil {
			return err
		}

		page.Data = append(page.Data, st)
	}

	if !output.IsTruncated {
		return typ.IterateDone
	}

	input.marker = output.NextMarker
	return nil
}
