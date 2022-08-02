package main

import (
	def "github.com/beyondstorage/go-storage/v5/definitions"
	"github.com/beyondstorage/go-storage/v5/types"
)

var Metadata = def.Metadata{
	Name:  "tar",
	Pairs: []def.Pair{},
	Infos: []def.Info{},
	Factory: []def.Pair{
		def.PairEndpoint,
		def.PairWorkDir,
	},
	Service: def.Service{
		Features: types.ServiceFeatures{},
	},
	Storage: def.Storage{
		Features: types.StorageFeatures{
			List: true,
			Read: true,
			Stat: true,
		},

		List: []def.Pair{
			def.PairListMode,
		},
		Read: []def.Pair{
			def.PairOffset,
			def.PairIoCallback,
			def.PairSize,
		},
		Stat: []def.Pair{
			def.PairObjectMode,
		},
	},
}
