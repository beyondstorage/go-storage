package kodo

import (
	"context"
	"fmt"

	qs "github.com/qiniu/api.v7/v7/storage"

	"github.com/Xuanwo/storage"
	ps "github.com/Xuanwo/storage/types/pairs"
)

func (s *Service) create(ctx context.Context, name string, opt *pairServiceCreate) (store storage.Storager, err error) {
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

	st, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) delete(ctx context.Context, name string, opt *pairServiceDelete) (err error) {
	err = s.service.DropBucket(name)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) get(ctx context.Context, name string, opt *pairServiceGet) (store storage.Storager, err error) {
	st, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) list(ctx context.Context, opt *pairServiceList) (err error) {
	buckets, err := s.service.Buckets(false)
	for _, v := range buckets {
		store, err := s.newStorage(ps.WithName(v))
		if err != nil {
			return err
		}
		opt.StoragerFunc(store)
	}
	return
}
