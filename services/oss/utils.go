package oss

import (
	"errors"
	"fmt"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// New will create a new aliyun oss service.
func New(pairs ...*types.Pair) (_ storage.Servicer, _ storage.Storager, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Type: Type, Err: err, Pairs: pairs}
		}
	}()

	srv := &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, nil, err
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}
	ep := opt.Endpoint.Value()

	srv.service, err = oss.New(ep.String(), cred[0], cred[1])
	if err != nil {
		return nil, nil, err
	}

	store, err := srv.newStorage(pairs...)
	if err != nil {
		if e := services.NewPairRequiredError(); errors.As(err, &e) {
			return srv, nil, nil
		}
		return nil, nil, err
	}
	return srv, store, nil
}

const (
	// ref: https://www.alibabacloud.com/help/doc-detail/31984.htm
	storageClassHeader = "x-oss-storage-class"

	// ref: https://www.alibabacloud.com/help/doc-detail/51374.htm
	storageClassStandard = "STANDARD"
	storageClassIA       = "IA"
	storageClassArchive  = "Archive"
)

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (string, error) {
	switch in {
	case storageclass.Hot:
		return storageClassStandard, nil
	case storageclass.Warm:
		return storageClassIA, nil
	case storageclass.Cold:
		return storageClassArchive, nil
	default:
		return "", services.NewPairUnsupportedError(ps.WithStorageClass(in))
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in string) storageclass.Type {
	switch in {
	case storageClassStandard:
		return storageclass.Hot
	case storageClassIA:
		return storageclass.Warm
	case storageClassArchive:
		return storageclass.Cold
	default:
		return ""
	}
}

func formatError(err error) error {
	switch e := err.(type) {
	case oss.ServiceError:
		switch e.Code {
		case "":
			switch e.StatusCode {
			case 404:
				return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
			default:
				return err
			}
		case "NoSuchKey":
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		case "AccessDenied":
			return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
		}
	case oss.UnexpectedStatusCodeError:
		switch e.Got() {
		case 404:
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		case 403:
			return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
		}
	}

	return err
}
