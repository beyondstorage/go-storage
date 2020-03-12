package s3

import (
	"errors"
	"fmt"

	"github.com/Xuanwo/storage/types/metadata"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// New will create a new s3 service.
func New(pairs ...*types.Pair) (_ storage.Servicer, _ storage.Storager, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Type: Type, Err: err, Pairs: pairs}
		}
	}()

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, nil, err
	}

	cfg := aws.NewConfig()

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	switch credProtocol {
	case credential.ProtocolHmac:
		cfg = cfg.WithCredentials(credentials.NewStaticCredentials(cred[0], cred[1], ""))
	case credential.ProtocolEnv:
		cfg = cfg.WithCredentials(credentials.NewEnvCredentials())
	default:
		return nil, nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, nil, err
	}

	srv := &Service{service: s3.New(sess)}

	store, err := srv.newStorage(pairs...)
	if err != nil {
		if e := services.NewPairRequiredError(); errors.As(err, &e) {
			if len(e.Keys) == 1 && e.Keys[0] == ps.Name {
				return srv, nil, nil
			}
		}
		return nil, nil, err
	}
	return srv, store, nil
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
		return "", services.NewPairUnsupportedError(ps.WithStorageClass(in))
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
		return "", services.NewMetadataNotRecognizedError(metadata.ObjectMetaStorageClass, in)
	}

}

func formatError(err error) error {
	e, ok := err.(awserr.Error)
	if !ok {
		return err
	}

	switch e.Code() {
	case "NoSuchKey":
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case "AccessDenied":
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	}

	return err
}
