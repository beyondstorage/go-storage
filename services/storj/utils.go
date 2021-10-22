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

// Storage is the example client.
type Storage struct {
	project      *uplink.Project
	defaultPairs DefaultStoragePairs
	features     StorageFeatures

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
	var accessGrant string
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return nil, err
	}

	st := &Storage{
		name:    opt.Name,
		workDir: "/",
	}
	if opt.HasWorkDir {
		st.workDir = opt.WorkDir
	}
	if opt.HasDefaultStoragePairs {
		st.defaultPairs = opt.DefaultStoragePairs
	}
	if opt.HasStorageFeatures {
		st.features = opt.StorageFeatures
	}
	cp, err := credential.Parse(opt.Credential)
	if err != nil {
		return nil, err
	}
	switch cp.Protocol() {
	case credential.ProtocolAPIKey:
		accessGrant = cp.APIKey()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(opt.Credential)}
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
