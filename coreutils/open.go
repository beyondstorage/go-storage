package coreutils

import (
	"errors"

	"github.com/Xuanwo/storage"
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
	// TODO: namespace has been removed, this error should be changed.
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
	oss.Type:      oss.New,
	qingstor.Type: qingstor.New,
	s3.Type:       s3.New,
	uss.Type:      uss.New,
}

// Open will parse config string and return valid Servicer and Storager.
//
// Depends on config string's service type, Servicer could be nil.
// Depends on config string's content, Storager could be nil if namespace not given.
func Open(t string, opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error) {
	fn, ok := opener[t]
	if !ok {
		return nil, nil, &OpenError{ErrServiceNotSupported, t, opt}
	}
	srv, store, err = fn(opt...)
	if err != nil {
		return nil, nil, &OpenError{err, t, opt}
	}
	return
}

// OpenServicer will open a servicer from config string.
func OpenServicer(t string, opt ...*types.Pair) (srv storage.Servicer, err error) {
	srv, _, err = Open(t, opt...)
	if err != nil {
		return nil, err
	}
	if srv == nil {
		return nil, &OpenError{ErrServiceNotImplemented, t, opt}
	}
	return
}

// OpenStorager will open a storager from config string.
func OpenStorager(t string, opt ...*types.Pair) (store storage.Storager, err error) {
	_, store, err = Open(t, opt...)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, &OpenError{ErrServiceNamespaceNotGiven, t, opt}
	}
	return
}
