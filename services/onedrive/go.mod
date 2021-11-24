module go.beyondstorage.io/services/onedrive

go 1.16

require (
	github.com/goh-chunlin/go-onedrive v1.1.1
	github.com/google/uuid v1.3.0
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/v5 v5.0.0
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
)

replace go.beyondstorage.io/v5 => ../../
