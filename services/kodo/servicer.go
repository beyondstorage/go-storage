package kodo

import (
	"errors"
	"fmt"

	"github.com/qiniu/api.v7/auth/qbox"
	qs "github.com/qiniu/api.v7/storage"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/types"
)

// Service is the kodo config.
type Service struct {
	service *qs.BucketManager
}

// New will create a new kodo service.
func New(pairs ...*types.Pair) (s *Service, err error) {
	const errorMessage = "%s New: %w"

	s = &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	mac := qbox.NewMac(cred[0], cred[1])
	cfg := &qs.Config{}
	s.service = qs.NewBucketManager(mac, cfg)
	return
}

// String implements Service.String
func (s Service) String() string {
	return fmt.Sprintf("Servicer kodo")
}

// List implements Service.List
func (s Service) List(pairs ...*types.Pair) (err error) {
	const errorMessage = "%s List: %w"

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, err)
	}

	buckets, err := s.service.Buckets(false)
	for _, v := range buckets {
		store := newStorage(s.service, v)
		opt.StoragerFunc(store)
	}
	return
}

// Get implements Service.Get
func (s Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const _ = "%s Get [%s]: %w"

	c := newStorage(s.service, name)
	return c, nil
}

// Create implements Service.Create
func (s Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Create [%s]: %w"

	// It looks like kodo go sdk doesn't provide bucket create API.
	return nil, fmt.Errorf(errorMessage, s, name, errors.New("not supported action"))
}

// Delete implements Service.Delete
func (s Service) Delete(name string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	// It looks like kodo go sdk doesn't provide bucket create API.
	return fmt.Errorf(errorMessage, s, name, errors.New("not supported action"))
}
