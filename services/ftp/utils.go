package ftp

import (
	"fmt"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	mime "github.com/qingstor/go-mime"

	credential "github.com/beyondstorage/go-storage/credential"
	endpoint "github.com/beyondstorage/go-storage/endpoint"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/services"
	"github.com/beyondstorage/go-storage/v5/types"
)

// Storage is the example client.
type Storage struct {
	f          Factory
	connection *ftp.ServerConn
	user       string
	password   string
	url        string
	workDir    string

	defaultPairs types.DefaultStoragePairs
	features     types.StorageFeatures

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager ftp {URL: %s, User: %s, WorkDir: %s}", s.url, s.user, s.workDir)
}

// Service is the ftp config.
// It is not usable, only for generate code
type Service struct {
	f Factory

	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	types.UnimplementedServicer
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

// NewStorager will create Storager only.
func NewStorager(pairs ...types.Pair) (types.Storager, error) {
	f := Factory{}
	err := f.WithPairs(pairs...)
	if err != nil {
		return nil, err
	}
	return f.newStorage()
}

func (f *Factory) newStorage(pairs ...types.Pair) (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
		}
	}()

	store = &Storage{
		f:        *f,
		features: f.storageFeatures(),
		user:     "anonymous",
		password: "anonymous",
		url:      "localhost:21",
		workDir:  "/",
	}

	if f.WorkDir != "" {
		store.workDir = f.WorkDir
	}

	if f.Endpoint != "" {
		ep, err := endpoint.Parse(f.Endpoint)
		if err != nil {
			return nil, err
		}
		var host string
		var port int
		switch ep.Protocol() {
		case endpoint.ProtocolTCP:
			_, host, port = ep.TCP()
		default:
			return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(f.Endpoint)}
		}
		url := fmt.Sprintf("%s:%d", host, port)
		store.url = url
	}

	if f.Credential != "" {
		cp, err := credential.Parse(f.Credential)
		if err != nil {
			return nil, err
		}
		switch cp.Protocol() {
		case credential.ProtocolBasic:
			user, pass := cp.Basic()
			store.password = pass
			store.user = user
		default:
			return nil, services.PairUnsupportedError{Pair: ps.WithCredential(f.Credential)}
		}
	}

	err = store.connect()
	if err != nil {
		return nil, err
	}
	return
}

func (s *Storage) connect() error {
	c, err := ftp.Dial(s.url, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	err = c.Login(s.user, s.password)
	if err != nil {
		return err
	}

	err = c.ChangeDir(s.workDir)
	if err != nil {
		return err
	}

	s.connection = c
	return nil
}

func (s *Storage) makeDir(path string) error {
	if path == s.workDir || path == "." {
		return nil
	}
	rp := s.getAbsPath(path)
	return s.connection.MakeDir(rp)
}

// getAbsPath will calculate object storage's abs path(include workDir).
func (s *Storage) getAbsPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	absPath := filepath.Join(s.workDir, path)

	// Join will clean the trailng "/", we need to append it back.
	if strings.HasSuffix(path, string(filepath.Separator)) {
		absPath += string(filepath.Separator)
	}
	return absPath
}

// getRelPath will get object storage's rel path(exclude workDir).
func (s *Storage) getRelPath(path string) string {
	path = strings.TrimPrefix(path, string(filepath.Separator))
	return strings.TrimPrefix(path, s.workDir)
}

func (s *Storage) getNameList(path string) (namelist []string, err error) {
	namelist, err = s.connection.NameList(s.getAbsPath(path))
	if err != nil {
		return nil, err
	}
	return
}

func (s *Storage) newObject(done bool) *types.Object {
	return types.NewObject(s, done)
}

func (s *Storage) mapMode(fet ftp.EntryType) types.ObjectMode {
	switch fet {
	case ftp.EntryTypeFile:
		return types.ModeRead
	case ftp.EntryTypeFolder:
		return types.ModeDir
	case ftp.EntryTypeLink:
		return types.ModeLink
	}
	return types.ModeRead
}

func (s *Storage) formatFileObject(fe *ftp.Entry, parent string) (obj *types.Object, err error) {
	path := filepath.Join(parent, fe.Name)
	obj = types.NewObject(s, false)
	obj.SetID(path)
	obj.SetMode(s.mapMode(fe.Type))
	obj.SetPath(s.getRelPath(path))
	if fe.Type == ftp.EntryTypeFile {
		obj.SetContentLength(int64(fe.Size))
		obj.SetContentType(mime.DetectFilePath(path))
	}
	obj.SetLastModified(fe.Time)
	return
}

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}
	switch errX := err.(type) {
	case *textproto.Error:
		switch errX.Code {
		case ftp.StatusInvalidCredentials,
			ftp.StatusLoginNeedAccount,
			ftp.StatusStorNeedAccount:
			return fmt.Errorf("%w, %v", services.ErrPermissionDenied, err)
		case ftp.StatusFileUnavailable,
			ftp.StatusFileActionIgnored:
			return fmt.Errorf("%w, %v", services.ErrObjectNotExist, err)
		default:
			return fmt.Errorf("%w, %v", services.ErrServiceInternal, err)
		}
	}
	return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
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
