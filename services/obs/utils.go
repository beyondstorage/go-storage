package obs

import (
	"fmt"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"

	"github.com/beyondstorage/go-storage/credential"
	"github.com/beyondstorage/go-storage/endpoint"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/services"
	"github.com/beyondstorage/go-storage/v5/types"
)

type Service struct {
	f Factory

	service *obs.ObsClient

	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	types.UnimplementedServicer
}

func (s *Service) String() string {
	return fmt.Sprintf("Servicer obs")
}

// Storage is the obs client.
type Storage struct {
	f Factory

	client *obs.ObsClient

	bucket  string
	workDir string

	defaultPairs types.DefaultStoragePairs
	features     types.StorageFeatures

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager obs {Name: %s, WorkDir: %s", s.bucket, s.workDir)
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

	var url string
	switch ep.Protocol() {
	case endpoint.ProtocolHTTP:
		url, _, _ = ep.HTTP()
	case endpoint.ProtocolHTTPS:
		url, _, _ = ep.HTTPS()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(f.Endpoint)}
	}

	srv.service, err = obs.New(ak, sk, url)
	if err != nil {
		return nil, err
	}

	return
}

func (f *Factory) newStorage() (store *Storage, err error) {
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

const (
	// writeSizeMaximum is the maximum size for write operation, 5GB.
	// ref: https://support.huaweicloud.com/sdk-go-devg-obs/obs_23_0402.html
	writeSizeMaximum = 5 * 1024 * 1024 * 1024
)

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

// formatError converts errors returned by SDK into errors defined in go-storage and go-service-*.
// The original error SHOULD NOT be wrapped.
func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	e, ok := err.(obs.ObsError)
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

func (s *Storage) formatFileObject(v obs.Content) (o *types.Object, err error) {
	o = s.newObject(false)
	o.ID = v.Key
	o.Path = s.getRelPath(v.Key)
	o.Mode |= types.ModeRead

	o.SetContentLength(v.Size)
	o.SetLastModified(v.LastModified)

	if v.ETag != "" {
		o.SetEtag(v.ETag)
	}

	var sm ObjectSystemMetadata
	if value := v.StorageClass; value != "" {
		sm.StorageClass = string(value)
	}
	o.SetSystemMetadata(sm)

	return
}

func (s *Storage) newObject(done bool) *types.Object {
	return types.NewObject(s, done)
}
