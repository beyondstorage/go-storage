module go.beyondstorage.io/services/oss/v3

go 1.16

require (
	github.com/aliyun/aliyun-oss-go-sdk v2.2.0+incompatible
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/google/uuid v1.3.0
	github.com/satori/go.uuid v1.2.0 // indirect
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/endpoint v1.2.0
	go.beyondstorage.io/v5 v5.0.0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
)

replace go.beyondstorage.io/v5 => ../../
