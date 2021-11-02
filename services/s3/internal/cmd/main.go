package main

import (
	def "go.beyondstorage.io/v5/definitions"
	"go.beyondstorage.io/v5/types"
)

func main() {
	def.GenerateService(metadata, "generated.go")
}

var metadata = def.Metadata{
	Name: "s3",
	Pairs: []def.Pair{
		pairForcePathStyle,
		pairDisable100Continue,
		pairUseAccelerate,
		pairUseArnRegion,
		pairExpectedBucketOwner,
		pairStorageClass,
		pairServerSideEncryption,
		pairServerSideEncryptionBucketKeyEnable,
		pairServerSideEncryptionCustomerAlgorithm,
		pairServerSideEncryptionCustomerKey,
		pairServerSideEncryptionAwsKmsKeyId,
		pairServerSideEncryptionContext,
	},
	Infos: []def.Info{
		infoObjectMetaStorageClass,
		infoObjectMetaServerSideEncryption,
		infoObjectMetaServerSideEncryptionCustomerAlgorithm,
		infoObjectMetaServerSideEncryptionAwsKmsKeyId,
		infoObjectMetaServerSideEncryptionContext,
		infoObjectMetaServerSideEncryptionCustomerKeyMd5,
		infoObjectMetaServerSideEncryptionBucketKeyEnabled,
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
			pairExpectedBucketOwner,
		},
		Get: []def.Pair{
			def.PairLocation,
		},
	},
	Storage: def.Storage{
		Features: types.StorageFeatures{
			VirtualDir:  true,
			VirtualLink: true,

			CompleteMultipart:           true,
			Create:                      true,
			CreateDir:                   true,
			CreateLink:                  true,
			CreateMultipart:             true,
			Delete:                      true,
			List:                        true,
			ListMultipart:               true,
			Metadata:                    true,
			QuerySignHTTPRead:           true,
			QuerySignHTTPWrite:          true,
			QuerySignHTTPWriteMultipart: true,
			Read:                        true,
			Stat:                        true,
			Write:                       true,
			WriteMultipart:              true,
		},

		Create: []def.Pair{
			def.PairMultipartID,
			def.PairObjectMode,
		},
		CreateDir: []def.Pair{
			pairExpectedBucketOwner,
			pairStorageClass,
		},
		Delete: []def.Pair{
			pairExpectedBucketOwner,
			def.PairMultipartID,
			def.PairObjectMode,
		},
		List: []def.Pair{
			def.PairListMode,
			pairExpectedBucketOwner,
		},
		Read: []def.Pair{
			def.PairOffset,
			def.PairIoCallback,
			def.PairSize,
			pairExpectedBucketOwner,
			pairServerSideEncryptionCustomerAlgorithm,
			pairServerSideEncryptionCustomerKey,
		},
		Write: []def.Pair{
			def.PairContentMD5,
			def.PairContentType,
			def.PairIoCallback,
			pairStorageClass,
			pairExpectedBucketOwner,
			pairServerSideEncryptionBucketKeyEnable,
			pairServerSideEncryptionCustomerAlgorithm,
			pairServerSideEncryptionCustomerKey,
			pairServerSideEncryptionAwsKmsKeyId,
			pairServerSideEncryptionContext,
			pairServerSideEncryption,
		},
		Stat: []def.Pair{
			def.PairMultipartID,
			def.PairObjectMode,
			pairExpectedBucketOwner,
			pairServerSideEncryptionCustomerAlgorithm,
			pairServerSideEncryptionCustomerKey,
		},
		CreateMultipart: []def.Pair{
			pairExpectedBucketOwner,
			pairServerSideEncryptionBucketKeyEnable,
			pairServerSideEncryptionCustomerAlgorithm,
			pairServerSideEncryptionCustomerKey,
			pairServerSideEncryptionAwsKmsKeyId,
			pairServerSideEncryptionContext,
			pairServerSideEncryption,
		},
		WriteMultipart: []def.Pair{
			def.PairIoCallback,
			pairExpectedBucketOwner,
			pairServerSideEncryptionCustomerAlgorithm,
			pairServerSideEncryptionCustomerKey,
		},
		ListMultipart: []def.Pair{
			pairExpectedBucketOwner,
		},
		CompleteMultipart: []def.Pair{
			pairExpectedBucketOwner,
		},
		QuerySignHTTPRead: []def.Pair{
			def.PairOffset,
			def.PairSize,
			pairExpectedBucketOwner,
			pairServerSideEncryptionCustomerAlgorithm,
			pairServerSideEncryptionCustomerKey,
		},
		QuerySignHTTPWrite: []def.Pair{
			def.PairContentMD5,
			def.PairContentType,
			pairStorageClass,
			pairExpectedBucketOwner,
			pairServerSideEncryptionBucketKeyEnable,
			pairServerSideEncryptionCustomerAlgorithm,
			pairServerSideEncryptionCustomerKey,
			pairServerSideEncryptionAwsKmsKeyId,
			pairServerSideEncryptionContext,
			pairServerSideEncryption,
		},
		QuerySignHTTPDelete: []def.Pair{
			def.PairMultipartID,
			def.PairObjectMode,
			pairExpectedBucketOwner,
		},
	},
}

var pairForcePathStyle = def.Pair{
	Name:        "force_path_style",
	Type:        def.Type{Name: "bool"},
	Description: "see http://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html for Amazon S3: Virtual Hosting of Buckets",
}
var pairDisable100Continue = def.Pair{
	Name:        "disable_100_continue",
	Type:        def.Type{Name: "bool"},
	Description: "set this to `true` to disable the SDK adding the `Expect: 100-Continue` header to PUT requests over 2MB of content",
}
var pairUseAccelerate = def.Pair{
	Name:        "use_accelerate",
	Type:        def.Type{Name: "bool"},
	Description: "set this to `true` to enable S3 Accelerate feature",
}
var pairUseArnRegion = def.Pair{
	Name:        "use_arn_region",
	Type:        def.Type{Name: "bool"},
	Description: "set this to `true` to have the S3 service client to use the region specified in the ARN, when an ARN is provided as an argument to a bucket parameter",
}
var pairExpectedBucketOwner = def.Pair{
	Name:        "expected_bucket_owner",
	Type:        def.Type{Name: "string"},
	Description: "the account ID of the expected bucket owner",
}

var pairStorageClass = def.Pair{
	Name:        "storage_class",
	Type:        def.Type{Name: "string"},
	Defaultable: true,
}
var pairServerSideEncryption = def.Pair{
	Name:        "server_side_encryption",
	Type:        def.Type{Name: "string"},
	Description: "the server-side encryption algorithm used when storing this object in Amazon",
}
var pairServerSideEncryptionBucketKeyEnable = def.Pair{
	Name:        "server_side_encryption_bucket_key_enabled",
	Type:        def.Type{Name: "bool"},
	Description: "specifies whether Amazon S3 should use an S3 Bucket Key for object encryption with server-side encryption using AWS KMS (SSE-KMS)",
}
var pairServerSideEncryptionCustomerAlgorithm = def.Pair{
	Name:        "server_side_encryption_customer_algorithm",
	Type:        def.Type{Name: "string"},
	Description: "specifies the algorithm to use to when encrypting the object. The header value must be `AES256`.",
}
var pairServerSideEncryptionCustomerKey = def.Pair{
	Name:        "server_side_encryption_customer_key",
	Type:        def.Type{Name: "[]byte"},
	Description: "specifies the customer-provided encryption key for Amazon S3 to use to encrypt/decrypt the source object. It must be 32-byte AES-256 key.",
}
var pairServerSideEncryptionAwsKmsKeyId = def.Pair{
	Name:        "server_side_encryption_aws_kms_key_id",
	Type:        def.Type{Name: "string"},
	Description: "specifies the AWS KMS key ID to use for object encryption",
}
var pairServerSideEncryptionContext = def.Pair{
	Name:        "server_side_encryption_context",
	Type:        def.Type{Name: "string"},
	Description: "specifies the AWS KMS Encryption Context to use for object encryption. The value of this header is a base64-encoded UTF-8 string holding JSON with the encryption context key-value pairs.",
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
var infoObjectMetaServerSideEncryptionAwsKmsKeyId = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "server_side_encryption_aws_kms_key_id",
	Type:      def.Type{Name: "string"},
}
var infoObjectMetaServerSideEncryptionContext = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "server_side_encryption_context",
	Type:      def.Type{Name: "string"},
}
var infoObjectMetaServerSideEncryptionCustomerKeyMd5 = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "server_side_encryption_customer_key_md5",
	Type:      def.Type{Name: "string"},
}
var infoObjectMetaServerSideEncryptionBucketKeyEnabled = def.Info{
	Namespace: def.NamespaceObject,
	Category:  def.CategoryMeta,
	Name:      "server_side_encryption_bucket_key_enabled",
	Type:      def.Type{Name: "bool"},
}
