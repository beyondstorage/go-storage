package s3

import (
	"fmt"

	ps "github.com/Xuanwo/storage/types/pairs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
)

// Service is the s3 service config.
type Service struct {
	service s3iface.S3API
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer s3")
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) (err error) {
	const errorMessage = "%s List: %w"

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, err)
	}

	input := &s3.ListBucketsInput{}

	output, err := s.service.ListBuckets(input)
	if err != nil {
		err = handleS3Error(err)
		return fmt.Errorf(errorMessage, s, err)
	}

	for _, v := range output.Buckets {
		store, err := s.newStorage(ps.WithName(*v.Name))
		if err != nil {
			return fmt.Errorf(errorMessage, s, err)
		}
		opt.StoragerFunc(store)
	}
	return nil
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s get [%s]: %w"

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return store, nil
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Create [%s]: %w"

	opt, err := parseServicePairCreate(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(opt.Location),
		},
	}

	_, err = s.service.CreateBucket(input)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return store, nil
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	_, err = parseServicePairDelete(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}

	input := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err = s.service.DeleteBucket(input)
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}
	return nil
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...*types.Pair) (*Storage, error) {
	const errorMessage = "s3 new_storage: %w"

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	c := &Storage{
		service: s.service,
		name:    opt.Name,
		workDir: opt.WorkDir,
	}
	return c, nil
}
