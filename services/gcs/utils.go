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
			err = &services.PairError{Op: "new gcs", Err: err, Pairs: pairs}
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
		return nil, nil, services.ErrCredentialProtocolNotSupported
	}

	client, err := gs.NewClient(opt.Context, options...)

	if err != nil {
		return nil, nil, err
	}

	srv.service = client
	srv.projectID = opt.Project

	store, err := srv.newStorage(pairs...)
	if err != nil && errors.Is(err, services.ErrPairRequired) {
		return srv, nil, nil
	}
	if err != nil {
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
		return "", &services.PairError{
			Op:    "parse storage class",
			Err:   services.ErrStorageClassNotSupported,
			Pairs: []*types.Pair{{Key: ps.StorageClass, Value: in}},
		}
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in string) (storageclass.Type, error) {
	switch in {
	case storageClassStandard:
		return storageclass.Hot, nil
	case storageClassNearLine:
		return storageclass.Warm, nil
	case storageClassColdLine:
		return storageclass.Cold, nil
	default:
		return "", &services.PairError{
			Op:    "format storage class",
			Err:   services.ErrStorageClassNotSupported,
			Pairs: []*types.Pair{{Key: ps.StorageClass, Value: in}},
		}
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
