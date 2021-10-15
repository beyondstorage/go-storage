package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	ps "go.beyondstorage.io/v5/pairs"
	. "go.beyondstorage.io/v5/types"
)

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store Storager, err error) {
	pairs := append(opt.pairs, ps.WithName(name))
	st, err := s.newStorage(pairs...)
	if err != nil {
		return nil, err
	}
	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &s3types.CreateBucketConfiguration{
			LocationConstraint: s3types.BucketLocationConstraint(opt.Location),
		},
	}
	_, err = s.service.CreateBucket(ctx, input)
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) delete(ctx context.Context, name string, opt pairServiceDelete) (err error) {
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}
	if opt.HasExceptedBucketOwner {
		input.ExpectedBucketOwner = &opt.ExceptedBucketOwner
	}
	_, err = s.service.DeleteBucket(ctx, input)
	if err != nil {
		return err
	}
	return
}

func (s *Service) get(ctx context.Context, name string, opt pairServiceGet) (store Storager, err error) {
	pairs := append(opt.pairs, ps.WithName(name))
	st, err := s.newStorage(pairs...)
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) list(ctx context.Context, opt pairServiceList) (it *StoragerIterator, err error) {
	input := &storagePageStatus{}
	return NewStoragerIterator(ctx, s.nextStoragePage, input), nil
}

func (s *Service) nextStoragePage(ctx context.Context, page *StoragerPage) error {
	output, err := s.service.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}
	for _, v := range output.Buckets {
		store, err := s.newStorage(ps.WithName(*v.Name))
		if err != nil {
			return err
		}
		page.Data = append(page.Data, store)
	}
	return IterateDone
}
