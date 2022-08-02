package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	"github.com/beyondstorage/go-storage/endpoint"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/services"
	"github.com/beyondstorage/go-storage/v5/types"
)

// Service is the tar config.
// It is not usable, only for generate code
type Service struct {
	f Factory

	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	types.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer tar")
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

	path string
	file *os.File
	r    *tar.Reader

	objects       []*types.Object
	objectsIndex  map[string]uint  // path -> index map.
	objectsOffset map[string]int64 // path -> object offset map.

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager tar {Path: %s}", s.path,
	)
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

func (f *Factory) newStorage() (store types.Storager, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err)}
		}
	}()

	ep, err := endpoint.Parse(f.Endpoint)
	if err != nil {
		return
	}
	if ep.Protocol() != endpoint.ProtocolFile {
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(f.Endpoint)}
	}

	file, err := os.Open(ep.File())
	if err != nil {
		return
	}

	s := &Storage{
		f:        *f,
		features: f.storageFeatures(),
		path:     f.Endpoint,
		file:     file,
		r:        tar.NewReader(file),

		objectsIndex:  make(map[string]uint),
		objectsOffset: make(map[string]int64),
	}
	err = s.parse()
	if err != nil {
		return
	}

	return s, nil
}

// formatError converts errors returned by SDK into errors defined in go-storage and go-service-*.
// The original error SHOULD NOT be wrapped.
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

func (s *Storage) parse() (err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "parse", Type: Type, Err: formatError(err)}
		}
	}()

	index := uint(0)

	for {
		h, err := s.r.Next()
		if err != nil && err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		o, err := s.formatTarHeader(h)
		if err != nil {
			return err
		}

		s.objects = append(s.objects, o)
		s.objectsIndex[o.Path] = index
		s.objectsOffset[o.Path], err = s.file.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		index += 1
	}
}

func (s *Storage) formatTarHeader(h *tar.Header) (o *types.Object, err error) {
	o = s.newObject(true)
	o.Path = h.Name
	o.ID = h.Name

	o.SetContentLength(h.Size)
	return
}

func (s *Storage) newObject(done bool) *types.Object {
	return types.NewObject(s, done)
}

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
}
