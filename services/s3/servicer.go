package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// Service is the s3 service config.
type Service struct {
	sess    *session.Session
	service *s3.S3
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer s3")
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

	input := &s3.ListBucketsInput{}

	output, err := s.service.ListBuckets(input)
	if err != nil {
		return err
	}

	for _, v := range output.Buckets {
		store, err := s.newStorage(ps.WithName(*v.Name))
		if err != nil {
			return err
		}
		opt.StoragerFunc(store)
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

	opt, err := s.parsePairCreate(pairs...)
	if err != nil {
		return nil, err
	}

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}

	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(opt.Location),
		},
	}

	_, err = s.service.CreateBucket(input)
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

	_, err = s.parsePairDelete(pairs...)
	if err != nil {
		return err
	}

	input := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err = s.service.DeleteBucket(input)
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

	st = &Storage{
		service: s3.New(s.sess, aws.NewConfig().WithRegion(opt.Location)),

		name:    opt.Name,
		workDir: "/",
	}

	if opt.HasWorkDir {
		st.workDir = opt.WorkDir
	}
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
