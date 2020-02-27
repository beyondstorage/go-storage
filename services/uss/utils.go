package uss

import (
	"fmt"

	"github.com/upyun/go-sdk/upyun"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/types"
)

// New will create a new uss service.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	const errorMessage = "uss New: %w"

	store := &Storage{}

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, nil, fmt.Errorf(errorMessage, credential.ErrUnsupportedProtocol)
	}

	cfg := &upyun.UpYunConfig{
		Bucket:   opt.Name,
		Operator: cred[0],
		Password: cred[1],
	}
	store.bucket = upyun.NewUpYun(cfg)
	store.name = opt.Name
	store.workDir = opt.WorkDir
	return nil, store, nil
}
