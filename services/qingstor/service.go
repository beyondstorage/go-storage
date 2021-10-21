package qingstor

import (
	"context"

	"github.com/qingstor/qingstor-sdk-go/v4/service"

	ps "go.beyondstorage.io/v5/pairs"
	. "go.beyondstorage.io/v5/types"
)

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store Storager, err error) {
	// ServicePairCreate requires location, so we don't need to add location into pairs
	pairs := append(opt.pairs, ps.WithName(name))

	st, err := s.newStorage(pairs...)
	if err != nil {
		return
	}

	_, err = st.bucket.PutWithContext(ctx)
	if err != nil {
		return
	}
	return st, nil
}

func (s *Service) delete(ctx context.Context, name string, opt pairServiceDelete) (err error) {
	pairs := append(opt.pairs, ps.WithName(name))

	store, err := s.newStorage(pairs...)
	if err != nil {
		return
	}
	_, err = store.bucket.DeleteWithContext(ctx)
	if err != nil {
		return
	}
	return nil
}

func (s *Service) get(ctx context.Context, name string, opt pairServiceGet) (store Storager, err error) {
	pairs := append(opt.pairs, ps.WithName(name))

	store, err = s.newStorage(pairs...)
	if err != nil {
		return
	}
	return
}

func (s *Service) list(ctx context.Context, opt pairServiceList) (it *StoragerIterator, err error) {
	input := &storagePageStatus{}

	if opt.HasLocation {
		input.location = opt.Location
	}

	return NewStoragerIterator(ctx, s.nextStoragePage, input), nil
}

func (s *Service) nextStoragePage(ctx context.Context, page *StoragerPage) error {
	input := page.Status.(*storagePageStatus)

	serviceInput := &service.ListBucketsInput{
		Limit:  &input.offset,
		Offset: &input.limit,
	}
	if input.location != "" {
		serviceInput.Location = &input.location
	}

	output, err := s.service.ListBucketsWithContext(ctx, serviceInput)
	if err != nil {
		return err
	}

	for _, v := range output.Buckets {
		store, err := s.newStorage(ps.WithName(*v.Name), ps.WithLocation(*v.Location))
		if err != nil {
			return err
		}
		page.Data = append(page.Data, store)
	}

	input.offset += len(output.Buckets)
	if input.offset >= service.IntValue(output.Count) {
		return IterateDone
	}

	return nil
}
