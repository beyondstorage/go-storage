package main

import (
	def "github.com/beyondstorage/go-storage/v5/definitions"
	"github.com/beyondstorage/go-storage/v5/types"
)

var Metadata = def.Metadata{
	Name:  "zip",
	Pairs: []def.Pair{},
	Infos: []def.Info{},
	Factory: []def.Pair{
		def.PairWorkDir,
	},
	Service: def.Service{},
	Storage: def.Storage{
		Features: types.StorageFeatures{},
	},
}
