package storj

import (
	"context"
	"fmt"
	"path/filepath"

	"storj.io/uplink"

	"go.beyondstorage.io/credential"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

// Service is the storj config.
// It is not usable, only for generate code
type Service struct {
	f Factory

	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	types.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer storj")
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
	srv = &Service{
		f:        *f,
		features: f.serviceFeatures(),
	}
	return
}

// Storage is the example client.
type Storage struct {
	f Factory

	project      *uplink.Project
	defaultPairs types.DefaultStoragePairs
	features     types.StorageFeatures

	name    string
	workDir string
	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager storj {Name: %s, WorkDir: %s}",
		s.name, s.workDir,
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

// NewStorage will create Storager only.
func (f *Factory) newStorage() (types.Storager, error) {
	var accessGrant string

	st := &Storage{
		f:        *f,
		features: f.storageFeatures(),
		name:     f.Name,
		workDir:  "/",
	}
	if f.WorkDir != "" {
		st.workDir = f.WorkDir
	}
	cp, err := credential.Parse(f.Credential)
	if err != nil {
		return nil, err
	}
	switch cp.Protocol() {
	case credential.ProtocolAPIKey:
		accessGrant = cp.APIKey()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(f.Credential)}
	}
	access, err := uplink.ParseAccess(accessGrant)
	if err != nil {
		return nil, err
	}
	st.project, err = uplink.OpenProject(context.Background(), access)
	if err != nil {
		return nil, err
	}
	return st, nil
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

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	return err
}

func (s *Storage) getAbsPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return s.workDir + path
}
