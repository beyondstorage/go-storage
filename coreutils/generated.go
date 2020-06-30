// Code generated by go generate via internal/cmd/coreutils; DO NOT EDIT.
package coreutils

import (
	"github.com/aos-dev/go-storage/v2"
	"github.com/aos-dev/go-storage/v2/services/azblob"
	"github.com/aos-dev/go-storage/v2/services/cos"
	"github.com/aos-dev/go-storage/v2/services/dropbox"
	"github.com/aos-dev/go-storage/v2/services/fs"
	"github.com/aos-dev/go-storage/v2/services/gcs"
	"github.com/aos-dev/go-storage/v2/services/kodo"
	"github.com/aos-dev/go-storage/v2/services/oss"
	"github.com/aos-dev/go-storage/v2/services/qingstor"
	"github.com/aos-dev/go-storage/v2/services/s3"
	"github.com/aos-dev/go-storage/v2/services/uss"
	"github.com/aos-dev/go-storage/v2/types"
)

type openFunc func(opt ...*types.Pair) (srv storage.Servicer, store storage.Storager, err error)

var openFuncMap = map[string]openFunc{
	azblob.Type:   azblob.New,
	cos.Type:      cos.New,
	gcs.Type:      gcs.New,
	kodo.Type:     kodo.New,
	oss.Type:      oss.New,
	qingstor.Type: qingstor.New,
	s3.Type:       s3.New,
}

type openServicerFunc func(opt ...*types.Pair) (srv storage.Servicer, err error)

var openServicerFuncMap = map[string]openServicerFunc{
	azblob.Type:   azblob.NewServicer,
	cos.Type:      cos.NewServicer,
	gcs.Type:      gcs.NewServicer,
	kodo.Type:     kodo.NewServicer,
	oss.Type:      oss.NewServicer,
	qingstor.Type: qingstor.NewServicer,
	s3.Type:       s3.NewServicer,
}

type openStoragerFunc func(opt ...*types.Pair) (store storage.Storager, err error)

var openStoragerFuncMap = map[string]openStoragerFunc{
	azblob.Type:   azblob.NewStorager,
	cos.Type:      cos.NewStorager,
	dropbox.Type:  dropbox.NewStorager,
	fs.Type:       fs.NewStorager,
	gcs.Type:      gcs.NewStorager,
	kodo.Type:     kodo.NewStorager,
	oss.Type:      oss.NewStorager,
	qingstor.Type: qingstor.NewStorager,
	s3.Type:       s3.NewStorager,
	uss.Type:      uss.NewStorager,
}
