package cos

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Xuanwo/storage/services"
	ps "github.com/Xuanwo/storage/types/pairs"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/types"
)

// New will create a new Tencent oss service.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	const errorMessage = "cos New: %w"

	srv := &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, nil, fmt.Errorf(errorMessage, credential.ErrUnsupportedProtocol)
	}

	srv.client = &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cred[0],
			SecretKey: cred[1],
		},
		Timeout: 100 * time.Second,
	}
	srv.service = cos.NewClient(nil, srv.client)

	store, err := srv.newStorage(pairs...)
	if err != nil && errors.Is(err, services.ErrPairRequired) {
		return srv, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}
	return srv, store, nil
}

const (
	// ref: https://cloud.tencent.com/document/product/436/7745
	storageClassHeader = "x-cos-storage-class"

	storageClassStandard   = "STANDARD"
	storageClassStandardIA = "STANDARD_IA"
	storageClassArchive    = "ARCHIVE"
)

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (string, error) {
	switch in {
	case storageclass.Cold:
		return storageClassArchive, nil
	case storageclass.Hot:
		return storageClassStandard, nil
	case storageclass.Warm:
		return storageClassStandardIA, nil
	default:
		return "", &services.PairError{
			Op:    "parse storage class",
			Err:   services.ErrStorageClassNotSupported,
			Key:   ps.StorageClass,
			Value: in,
		}
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in string) (storageclass.Type, error) {
	switch in {
	case storageClassArchive:
		return storageclass.Cold, nil
	case storageClassStandardIA:
		return storageclass.Warm, nil
	// cos only return storage class while not standard, we should handle empty string
	case storageClassStandard, "":
		return storageclass.Hot, nil
	default:
		return "", &services.PairError{
			Op:    "format storage class",
			Err:   services.ErrStorageClassNotSupported,
			Key:   ps.StorageClass,
			Value: in,
		}
	}
}

// ref: https://www.qcloud.com/document/product/436/7730
func formatCosError(err *cos.ErrorResponse) error {
	switch err.Code {
	case "NoSuchKey":
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case "AccessDenied":
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
	}
}
