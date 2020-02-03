package gcs

import (
	"errors"
	"fmt"

	gs "cloud.google.com/go/storage"
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/types"
	"google.golang.org/api/option"
)

// New will create a new aliyun oss service.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	const errorMessage = "gcs New: %w"

	srv := &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	options := make([]option.ClientOption, 0)

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	switch credProtocol {
	case credential.ProtocolAPIKey:
		options = append(options, option.WithAPIKey(cred[0]))
	case credential.ProtocolFile:
		options = append(options, option.WithCredentialsFile(cred[0]))
	default:
		return nil, nil, fmt.Errorf(errorMessage, credential.ErrUnsupportedProtocol)
	}

	client, err := gs.NewClient(opt.Context, options...)

	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	srv.service = client
	srv.projectID = opt.Project

	store, err := srv.newStorage(pairs...)
	if err != nil && errors.Is(err, types.ErrPairRequired) {
		return srv, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
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
		return "", types.ErrStorageClassNotSupported
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
		return "", types.ErrStorageClassNotSupported
	}
}
