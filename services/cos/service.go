package cos

import (
	"context"

	"github.com/aos-dev/go-storage/v2"
	ps "github.com/aos-dev/go-storage/v2/types/pairs"
)

func (s *Service) create(ctx context.Context, name string, opt *pairServiceCreate) (store storage.Storager, err error) {
	st, err := s.newStorage(ps.WithName(name), ps.WithLocation(opt.Location))
	if err != nil {
		return nil, err
	}
	_, err = st.bucket.Put(ctx, nil)
	if err != nil {
		return nil, err
	}
	return st, nil
}
func (s *Service) delete(ctx context.Context, name string, opt *pairServiceDelete) (err error) {
	store, err := s.newStorage(ps.WithName(name), ps.WithLocation(opt.Location))
	if err != nil {
		return err
	}
	_, err = store.bucket.Delete(ctx)
	if err != nil {
		return err
	}
	return
}
func (s *Service) get(ctx context.Context, name string, opt *pairServiceGet) (store storage.Storager, err error) {
	st, err := s.newStorage(ps.WithName(name), ps.WithLocation(opt.Location))
	if err != nil {
		return nil, err
	}
	return st, nil
}
func (s *Service) list(ctx context.Context, opt *pairServiceList) (err error) {
	output, _, err := s.service.Service.Get(ctx)
	if err != nil {
		return err
	}
	for _, v := range output.Buckets {
		store, err := s.newStorage(ps.WithName(v.Name), ps.WithLocation(v.Region))
		if err != nil {
			return err
		}
		opt.StoragerFunc(store)
	}
	return
}
