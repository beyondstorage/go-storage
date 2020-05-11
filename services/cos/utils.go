package cos

import (
	"fmt"

	"github.com/Xuanwo/storage/pkg/httpclient"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// New will create both Servicer and Storager.
func New(pairs ...*types.Pair) (_ storage.Servicer, _ storage.Storager, err error) {
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

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	httpClient := httpclient.New(opt.HTTPClientOptions)
	httpClient.Transport = &cos.AuthorizationTransport{
		Transport: httpClient.Transport,
		SecretID:  cred[0],
		SecretKey: cred[1],
	}

	srv.client = httpClient
	srv.service = cos.NewClient(nil, srv.client)
	return
}

// newServicerAndStorager will create a new Tencent oss service.
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
		return
	}
	return
}

// All available storage classes are listed here.
const (
	// ref: https://cloud.tencent.com/document/product/436/7745
	storageClassHeader = "x-cos-storage-class"

	StorageClassStandard   = "STANDARD"
	StorageClassStandardIA = "STANDARD_IA"
	StorageClassArchive    = "ARCHIVE"
)

// ref: https://www.qcloud.com/document/product/436/7730
func formatError(err error) error {
	// Handle errors returned by cos.
	e, ok := err.(*cos.ErrorResponse)
	if !ok {
		return err
	}

	switch e.Code {
	case "":
		switch e.Response.StatusCode {
		case 404:
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		default:
			return err
		}
	case "NoSuchKey":
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case "AccessDenied":
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
	}
}
