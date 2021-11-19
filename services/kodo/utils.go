package kodo

import (
	"fmt"
	"net/http"

	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	qc "github.com/qiniu/go-sdk/v7/client"
	qs "github.com/qiniu/go-sdk/v7/storage"

	"go.beyondstorage.io/credential"
	"go.beyondstorage.io/endpoint"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	typ "go.beyondstorage.io/v5/types"
)

// Service is the kodo config.
type Service struct {
	f Factory

	service *qs.BucketManager

	defaultPairs typ.DefaultServicePairs
	features     typ.ServiceFeatures

	typ.UnimplementedServicer
}

// String implements Service.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer kodo")
}

// Storage is the gcs service client.
type Storage struct {
	f Factory

	bucket    *qs.BucketManager
	domain    string
	putPolicy qs.PutPolicy // kodo need PutPolicy to generate upload token.

	name    string
	workDir string

	defaultPairs typ.DefaultStoragePairs
	features     typ.StorageFeatures

	typ.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager kodo {Name: %s, WorkDir: %s}",
		s.name, s.workDir,
	)
}

// New will create both Servicer and Storager.
func New(pairs ...typ.Pair) (typ.Servicer, typ.Storager, error) {
	f := Factory{}
	err := f.WithPairs(pairs...)
	if err != nil {
		return nil, nil, err
	}
	srv, err := f.NewServicer()
	if err != nil {
		return nil, nil, err
	}
	sto, err := f.NewStorager()
	if err != nil {
		return nil, nil, err
	}
	return srv, sto, nil
}

// NewServicer will create Servicer only.
func NewServicer(pairs ...typ.Pair) (typ.Servicer, error) {
	f := Factory{}
	err := f.WithPairs(pairs...)
	if err != nil {
		return nil, err
	}
	return f.NewServicer()
}

// NewStorager will create Storager only.
func NewStorager(pairs ...typ.Pair) (typ.Storager, error) {
	f := Factory{}
	err := f.WithPairs(pairs...)
	if err != nil {
		return nil, err
	}
	return f.newStorage()
}

func (f *Factory) newService() (srv *Service, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_servicer", Type: Type, Err: formatError(err)}
		}
	}()

	srv = &Service{
		f:        *f,
		features: f.serviceFeatures(),
	}

	cp, err := credential.Parse(f.Credential)
	if err != nil {
		return nil, err
	}
	if cp.Protocol() != credential.ProtocolHmac {
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(f.Credential)}
	}
	ak, sk := cp.Hmac()

	mac := qbox.NewMac(ak, sk)
	cfg := &qs.Config{}
	srv.service = qs.NewBucketManager(mac, cfg)
	srv.service.Client.Client = &http.Client{}

	return
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
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	e, ok := err.(*qc.ErrorInfo)
	if !ok {
		return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
	}

	// error code returned by kodo looks like http status code, but it's not.
	// kodo could return 6xx or 7xx for their costumed errors.
	switch e.Code {
	case responseCodeResourceNotExist:
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case responseCodePermissionDenied:
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
	}
}

// Error code returned by kodo.
//
// ref: https://developer.qiniu.com/kodo/api/3928/error-responses
const (
	// responseCodeResourceNotExist is an error code that is returned if insufficient permissions and access denied.
	responseCodePermissionDenied = 403
	// responseCodeResourceNotExist is an error code that is returned if the specified resource does not exist or has been deleted.
	responseCodeResourceNotExist = 612
)

func checkError(err error, code int) bool {
	e, ok := err.(*qc.ErrorInfo)
	if !ok {
		return false
	}

	return e.Code == code
}

// newStorage will create a new client.
func (f *Factory) newStorage() (store *Storage, err error) {
	s, err := f.newService()
	if err != nil {
		return
	}

	ep, err := endpoint.Parse(f.Endpoint)
	if err != nil {
		return nil, err
	}

	var url string
	switch ep.Protocol() {
	case endpoint.ProtocolHTTPS:
		url, _, _ = ep.HTTPS()
	case endpoint.ProtocolHTTP:
		url, _, _ = ep.HTTP()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(f.Endpoint)}
	}

	store = &Storage{
		f:        *f,
		features: f.storageFeatures(),
		bucket:   s.service,
		domain:   url,
		putPolicy: qs.PutPolicy{
			Scope: f.Name,
		},

		name:    f.Name,
		workDir: "/",
	}

	if f.WorkDir != "" {
		store.workDir = f.WorkDir
	}
	return store, nil
}

func (s *Service) formatError(op string, err error, name string) error {
	if err == nil {
		return nil
	}

	return services.ServiceError{
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

	return services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}

func (s *Storage) formatFileObject(v qs.ListItem) (o *typ.Object, err error) {
	o = s.newObject(false)
	o.ID = v.Key
	o.Path = s.getRelPath(v.Key)
	o.Mode |= typ.ModeRead

	o.SetContentLength(v.Fsize)
	o.SetLastModified(convertUnixTimestampToTime(v.PutTime))

	if v.MimeType != "" {
		o.SetContentType(v.MimeType)
	}
	if v.Hash != "" {
		o.SetEtag(v.Hash)
	}

	var sm ObjectSystemMetadata
	sm.StorageClass = v.Type
	o.SetSystemMetadata(sm)

	return
}

func (s *Storage) newObject(done bool) *typ.Object {
	return typ.NewObject(s, done)
}
