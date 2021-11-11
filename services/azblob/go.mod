module go.beyondstorage.io/services/azblob/v3

go 1.16

require (
	github.com/Azure/azure-storage-blob-go v0.14.0
	github.com/google/uuid v1.3.0
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/endpoint v1.2.0
	go.beyondstorage.io/v5 v5.0.0
)

replace go.beyondstorage.io/v5 => ../../
