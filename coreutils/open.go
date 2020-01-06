package coreutils

import (
	"errors"
	"fmt"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/config"
	"github.com/Xuanwo/storage/services/azblob"
	"github.com/Xuanwo/storage/services/fs"
	"github.com/Xuanwo/storage/services/oss"
	"github.com/Xuanwo/storage/services/qingstor"
	"github.com/Xuanwo/storage/services/s3"
	"github.com/Xuanwo/storage/types/pairs"
)

var (
	// ErrServiceNotSupported will return when service not supported.
	ErrServiceNotSupported = errors.New("service not supported")
	// ErrServiceNotImplemented will return when service doesn't implement Servicer.
	ErrServiceNotImplemented = errors.New("service not implemented")
	// ErrServiceNamespaceNotGiven will return when service namespace not given.
	ErrServiceNamespaceNotGiven = errors.New("service namespace not given")
)

// Open will parse config string and return valid Servicer and Storager.
//
// Depends on config string's service type, Servicer could be nil.
// Depends on config string's content, Storager could be nil if namespace not given.
func Open(cfg string) (srv storage.Servicer, store storage.Storager, err error) {
	errorMessage := "coreutils Open [%s]: <%w>"

	t, namespace, opt, err := config.Parse(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, cfg, err)
	}

	switch t {
	case fs.Type:
		store = fs.New()
		path := fs.ParseNamespace(namespace)
		err = store.Init(pairs.WithWorkDir(path))
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		return
	case qingstor.Type:
		srv, err = qingstor.New(opt...)
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		name, prefix := qingstor.ParseNamespace(namespace)
		if name == "" {
			return
		}
		store, err = srv.Get(name)
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		err = store.Init(pairs.WithWorkDir(prefix))
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		return
	case s3.Type:
		srv, err = s3.New(opt...)
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		// FIXME: this util function should move to config package.
		name, prefix := qingstor.ParseNamespace(namespace)
		if name == "" {
			return
		}
		store, err = srv.Get(name)
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		err = store.Init(pairs.WithWorkDir(prefix))
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		return
	case oss.Type:
		srv, err = oss.New(opt...)
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		// FIXME: this util function should move to config package.
		name, prefix := qingstor.ParseNamespace(namespace)
		if name == "" {
			return
		}
		store, err = srv.Get(name)
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		err = store.Init(pairs.WithWorkDir(prefix))
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		return
	case azblob.Type:
		srv, err = azblob.New(opt...)
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		// FIXME: this util function should move to config package.
		name, prefix := qingstor.ParseNamespace(namespace)
		if name == "" {
			return
		}
		store, err = srv.Get(name)
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		err = store.Init(pairs.WithWorkDir(prefix))
		if err != nil {
			err = fmt.Errorf(errorMessage, cfg, err)
			return
		}
		return
	default:
		err = fmt.Errorf(errorMessage, cfg, ErrServiceNotSupported)
		return nil, nil, err
	}
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
