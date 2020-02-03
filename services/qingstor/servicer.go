package qingstor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
	"github.com/yunify/qingstor-sdk-go/v3/config"
	iface "github.com/yunify/qingstor-sdk-go/v3/interface"
	"github.com/yunify/qingstor-sdk-go/v3/service"
)

// Service is the qingstor service config.
type Service struct {
	config  *config.Config
	service iface.Service

	noRedirectClient *http.Client
}

// String implements Service.String
func (s *Service) String() string {
	if s.config == nil {
		return fmt.Sprintf("Servicer qingstor")
	}
	return fmt.Sprintf("Servicer qingstor {Host: %s, Port: %d, Protocol: %s, AccessKey: %s}", s.config.Host, s.config.Port, s.config.Protocol, s.config.AccessKeyID)
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
		store, err := s.newStorage(ps.WithName(*v.Name), ps.WithLocation(*v.Location))
		if err != nil {
			return fmt.Errorf(errorMessage, s, err)
		}
		opt.StoragerFunc(store)
	}
	return nil
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Get [%s]: %w"

	store, err := s.newStorage(append(pairs, ps.WithName(name))...)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return store, nil
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Create [%s]: %w"

	store, err := s.newStorage(append(pairs, ps.WithName(name))...)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	_, err = store.bucket.Put()
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return store, nil
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	store, err := s.newStorage(append(pairs, ps.WithName(name))...)
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, name, err)
	}
	_, err = store.bucket.Delete()
	if err != nil {
		err = handleQingStorError(err)
		return fmt.Errorf(errorMessage, s, name, err)
	}
	return nil
}

func (s *Service) newStorage(pairs ...*types.Pair) (*Storage, error) {
	const errorMessage = "qingstor new_storage: %w"

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	if !IsBucketNameValid(opt.Name) {
		err := handleQingStorError(ErrInvalidBucketName)
		return nil, fmt.Errorf(errorMessage, err)
	}

	if opt.HasLocation {
		bucket, err := s.service.Bucket(opt.Name, opt.Location)
		if err != nil {
			err = handleQingStorError(err)
			return nil, fmt.Errorf(errorMessage, err)
		}
		return &Storage{
			bucket:     bucket,
			config:     bucket.Config,
			properties: bucket.Properties,
			segments:   make(map[string]*segment.Segment),
		}, nil
	}

	url := fmt.Sprintf("%s://%s.%s:%d", s.config.Protocol, opt.Name, s.config.Host, s.config.Port)

	r, err := s.noRedirectClient.Head(url)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, err)
	}
	if r.StatusCode != http.StatusTemporaryRedirect {
		err = fmt.Errorf("head status is %d instead of %d", r.StatusCode, http.StatusTemporaryRedirect)
		return nil, fmt.Errorf(errorMessage, handleQingStorError(err))
	}

	// Example URL: https://bucket.zone.qingstor.com
	location := strings.Split(r.Header.Get("Location"), ".")[1]
	bucket, err := s.service.Bucket(opt.Name, location)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, err)
	}
	return &Storage{
		bucket:     bucket,
		config:     bucket.Config,
		properties: bucket.Properties,
		segments:   make(map[string]*segment.Segment),
	}, nil
}
