package dropbox

import (
	"fmt"
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"

	"go.beyondstorage.io/credential"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
	typ "go.beyondstorage.io/v5/types"
)

// Service is the dropbox config.
// It is not usable, only for generate code
type Service struct {
	f Factory

	defaultPairs typ.DefaultServicePairs
	features     typ.ServiceFeatures

	typ.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer azfile")
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

// Storage is the dropbox client.
type Storage struct {
	f Factory

	client files.Client

	workDir string

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	typ.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager dropbox {WorkDir: %s}",
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

// New will create a new client.
func (f *Factory) newStorage() (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err)}
		}
	}()

	cfg := dropbox.Config{
		// Client: httpclient.New(opt.HTTPClientOptions),
	}

	cred, err := credential.Parse(f.Credential)
	if err != nil {
		return nil, err
	}

	switch cred.Protocol() {
	case credential.ProtocolAPIKey:
		cfg.Token = cred.APIKey()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(f.Credential)}
	}

	store = &Storage{
		f:        *f,
		features: f.storageFeatures(),
		client:   files.New(cfg),

		workDir: "/",
	}

	if f.WorkDir != "" {
		store.workDir = f.WorkDir
	}
	return
}

// ref: https://www.dropbox.com/developers/documentation/http/documentation
//
// FIXME: I don't know how to handle dropbox's API error correctly, please give me some help.
func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	fn := func(errorSummary, s string) bool {
		return strings.Contains(errorSummary, s)
	}

	switch e := err.(type) {
	case files.GetMetadataAPIError:
		if fn(e.ErrorSummary, "not_found") {
			err = fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		}
	case files.DownloadAPIError:
		if fn(e.ErrorSummary, "not_found") {
			err = fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		}
	case auth.AccessAPIError:
		err = fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		err = fmt.Errorf("%w, %v", services.ErrUnexpected, err)
	}
	return err
}

func checkError(err error, codes ...string) bool {
	var s strings.Builder
	for _, code := range codes {
		s.WriteString(code)
		s.WriteString("/")
	}
	return strings.Contains(err.Error(), s.String())
}

func (s *Storage) getAbsPath(path string) string {
	// Return workDir while input path is empty.
	if path == "" {
		return s.workDir
	}
	// Return directly if input path is already an absolute path.
	if strings.HasPrefix(path, "/") {
		return path
	}
	return s.workDir + "/" + path
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

// getRelPath will get object storage's rel path.
func (s *Storage) getRelPath(path string) string {
	prefix := s.workDir + "/"
	return strings.TrimPrefix(path, prefix)
}

func (s *Storage) newObject(done bool) *typ.Object {
	return typ.NewObject(s, done)
}

func (s *Storage) formatFolderObject(path string, v *files.FolderMetadata) (o *typ.Object) {
	o = s.newObject(true)
	o.ID = v.Id
	o.Path = s.getRelPath(path)
	o.Mode |= typ.ModeDir

	return o
}

func (s *Storage) formatFileObject(path string, v *files.FileMetadata) (o *typ.Object) {
	o = s.newObject(true)
	o.ID = v.Id
	o.Path = s.getRelPath(path) + v.Metadata.Name
	o.Mode |= typ.ModeRead

	o.SetContentLength(int64(v.Size))
	o.SetLastModified(v.ServerModified)
	o.SetEtag(v.ContentHash)

	return o
}

const (
	// WriteSizeMaximum is the maximum size for write operation, 150MB.
	// ref: https://www.dropbox.com/developers/documentation/http/documentation#files-upload
	writeSizeMaximum = 150 * 1024 * 1024
	// AppendTotalSizeMaximum is the max append total size in append operation, 350GB.
	// ref: https://www.dropbox.com/developers/documentation/http/documentation#files-upload_session-append
	appendTotalSizeMaximum = 350 * 1024 * 1024 * 1024
)
