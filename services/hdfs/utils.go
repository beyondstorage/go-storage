package hdfs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/colinmarc/hdfs/v2"

	"go.beyondstorage.io/endpoint"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

// Service is the hdfs config.
// It is not usable, only for generate code
type Service struct {
	f Factory

	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	types.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer hdfs")
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

	hdfs *hdfs.Client

	defaultPairs types.DefaultStoragePairs
	features     types.StorageFeatures

	workDir string

	types.UnimplementedStorager
	types.UnimplementedDirer
	types.UnimplementedMover
	types.UnimplementedAppender
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager hdfs {WorkDir: %s}", s.workDir)
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

func (f *Factory) newStorage() (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err)}
		}
	}()

	store = &Storage{
		f:        *f,
		features: f.storageFeatures(),
		workDir:  "/",
	}
	if f.WorkDir != "" {
		store.workDir = f.WorkDir
	}

	ep, err := endpoint.Parse(f.Endpoint)
	if err != nil {
		return nil, err
	}

	var addr string

	switch ep.Protocol() {
	case endpoint.ProtocolTCP:
		addr, _, _ = ep.TCP()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(f.Endpoint)}
	}
	store.hdfs, err = hdfs.New(addr)
	if err != nil {
		return nil, errors.New("hdfs address is not exist")
	}

	return store, nil
}

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	switch {
	case errors.Is(err, os.ErrNotExist):
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case errors.Is(err, os.ErrPermission):
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
	}
}

func (s *Storage) getAbsPath(fp string) string {
	if filepath.IsAbs(fp) {
		return fp
	}
	return path.Join(s.workDir, fp)
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

func (s *Storage) newObject(done bool) *types.Object {
	return types.NewObject(s, done)
}
