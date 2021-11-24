package memory

import (
	"fmt"
	"path"
	"strings"

	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

// Service is the memory config.
// It is not usable, only for generate code
type Service struct {
	f Factory

	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	types.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer memory")
}

// NewServicer is not usable, only for generate code
func NewServicer(pairs ...types.Pair) (types.Servicer, error) {
	f := Factory{}
	err := f.WithPairs(pairs...)
	if err != nil {
		return nil, err
	}
	return f.NewServicer()
}

// newService is not usable, only for generate code
func (f *Factory) newService() (srv *Service, err error) {
	srv = &Service{}
	return
}

// Storage is the example client.
type Storage struct {
	f Factory

	defaultPairs types.DefaultStoragePairs
	features     types.StorageFeatures

	workDir string
	root    *object

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return "memory"
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

func (f *Factory) newStorage() (st *Storage, err error) {
	root := newObject("", nil, types.ModeDir)
	root.parent = root

	return &Storage{
		f:        *f,
		features: f.storageFeatures(),
		root:     root,
		workDir:  "/",
	}, nil
}

// formatError converts errors returned by SDK into errors defined in go-storage and go-service-*.
// The original error SHOULD NOT be wrapped.
func (s *Storage) formatError(op string, err error, path ...string) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	return err
}

func (s *Storage) absPath(p string) string {
	p = path.Clean(p)
	if path.IsAbs(p) {
		return p
	}

	return path.Join(s.workDir, p)
}

func (s *Storage) relPath(p string) string {
	return strings.TrimPrefix(p, s.workDir)
}
