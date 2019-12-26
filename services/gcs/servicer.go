package gcs

import (
	"context"
	"fmt"

	gs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/types"
)

// Service is the gcs config.
type Service struct {
	service   *gs.Client
	projectID string
}

// New will create a new aliyun oss service.
func New(pairs ...*types.Pair) (s *Service, err error) {
	const errorMessage = "%s New: %w"

	s = &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	ctx := context.Background()

	options := make([]option.ClientOption, 0)

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	switch credProtocol {
	case credential.ProtocolAPIKey:
		options = append(options, option.WithAPIKey(cred[0]))
	case credential.ProtocolFile:
		options = append(options, option.WithCredentialsFile(cred[0]))
	default:
		return nil, fmt.Errorf(errorMessage, s, credential.ErrUnsupportedProtocol)
	}

	client, err := gs.NewClient(ctx, options...)

	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	s.service = client
	s.projectID = opt.Project
	return
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer gcs")
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) (err error) {
	const errorMessage = "%s List: %w"

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, err)
	}

	it := s.service.Buckets(context.TODO(), s.projectID)
	for {
		bucketAttr, err := it.Next()
		// Next will return iterator.Done if there is no more items.
		if err != nil && err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf(errorMessage, s, err)
		}
		bucket := s.service.Bucket(bucketAttr.Name)
		c := newStorage(bucket, bucketAttr.Name)
		opt.StoragerFunc(c)
	}
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const _ = "%s Get [%s]: %w"

	bucket := s.service.Bucket(name)
	c := newStorage(bucket, name)
	return c, nil
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Create [%s]: %w"

	bucket := s.service.Bucket(name)

	err := bucket.Create(context.TODO(), s.projectID, nil)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	c := newStorage(bucket, name)
	return c, nil
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	bucket := s.service.Bucket(name)

	err = bucket.Delete(context.TODO())
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}
	return nil
}
