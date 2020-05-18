package gcs

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	gs "cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/httpclient"
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
			err = &services.InitError{Op: services.OpNewServicer, Type: Type, Err: err, Pairs: pairs}
		}
	}()

	srv = &Service{}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, err
	}

	hc := httpclient.New(opt.HTTPClientOptions)

	var credJSON []byte

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	switch credProtocol {
	case credential.ProtocolFile:
		credJSON, err = ioutil.ReadFile(cred[0])
		if err != nil {
			return nil, err
		}
	case credential.ProtocolBase64:
		credJSON, err = base64.StdEncoding.DecodeString(cred[0])
		if err != nil {
			return nil, err
		}
	default:
		return nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	// Loading token source from binary data.
	creds, err := google.CredentialsFromJSON(opt.Context, credJSON, gs.ScopeFullControl)
	if err != nil {
		return nil, err
	}
	ot := &oauth2.Transport{
		Source: creds.TokenSource,
		Base:   hc.Transport,
	}
	hc.Transport = ot

	client, err := gs.NewClient(opt.Context, option.WithHTTPClient(hc))
	if err != nil {
		return nil, err
	}

	srv.service = client
	srv.projectID = opt.Project

	return
}

// New will create a new aliyun oss service.
func newServicerAndStorager(pairs ...*types.Pair) (srv *Service, store *Storage, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Op: services.OpNewStorager, Type: Type, Err: err, Pairs: pairs}
		}
	}()

	srv, err = newServicer(pairs...)
	if err != nil {
		return
	}

	store, err = srv.newStorage(pairs...)
	if err != nil {
		return nil, nil, err
	}
	return srv, store, nil
}

// All available storage classes are listed here.
const (
	StorageClassStandard = "STANDARD"
	StorageClassNearLine = "NEARLINE"
	StorageClassColdLine = "COLDLINE"
	StorageClassArchive  = "ARCHIVE"
)

// ref: https://cloud.google.com/storage/docs/json_api/v1/status-codes
func formatError(err error) error {
	// gcs sdk could return explicit error, we should handle them.
	if errors.Is(err, gs.ErrObjectNotExist) {
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	}

	e, ok := err.(*googleapi.Error)
	if !ok {
		return err
	}

	switch e.Code {
	case http.StatusNotFound:
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case http.StatusForbidden:
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
	}
}
