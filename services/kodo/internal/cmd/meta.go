package main

import (
	def "go.beyondstorage.io/v5/definitions"
	"go.beyondstorage.io/v5/types"
)

var Metadata = def.Metadata{
	Name: "kodo",
	Pairs: []def.Pair{
		pairStorageClass,
	},
	Infos: []def.Info{
		infoObjectMetaStorageClass,
	},
	Factory: []def.Pair{
		def.PairCredential,
		def.PairEndpoint,
		def.PairName,
		def.PairWorkDir,
	},
	Service: def.Service{
		Features: types.ServiceFeatures{
			Create: true,
			Delete: true,
			Get:    true,
			List:   true,
		},
		Create: []def.Pair{
			def.PairLocation,
		},
	},
	Storage: def.Storage{
		Features: types.StorageFeatures{
			VirtualDir:       true,
			WriteEmptyObject: true,

			Create:    true,
			CreateDir: true,
			Delete:    true,
			List:      true,
			Metadata:  true,
			Read:      true,
			Stat:      true,
			Write:     true,
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
			pairStorageClass,
		},
		Stat: []def.Pair{
			def.PairObjectMode,
		},
	},
}

var pairStorageClass = def.Pair{
	Name:        "storage_class",
	Type:        def.Type{Name: "int"},
	Defaultable: true,
}

var infoObjectMetaStorageClass = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "storage_class",
	Type:      def.Type{Name: "int"},
}
