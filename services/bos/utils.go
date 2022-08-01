package bos

import (
	"fmt"
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"

	"github.com/beyondstorage/go-storage/credential"
	"github.com/beyondstorage/go-storage/endpoint"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/services"
	"github.com/beyondstorage/go-storage/v5/types"
)

// Service is the bos service
type Service struct {
	f Factory

	service *bos.Client

	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	types.UnimplementedServicer
}

func (s *Service) String() string {
	return fmt.Sprintf("Servicer bos")
}

// Storage is the bos client
type Storage struct {
	f Factory

	client *bos.Client

	bucket  string
	workDir string

	defaultPairs types.DefaultStoragePairs
	features     types.StorageFeatures

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager bos {Name: %s, WorkDir: %s}",
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
			err = services.InitError{
				Op:   "new_servicer",
				Type: Type,
				Err:  formatError(err),
			}
		}
	}()

	srv = &Service{}

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

	var url string
	switch ep.Protocol() {
	case endpoint.ProtocolHTTP:
		url, _, _ = ep.HTTP()
	case endpoint.ProtocolHTTPS:
		url, _, _ = ep.HTTPS()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(f.Endpoint)}
	}

	srv.service, err = bos.NewClient(ak, sk, url)
	if err != nil {
		return nil, err
	}

	return
}

const (
	// writeSizeMaximum is the maximum size for write operation, 5GB.
	// ref: https://cloud.baidu.com/doc/BOS/s/Ikc5nv3wc
	writeSizeMaximum = 5 * 1024 * 1024 * 1024
)

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

func (s *Service) formatError(op string, err error, name string) error {
	if err == nil {
		return nil
	}

	return services.ServiceError{
		Op:       op,
		Err:      err,
		Servicer: s,
		Name:     name,
	}
}

// formatError converts errors returned by SDK into errors defined in go-storage and go-service-*.
// The original error SHOULD NOT be wrapped.
func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	e, ok := err.(*bce.BceServiceError)
	if ok {
		switch e.Code {
		case "AccessDenied":
			return fmt.Errorf("%w, %v", services.ErrPermissionDenied, err)
		case "NoSuchKey":
			return fmt.Errorf("%w, %v", services.ErrObjectNotExist, err)
		default:
			return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
		}
	}

	return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
}

func (f *Factory) newStorage(pairs ...types.Pair) (store *Storage, err error) {
	s, err := f.newService()
	if err != nil {
		return nil, err
	}

	store = &Storage{
		f:        *f,
		features: f.storageFeatures(),
		client:   s.service,
		bucket:   f.Name,
		workDir:  "/",
	}

	if f.WorkDir != "" {
		store.workDir = f.WorkDir
	}

	return
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

func (s *Storage) formatFileObject(v api.ObjectSummaryType) (o *types.Object, err error) {
	o = s.newObject(false)
	o.ID = v.Key
	o.Path = s.getRelPath(v.Key)
	o.Mode |= types.ModeRead

	o.SetContentLength(int64(v.Size))
	// Last-Modified returns a format of :
	// 2009-10-12T17:50:30Z
	// ref:https://cloud.baidu.com/doc/BOS/s/Ekc4epj6m#%E7%A4%BA%E4%BE%8B
	lastModified, err := time.Parse(time.RFC3339, v.LastModified)
	if err != nil {
		return nil, err
	}
	o.SetLastModified(lastModified)

	if v.ETag != "" {
		o.SetEtag(v.ETag)
	}

	var sm ObjectSystemMetadata
	if value := v.StorageClass; value != "" {
		sm.StorageClass = value
	}
	o.SetSystemMetadata(sm)

	return
}

func (s *Storage) newObject(done bool) *types.Object {
	return types.NewObject(s, done)
}
