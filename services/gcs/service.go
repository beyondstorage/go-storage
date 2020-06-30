package gcs

import (
	"context"

	"google.golang.org/api/iterator"

	"github.com/cns-io/go-storage/v2"
	ps "github.com/cns-io/go-storage/v2/types/pairs"
)

func (s *Service) create(ctx context.Context, name string, opt *pairServiceCreate) (store storage.Storager, err error) {
	st, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return nil, err
	}
	err = st.bucket.Create(ctx, s.projectID, nil)
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) delete(ctx context.Context, name string, opt *pairServiceDelete) (err error) {
	store, err := s.newStorage(ps.WithName(name))
	if err != nil {
		return err
	}
	err = store.bucket.Delete(ctx)
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
	it := s.service.Buckets(ctx, s.projectID)

	for {
		bucketAttr, err := it.Next()
		// Next will return iterator.Done if there is no more items.
		if err != nil && err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		store, err := s.newStorage(ps.WithName(bucketAttr.Name))
		if err != nil {
			return err
		}
		opt.StoragerFunc(store)
	}
}
