package onedrive

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"go.beyondstorage.io/credential"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	typ "go.beyondstorage.io/v5/types"
)

// Storage is the example client.
type Storage struct {
	client  *onedriveClient
	workDir string

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

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
	return newStorager(pairs...)
}

// newStorager will create a new onedrive storager client.
func newStorager(pairs ...typ.Pair) (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
		}
	}()

	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return
	}

	var token []byte

	cp, err := credential.Parse(opt.Credential)
	if err != nil {
		return nil, err
	}
	switch cp.Protocol() {
	case credential.ProtocolFile:
		token, err = ioutil.ReadFile(cp.File())
		if err != nil {
			return nil, err
		}
	case credential.ProtocolBase64:
		token, err = base64.StdEncoding.DecodeString(cp.Base64())
		if err != nil {
			return nil, err
		}
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(opt.Credential)}
	}

	// create new onedrive client
	client := getClient(context.TODO(), string(token))

	// generate work dir
	workDir := "/"
	if opt.HasWorkDir {
		workDir = opt.WorkDir
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
