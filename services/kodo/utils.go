package kodo

import (
	"errors"
	"fmt"
	"time"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	qs "github.com/qiniu/api.v7/v7/storage"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/types"
)

// New will create a new kodo service.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	const errorMessage = "kodo New: %w"

	srv := &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	mac := qbox.NewMac(cred[0], cred[1])
	cfg := &qs.Config{}
	srv.service = qs.NewBucketManager(mac, cfg)

	store, err := srv.newStorage(pairs...)
	if err != nil && errors.Is(err, types.ErrPairRequired) {
		return srv, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}
	return srv, store, nil
}

func convertUnixTimestampToTime(v int64) time.Time {
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(v, 0)
}

const (
	// ref: https://developer.qiniu.com/kodo/api/3710/chtype
	storageClassStandard   = 0
	storageClassStandardIA = 1
	storageClassArchive    = 2
)

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (int, error) {
	switch in {
	case storageclass.Hot:
		return storageClassStandard, nil
	case storageclass.Warm:
		return storageClassStandardIA, nil
	case storageclass.Cold:
		return storageClassArchive, nil
	default:
		return 0, types.ErrStorageClassNotSupported
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in int) (storageclass.Type, error) {
	switch in {
	case 0:
		return storageclass.Hot, nil
	case 1:
		return storageclass.Warm, nil
	case 2:
		return storageclass.Cold, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}
