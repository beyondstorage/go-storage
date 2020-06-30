package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aos-dev/go-storage/v2"
	ps "github.com/aos-dev/go-storage/v2/types/pairs"
)

func (s *Service) create(ctx context.Context, name string, opt *pairServiceCreate) (store storage.Storager, err error) {
	pairs := append(opt.pairs, ps.WithName(name))

	st, err := s.newStorage(pairs...)
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
	return st, nil
}
func (s *Service) delete(ctx context.Context, name string, opt *pairServiceDelete) (err error) {
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err = s.service.DeleteBucket(input)
	if err != nil {
		return err
	}
	return
}
func (s *Service) get(ctx context.Context, name string, opt *pairServiceGet) (store storage.Storager, err error) {
	pairs := append(opt.pairs, ps.WithName(name))

	st, err := s.newStorage(pairs...)
	if err != nil {
		return nil, err
	}
	return st, nil
}
func (s *Service) list(ctx context.Context, opt *pairServiceList) (err error) {
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
