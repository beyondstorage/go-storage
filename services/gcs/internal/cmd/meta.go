package main

import (
	def "go.beyondstorage.io/v5/definitions"
	"go.beyondstorage.io/v5/types"
)

var Metadata = def.Metadata{
	Name: "gcs",
	Pairs: []def.Pair{
		pairEncryptionKey,
		pairKmsKeyName,
		pairProjectId,
		pairStorageClass,
	},
	Infos: []def.Info{
		infoObjectMetaStorageClass,
		infoObjectMetaEncryptionKeySha256,
	},
	Factory: []def.Pair{
		def.PairCredential,
		def.PairName,
		def.PairWorkDir,
		pairProjectId,
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
		CreateDir: []def.Pair{
			pairStorageClass,
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
		},
		Write: []def.Pair{
			def.PairContentMD5,
			def.PairContentType,
			def.PairIoCallback,
			pairStorageClass,
			pairEncryptionKey,
			pairKmsKeyName,
		},
		Stat: []def.Pair{
			def.PairObjectMode,
		},
	},
}

var pairEncryptionKey = def.Pair{
	Name:        "encryption_key",
	Type:        def.Type{Name: "[]byte"},
	Description: "is the customer's 32-byte AES-256 key",
}
var pairKmsKeyName = def.Pair{
	Name:        "kms_key_name",
	Type:        def.Type{Name: "string"},
	Description: "is the Cloud KMS key resource. For example, `projects/my-pet-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key`.\n\nRefer to https://cloud.google.com/storage/docs/encryption/using-customer-managed-keys#add-object-key for more details.",
}
var pairProjectId = def.Pair{
	Name: "project_id",
	Type: def.Type{Name: "string"},
}
var pairStorageClass = def.Pair{
	Name:        "storage_class",
	Type:        def.Type{Name: "string"},
	Defaultable: true,
}

var infoObjectMetaStorageClass = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "storage_class",
	Type:      def.Type{Name: "string"},
}
var infoObjectMetaEncryptionKeySha256 = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "encryption_key_sha256",
	Type:      def.Type{Name: "string"},
}
