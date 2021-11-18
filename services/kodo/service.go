package kodo

import (
	"context"
	"fmt"

	qs "github.com/qiniu/go-sdk/v7/storage"

	typ "go.beyondstorage.io/v5/types"
)

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store typ.Storager, err error) {
	// Check region ID.
	_, ok := qs.GetRegionByID(qs.RegionID(opt.Location))
	if !ok {
		err = fmt.Errorf("region %s is invalid", opt.Location)
		return nil, err
	}

	err = s.service.CreateBucket(name, qs.RegionID(opt.Location))
	if err != nil {
		return nil, err
	}

	f := s.f
	f.Name = name
	st, err := f.newStorage()
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) delete(ctx context.Context, name string, opt pairServiceDelete) (err error) {
	err = s.service.DropBucket(name)
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
	input := &storagePageStatus{}

	return typ.NewStoragerIterator(ctx, s.nextStoragePage, input), nil
}

func (s *Service) nextStoragePage(ctx context.Context, page *typ.StoragerPage) error {
	buckets, err := s.service.Buckets(true)
	if err != nil {
		return err
	}

	for _, v := range buckets {
		f := s.f
		f.Name = v
		st, err := f.newStorage()
		if err != nil {
			return err
		}

		page.Data = append(page.Data, st)
	}

	return typ.IterateDone
}
