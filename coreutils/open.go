package coreutils

import (
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
)

// Open will parse config string and return valid Servicer and Storager.
//
// Depends on config string's service type, Servicer could be nil.
// Depends on config string's content, Storager could be nil if namespace not given.
func Open(t string, opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error) {
	defer func() {
		if err != nil {
			err = &OpenError{err, t, opt}
		}
	}()

	fn, ok := openFuncMap[t]
	if !ok {
		return nil, nil, ErrServicerNotImplemented
	}
	srv, store, err = fn(opt...)
	if err != nil {
		return
	}
	return
}

// OpenServicer will open a servicer from config string.
func OpenServicer(t string, opt ...*types.Pair) (srv storage.Servicer, err error) {
	defer func() {
		if err != nil {
			err = &OpenError{err, t, opt}
		}
	}()

	fn, ok := openServicerFuncMap[t]
	if !ok {
		return nil, ErrServicerNotImplemented
	}
	srv, err = fn(opt...)
	if err != nil {
		return
	}
	return
}

// OpenStorager will open a storager from config string.
func OpenStorager(t string, opt ...*types.Pair) (store storage.Storager, err error) {
	defer func() {
		if err != nil {
			err = &OpenError{err, t, opt}
		}
	}()

	fn, ok := openStoragerFuncMap[t]
	if !ok {
		return nil, ErrStoragerNotImplemented
	}
	store, err = fn(opt...)
	if err != nil {
		return
	}
	return
}
