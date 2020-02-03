package gcs

import (
	"fmt"

	gs "cloud.google.com/go/storage"
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
	"google.golang.org/api/iterator"
)

// Service is the gcs config.
type Service struct {
	service   *gs.Client
	projectID string
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer gcs")
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) (err error) {
	const errorMessage = "%s List: %w"

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, err)
	}

	it := s.service.Buckets(opt.Context, s.projectID)
	for {
		bucketAttr, err := it.Next()
		// Next will return iterator.Done if there is no more items.
		if err != nil && err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf(errorMessage, s, err)
		}
		store, err := s.newStorage(ps.WithName(bucketAttr.Name))
		if err != nil {
			return fmt.Errorf(errorMessage, s, err)
		}
		opt.StoragerFunc(store)
	}
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Get [%s]: %w"

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return store, nil
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Create [%s]: %w"

	opt, err := parseServicePairCreate(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	err = store.bucket.Create(opt.Context, s.projectID, nil)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return store, nil
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	opt, err := parseServicePairDelete(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}
	err = store.bucket.Delete(opt.Context)
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}
	return nil
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...*types.Pair) (*Storage, error) {
	const errorMessage = "gcs new_storage: %w"

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	bucket := s.service.Bucket(opt.Name)

	store := &Storage{
		bucket: bucket,
		name:   opt.Name,
	}
	return store, nil
}
