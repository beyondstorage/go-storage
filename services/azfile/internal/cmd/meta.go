package main

import (
	def "go.beyondstorage.io/v5/definitions"
	"go.beyondstorage.io/v5/types"
)

var Metadata = def.Metadata{
	Name: "azfile",
	Infos: []def.Info{
		infoObjectMetaServerEncrypted,
	},
	Pairs: []def.Pair{},
	Factory: []def.Pair{
		def.PairCredential,
		def.PairEndpoint,
		def.PairName,
		def.PairWorkDir,
	},
	Service: def.Service{
		Features: types.ServiceFeatures{},
	},
	Storage: def.Storage{
		Features: types.StorageFeatures{
			WriteEmptyObject: true,

			Create: true,
			Delete: true,
			List:   true,
			Read:   true,
			Stat:   true,
			Write:  true,
		},

		Create: []def.Pair{
			def.PairObjectMode,
		},
		Delete: []def.Pair{
			def.PairObjectMode,
		},
		List: []def.Pair{
			def.PairListMode,
		},
		Read: []def.Pair{
			def.PairOffset,
			def.PairIoCallback,
			def.PairSize,
		},
		Write: []def.Pair{
			def.PairContentMD5,
			def.PairContentType,
			def.PairIoCallback,
		},
		Stat: []def.Pair{
			def.PairObjectMode,
		},
	},
}

var infoObjectMetaServerEncrypted = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "server_encrypted",
	Type:      def.Type{Name: "bool"},
}
