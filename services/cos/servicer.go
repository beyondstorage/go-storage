package cos

import (
	"fmt"
	"net/http"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/services"
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
	defer func() {
		err = s.formatError("list", err, "")
	}()

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return err
	}

	output, _, err := s.service.Service.Get(opt.Context)
	if err != nil {
		return err
	}
	for _, v := range output.Buckets {
		store := newStorage(v.Name, v.Region, s.client)
		opt.StoragerFunc(store)
	}
	return
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (st storage.Storager, err error) {
	defer func() {
		err = s.formatError("get", err, name)
	}()

	opt, err := parseServicePairGet(pairs...)
	if err != nil {
		return nil, err
	}

	store := newStorage(name, opt.Location, s.client)
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

	store := newStorage(name, opt.Location, s.client)
	_, err = store.bucket.Put(opt.Context, nil)
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

	store := newStorage(name, opt.Location, s.client)
	_, err = store.bucket.Delete(opt.Context)
	if err != nil {
		return err
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
	store.workDir = opt.WorkDir
	return store, nil
}

func (s *Service) formatError(op string, err error, name string) error {
	if err == nil {
		return nil
	}

	// Handle errors returned by cos.
	e, ok := err.(*cos.ErrorResponse)
	if ok {
		err = formatCosError(e)
	}

	return &services.ServiceError{
		Op:       op,
		Err:      err,
		Servicer: s,
		Name:     name,
	}
}
