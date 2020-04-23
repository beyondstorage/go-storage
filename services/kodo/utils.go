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
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// New will create both Servicer and Storager.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	return newServicerAndStorager(pairs...)
}

// NewServicer will create Servicer only.
func NewServicer(pairs ...*types.Pair) (storage.Servicer, error) {
	return newServicer(pairs...)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...*types.Pair) (storage.Storager, error) {
	_, store, err := newServicerAndStorager(pairs...)
	return store, err
}

func newServicer(pairs ...*types.Pair) (srv *Service, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Type: Type, Err: err, Pairs: pairs}
		}
	}()

	srv = &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, err
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	mac := qbox.NewMac(cred[0], cred[1])
	cfg := &qs.Config{}
	srv.service = qs.NewBucketManager(mac, cfg)

	return
}

func newServicerAndStorager(pairs ...*types.Pair) (srv *Service, store *Storage, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Type: Type, Err: err, Pairs: pairs}
		}
	}()

	srv, err = newServicer(pairs...)
	if err != nil {
		return
	}

	store, err = srv.newStorage(pairs...)
	if err != nil {
		if e := services.NewPairRequiredError(); errors.As(err, &e) {
			return srv, nil, nil
		}
		return nil, nil, err
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
		return 0, services.NewPairUnsupportedError(ps.WithStorageClass(in))
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in int) storageclass.Type {
	switch in {
	case 0:
		return storageclass.Hot
	case 1:
		return storageclass.Warm
	case 2:
		return storageclass.Cold
	default:
		return ""
	}
}

// ref: https://developer.qiniu.com/kodo/api/3928/error-responses
func formatError(err error) error {
	e, ok := err.(*qs.ErrorInfo)
	if !ok {
		return err
	}

	// error code returned by kodo looks like http status code, but it's not.
	// kodo could return 6xx or 7xx for their costumed errors, so we use untyped int directly.
	switch e.Errno {
	case 0:
		switch e.Code {
		case 404:
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		default:
			return err
		}
	case 404:
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case 403:
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
	}
}
