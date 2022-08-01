package cos

import "github.com/beyondstorage/go-storage/v5/services"

var (
	// ErrServerSideEncryptionCustomerKeyInvalid will be returned while server-side encryption customer key is invalid.
	ErrServerSideEncryptionCustomerKeyInvalid = services.NewErrorCode("invalid server-side encryption customer key")
)
