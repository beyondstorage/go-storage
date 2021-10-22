package qingstor

import "go.beyondstorage.io/v5/services"

var (
	// ErrBucketNameInvalid will be returned while bucket name is invalid.
	ErrBucketNameInvalid = services.NewErrorCode("invalid bucket name")

	// ErrWorkDirInvalid will be returned while work dir is invalid.
	// Work dir must start and end with only one '/'
	ErrWorkDirInvalid = services.NewErrorCode("invalid work dir")

	// ErrEncryptionCustomerKeyInvalid will be returned while encryption customer key is invalid.
	// Encryption key must be a 32-byte AES-256 key.
	ErrEncryptionCustomerKeyInvalid = services.NewErrorCode("invalid encryption customer key")

	// ErrAppendNextPositionEmpty will be returned while next append position is empty.
	ErrAppendNextPositionEmpty = services.NewErrorCode("next append position is empty")

	// ErrPartNumberInvalid will be returned while part number is out of range [0, 10000] when uploading multipart.
	ErrPartNumberInvalid = services.NewErrorCode("part number is out of range [0, 10000]")
)
