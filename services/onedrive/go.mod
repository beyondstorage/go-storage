module github.com/beyondstorage/go-storage/services/onedrive

go 1.16

require (
	github.com/beyondstorage/go-storage/credential v1.0.0
	github.com/beyondstorage/go-storage/v5 v5.0.0
	github.com/goh-chunlin/go-onedrive v1.1.1
	github.com/google/uuid v1.4.0
	github.com/kr/pretty v0.3.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	golang.org/x/oauth2 v0.20.0
	golang.org/x/sys v0.15.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace (
	github.com/beyondstorage/go-storage/credential => ../../credential
	github.com/beyondstorage/go-storage/v5 => ../../
)
