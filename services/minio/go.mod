module go.beyondstorage.io/services/minio

go 1.16

require (
	github.com/google/uuid v1.3.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/minio/minio-go/v7 v7.0.28
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/endpoint v1.2.0
	go.beyondstorage.io/v5 v5.0.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace go.beyondstorage.io/v5 => ../../
