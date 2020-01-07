package azblob

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/types"
)

// Service is the azblob config.
type Service struct {
	service azblob.ServiceURL
}

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
func New(pairs ...*types.Pair) (s *Service, err error) {
	const errorMessage = "%s New: %w"

	s = &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	primaryURL, _ := url.Parse(opt.Endpoint.Value().String())

	credProtocol, credValue := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, fmt.Errorf(errorMessage, s, credential.ErrUnsupportedProtocol)
	}

	cred, err := azblob.NewSharedKeyCredential(credValue[0], credValue[1])
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	p := azblob.NewPipeline(cred, azblob.PipelineOptions{})
	s.service = azblob.NewServiceURL(*primaryURL, p)
	return
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer azblob")
}

// List implements Servicer.List
func (s *Service) List(pairs ...*types.Pair) (err error) {
	const errorMessage = "%s List: %w"

	opt, err := parseServicePairList(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, err)
	}

	marker := azblob.Marker{}
	var output *azblob.ListContainersSegmentResponse
	for {
		output, err = s.service.ListContainersSegment(context.TODO(),
			marker, azblob.ListContainersSegmentOptions{})
		if err != nil {
			return fmt.Errorf(errorMessage, s, err)
		}

		for _, v := range output.ContainerItems {
			bucket := s.service.NewContainerURL(v.Name)
			opt.StoragerFunc(newStorage(bucket, v.Name))
		}

		marker = output.NextMarker
		if !marker.NotDone() {
			break
		}
	}
	return
}

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const _ = "%s Get [%s]: %w"

	bucket := s.service.NewContainerURL(name)
	return newStorage(bucket, name), nil
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Create [%s]: %w"

	bucket := s.service.NewContainerURL(name)
	_, err := bucket.Create(context.TODO(), azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return newStorage(bucket, name), nil
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	bucket := s.service.NewContainerURL(name)
	_, err = bucket.Delete(context.TODO(), azblob.ContainerAccessConditions{})
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}
	return nil
}
