module github.com/beyondstorage/go-storage/services/gdrive

go 1.16

require (
	github.com/beyondstorage/go-storage/credential v1.0.0
	github.com/beyondstorage/go-storage/v5 v5.0.0
	github.com/dgraph-io/ristretto v0.1.1
	github.com/google/uuid v1.3.0
	golang.org/x/oauth2 v0.11.0
	google.golang.org/api v0.134.0
)

replace (
	github.com/beyondstorage/go-storage/credential => ../../credential
	github.com/beyondstorage/go-storage/v5 => ../../
)
