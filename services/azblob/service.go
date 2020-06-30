package azblob

import (
	"context"

	"github.com/Azure/azure-storage-blob-go/azblob"

	"github.com/aos-dev/go-storage/v2"
	ps "github.com/aos-dev/go-storage/v2/types/pairs"
)

func (s *Service) create(ctx context.Context, name string, opt *pairServiceCreate) (store storage.Storager, err error) {
	st, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	_, err = st.bucket.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		return nil, err
	}
	return st, nil
}
func (s *Service) delete(ctx context.Context, name string, opt *pairServiceDelete) (err error) {
	bucket := s.service.NewContainerURL(name)
	_, err = bucket.Delete(ctx, azblob.ContainerAccessConditions{})
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) get(ctx context.Context, name string, opt *pairServiceGet) (store storage.Storager, err error) {
	st, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}

	return st, nil
}
func (s *Service) list(ctx context.Context, opt *pairServiceList) (err error) {
	marker := azblob.Marker{}
	var output *azblob.ListContainersSegmentResponse
	for {
		output, err = s.service.ListContainersSegment(ctx,
			marker, azblob.ListContainersSegmentOptions{})
		if err != nil {
			return err
		}

		for _, v := range output.ContainerItems {
			store, err := s.newStorage(ps.WithName(v.Name))
			if err != nil {
				return err
			}
			opt.StoragerFunc(store)
		}

		marker = output.NextMarker
		if !marker.NotDone() {
			break
		}
	}
	return nil
}
