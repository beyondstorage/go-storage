package hdfs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"go.beyondstorage.io/endpoint"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"

	"github.com/colinmarc/hdfs/v2"
)

// Storage is the example client.
type Storage struct {
	hdfs *hdfs.Client

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

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
	return newStorager(pairs...)
}

func newStorager(pairs ...types.Pair) (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
		}
	}()

	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return nil, err
	}

	store = &Storage{
		workDir: "/",
	}
	if opt.HasWorkDir {
		store.workDir = opt.WorkDir
	}

	ep, err := endpoint.Parse(opt.Endpoint)
	if err != nil {
		return nil, err
	}

	var addr string

	switch ep.Protocol() {
	case endpoint.ProtocolTCP:
		addr, _, _ = ep.TCP()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(opt.Endpoint)}
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
