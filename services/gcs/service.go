package gcs

import (
	"context"

	"google.golang.org/api/iterator"

	typ "go.beyondstorage.io/v5/types"
)

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store typ.Storager, err error) {
	f := s.f
	f.Name = name
	st, err := f.newStorage()
	if err != nil {
		return nil, err
	}
	err = st.bucket.Create(ctx, s.projectID, nil)
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) delete(ctx context.Context, name string, opt pairServiceDelete) (err error) {
	f := s.f
	f.Name = name
	st, err := f.newStorage()
	if err != nil {
		return err
	}
	err = st.bucket.Delete(ctx)
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
	return typ.NewStoragerIterator(ctx, s.nextStoragePage, nil), nil
}

func (s *Service) nextStoragePage(ctx context.Context, page *typ.StoragerPage) error {
	it := s.service.Buckets(ctx, s.projectID)
	for {
		bucket, err := it.Next()
		if err == iterator.Done {
			return typ.IterateDone
		}
		if err != nil {
			return err
		}
		f := s.f
		f.Name = bucket.Name
		st, err := f.newStorage()
		if err != nil {
			return err
		}
		page.Data = append(page.Data, st)
	}
}
