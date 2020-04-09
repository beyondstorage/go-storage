package gcs

import (
	"errors"
	"fmt"
	"net/http"

	gs "cloud.google.com/go/storage"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
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

	options := make([]option.ClientOption, 0)

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	switch credProtocol {
	case credential.ProtocolAPIKey:
		options = append(options, option.WithAPIKey(cred[0]))
	case credential.ProtocolFile:
		options = append(options, option.WithCredentialsFile(cred[0]))
	default:
		return nil, nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	client, err := gs.NewClient(opt.Context, options...)

	if err != nil {
		return nil, nil, err
	}

	srv.service = client
	srv.projectID = opt.Project

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
	storageClassStandard = "STANDARD"
	storageClassNearLine = "NEARLINE"
	storageClassColdLine = "COLDLINE"
)

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (string, error) {
	switch in {
	case storageclass.Hot:
		return storageClassStandard, nil
	case storageclass.Warm:
		return storageClassNearLine, nil
	case storageclass.Cold:
		return storageClassColdLine, nil
	default:
		return "", services.NewPairUnsupportedError(ps.WithStorageClass(in))
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in string) storageclass.Type {
	switch in {
	case storageClassStandard:
		return storageclass.Hot
	case storageClassNearLine:
		return storageclass.Warm
	case storageClassColdLine:
		return storageclass.Cold
	default:
		return ""
	}
}

// ref: https://cloud.google.com/storage/docs/json_api/v1/status-codes
func formatError(err error) error {
	// gcs sdk could return explicit error, we should handle them.
	if errors.Is(err, gs.ErrObjectNotExist) {
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	}

	e, ok := err.(*googleapi.Error)
	if !ok {
		return err
	}

	switch e.Code {
	case http.StatusNotFound:
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case http.StatusForbidden:
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
	}
}
