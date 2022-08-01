module github.com/beyondstorage/go-storage/services/azblob/v3

go 1.16

require (
	github.com/Azure/azure-storage-blob-go v0.14.0
	github.com/beyondstorage/go-storage/credential v1.0.0
	github.com/beyondstorage/go-storage/endpoint v1.2.0
	github.com/beyondstorage/go-storage/v5 v5.0.0
	github.com/google/uuid v1.3.0
)

replace (
	github.com/beyondstorage/go-storage/credential => ../../credential
	github.com/beyondstorage/go-storage/endpoint => ../../endpoint
	github.com/beyondstorage/go-storage/v5 => ../../
)
