package main

import (
	"go.beyondstorage.io/v5/types"
)

type Storage struct {
	defaultPairs types.DefaultStoragePairs
	features     types.StorageFeatures

	objects []*types.Object

	Pairs []types.Pair

	types.UnimplementedStorager
}
