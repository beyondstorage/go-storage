module go.beyondstorage.io/services/obs/v2

go 1.16

require (
	github.com/google/uuid v1.3.0
	github.com/huaweicloud/huaweicloud-sdk-go-obs v3.21.12+incompatible
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/endpoint v1.2.0
	go.beyondstorage.io/v5 v5.0.0
)

replace go.beyondstorage.io/v5 => ../../
