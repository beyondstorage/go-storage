package kodo

import (
	"errors"
	"fmt"

	qs "github.com/qiniu/api.v7/v7/storage"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// Service is the kodo config.
type Service struct {
	service *qs.BucketManager

	loose bool
}

// String implements Service.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer kodo")
}

// List implements Service.List
func (s *Service) List(pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list", err, "")
	}()

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return err
	}

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

// Get implements Service.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (st storage.Storager, err error) {
	defer func() {
		err = s.formatError("get", err, name)
	}()

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	return store, nil
}

// Create implements Service.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (st storage.Storager, err error) {
	defer func() {
		err = s.formatError("create", err, name)
	}()

	opt, err := parseServicePairCreate(pairs...)
	if err != nil {
		return nil, err
	}

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

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	return store, nil
}

// Delete implements Service.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("delete", err, name)
	}()

	err = s.service.DropBucket(name)
	if err != nil {
		return err
	}
	return nil
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...*types.Pair) (store *Storage, err error) {
	defer func() {
		err = s.formatError("new storage", err, "")
	}()

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, err
	}

	// Get bucket's domain.
	domains, err := s.service.ListBucketDomains(opt.Name)
	if err != nil {
		return nil, err
	}
	// TODO: we need to choose user's production domain.
	if len(domains) == 0 {
		return nil, errors.New("no available domains")
	}

	store = &Storage{
		bucket: s.service,
		domain: domains[0].Domain,
		putPolicy: qs.PutPolicy{
			Scope: opt.Name,
		},

		name:    opt.Name,
		workDir: opt.WorkDir,
		loose:   opt.Loose || s.loose,
	}
	return store, nil
}

func (s *Service) formatError(op string, err error, name string) error {
	if err == nil {
		return nil
	}

	if s.loose && errors.Is(err, services.ErrCapabilityInsufficient) {
		return nil
	}

	return &services.ServiceError{
		Op:       op,
		Err:      formatError(err),
		Servicer: s,
		Name:     name,
	}
}
