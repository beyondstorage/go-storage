package oss

import (
	"errors"
	"fmt"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/types"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// New will create a new aliyun oss service.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	const errorMessage = "oss New: %w"

	srv := &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, nil, fmt.Errorf(errorMessage, credential.ErrUnsupportedProtocol)
	}
	ep := opt.Endpoint.Value()

	srv.service, err = oss.New(ep.String(), cred[0], cred[1])
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

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
		return "", types.ErrStorageClassNotSupported
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in string) (storageclass.Type, error) {
	switch in {
	case storageClassStandard:
		return storageclass.Hot, nil
	case storageClassIA:
		return storageclass.Warm, nil
	case storageClassArchive:
		return storageclass.Cold, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}
