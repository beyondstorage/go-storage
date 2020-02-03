package s3

import (
	"errors"
	"fmt"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// New will create a new s3 service.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	const errorMessage = "s3 New: %w"

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	cfg := aws.NewConfig()

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	switch credProtocol {
	case credential.ProtocolHmac:
		cfg = cfg.WithCredentials(credentials.NewStaticCredentials(cred[0], cred[1], ""))
	case credential.ProtocolEnv:
		cfg = cfg.WithCredentials(credentials.NewEnvCredentials())
	default:
		return nil, nil, fmt.Errorf(errorMessage, credential.ErrUnsupportedProtocol)
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	srv := &Service{service: s3.New(sess)}

	store, err := srv.newStorage(pairs...)
	if err != nil && errors.Is(err, types.ErrPairRequired) {
		return srv, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}
	return srv, store, nil
}

func handleS3Error(err error) error {
	if err == nil {
		panic("error must not be nil")
	}

	var e awserr.Error
	e, ok := err.(awserr.Error)
	if !ok {
		return fmt.Errorf("%w: %v", types.ErrUnhandledError, err)
	}

	switch e.Code() {
	default:
		return fmt.Errorf("%w: %v", types.ErrUnhandledError, err)
	}
}

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (string, error) {
	switch in {
	case storageclass.Hot:
		return s3.ObjectStorageClassStandard, nil
	case storageclass.Warm:
		return s3.ObjectStorageClassStandardIa, nil
	case storageclass.Cold:
		return s3.ObjectStorageClassGlacier, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in string) (storageclass.Type, error) {
	switch in {
	case s3.ObjectStorageClassStandard:
		return storageclass.Hot, nil
	case s3.ObjectStorageClassStandardIa:
		return storageclass.Warm, nil
	case s3.ObjectStorageClassGlacier:
		return storageclass.Cold, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}

}
