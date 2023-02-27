module github.com/beyondstorage/go-storage/services/tar

go 1.16

require (
	github.com/beyondstorage/go-storage/endpoint v1.2.0
	github.com/beyondstorage/go-storage/v5 v5.0.0
	github.com/stretchr/testify v1.8.2
)

replace (
	github.com/beyondstorage/go-storage/endpoint => ../../endpoint
	github.com/beyondstorage/go-storage/v5 => ../../
)
