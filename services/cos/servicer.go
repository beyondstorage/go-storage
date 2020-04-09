package cos

import (
	"fmt"
	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
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

	opt, err := s.parsePairList(pairs...)
	if err != nil {
		return err
	}

	output, _, err := s.service.Service.Get(opt.Context)
	if err != nil {
		return err
	}
	for _, v := range output.Buckets {
		store, err := s.newStorage(ps.WithName(v.Name), ps.WithLocation(v.Region))
		if err != nil {
			return err
		}
		opt.StoragerFunc(store)
	}
	return
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (st storage.Storager, err error) {
	defer func() {
		err = s.formatError("get", err, name)
	}()

	opt, err := s.parsePairGet(pairs...)
	if err != nil {
		return nil, err
	}

	store, err := s.newStorage(ps.WithName(name), ps.WithLocation(opt.Location))
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

	opt, err := s.parsePairCreate(pairs...)
	if err != nil {
		return nil, err
	}

	store, err := s.newStorage(ps.WithName(name), ps.WithLocation(opt.Location))
	if err != nil {
		return nil, err
	}
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

	opt, err := s.parsePairDelete(pairs...)
	if err != nil {
		return err
	}

	store, err := s.newStorage(ps.WithName(name), ps.WithLocation(opt.Location))
	if err != nil {
		return err
	}
	_, err = store.bucket.Delete(opt.Context)
	if err != nil {
		return err
	}
	return
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...*types.Pair) (st *Storage, err error) {
	defer func() {
		err = s.formatError("new_storage", err, "")
	}()

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, err
	}

	st = &Storage{}

	url := cos.NewBucketURL(opt.Name, opt.Location, true)
	c := cos.NewClient(&cos.BaseURL{BucketURL: url}, s.client)

	st.bucket = c.Bucket
	st.object = c.Object
	st.name = opt.Name
	st.location = opt.Location
	st.workDir = opt.WorkDir
	return st, nil
}

func (s *Service) formatError(op string, err error, name string) error {
	if err == nil {
		return nil
	}

	return &services.ServiceError{
		Op:       op,
		Err:      formatError(err),
		Servicer: s,
		Name:     name,
	}
}
