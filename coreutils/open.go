package coreutils

import (
	"errors"
	"fmt"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/config"
	"github.com/Xuanwo/storage/services/fs"
	"github.com/Xuanwo/storage/services/qingstor"
	"github.com/Xuanwo/storage/types/pairs"
)

var (
	// ErrNotSupportedServiceType will be return when service not supported.
	ErrNotSupportedServiceType = errors.New("not_supported_service_type")
)

// Open will parse config string and return valid Servicer and Storager.
//
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
	default:
		err = fmt.Errorf(errorMessage, cfg, ErrNotSupportedServiceType)
		return nil, nil, err
	}
}
