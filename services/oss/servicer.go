package oss

import (
	"fmt"

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
	const errorMessage = "%s List: %w"

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, err)
	}

	marker := ""
	var output oss.ListBucketsResult
	for {
		output, err = s.service.ListBuckets(
			oss.Marker(marker),
			oss.MaxKeys(1000),
		)
		if err != nil {
			return fmt.Errorf(errorMessage, s, err)
		}

		for _, v := range output.Buckets {
			store, err := s.newStorage(ps.WithName(v.Name))
			if err != nil {
				return fmt.Errorf(errorMessage, s, err)
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

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	err = s.service.CreateBucket(name)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return store, nil
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	err = s.service.DeleteBucket(name)
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}
	return nil
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...*types.Pair) (*Storage, error) {
	const errorMessage = "oss new_storage: %w"

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	bucket, err := s.service.Bucket(opt.Name)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	store := &Storage{
		bucket: bucket,

		workDir: opt.WorkDir,
	}
	return store, nil
}
