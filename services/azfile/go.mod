module github.com/beyondstorage/go-storage/services/azfile

go 1.16

require (
	github.com/Azure/azure-storage-file-go v0.8.0
	github.com/beyondstorage/go-storage/credential v1.0.0
	github.com/beyondstorage/go-storage/endpoint v1.2.0
	github.com/beyondstorage/go-storage/v5 v5.0.0
	github.com/google/uuid v1.6.0
	github.com/pkg/errors v0.9.1 // indirect
)

replace (
	github.com/beyondstorage/go-storage/credential => ../../credential
	github.com/beyondstorage/go-storage/endpoint => ../../endpoint
	github.com/beyondstorage/go-storage/v5 => ../../
)
