package main

import (
	"context"
	"go.beyondstorage.io/v5/types"
	"io"
)

type Storage struct {
	defaultPairs types.DefaultStoragePairs
	features     types.StorageFeatures

	Pairs []types.Pair

	types.UnimplementedStorager
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	panic("not implemented")
}
