module github.com/beyondstorage/go-storage/services/gcs/v3

go 1.16

require (
	cloud.google.com/go/storage v1.27.0
	github.com/beyondstorage/go-storage/credential v1.0.0
	github.com/beyondstorage/go-storage/v5 v5.0.0
	github.com/google/uuid v1.3.0
	golang.org/x/oauth2 v0.0.0-20220909003341-f21342109be1
	google.golang.org/api v0.97.0
)

replace (
	github.com/beyondstorage/go-storage/credential => ../../credential
	github.com/beyondstorage/go-storage/v5 => ../../
)
