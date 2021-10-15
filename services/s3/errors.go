package s3

import (
	"go.beyondstorage.io/v5/services"
)

var (
	// ErrServerSideEncryptionCustomerKeyInvalid will be returned while server-side encryption customer key is invalid.
	ErrServerSideEncryptionCustomerKeyInvalid = services.NewErrorCode("invalid server-side encryption customer key")
)
