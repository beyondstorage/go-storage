package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/httpclient"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// New will create both Servicer and Storager.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	return newServicerAndStorager(pairs...)
}

// NewServicer will create Servicer only.
func NewServicer(pairs ...*types.Pair) (storage.Servicer, error) {
	return newServicer(pairs...)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...*types.Pair) (storage.Storager, error) {
	_, store, err := newServicerAndStorager(pairs...)
	return store, err
}

func newServicer(pairs ...*types.Pair) (srv *Service, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Op: services.OpNewServicer, Type: Type, Err: err, Pairs: pairs}
		}
	}()

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, err
	}

	cfg := aws.NewConfig()

	// Set s3 config's http client
	cfg.HTTPClient = httpclient.New(opt.HTTPClientOptions)

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	switch credProtocol {
	case credential.ProtocolHmac:
		cfg = cfg.WithCredentials(credentials.NewStaticCredentials(cred[0], cred[1], ""))
	case credential.ProtocolEnv:
		cfg = cfg.WithCredentials(credentials.NewEnvCredentials())
	default:
		return nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	srv = &Service{
		sess:    sess,
		service: newS3Service(sess),
	}
	return
}

// New will create a new s3 service.
func newServicerAndStorager(pairs ...*types.Pair) (srv *Service, store *Storage, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Op: services.OpNewStorager, Type: Type, Err: err, Pairs: pairs}
		}
	}()

	srv, err = newServicer(pairs...)
	if err != nil {
		return
	}

	store, err = srv.newStorage(pairs...)
	if err != nil {
		return
	}
	return
}

// All available storage classes are listed here.
const (
	StorageClassStandard           = s3.ObjectStorageClassStandard
	StorageClassReducedRedundancy  = s3.ObjectStorageClassReducedRedundancy
	StorageClassGlacier            = s3.ObjectStorageClassGlacier
	StorageClassStandardIa         = s3.ObjectStorageClassStandardIa
	StorageClassOnezoneIa          = s3.ObjectStorageClassOnezoneIa
	StorageClassIntelligentTiering = s3.ObjectStorageClassIntelligentTiering
	StorageClassDeepArchive        = s3.ObjectStorageClassDeepArchive
)

func formatError(err error) error {
	e, ok := err.(awserr.RequestFailure)
	if !ok {
		return err
	}

	switch e.Code() {
	// AWS SDK will use status code to generate awserr.Error, so "NotFound" should also be supported.
	case "NoSuchKey", "NotFound":
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case "AccessDenied":
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	}

	return err
}

func newS3Service(sess *session.Session, cfgs ...*aws.Config) (srv *s3.S3) {
	srv = s3.New(sess, cfgs...)

	// S3 will calculate payload's content-sha256 by default, we change this behavior for following reasons:
	// - To support uploading content without seek support: stdin, bytes.Reader
	// - To allow user decide when to calculate the hash, especially for big files
	srv.Handlers.Sign.SwapNamed(v4.BuildNamedHandler(v4.SignRequestHandler.Name, func(s *v4.Signer) {
		s.DisableURIPathEscaping = true
		// With UnsignedPayload set to true, signer will set "X-Amz-Content-Sha256" to "UNSIGNED-PAYLOAD"
		s.UnsignedPayload = true
	}))
	return
}
