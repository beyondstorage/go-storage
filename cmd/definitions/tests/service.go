package tests

import (
	"context"

	"go.beyondstorage.io/v5/types"
)

type Service struct {
	defaultPairs DefaultServicePairs
	features     ServiceFeatures //nolint:structcheck,unused

	Pairs []types.Pair

	types.UnimplementedServicer
}

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store types.Storager, err error) {
	panic("not implemented")
}

func (s *Service) delete(ctx context.Context, name string, opt pairServiceDelete) (err error) {
	panic("not implemented")
}

func (s *Service) get(ctx context.Context, name string, opt pairServiceGet) (store types.Storager, err error) {
	panic("not implemented")
}

func (s *Service) list(ctx context.Context, opt pairServiceList) (sti *types.StoragerIterator, err error) {
	panic("not implemented")
}
