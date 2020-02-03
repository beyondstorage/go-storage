package azblob

import (
	"fmt"

	"github.com/Azure/azure-storage-blob-go/azblob"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// Service is the azblob config.
type Service struct {
	service azblob.ServiceURL
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
		output, err = s.service.ListContainersSegment(opt.Context,
			marker, azblob.ListContainersSegmentOptions{})
		if err != nil {
			return fmt.Errorf(errorMessage, s, err)
		}

		for _, v := range output.ContainerItems {
			store, err := s.newStorage(ps.WithName(v.Name))
			if err != nil {
				return fmt.Errorf(errorMessage, s, err)
			}
			opt.StoragerFunc(store)
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
	const errorMessage = "%s Get [%s]: %w"

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	return store, nil
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (storage.Storager, error) {
	const errorMessage = "%s Create [%s]: %w"

	opt, err := parseServicePairCreate(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}

	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	_, err = store.bucket.Create(opt.Context, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, name, err)
	}
	return store, nil
}

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	opt, err := parseServicePairDelete(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}

	bucket := s.service.NewContainerURL(name)
	_, err = bucket.Delete(opt.Context, azblob.ContainerAccessConditions{})
	if err != nil {
		return fmt.Errorf(errorMessage, s, name, err)
	}
	return nil
}
