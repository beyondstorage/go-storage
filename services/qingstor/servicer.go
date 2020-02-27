package qingstor

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
	"github.com/yunify/qingstor-sdk-go/v3/config"
	iface "github.com/yunify/qingstor-sdk-go/v3/interface"
	qserror "github.com/yunify/qingstor-sdk-go/v3/request/errors"
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
	return fmt.Sprintf("Servicer qingstor {AccessKey: %s}", s.config.AccessKeyID)
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list", err, "")
	}()

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return
	}

	input := &service.ListBucketsInput{}
	if opt.HasLocation {
		input.Location = &opt.Location
	}

	// FIXME: List buckets could be incomplete.
	output, err := s.service.ListBuckets(input)
	if err != nil {
		return
	}

	for _, v := range output.Buckets {
		store, err := s.newStorage(ps.WithName(*v.Name), ps.WithLocation(*v.Location))
		if err != nil {
			return err
		}
		opt.StoragerFunc(store)
	}
	return nil
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (store storage.Storager, err error) {
	defer func() {
		err = s.formatError("get", err, name)
	}()

	opt, err := parseServicePairGet(pairs...)
	if err != nil {
		return
	}

	location := opt.Location
	if !opt.HasLocation {
		location, err = s.detectLocation(name)
		if err != nil {
			return
		}
	}
	pairs = append(pairs, ps.WithName(name), ps.WithLocation(location))

	store, err = s.newStorage(append(pairs, ps.WithName(name))...)
	if err != nil {
		return
	}
	return
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (store storage.Storager, err error) {
	defer func() {
		err = s.formatError("create", err, name)
	}()

	_, err = parseServicePairCreate(pairs...)
	if err != nil {
		return
	}

	// ServicePairCreate requires location, so we don't need to add location into pairs
	pairs = append(pairs, ps.WithName(name))

	st, err := s.newStorage(pairs...)
	if err != nil {
		return
	}

	_, err = st.bucket.Put()
	if err != nil {
		return
	}
	return st, nil
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("create", err, name)
	}()

	opt, err := parseServicePairDelete(pairs...)
	if err != nil {
		return
	}

	location := opt.Location
	if !opt.HasLocation {
		location, err = s.detectLocation(name)
		if err != nil {
			return
		}
	}
	pairs = append(pairs, ps.WithName(name), ps.WithLocation(location))

	store, err := s.newStorage(pairs...)
	if err != nil {
		return
	}
	_, err = store.bucket.Delete()
	if err != nil {
		return
	}
	return nil
}

func (s *Service) newStorage(pairs ...*types.Pair) (store *Storage, err error) {
	defer func() {
		err = s.formatError("new storage", err, "")
	}()

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return
	}

	if !IsBucketNameValid(opt.Name) {
		err = ErrInvalidBucketName
		return
	}

	bucket, err := s.service.Bucket(opt.Name, opt.Location)
	if err != nil {
		return
	}
	return &Storage{
		bucket:     bucket,
		config:     bucket.Config,
		properties: bucket.Properties,
		segments:   make(map[string]*segment.Segment),

		workDir: opt.WorkDir,
	}, nil
}

func (s *Service) detectLocation(name string) (location string, err error) {
	defer func() {
		err = s.formatError("detect location", err, "")
	}()

	url := fmt.Sprintf("%s://%s.%s:%d", s.config.Protocol, name, s.config.Host, s.config.Port)

	r, err := s.noRedirectClient.Head(url)
	if err != nil {
		return
	}
	if r.StatusCode != http.StatusTemporaryRedirect {
		err = fmt.Errorf("head status is %d instead of %d", r.StatusCode, http.StatusTemporaryRedirect)
		return
	}

	// Example URL: https://bucket.zone.qingstor.com
	location = strings.Split(r.Header.Get("Location"), ".")[1]
	return
}

func (s *Service) formatError(op string, err error, name string) error {
	if err == nil {
		return nil
	}

	// Handle errors returned by qingstor.
	var e qserror.QingStorError
	if errors.As(err, &e) {
		err = formatQingStorError(&e)
	}

	return &services.ServiceError{
		Op:       op,
		Err:      err,
		Servicer: s,
		Name:     name,
	}
}
