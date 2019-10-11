package qingstor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yunify/qingstor-sdk-go/v3/config"
	iface "github.com/yunify/qingstor-sdk-go/v3/interface"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
)

// Service is the qingstor service config.
type Service struct {
	config  *config.Config
	service iface.Service

	noRedirectClient *http.Client
}

// New will create a new qingstor service.
func New() *Service {
	return &Service{
		noRedirectClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// Init implements Servicer.Init
func (s *Service) Init(pairs ...*types.Pair) (err error) {
	errorMessage := "init qingstor service failed: %w"

	opt, err := parseServicePairInit(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	cfg, err := config.New(opt.AccessKey, opt.SecretKey)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	if opt.HasHost {
		cfg.Host = opt.Host
	}
	if opt.HasPort {
		cfg.Port = opt.Port
	}
	if opt.HasProtocol {
		cfg.Protocol = opt.Protocol
	}

	s.config = cfg
	s.service, _ = service.Init(cfg)
	return
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	errorMessage := "get qingstor storager [%s] failed: %w"

	opt, err := parseServicePairGet(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, name, err)
	}

	// TODO: add bucket name check here.

	if opt.HasLocation {
		bucket, err := s.service.Bucket(name, opt.Location)
		if err != nil {
			err = handleQingStorError(err)
			return nil, fmt.Errorf(errorMessage, name, err)
		}
		return &Client{
			bucket:   bucket,
			segments: make(map[string]*segment.Segment),
		}, nil
	}

	url := fmt.Sprintf("%s://%s.%s:%d", s.config.Protocol, name, s.config.Host, s.config.Port)

	r, err := s.noRedirectClient.Head(url)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, name, err)
	}
	if r.StatusCode != http.StatusTemporaryRedirect {
		err = fmt.Errorf("head status is %d instead of %d", r.StatusCode, http.StatusTemporaryRedirect)
		return nil, fmt.Errorf(errorMessage, name, handleQingStorError(err))
	}

	// Example URL: https://bucket.zone.qingstor.com
	location := strings.Split(r.Header.Get("Location"), ".")[1]
	bucket, err := s.service.Bucket(name, location)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, name, err)
	}
	return &Client{
		bucket:   bucket,
		segments: make(map[string]*segment.Segment),
	}, nil
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	errorMessage := "create qingstor storager [%s] failed: %w"

	opt, err := parseServicePairCreate(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, name, err)
	}

	// TODO: check bucket name here.

	bucket, err := s.service.Bucket(name, opt.Location)
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, err)
	}

	_, err = bucket.Put()
	if err != nil {
		err = handleQingStorError(err)
		return nil, fmt.Errorf(errorMessage, err)
	}
	return &Client{
		bucket:   bucket,
		segments: make(map[string]*segment.Segment),
	}, nil
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) ([]storage.Storager, error) {
	panic("not supported")
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}
