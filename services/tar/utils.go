package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	"go.beyondstorage.io/endpoint"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

// Storage is the example client.
type Storage struct {
	path string
	f    *os.File
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
	return newStorager(pairs...)
}

func newStorager(pairs ...types.Pair) (store types.Storager, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
		}
	}()

	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return nil, err
	}

	ep, err := endpoint.Parse(opt.Endpoint)
	if err != nil {
		return
	}
	if ep.Protocol() != endpoint.ProtocolFile {
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(opt.Endpoint)}
	}

	f, err := os.Open(ep.File())
	if err != nil {
		return
	}

	s := &Storage{
		path: opt.Endpoint,
		f:    f,
		r:    tar.NewReader(f),

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
		s.objectsOffset[o.Path], err = s.f.Seek(0, io.SeekCurrent)
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
