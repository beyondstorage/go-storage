package main

import (
	def "go.beyondstorage.io/v5/definitions"
	"go.beyondstorage.io/v5/types"
)

var Metadata = def.Metadata{
	Name: "ipfs",
	Pairs: []def.Pair{
		pairGateway,
	},
	Infos: []def.Info{
		infoObjectMetaHash,
		infoObjectMetaBlocks,
		infoObjectMetaCumulativeSize,
		infoObjectMetaLocal,
		infoObjectMetaWithLocality,
		infoObjectMetaSizeLocal,
	},
	Factory: []def.Pair{
		def.PairEndpoint,
		def.PairWorkDir,
		pairGateway,
	},
	Service: def.Service{},
	Storage: def.Storage{
		Features: types.StorageFeatures{
			WriteEmptyObject: true,

			Create:            true,
			CreateDir:         true,
			Copy:              true,
			Delete:            true,
			List:              true,
			Metadata:          true,
			Move:              true,
			QuerySignHTTPRead: true,
			Read:              true,
			Stat:              true,
			Write:             true,
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

var pairGateway = def.Pair{
	Name:        "gateway",
	Type:        def.Type{Name: "string"},
	Description: "set storage gateway, for http(s) request purpose.",
}
var infoObjectMetaHash = def.Info{
	Namespace:   def.NamespaceObject,
	Category:    def.CategoryMeta,
	Name:        "hash",
	Type:        def.Type{Name: "string"},
	Description: "the CID of the file or directory",
}
var infoObjectMetaBlocks = def.Info{
	Namespace:   def.NamespaceObject,
	Category:    def.CategoryMeta,
	Name:        "blocks",
	Type:        def.Type{Name: "int"},
	Description: "the number of files in the directory or the number of blocks that make up the file",
}
var infoObjectMetaCumulativeSize = def.Info{
	Namespace:   def.NamespaceObject,
	Category:    def.CategoryMeta,
	Name:        "cumulative_size",
	Type:        def.Type{Name: "uint64"},
	Description: "the size of the DAGNodes making up the file in Bytes, or the sum of the sizes of all files in the directory",
}
var infoObjectMetaLocal = def.Info{
	Namespace:   def.NamespaceObject,
	Category:    def.CategoryMeta,
	Name:        "local",
	Type:        def.Type{Name: "bool"},
	Description: "whether the file`s dags is fully present locally",
}
var infoObjectMetaWithLocality = def.Info{
	Namespace:   def.NamespaceObject,
	Category:    def.CategoryMeta,
	Name:        "with_locality;",
	Type:        def.Type{Name: "bool"},
	Description: "whether the locality information is present",
}
var infoObjectMetaSizeLocal = def.Info{
	Namespace:   def.NamespaceObject,
	Category:    def.CategoryMeta,
	Name:        "size_local;",
	Type:        def.Type{Name: "uint64"},
	Description: "the cumulative size of the data present locally",
}
