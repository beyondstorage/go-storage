package qingstor

import (
	"context"

	"github.com/qingstor/qingstor-sdk-go/v4/service"

	"github.com/aos-dev/go-storage/v2"
	ps "github.com/aos-dev/go-storage/v2/types/pairs"
)

func (s *Service) create(ctx context.Context, name string, opt *pairServiceCreate) (store storage.Storager, err error) {
	// ServicePairCreate requires location, so we don't need to add location into pairs
	pairs := append(opt.pairs, ps.WithName(name))

	st, err := s.newStorage(pairs...)
	if err != nil {
		return
	}

	_, err = st.bucket.Put()
	if err != nil {
		return
	}
	return st, nil
}
func (s *Service) delete(ctx context.Context, name string, opt *pairServiceDelete) (err error) {
	pairs := append(opt.pairs, ps.WithName(name))

	store, err := s.newStorage(pairs...)
	if err != nil {
		return
	}
	_, err = store.bucket.Delete()
	if err != nil {
		return
	}
	return nil
}
func (s *Service) get(ctx context.Context, name string, opt *pairServiceGet) (store storage.Storager, err error) {
	pairs := append(opt.pairs, ps.WithName(name))

	store, err = s.newStorage(pairs...)
	if err != nil {
		return
	}
	return
}
func (s *Service) list(ctx context.Context, opt *pairServiceList) (err error) {
	input := &service.ListBucketsInput{}
	if opt.HasLocation {
		input.Location = &opt.Location
	}

	offset := 0
	var output *service.ListBucketsOutput
	for {
		input.Offset = service.Int(offset)

		output, err = s.service.ListBuckets(input)
		if err != nil {
			return
		}

		for _, v := range output.Buckets {
			store, err := s.newStorage(ps.WithName(*v.Name), ps.WithLocation(*v.Location))
			if err != nil {
				return err
			}
			opt.StoragerFunc(store)
		}

		offset += len(output.Buckets)
		if offset >= service.IntValue(output.Count) {
			break
		}
	}
	return nil
}
