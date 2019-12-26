package azblob

import (
	"fmt"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/types"
)

// Service is the azblob config.
type Service struct {
	service pipeline.Pipeline
}

// New will create a new azblob oss service.
func New(pairs ...*types.Pair) (s *Service, err error) {
	const errorMessage = "%s New: %w"

	s = &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	credProtocol, credValue := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, fmt.Errorf(errorMessage, s, credential.ErrUnsupportedProtocol)
	}

	cred, err := azblob.NewSharedKeyCredential(credValue[0], credValue[1])
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	p := azblob.NewPipeline(cred, azblob.PipelineOptions{})
	s.service = p
	return
}

// String implements Servicer.String
func (s Service) String() string {
	panic("implement me")
}

// List implements Servicer.List
func (s Service) List(pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// Get implements Servicer.Get
func (s Service) Get(name string, pairs ...*types.Pair) (storage.Storager, error) {
	panic("implement me")
}

// Create implements Servicer.Create
func (s Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	panic("implement me")
}

// Delete implements Servicer.Delete
func (s Service) Delete(name string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}
