module go.beyondstorage.io/services/s3/v3

go 1.16

require (
	github.com/aws/aws-sdk-go-v2 v1.13.0
	github.com/aws/aws-sdk-go-v2/config v1.13.1
	github.com/aws/aws-sdk-go-v2/credentials v1.8.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.24.1
	github.com/aws/smithy-go v1.10.0
	github.com/google/uuid v1.3.0
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/endpoint v1.2.0
	go.beyondstorage.io/v5 v5.0.0
)

replace go.beyondstorage.io/v5 => ../../
