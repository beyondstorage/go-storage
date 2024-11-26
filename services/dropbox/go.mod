module github.com/beyondstorage/go-storage/services/dropbox/v3

go 1.16

require (
	github.com/beyondstorage/go-storage/credential v1.0.0
	github.com/beyondstorage/go-storage/v5 v5.0.0
	github.com/dropbox/dropbox-sdk-go-unofficial/v6 v6.0.5
	github.com/google/uuid v1.6.0
)

require golang.org/x/oauth2 v0.0.0-20210413134643-5e61552d6c78 // indirect

replace (
	github.com/beyondstorage/go-storage/credential => ../../credential
	github.com/beyondstorage/go-storage/v5 => ../../
)
