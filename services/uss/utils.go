package uss

import (
	"fmt"
	"strings"

	"github.com/upyun/go-sdk/upyun"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/services"
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

// ref: https://help.upyun.com/knowledge-base/errno/
func formatError(err error) error {
	fn := func(s string) bool {
		return strings.Contains(err.Error(), `"code": `+s)
	}

	switch {
	case fn("40400001"):
		// 40400001:	file or directory not found
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case fn("40100017"), fn("40100019"), fn("40300011"):
		// 40100017: user need permission
		// 40100019: account forbidden
		// 40300011: has no permission to delete
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
	}
}
