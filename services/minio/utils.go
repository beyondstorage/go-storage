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
	service *minio.Client

	defaultPairs DefaultServicePairs
	features     ServiceFeatures

	types.UnimplementedServicer
}

func (s *Service) String() string {
	return fmt.Sprintf("Servicer minio")
}

// Storage is the example client.
type Storage struct {
	client *minio.Client

	bucket  string
	workDir string

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	types.UnimplementedStorager
	types.UnimplementedCopier
	types.UnimplementedReacher
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager minio {Name: %s, WorkDir: %s}",
		s.bucket, s.workDir,
	)
}

func New(pairs ...types.Pair) (types.Servicer, types.Storager, error) {
	return newServicerAndStorager(pairs...)
}

func NewServicer(pairs ...types.Pair) (types.Servicer, error) {
	return newServicer(pairs...)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...types.Pair) (types.Storager, error) {
	_, store, err := newServicerAndStorager(pairs...)
	return store, err
}

func newServicer(pairs ...types.Pair) (srv *Service, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_servicer", Type: Type, Err: formatError(err), Pairs: pairs}
		}
	}()

	srv = &Service{}

	opt, err := parsePairServiceNew(pairs)
	if err != nil {
		return nil, err
	}

	cp, err := credential.Parse(opt.Credential)
	if err != nil {
		return nil, err
	}
	if cp.Protocol() != credential.ProtocolHmac {
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(opt.Credential)}
	}
	ak, sk := cp.Hmac()

	ep, err := endpoint.Parse(opt.Endpoint)
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
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(opt.Endpoint)}
	}
	url := fmt.Sprintf("%s:%d", host, port)

	srv.service, err = minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(ak, sk, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	if opt.HasDefaultServicePairs {
		srv.defaultPairs = opt.DefaultServicePairs
	}
	if opt.HasServiceFeatures {
		srv.features = opt.ServiceFeatures
	}

	return
}

func newServicerAndStorager(pairs ...types.Pair) (srv *Service, store *Storage, err error) {
	srv, err = newServicer(pairs...)
	if err != nil {
		return
	}

	store, err = srv.newStorage(pairs...)
	if err != nil {
		err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
		return nil, nil, err
	}
	return srv, store, nil
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

func (s *Service) newStorage(pairs ...types.Pair) (st *Storage, err error) {
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return nil, err
	}

	store := &Storage{
		client:  s.service,
		bucket:  opt.Name,
		workDir: "/",
	}

	if opt.HasWorkDir {
		if !strings.HasSuffix(opt.WorkDir, "/") {
			opt.WorkDir += "/"
		}
		store.workDir = opt.WorkDir
	}
	if opt.HasDefaultStoragePairs {
		store.defaultPairs = opt.DefaultStoragePairs
	}
	if opt.HasStorageFeatures {
		store.features = opt.StorageFeatures
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
