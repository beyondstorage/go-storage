module github.com/beyondstorage/go-storage/services/gdrive

go 1.16

require (
	github.com/beyondstorage/go-storage/credential v1.0.0
	github.com/beyondstorage/go-storage/v5 v5.0.0
	github.com/dgraph-io/ristretto v0.1.0
	github.com/google/uuid v1.3.0
	golang.org/x/oauth2 v0.0.0-20220822191816-0ebed06d0094
	google.golang.org/api v0.97.0
)

replace (
	github.com/beyondstorage/go-storage/credential => ../../credential
	github.com/beyondstorage/go-storage/v5 => ../../
)
