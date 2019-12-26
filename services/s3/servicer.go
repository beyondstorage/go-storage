package s3

import (
	"fmt"

	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
)

// Service is the s3 service config.
type Service struct {
	service s3iface.S3API
}

// New will create a new s3 service.
func New(pairs ...*types.Pair) (s *Service, err error) {
	const errorMessage = "%s New: %w"

	opt, err := parseServicePairInit(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	cfg := aws.NewConfig()

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	switch credProtocol {
	case credential.ProtocolHmac:
		cfg = cfg.WithCredentials(credentials.NewStaticCredentials(cred[0], cred[1], ""))
	case credential.ProtocolEnv:
		cfg = cfg.WithCredentials(credentials.NewEnvCredentials())
	default:
		return nil, fmt.Errorf(errorMessage, s, credential.ErrUnsupportedProtocol)
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	srv := s3.New(sess)

	s = &Service{service: srv}
	return s, nil
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer s3")
}

// List implements Servicer.List
func (s Service) List(pairs ...*types.Pair) (err error) {
	errorMessage := "list s3 storager: %w"

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	input := &s3.ListBucketsInput{}

	output, err := s.service.ListBuckets(input)
	if err != nil {
		err = handleS3Error(err)
		return fmt.Errorf(errorMessage, err)
	}

	for _, v := range output.Buckets {
		store, err := newStorage(s.service, *v.Name)
		if err != nil {
			return fmt.Errorf(errorMessage, err)
		}
		if opt.HasStoragerFunc {
			opt.StoragerFunc(store)
		}
	}
	return nil
}

// Get implements Servicer.Get
func (s Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	errorMessage := "get s3 storager: %w"

	store, err := newStorage(s.service, name)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}
	return store, nil
}

// Create implements Servicer.Create
func (s Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	errorMessage := "create s3 storager: %w"

	opt, err := parseServicePairCreate(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(opt.Location),
		},
	}

	_, err = s.service.CreateBucket(input)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	store, err := newStorage(s.service, name)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}
	return store, nil
}

// Delete implements Servicer.Delete
func (s Service) Delete(name string, pairs ...*types.Pair) (err error) {
	errorMessage := "delete s3 storager: %w"

	_, err = parseServicePairDelete(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	input := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err = s.service.DeleteBucket(input)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}
