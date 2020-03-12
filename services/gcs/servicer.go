package gcs

import (
	"errors"
	"fmt"

	gs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// Service is the gcs config.
type Service struct {
	service   *gs.Client
	projectID string

	loose bool
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer gcs")
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list", err, "")
	}()

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return err
	}

	it := s.service.Buckets(opt.Context, s.projectID)
	for {
		bucketAttr, err := it.Next()
		// Next will return iterator.Done if there is no more items.
		if err != nil && err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		store, err := s.newStorage(ps.WithName(bucketAttr.Name))
		if err != nil {
			return err
		}
		opt.StoragerFunc(store)
	}
}

// Get implements Servicer.Get
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

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (st storage.Storager, err error) {
	defer func() {
		err = s.formatError("create", err, name)
	}()

	opt, err := parseServicePairCreate(pairs...)
	if err != nil {
		return nil, err
	}

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	err = store.bucket.Create(opt.Context, s.projectID, nil)
	if err != nil {
		return nil, err
	}
	return store, nil
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("delete", err, name)
	}()

	opt, err := parseServicePairDelete(pairs...)
	if err != nil {
		return err
	}

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return err
	}
	err = store.bucket.Delete(opt.Context)
	if err != nil {
		return err
	}
	return nil
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...*types.Pair) (st *Storage, err error) {
	defer func() {
		err = s.formatError("new storage", err, "")
	}()

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, err
	}

	bucket := s.service.Bucket(opt.Name)

	store := &Storage{
		bucket: bucket,
		name:   opt.Name,

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
