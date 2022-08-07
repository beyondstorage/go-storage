module github.com/beyondstorage/go-storage/services/ftp

go 1.16

require (
	github.com/beyondstorage/go-storage/credential v1.0.0
	github.com/beyondstorage/go-storage/endpoint v1.2.0
	github.com/beyondstorage/go-storage/v5 v5.0.0
	github.com/jlaffaye/ftp v0.0.0-20210307004419-5d4190119067
	github.com/qingstor/go-mime v0.1.0
	go.beyondstorage.io/services/ftp v0.3.0
)

replace (
	github.com/beyondstorage/go-storage/credential => ../../credential
	github.com/beyondstorage/go-storage/endpoint => ../../endpoint
	github.com/beyondstorage/go-storage/v5 => ../../
)
