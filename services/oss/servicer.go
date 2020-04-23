package oss

import (
	"fmt"

	"github.com/Xuanwo/storage/services"
	ps "github.com/Xuanwo/storage/types/pairs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
)

// Service is the aliyun oss *Service config.
type Service struct {
	service *oss.Client
}

// String implements Servicer.String
func (s *Service) String() string {
	if s.service == nil {
		return fmt.Sprintf("Servicer oss")
	}
	return fmt.Sprintf("Servicer oss {AccessKey: %s}", s.service.Config.AccessKeyID)
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpList, err, "")
	}()

	opt, err := s.parsePairList(pairs...)
	if err != nil {
		return err
	}

	marker := ""
	var output oss.ListBucketsResult
	for {
		output, err = s.service.ListBuckets(
			oss.Marker(marker),
			oss.MaxKeys(1000),
		)
		if err != nil {
			return err
		}

		for _, v := range output.Buckets {
			store, err := s.newStorage(ps.WithName(v.Name))
			if err != nil {
				return err
			}
			opt.StoragerFunc(store)
		}

		marker = output.NextMarker
		if output.IsTruncated {
			break
		}
	}
	return nil
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (st storage.Storager, err error) {
	defer func() {
		err = s.formatError(services.OpGet, err, name)
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
		err = s.formatError(services.OpCreate, err, name)
	}()

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	err = s.service.CreateBucket(name)
	if err != nil {
		return nil, err
	}
	return store, nil
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpDelete, err, name)
	}()

	err = s.service.DeleteBucket(name)
	if err != nil {
		return err
	}
	return nil
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...*types.Pair) (st *Storage, err error) {
	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, err
	}

	bucket, err := s.service.Bucket(opt.Name)
	if err != nil {
		return nil, err
	}

	store := &Storage{
		bucket: bucket,

		workDir: "/",
	}

	if opt.HasWorkDir {
		store.workDir = opt.WorkDir
	}
	return store, nil
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
