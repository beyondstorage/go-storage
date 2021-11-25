module go.beyondstorage.io/services/minio

go 1.16

require (
	github.com/google/uuid v1.3.0
	github.com/minio/minio-go/v7 v7.0.16
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/endpoint v1.2.0
	go.beyondstorage.io/v5 v5.0.0
)

replace go.beyondstorage.io/v5 => ../../
