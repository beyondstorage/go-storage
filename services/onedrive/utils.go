package onedrive

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"go.beyondstorage.io/credential"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	typ "go.beyondstorage.io/v5/types"
)

// Service is the onedrive config.
// It is not usable, only for generate code
type Service struct {
	f Factory

	defaultPairs typ.DefaultServicePairs
	features     typ.ServiceFeatures

	typ.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer onedrive")
}

// NewServicer is not usable, only for generate code
func NewServicer(pairs ...typ.Pair) (typ.Servicer, error) {
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

	client  *onedriveClient
	workDir string

	defaultPairs typ.DefaultStoragePairs
	features     typ.StorageFeatures

	typ.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager onedrive {WorkDir: %s}",
		s.workDir,
	)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...typ.Pair) (typ.Storager, error) {
	f := Factory{}
	err := f.WithPairs(pairs...)
	if err != nil {
		return nil, err
	}
	return f.newStorage()
}

// newStorager will create a new onedrive storager client.
func (f *Factory) newStorage() (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err)}
		}
	}()

	var token []byte

	cp, err := credential.Parse(f.Credential)
	if err != nil {
		return nil, err
	}
	switch cp.Protocol() {
	case credential.ProtocolFile:
		token, err = os.ReadFile(cp.File())
		if err != nil {
			return nil, err
		}
	case credential.ProtocolBase64:
		token, err = base64.StdEncoding.DecodeString(cp.Base64())
		if err != nil {
			return nil, err
		}
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(f.Credential)}
	}

	// create new onedrive client
	client := getClient(context.TODO(), string(token))

	// generate work dir
	workDir := "/"
	if f.WorkDir != "" {
		workDir = f.WorkDir
	}
	store = &Storage{
		client:  &onedriveClient{client},
		workDir: workDir,
	}

	return
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

// formatError format all error into go-storage error
func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	if strings.Contains(err.Error(), "itemNotFound") {
		err = services.ErrObjectNotExist
	}

	return err
}

// getAbsPath return absolute path
func (s *Storage) getAbsPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	// join path
	absPath := filepath.Join(s.workDir, path)

	// append again
	if strings.HasSuffix(path, "/") {
		absPath += "/"
	}
	return absPath
}

// newObject new a object
func (s *Storage) newObject(done bool) *typ.Object {
	return typ.NewObject(s, done)
}

func (s *Storage) formatObject(v *Item, dir, rp string) *typ.Object {
	if v.Folder != nil {
		return s.formatFolderObject(v, dir, rp)
	}

	return s.formatFileObject(v, dir, rp)
}

// formatFolderObject format a onedrive folder object into go-storage object
func (s *Storage) formatFolderObject(v *Item, dir, rp string) (o *typ.Object) {
	o = s.newObject(true)

	folderName := path.Base(v.Name)

	o.ID = filepath.Join(rp, folderName)
	o.Path = path.Join(dir, folderName)

	o.Mode |= typ.ModeDir

	return o
}

// formatFolderObject format a onedrive file object into go-storage object
func (s *Storage) formatFileObject(v *Item, dir, rp string) (o *typ.Object) {
	o = s.newObject(true)

	fileName := path.Base(v.Name)

	o.ID = filepath.Join(rp, fileName)
	o.Path = path.Join(dir, fileName)

	o.Mode |= typ.ModeRead

	o.SetEtag(v.Etag)
	o.SetLastModified(v.LastModifiedDateTime)
	o.SetContentLength(v.Size)

	return o
}
