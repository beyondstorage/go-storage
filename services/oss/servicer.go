package oss

import (
	"fmt"

	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
)

// Service is the aliyun oss *Service config.
type Service struct {
	service *oss.Client
}

// New will create a new aliyun oss service.
func New(pairs ...*types.Pair) (s *Service, err error) {
	const errorMessage = "%s New: %w"

	s = &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, fmt.Errorf(errorMessage, s, credential.ErrUnsupportedProtocol)
	}
	ep := opt.Endpoint.Value()

	s.service, err = oss.New(ep.String(), cred[0], cred[1])
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}
	return
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
			bucket, err := s.service.Bucket(v.Name)
			if err != nil {
				return fmt.Errorf(errorMessage, s, err)
			}
			if opt.HasStoragerFunc {
				c := newStorage(bucket)
				opt.StoragerFunc(c)
			}
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

	bucket, err := s.service.Bucket(name)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return newStorage(bucket), nil
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Create [%s]: %w"

	err := s.service.CreateBucket(name)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	bucket, err := s.service.Bucket(name)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return newStorage(bucket), nil
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
