package main

import (
	def "github.com/beyondstorage/go-storage/v5/definitions"
	"github.com/beyondstorage/go-storage/v5/types"
)

var Metadata = def.Metadata{
	Name: "oss",
	Pairs: []def.Pair{
		pairStorageClass,
		pairServerSideEncryption,
		pairServerSideDataEncryption,
		pairServerSideEncryptionKeyId,
	},
	Infos: []def.Info{
		infoObjectMetaStorageClass,
		infoObjectMetaServerSideEncryption,
		infoObjectMetaServerSideEncryptionKeyId,
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

			CompleteMultipart: true,
			CommitAppend:      true,
			Create:            true,
			CreateAppend:      true,
			CreateDir:         true,
			CreateLink:        true,
			CreateMultipart:   true,
			Delete:            true,
			List:              true,
			ListMultipart:     true,
			Metadata:          true,
			Read:              true,
			Stat:              true,
			Write:             true,
			WriteAppend:       true,
			WriteMultipart:    true,
		},

		Create: []def.Pair{
			def.PairMultipartID,
			def.PairObjectMode,
		},
		CreateDir: []def.Pair{
			pairStorageClass,
		},
		CreateAppend: []def.Pair{
			def.PairContentType,
			pairStorageClass,
			pairServerSideEncryption,
		},
		CreateMultipart: []def.Pair{
			def.PairContentType,
			pairStorageClass,
			pairServerSideEncryption,
			pairServerSideDataEncryption,
			pairServerSideEncryptionKeyId,
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
		},
		Write: []def.Pair{
			def.PairContentMD5,
			def.PairContentType,
			def.PairIoCallback,
			pairStorageClass,
			pairServerSideEncryption,
			pairServerSideDataEncryption,
			pairServerSideEncryptionKeyId,
		},
		WriteAppend: []def.Pair{
			def.PairContentMD5,
			def.PairIoCallback,
		},
		WriteMultipart: []def.Pair{
			def.PairContentMD5,
		},
		Stat: []def.Pair{
			def.PairMultipartID,
			def.PairObjectMode,
		},
	},
}

var pairStorageClass = def.Pair{
	Name:        "storage_class",
	Type:        def.Type{Name: "string"},
	Defaultable: true,
}
var pairServerSideEncryption = def.Pair{
	Name:        "server_side_encryption",
	Type:        def.Type{Name: "string"},
	Description: "specifies the encryption algorithm. Can be AES256, KMS or SM4.\n\nFor Chinese users, refer to https://help.aliyun.com/document_detail/31871.html for details.\n\nFor global users, refer to https://www.alibabacloud.com/help/doc-detail/31871.htm for details, and double-check whether SM4 can be used.",
}
var pairServerSideDataEncryption = def.Pair{
	Name:        "server_side_data_encryption",
	Type:        def.Type{Name: "string"},
	Description: "specifies the encryption algorithm when server_side_encryption is KMS. Can only be set to SM4. If this is not set, AES256 will be used.\n\nFor Chinese users, refer to https://help.aliyun.com/document_detail/31871.html for details.\n\nFor global users, refer to https://www.alibabacloud.com/help/doc-detail/31871.htm for details, and double-check whether SM4 can be used.",
}
var pairServerSideEncryptionKeyId = def.Pair{
	Name:        "server_side_encryption_key_id",
	Type:        def.Type{Name: "string"},
	Description: "specifies the COS KMS key ID to use for object encryption.",
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
var infoObjectMetaServerSideEncryptionKeyId = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "server_side_encryption_key_id",
	Type:      def.Type{Name: "string"},
}
