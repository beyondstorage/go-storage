package coreutils

import (
	"errors"
	"fmt"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/config"
	"github.com/Xuanwo/storage/pkg/namespace"
	"github.com/Xuanwo/storage/services/azblob"
	"github.com/Xuanwo/storage/services/fs"
	"github.com/Xuanwo/storage/services/gcs"
	"github.com/Xuanwo/storage/services/kodo"
	"github.com/Xuanwo/storage/services/oss"
	"github.com/Xuanwo/storage/services/qingstor"
	"github.com/Xuanwo/storage/services/s3"
	"github.com/Xuanwo/storage/types"
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

type openFunc func(ns string, opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error)

var opener = map[string]openFunc{
	azblob.Type:   openAzblob,
	fs.Type:       openFs,
	gcs.Type:      openGCS,
	kodo.Type:     openKodo,
	oss.Type:      openOSS,
	qingstor.Type: openQingStor,
	s3.Type:       openS3,
}

// Open will parse config string and return valid Servicer and Storager.
//
// Depends on config string's service type, Servicer could be nil.
// Depends on config string's content, Storager could be nil if namespace not given.
func Open(cfg string) (srv storage.Servicer, store storage.Storager, err error) {
	errorMessage := "coreutils Open [%s]: <%w>"

	t, ns, opt, err := config.Parse(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, cfg, err)
	}

	fn, ok := opener[t]
	if !ok {
		err = fmt.Errorf(errorMessage, cfg, ErrServiceNotSupported)
		return nil, nil, err
	}
	srv, store, err = fn(ns, opt...)
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

func openObjectStorage(srv storage.Servicer, ns string) (store storage.Storager, err error) {
	name, prefix := namespace.ParseObjectStorage(ns)
	// name == "" means no bucket name input, return nil directly.
	if name == "" {
		return
	}
	store, err = srv.Get(name)
	if err != nil {
		return
	}
	err = store.Init(pairs.WithWorkDir(prefix))
	if err != nil {
		return
	}
	return
}

func openAzblob(ns string, opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error) {
	srv, err = azblob.New(opt...)
	if err != nil {
		return
	}
	store, err = openObjectStorage(srv, ns)
	return
}

func openFs(ns string, opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error) {
	store = fs.New()
	path := fs.ParseNamespace(ns)
	err = store.Init(pairs.WithWorkDir(path))
	if err != nil {
		return
	}
	return
}

func openGCS(ns string, opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error) {
	srv, err = gcs.New(opt...)
	if err != nil {
		return
	}
	store, err = openObjectStorage(srv, ns)
	return
}

func openKodo(ns string, opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error) {
	srv, err = kodo.New(opt...)
	if err != nil {
		return
	}
	store, err = openObjectStorage(srv, ns)
	return
}

func openOSS(ns string, opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error) {
	srv, err = oss.New(opt...)
	if err != nil {
		return
	}
	store, err = openObjectStorage(srv, ns)
	return
}

func openQingStor(ns string, opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error) {
	srv, err = qingstor.New(opt...)
	if err != nil {
		return
	}
	store, err = openObjectStorage(srv, ns)
	return
}

func openS3(ns string, opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error) {
	srv, err = s3.New(opt...)
	if err != nil {
		return
	}
	store, err = openObjectStorage(srv, ns)
	return
}
