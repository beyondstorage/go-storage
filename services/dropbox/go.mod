module go.beyondstorage.io/services/dropbox/v3

go 1.16

require (
	github.com/dropbox/dropbox-sdk-go-unofficial/v6 v6.0.3
	github.com/google/uuid v1.3.0
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/v5 v5.0.0
)

require golang.org/x/oauth2 v0.0.0-20210413134643-5e61552d6c78 // indirect

replace go.beyondstorage.io/v5 => ../../
