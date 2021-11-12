package main

import (
	def "go.beyondstorage.io/v5/definitions"
	"go.beyondstorage.io/v5/types"
)

var Metadata = def.Metadata{
	Name: "azblob",
	Pairs: []def.Pair{
		pairAccessTier,
		pairEncryptionKey,
		pairEncryptionScope,
	},
	Infos: []def.Info{
		infoObjectMetaAccessTier,
		infoObjectMetaEncryptionKeySha256,
		infoObjectMetaEncryptionScope,
		infoObjectMetaServerEncrypted,
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
	},
	Storage: def.Storage{
		Features: types.StorageFeatures{
			VirtualDir:       true,
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
			pairEncryptionKey,
			pairEncryptionScope,
		},
		Write: []def.Pair{
			def.PairContentMD5,
			def.PairContentType,
			def.PairIoCallback,
			pairAccessTier,
			pairEncryptionKey,
			pairEncryptionScope,
		},
		Stat: []def.Pair{
			def.PairObjectMode,
			pairEncryptionKey,
			pairEncryptionScope,
		},
		CreateAppend: []def.Pair{
			def.PairContentType,
			pairEncryptionKey,
			pairEncryptionScope,
		},
		WriteAppend: []def.Pair{
			def.PairContentMD5,
			pairEncryptionKey,
			pairEncryptionScope,
		},
		CreateDir: []def.Pair{
			pairAccessTier,
		},
	},
}

var pairAccessTier = def.Pair{
	Name:        "access_tier",
	Type:        def.Type{Name: "string"},
	Description: "See https://docs.microsoft.com/en-us/azure/storage/blobs/access-tiers-overview for details. Specifies the access tier.",
}
var pairEncryptionKey = def.Pair{
	Name:        "encryption_key",
	Type:        def.Type{Name: "[]byte"},
	Description: "is the customer's 32-byte AES-256 key",
}
var pairEncryptionScope = def.Pair{
	Name:        "encryption_scope",
	Type:        def.Type{Name: "string"},
	Description: "See https://docs.microsoft.com/en-us/azure/storage/blobs/encryption-scope-overview for details. Specifies the name of the encryption scope.",
}
var infoObjectMetaAccessTier = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "Access_tier",
	Type:      def.Type{Name: "string"},
}
var infoObjectMetaEncryptionKeySha256 = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "encryption_key_sha256",
	Type:      def.Type{Name: "string"},
}
var infoObjectMetaEncryptionScope = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "encryption_scope",
	Type:      def.Type{Name: "string"},
}
var infoObjectMetaServerEncrypted = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "server_encrypted",
	Type:      def.Type{Name: "bool"},
}
