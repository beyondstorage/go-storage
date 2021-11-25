package minio

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"go.beyondstorage.io/credential"
	"go.beyondstorage.io/endpoint"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

// Service is the minio service.
type Service struct {
	f Factory

	service *minio.Client

	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	types.UnimplementedServicer
}

func (s *Service) String() string {
	return fmt.Sprintf("Servicer minio")
}

// Storage is the example client.
type Storage struct {
	f Factory

	client *minio.Client

	bucket  string
	workDir string

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager minio {Name: %s, WorkDir: %s}",
		s.bucket, s.workDir,
	)
}

func New(pairs ...types.Pair) (types.Servicer, types.Storager, error) {
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

func NewServicer(pairs ...types.Pair) (types.Servicer, error) {
	f := Factory{}
	err := f.WithPairs(pairs...)
	if err != nil {
		return nil, err
	}
	return f.NewServicer()
}

// NewStorager will create Storager only.
func NewStorager(pairs ...types.Pair) (types.Storager, error) {
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

	ep, err := endpoint.Parse(f.Endpoint)
	if err != nil {
		return nil, err
	}

	var host string
	var port int
	var secure bool
	switch ep.Protocol() {
	case endpoint.ProtocolHTTP:
		_, host, port = ep.HTTP()
		secure = false
	case endpoint.ProtocolHTTPS:
		_, host, port = ep.HTTPS()
		secure = true
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(f.Endpoint)}
	}
	url := fmt.Sprintf("%s:%d", host, port)

	srv.service, err = minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(ak, sk, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	return
}

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	e, ok := err.(minio.ErrorResponse)
	if ok {
		switch e.Code {
		case "AccessDenied":
			return fmt.Errorf("%w, %v", services.ErrPermissionDenied, err)
		case "NoSuchKey":
			return fmt.Errorf("%w, %v", services.ErrObjectNotExist, err)
		case "InternalError":
			return fmt.Errorf("%w, %v", services.ErrServiceInternal, err)
		}

		switch e.StatusCode {
		case http.StatusTooManyRequests:
			return fmt.Errorf("%w, %v", services.ErrRequestThrottled, err)
		case http.StatusServiceUnavailable:
			return fmt.Errorf("%w, %v", services.ErrRequestThrottled, err)
		}
	}

	return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
}

func (f *Factory) newStorage() (st *Storage, err error) {
	s, err := f.newService()

	store := &Storage{
		f:        *f,
		features: f.storageFeatures(),

		client:  s.service,
		bucket:  f.Name,
		workDir: "/",
	}

	if f.WorkDir != "" {
		if !strings.HasSuffix(f.WorkDir, "/") {
			f.WorkDir += "/"
		}
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

// getAbsPath will calculate object storage's abs path
func (s *Storage) getAbsPath(path string) string {
	if strings.HasPrefix(path, s.workDir) {
		return strings.TrimPrefix(path, "/")
	}
	prefix := strings.TrimPrefix(s.workDir, "/")
	return prefix + path
}

// getRelPath will get object storage's rel path.
func (s *Storage) getRelPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return strings.TrimPrefix(path, prefix)
}

func (s *Storage) formatFileObject(v minio.ObjectInfo) (o *types.Object, err error) {
	o = s.newObject(true)
	if v.ETag == "" {
		o.Mode |= types.ModeDir
	} else {
		o.Mode |= types.ModeRead
	}

	o.SetID(v.Key)
	o.SetPath(s.getRelPath(v.Key))
	o.SetEtag(v.ETag)
	o.SetContentLength(v.Size)
	o.SetContentType(v.ContentType)
	o.SetLastModified(v.LastModified)
	o.SetUserMetadata(v.UserMetadata)
	o.SetSystemMetadata(ObjectSystemMetadata{
		StorageClass: v.StorageClass,
	})

	return
}

func (s *Storage) newObject(done bool) *types.Object {
	return types.NewObject(s, done)
}
