package ipfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ipfs "github.com/ipfs/go-ipfs-api"
	cmds "github.com/ipfs/go-ipfs-cmds"

	"go.beyondstorage.io/endpoint"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

// Service is the ipfs config.
// It is not usable, only for generate code
type Service struct {
	f Factory

	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	types.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer ipfs")
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

	ipfs *ipfs.Shell

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	workDir string
	gateway string

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager IPFS {WorkDir: %s}", s.workDir)
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

// NewStorager will create Storager only.
func (f *Factory) newStorage() (st *Storage, err error) {
	st = &Storage{
		f:        *f,
		features: f.storageFeatures(),
		workDir:  "/",
	}
	if f.WorkDir != "" {
		if !strings.HasSuffix(f.WorkDir, "/") {
			f.WorkDir += "/"
		}
		st.workDir = f.WorkDir
	}

	ep, err := endpoint.Parse(f.Endpoint)
	if err != nil {
		return nil, err
	}
	var e string
	switch ep.Protocol() {
	case endpoint.ProtocolHTTP:
		e, _, _ = ep.HTTP()
	case endpoint.ProtocolHTTPS:
		e, _, _ = ep.HTTPS()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(f.Endpoint)}
	}

	gate, err := endpoint.Parse(f.Gateway)
	if err != nil {
		return nil, err
	}
	switch gate.Protocol() {
	case endpoint.ProtocolHTTP:
		st.gateway, _, _ = gate.HTTP()
	case endpoint.ProtocolHTTPS:
		st.gateway, _, _ = gate.HTTPS()
	default:
		return nil, services.PairUnsupportedError{Pair: WithGateway(f.Gateway)}
	}

	sh := ipfs.NewShell(e)
	if !sh.IsUp() {
		return nil, errors.New("ipfs not online")
	}
	st.ipfs = sh

	return st, nil
}

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	e, ok := err.(*ipfs.Error)
	if !ok {
		return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
	}

	switch e.Message {
	case os.ErrNotExist.Error():
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	}

	// ref: https://github.com/ipfs/go-ipfs-cmds/blob/4ade007405e5d3befb14184290576c63cc43a6a3/error.go#L31
	switch e.Code {
	case int(cmds.ErrRateLimited):
		return fmt.Errorf("%w: %v", services.ErrRequestThrottled, err)
	case int(cmds.ErrImplementation):
		return fmt.Errorf("%w: %v", services.ErrServiceInternal, err)
	}

	return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
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
	path = strings.ReplaceAll(path, "\\", "/")

	if filepath.IsAbs(path) {
		return path
	}

	return s.workDir + path
}
