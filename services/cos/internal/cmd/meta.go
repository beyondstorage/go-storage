package main

import (
	def "go.beyondstorage.io/v5/definitions"
	"go.beyondstorage.io/v5/types"
)

var Metadata = def.Metadata{
	Name: "cos",
	Pairs: []def.Pair{
		pairStorageClass,
		pairServerSideEncryption,
		pairServerSideEncryptionCustomerAlgorithm,
		pairServerSideEncryptionCustomerKey,
		pairServerSideEncryptionContext,
	},
	Infos: []def.Info{
		infoObjectMetaStorageClass,
		infoObjectMetaServerSideEncryption,
		infoObjectMetaServerSideEncryptionCosKmsKeyId,
		infoObjectMetaServerSideEncryptionCustomerAlgorithm,
		infoObjectMetaServerSideEncryptionCustomerKeyMd5,
	},
	Factory: []def.Pair{
		def.PairCredential,
		def.PairEndpoint,
		def.PairName,
		def.PairLocation,
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
		Delete: []def.Pair{
			def.PairLocation,
		},
		Get: []def.Pair{
			def.PairLocation,
		},
	},
	Storage: def.Storage{
		Features: types.StorageFeatures{
			VirtualDir:       true,
			WriteEmptyObject: true,

			CompleteMultipart: true,
			Create:            true,
			CreateDir:         true,
			CreateMultipart:   true,
			Delete:            true,
			List:              true,
			ListMultipart:     true,
			Metadata:          true,
			Read:              true,
			Stat:              true,
			Write:             true,
			WriteMultipart:    true,
		},

		Create: []def.Pair{
			def.PairMultipartID,
			def.PairObjectMode,
		},
		Delete: []def.Pair{
			def.PairMultipartID,
			def.PairObjectMode,
		},
		List: []def.Pair{
			def.PairListMode,
		},
		Read: []def.Pair{
			def.PairOffset,
			def.PairIoCallback,
			def.PairSize,
			pairServerSideEncryptionCustomerAlgorithm,
			pairServerSideEncryptionCustomerKey,
		},
		Write: []def.Pair{
			def.PairContentMD5,
			def.PairContentType,
			def.PairIoCallback,
			pairStorageClass,
			pairServerSideEncryption,
			pairServerSideEncryptionCustomerAlgorithm,
			pairServerSideEncryptionCustomerKey,
			pairServerSideEncryptionCosKmsKeyId,
			pairServerSideEncryptionContext,
		},
		Stat: []def.Pair{
			def.PairMultipartID,
			def.PairObjectMode,
			pairServerSideEncryptionCustomerAlgorithm,
			pairServerSideEncryptionCosKmsKeyId,
			pairServerSideEncryptionCustomerKey,
		},
		CreateMultipart: []def.Pair{
			def.PairContentType,
			pairStorageClass,
			pairServerSideEncryption,
			pairServerSideEncryptionCustomerKey,
			pairServerSideEncryptionCustomerAlgorithm,
			pairServerSideEncryptionCosKmsKeyId,
			pairServerSideEncryptionContext,
		},
		WriteMultipart: []def.Pair{
			def.PairContentMD5,
		},
		CreateDir: []def.Pair{
			pairStorageClass,
		},
	},
}

var pairStorageClass = def.Pair{
	Name:        "storage_class",
	Type:        def.Type{Name: "string"},
	Defaultable: true,
}
var pairServerSideEncryptionCustomerAlgorithm = def.Pair{
	Name:        "server_side_encryption_customer_algorithm",
	Type:        def.Type{Name: "string"},
	Description: "specifies the algorithm to use to when encrypting the object. Now only `AES256` is supported.",
}
var pairServerSideEncryptionCustomerKey = def.Pair{
	Name:        "server_side_encryption_customer_key",
	Type:        def.Type{Name: "[]byte"},
	Description: "specifies the customer-provided encryption key to encrypt/decrypt the source object. It must be a 32-byte AES-256 key.",
}
var pairServerSideEncryptionCosKmsKeyId = def.Pair{
	Name:        "server_side_encryption_cos_kms_key_id",
	Type:        def.Type{Name: "string"},
	Description: "specifies the COS KMS key ID to use for object encryption.",
}
var pairServerSideEncryptionContext = def.Pair{
	Name:        "server_side_encryption_context",
	Type:        def.Type{Name: "string"},
	Description: "specifies the COS KMS Encryption Context to use for object encryption. The value of this header is a base64-encoded UTF-8 string holding JSON with the encryption context key-value pairs.",
}
var pairServerSideEncryption = def.Pair{
	Name:        "server_side_encryption",
	Type:        def.Type{Name: "string"},
	Description: "the server-side encryption algorithm used when storing this object. It can be `AES-256` for SSE-COS, and `cos/kms` for SSE-KMS.",
}

var infoObjectMetaStorageClass = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "storage_class",
	Type:      def.Type{Name: "string"},
}
var infoObjectMetaServerSideEncryption = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "server_side_encryption",
	Type:      def.Type{Name: "string"},
}
var infoObjectMetaServerSideEncryptionCustomerAlgorithm = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "server_side_encryption_customer_algorithm",
	Type:      def.Type{Name: "string"},
}
var infoObjectMetaServerSideEncryptionCosKmsKeyId = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "server_side_encryption_cos_kms_key_id",
	Type:      def.Type{Name: "string"},
}
var infoObjectMetaServerSideEncryptionCustomerKeyMd5 = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "server_side_encryption_customer_key_md5",
	Type:      def.Type{Name: "string"},
}
