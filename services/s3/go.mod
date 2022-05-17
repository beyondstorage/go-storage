module go.beyondstorage.io/services/s3/v3

go 1.16

require (
	github.com/aws/aws-sdk-go-v2 v1.16.3
	github.com/aws/aws-sdk-go-v2/config v1.15.6
	github.com/aws/aws-sdk-go-v2/credentials v1.12.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.26.9
	github.com/aws/smithy-go v1.11.2
	github.com/google/uuid v1.3.0
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/endpoint v1.2.0
	go.beyondstorage.io/v5 v5.0.0
)

replace go.beyondstorage.io/v5 => ../../
