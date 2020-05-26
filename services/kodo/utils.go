package kodo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Xuanwo/storage/types/info"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	qc "github.com/qiniu/api.v7/v7/client"
	qs "github.com/qiniu/api.v7/v7/storage"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/httpclient"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// Service is the kodo config.
type Service struct {
	service *qs.BucketManager
}

// String implements Service.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer kodo")
}

// Storage is the gcs service client.
type Storage struct {
	bucket    *qs.BucketManager
	domain    string
	putPolicy qs.PutPolicy // kodo need PutPolicy to generate upload token.

	name    string
	workDir string
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager kodo {Name: %s, WorkDir: %s}",
		s.name, s.workDir,
	)
}

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

	opt, err := parsePairServiceNew(pairs)
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
	srv.service.Client.Client = httpclient.New(opt.HTTPClientOptions)
	return
}

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

// All available storage classes are listed here.
const (
	// ref: https://developer.qiniu.com/kodo/api/3710/chtype
	StorageClassStandard   = 0
	StorageClassStandardIA = 1
	StorageClassArchive    = 2
)

// ref: https://developer.qiniu.com/kodo/api/3928/error-responses
func formatError(err error) error {
	e, ok := err.(*qc.ErrorInfo)
	if !ok {
		return err
	}

	// error code returned by kodo looks like http status code, but it's not.
	// kodo could return 6xx or 7xx for their costumed errors, so we use untyped int directly.
	switch e.Code {
	case 612:
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case 403:
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
	}
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...*types.Pair) (store *Storage, err error) {
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return nil, err
	}

	store = &Storage{
		bucket: s.service,
		domain: opt.Endpoint.Value().String(),
		putPolicy: qs.PutPolicy{
			Scope: opt.Name,
		},

		name:    opt.Name,
		workDir: "/",
	}

	if opt.HasWorkDir {
		store.workDir = opt.WorkDir
	}
	return store, nil
}

func (s *Service) formatError(op string, err error, name string) error {
	if err == nil {
		return nil
	}

	return &services.ServiceError{
		Op:       op,
		Err:      formatError(err),
		Servicer: s,
		Name:     name,
	}
}

// getAbsPath will calculate object storage's abs path
func (s *Storage) getAbsPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return prefix + path
}

// getRelPath will get object storage's rel path.
func (s *Storage) getRelPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return strings.TrimPrefix(path, prefix)
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	if err == nil {
		return nil
	}

	return &services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}

func (s *Storage) formatFileObject(v qs.ListItem) (o *types.Object, err error) {
	o = &types.Object{
		ID:         v.Key,
		Name:       s.getRelPath(v.Key),
		Type:       types.ObjectTypeFile,
		Size:       v.Fsize,
		UpdatedAt:  convertUnixTimestampToTime(v.PutTime),
		ObjectMeta: info.NewObjectMeta(),
	}

	if v.MimeType != "" {
		o.SetContentType(v.MimeType)
	}
	if v.Hash != "" {
		o.SetETag(v.Hash)
	}

	setStorageClass(o.ObjectMeta, v.Type)
	return
}
