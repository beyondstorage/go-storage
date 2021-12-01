package main

import (
	def "go.beyondstorage.io/v5/definitions"
	"go.beyondstorage.io/v5/types"
)

var Metadata = def.Metadata{
	Name:  "dropbox",
	Pairs: []def.Pair{},
	Infos: []def.Info{
		infoObjectMetaUploadSessionId,
	},
	Factory: []def.Pair{
		def.PairCredential,
		def.PairWorkDir,
	},
	Service: def.Service{
		Features: types.ServiceFeatures{},
	},
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

var infoObjectMetaUploadSessionId = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "upload_session_id",
	Type:      def.Type{Name: "string"},
}
