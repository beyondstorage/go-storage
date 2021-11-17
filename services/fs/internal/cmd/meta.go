package main

import (
	def "go.beyondstorage.io/v5/definitions"
	"go.beyondstorage.io/v5/types"
)

var Metadata = def.Metadata{
	Name: "fs",
	Pairs: []def.Pair{
		pairContinuationToken,
	},
	Infos: []def.Info{},
	Factory: []def.Pair{
		def.PairWorkDir,
	},
	Service: def.Service{},
	Storage: def.Storage{
		Features: types.StorageFeatures{
			WriteEmptyObject: true,

			Create:       true,
			CreateDir:    true,
			CreateAppend: true,
			Delete:       true,
			List:         true,
			Metadata:     true,
			Read:         true,
			Stat:         true,
			Write:        true,
			WriteAppend:  true,
			CommitAppend: true,
		},

		Create: []def.Pair{
			def.PairObjectMode,
		},
		Delete: []def.Pair{
			def.PairObjectMode,
		},
		List: []def.Pair{
			def.PairListMode,
			pairContinuationToken,
		},
		Read: []def.Pair{
			def.PairOffset,
			def.PairIoCallback,
			def.PairSize,
		},
		Write: []def.Pair{
			def.PairContentMD5,
			def.PairContentType,
			def.PairOffset,
			def.PairIoCallback,
		},
		Stat: []def.Pair{
			def.PairObjectMode,
		},
	},
}

var pairContinuationToken = def.Pair{
	Name:        "continuation_token",
	Type:        def.Type{Name: "string"},
	Description: "specify the continuation token for list.",
}