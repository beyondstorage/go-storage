package main

import (
	"go.beyondstorage.io/v5/types"
)

type Service struct {
	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	Pairs []types.Pair

	types.UnimplementedServicer
}
