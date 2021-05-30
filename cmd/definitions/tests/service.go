package tests

import (
	"context"

	. "github.com/beyondstorage/go-storage/v4/types"
)

type Service struct {
	defaultPairs DefaultServicePairs
	features     ServiceFeatures

	UnimplementedServicer
}

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store Storager, err error) {
	panic("not implemented")
}

func (s *Service) delete(ctx context.Context, name string, opt pairServiceDelete) (err error) {
	panic("not implemented")
}

func (s *Service) get(ctx context.Context, name string, opt pairServiceGet) (store Storager, err error) {
	panic("not implemented")
}

func (s *Service) list(ctx context.Context, opt pairServiceList) (sti *StoragerIterator, err error) {
	panic("not implemented")
}
