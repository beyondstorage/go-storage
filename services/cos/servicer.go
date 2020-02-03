package cos

import (
	"fmt"
	"net/http"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// Service is the Tencent oss *Service config.
type Service struct {
	service *cos.Client
	client  *http.Client
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer cos")
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) (err error) {
	const errorMessage = "%s List: %w"

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, err)
	}

	output, _, err := s.service.Service.Get(opt.Context)
	if err != nil {
		return fmt.Errorf(errorMessage, s, err)
	}
	for _, v := range output.Buckets {
		store := newStorage(v.Name, v.Region, s.client)
		opt.StoragerFunc(store)
	}
	return
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Get [%s]: %w"

	opt, err := parseServicePairGet(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	store := newStorage(name, opt.Location, s.client)
	return store, nil
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Create [%s]: %w"

	opt, err := parseServicePairCreate(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	store := newStorage(name, opt.Location, s.client)
	_, err = store.bucket.Put(opt.Context, nil)
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

	store := newStorage(name, opt.Location, s.client)
	_, err = store.bucket.Delete(opt.Context)
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}
	return
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...*types.Pair) (*Storage, error) {
	const errorMessage = "cos new_storage: %w"

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	store := &Storage{}

	url := cos.NewBucketURL(opt.Name, opt.Location, true)
	c := cos.NewClient(&cos.BaseURL{BucketURL: url}, s.client)
	store.bucket = c.Bucket
	store.object = c.Object
	store.name = opt.Name
	store.location = opt.Location
	return store, nil
}
