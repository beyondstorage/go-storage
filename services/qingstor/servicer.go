package qingstor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yunify/qingstor-sdk-go/v3/config"
	iface "github.com/yunify/qingstor-sdk-go/v3/interface"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
)

// Service is the qingstor service config.
type Service struct {
	config  *config.Config
	service iface.Service

	noRedirectClient *http.Client
}

// New will create a new qingstor service.
func New(pairs ...*types.Pair) (s *Service, err error) {
	const errorMessage = "%s New: %w"

	s = &Service{
		noRedirectClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	cred := opt.Credential.Value()
	cfg, err := config.New(cred.AccessKey, cred.SecretKey)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}
	if opt.HasEndpoint {
		ep := opt.Endpoint.Value()
		cfg.Host = ep.Host
		cfg.Port = ep.Port
		cfg.Protocol = ep.Protocol
	}

	s.config = cfg
	s.service, _ = service.Init(cfg)
	return
}

// String implements Service.String
func (s *Service) String() string {
	if s.config == nil {
		return fmt.Sprintf("Servicer qingstor")
	}
	return fmt.Sprintf("Servicer qingstor {Host: %s, Port: %d, Protocol: %s, AccessKey: %s}", s.config.Host, s.config.Port, s.config.Protocol, s.config.AccessKeyID)
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Create [%s]: %w"

	opt, err := parseServicePairCreate(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	// TODO: check bucket name here.

	bucket, err := s.service.Bucket(name, opt.Location)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	_, err = bucket.Put()
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return newStorage(bucket)
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	opt, err := parseServicePairDelete(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}
	bucket, err := s.get(name, opt.Location)
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, name, err)
	}
	_, err = bucket.Delete()
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, name, err)
	}
	return nil
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Get [%s]: %w"

	opt, err := parseServicePairGet(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	bucket, err := s.get(name, opt.Location)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return newStorage(bucket)
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) (err error) {
	const errorMessage = "%s List: %w"

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, err)
	}

	input := &service.ListBucketsInput{}
	if opt.HasLocation {
		input.Location = &opt.Location
	}

	// FIXME: List buckets could be incomplete.
	output, err := s.service.ListBuckets(input)
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, err)
	}

	for _, v := range output.Buckets {
		store, err := s.get(*v.Name, *v.Location)
		if err != nil {
			return fmt.Errorf(errorMessage, s, err)
		}
		if opt.HasStoragerFunc {
			c, err := newStorage(store)
			if err != nil {
				return fmt.Errorf(errorMessage, s, err)
			}
			opt.StoragerFunc(c)
		}
	}
	return nil
}

func (s *Service) get(name, location string) (*service.Bucket, error) {
	const errorMessage = "%s get [%s]: %w"

	if !IsBucketNameValid(name) {
		err := handleQingStorError(ErrInvalidBucketName)
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	// TODO: add bucket name check here.
	if location != "" {
		bucket, err := s.service.Bucket(name, location)
		if err != nil {
			err = handleQingStorError(err)
			return nil, fmt.Errorf(errorMessage, s, name, err)
		}
		return bucket, nil
	}

	url := fmt.Sprintf("%s://%s.%s:%d", s.config.Protocol, name, s.config.Host, s.config.Port)

	r, err := s.noRedirectClient.Head(url)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	if r.StatusCode != http.StatusTemporaryRedirect {
		err = fmt.Errorf("head status is %d instead of %d", r.StatusCode, http.StatusTemporaryRedirect)
		return nil, fmt.Errorf(errorMessage, s, name, handleQingStorError(err))
	}

	// Example URL: https://bucket.zone.qingstor.com
	location = strings.Split(r.Header.Get("Location"), ".")[1]
	bucket, err := s.service.Bucket(name, location)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return bucket, nil
}
