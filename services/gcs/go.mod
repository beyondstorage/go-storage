module github.com/beyondstorage/go-storage/services/gcs/v3

go 1.16

require (
	cloud.google.com/go/storage v1.30.1
	github.com/beyondstorage/go-storage/credential v1.0.0
	github.com/beyondstorage/go-storage/v5 v5.0.0
	github.com/google/uuid v1.3.0
	golang.org/x/oauth2 v0.10.0
	google.golang.org/api v0.135.0
)

replace (
	github.com/beyondstorage/go-storage/credential => ../../credential
	github.com/beyondstorage/go-storage/v5 => ../../
)
