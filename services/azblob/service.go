package azblob

import (
	"context"

	"github.com/Azure/azure-storage-blob-go/azblob"

	typ "go.beyondstorage.io/v5/types"
)

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store typ.Storager, err error) {
	f := s.f
	f.Name = name
	st, err := f.newStorage()
	if err != nil {
		return nil, err
	}
	_, err = st.bucket.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) delete(ctx context.Context, name string, opt pairServiceDelete) (err error) {
	bucket := s.service.NewContainerURL(name)
	_, err = bucket.Delete(ctx, azblob.ContainerAccessConditions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) get(ctx context.Context, name string, opt pairServiceGet) (store typ.Storager, err error) {
	f := s.f
	f.Name = name
	st, err := f.newStorage()
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (s *Service) list(ctx context.Context, opt pairServiceList) (it *typ.StoragerIterator, err error) {
	input := &storagePageStatus{
		maxResults: 200,
	}

	return typ.NewStoragerIterator(ctx, s.nextStoragePage, input), nil
}

func (s *Service) nextStoragePage(ctx context.Context, page *typ.StoragerPage) error {
	input := page.Status.(*storagePageStatus)

	output, err := s.service.ListContainersSegment(ctx, input.marker, azblob.ListContainersSegmentOptions{
		MaxResults: input.maxResults,
	})
	if err != nil {
		return err
	}

	for _, v := range output.ContainerItems {
		f := s.f
		f.Name = v.Name
		store, err := f.newStorage()
		if err != nil {
			return err
		}

		page.Data = append(page.Data, store)
	}

	if !output.NextMarker.NotDone() {
		return typ.IterateDone
	}

	input.marker = output.NextMarker
	return nil
}
