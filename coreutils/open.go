package coreutils

import (
	"errors"
	"fmt"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/config"
	"github.com/Xuanwo/storage/services/azblob"
	"github.com/Xuanwo/storage/services/cos"
	"github.com/Xuanwo/storage/services/dropbox"
	"github.com/Xuanwo/storage/services/fs"
	"github.com/Xuanwo/storage/services/gcs"
	"github.com/Xuanwo/storage/services/kodo"
	"github.com/Xuanwo/storage/services/oss"
	"github.com/Xuanwo/storage/services/qingstor"
	"github.com/Xuanwo/storage/services/s3"
	"github.com/Xuanwo/storage/services/uss"
	"github.com/Xuanwo/storage/types"
)

var (
	// ErrServiceNotSupported will return when service not supported.
	ErrServiceNotSupported = errors.New("service not supported")
	// ErrServiceNotImplemented will return when service doesn't implement Servicer.
	ErrServiceNotImplemented = errors.New("service not implemented")
	// ErrServiceNamespaceNotGiven will return when service namespace not given.
	ErrServiceNamespaceNotGiven = errors.New("service namespace not given")
)

type openFunc func(opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error)

var opener = map[string]openFunc{
	azblob.Type:   azblob.New,
	cos.Type:      cos.New,
	dropbox.Type:  dropbox.New,
	fs.Type:       fs.New,
	gcs.Type:      gcs.New,
	kodo.Type:     kodo.New,
	oss.Type:      nil,
	qingstor.Type: nil,
	s3.Type:       nil,
	uss.Type:      nil,
}

// Open will parse config string and return valid Servicer and Storager.
//
// Depends on config string's service type, Servicer could be nil.
// Depends on config string's content, Storager could be nil if namespace not given.
func Open(cfg string) (srv storage.Servicer, store storage.Storager, err error) {
	errorMessage := "coreutils Open [%s]: <%w>"

	t, opt, err := config.Parse(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, cfg, err)
	}

	fn, ok := opener[t]
	if !ok {
		err = fmt.Errorf(errorMessage, cfg, ErrServiceNotSupported)
		return nil, nil, err
	}
	srv, store, err = fn(opt...)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, cfg, err)
	}
	return
}

// OpenServicer will open a servicer from config string.
func OpenServicer(cfg string) (srv storage.Servicer, err error) {
	errorMessage := "coreutils OpenServicer [%s]: <%w>"

	srv, _, err = Open(cfg)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, cfg, err)
	}
	if srv == nil {
		return nil, fmt.Errorf(errorMessage, cfg, ErrServiceNotImplemented)
	}
	return
}

// OpenStorager will open a storager from config string.
func OpenStorager(cfg string) (store storage.Storager, err error) {
	errorMessage := "coreutils OpenStorager [%s]: <%w>"

	_, store, err = Open(cfg)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, cfg, err)
	}
	if store == nil {
		return nil, fmt.Errorf(errorMessage, cfg, ErrServiceNamespaceNotGiven)
	}
	return
}
