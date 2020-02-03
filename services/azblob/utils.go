package azblob

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/types"
)

// New will create a new azblob oss service.
//
// azblob use different URL to represent different sub services.
// - ServiceURL's          methods perform operations on a storage account.
//   - ContainerURL's     methods perform operations on an account's container.
//      - BlockBlobURL's  methods perform operations on a container's block blob.
//      - AppendBlobURL's methods perform operations on a container's append blob.
//      - PageBlobURL's   methods perform operations on a container's page blob.
//      - BlobURL's       methods perform operations on a container's blob regardless of the blob's type.
//
// Our Service will store a ServiceURL for operation.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	const errorMessage = "azblob New: %w"

	srv := &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	primaryURL, _ := url.Parse(opt.Endpoint.Value().String())

	credProtocol, credValue := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, nil, fmt.Errorf(errorMessage, credential.ErrUnsupportedProtocol)
	}

	cred, err := azblob.NewSharedKeyCredential(credValue[0], credValue[1])
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	p := azblob.NewPipeline(cred, azblob.PipelineOptions{})
	srv.service = azblob.NewServiceURL(*primaryURL, p)

	store, err := srv.newStorage(pairs...)
	if err != nil && errors.Is(err, types.ErrPairRequired) {
		return srv, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}
	return srv, store, nil
}

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (azblob.AccessTierType, error) {
	switch in {
	case storageclass.Cold:
		return azblob.AccessTierArchive, nil
	case storageclass.Hot:
		return azblob.AccessTierHot, nil
	case storageclass.Warm:
		return azblob.AccessTierCool, nil
	default:
		return azblob.AccessTierType(""), types.ErrStorageClassNotSupported
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in azblob.AccessTierType) (storageclass.Type, error) {
	switch in {
	case azblob.AccessTierArchive:
		return storageclass.Cold, nil
	case azblob.AccessTierCool:
		return storageclass.Warm, nil
	case azblob.AccessTierHot:
		return storageclass.Hot, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}
