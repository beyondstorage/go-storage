package main

import (
	def "go.beyondstorage.io/v5/definitions"
	"go.beyondstorage.io/v5/types"
)

var Metadata = def.Metadata{
	Name:  "us3",
	Infos: []def.Info{},
	Pairs: []def.Pair{},
	Factory: []def.Pair{
		def.PairWorkDir,
	},
	Service: def.Service{
		Features: types.ServiceFeatures{},
	},
	Storage: def.Storage{
		Features: types.StorageFeatures{},
	},
}
