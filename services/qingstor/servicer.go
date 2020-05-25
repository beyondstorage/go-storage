package qingstor

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/qingstor/qingstor-sdk-go/v4/config"
	iface "github.com/qingstor/qingstor-sdk-go/v4/interface"
	"github.com/qingstor/qingstor-sdk-go/v4/service"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/headers"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// Get implements Servicer.Get
func (s *Service) Get(name string, pairs ...*types.Pair) (store storage.Storager, err error) {
	defer func() {
		err = s.formatError(services.OpGet, err, name)
	}()

	_, err = s.parsePairGet(pairs...)
	if err != nil {
		return
	}

	pairs = append(pairs, ps.WithName(name))

	store, err = s.newStorage(pairs...)
	if err != nil {
		return
	}
	return
}

// Create implements Servicer.Create
func (s *Service) Create(name string, pairs ...*types.Pair) (store storage.Storager, err error) {
	defer func() {
		err = s.formatError(services.OpCreate, err, name)
	}()

	_, err = s.parsePairCreate(pairs...)
	if err != nil {
		return
	}

	// ServicePairCreate requires location, so we don't need to add location into pairs
	pairs = append(pairs, ps.WithName(name))

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

// Delete implements Servicer.Delete
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpDelete, err, name)
	}()

	_, err = s.parsePairDelete(pairs...)
	if err != nil {
		return
	}

	pairs = append(pairs, ps.WithName(name))

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

func (s *Service) create(ctx context.Context, name string, opt *pairServiceCreate) (store storage.Storager, err error) {
	panic("implement it")
}
func (s *Service) delete(ctx context.Context, name string, opt *pairServiceDelete) (err error) {
	panic("implement it")
}
func (s *Service) get(ctx context.Context, name string, opt *pairServiceGet) (store storage.Storager, err error) {
	panic("implement it")
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
