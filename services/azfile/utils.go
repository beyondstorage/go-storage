package azfile

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-storage-file-go/azfile"

	"go.beyondstorage.io/credential"
	"go.beyondstorage.io/endpoint"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

// Service is the azfile config.
// It is not usable, only for generate code
type Service struct {
	f Factory

	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	types.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer azblob")
}

// Storage is the azfile client.
type Storage struct {
	f      Factory
	client azfile.DirectoryURL

	workDir string

	defaultPairs types.DefaultStoragePairs
	features     types.StorageFeatures

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager azfile {WorkDir: %s}", s.workDir)
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

// NewStorager will create Storager only.
func NewStorager(pairs ...types.Pair) (types.Storager, error) {
	f := Factory{}
	err := f.WithPairs(pairs...)
	if err != nil {
		return nil, err
	}
	return f.newStorage()
}

// newStorager will create a storage client.
func (f *Factory) newStorage() (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err)}
		}
	}()

	store = &Storage{
		workDir: "/",
	}

	if f.WorkDir == "" {
		store.workDir = f.WorkDir
	}

	ep, err := endpoint.Parse(f.Endpoint)
	if err != nil {
		return nil, err
	}

	var uri string
	switch ep.Protocol() {
	case endpoint.ProtocolHTTP:
		uri, _, _ = ep.HTTP()
	case endpoint.ProtocolHTTPS:
		uri, _, _ = ep.HTTPS()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(f.Endpoint)}
	}

	primaryURL, _ := url.Parse(uri)

	cred, err := credential.Parse(f.Credential)
	if err != nil {
		return nil, err
	}
	if cred.Protocol() != credential.ProtocolHmac {
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(f.Credential)}
	}

	credValue, err := azfile.NewSharedKeyCredential(cred.Hmac())
	if err != nil {
		return nil, err
	}

	p := azfile.NewPipeline(credValue, azfile.PipelineOptions{
		Retry: azfile.RetryOptions{},
	})

	serviceURL := azfile.NewServiceURL(*primaryURL, p)

	ctx := context.Background()
	shareURL := serviceURL.NewShareURL(f.Name)

	workDir := strings.TrimPrefix(store.workDir, "/")
	store.client = shareURL.NewDirectoryURL(workDir)
	_, err = store.client.Create(ctx, azfile.Metadata{}, azfile.SMBProperties{})
	if err != nil {
		return nil, err
	}

	return store, nil
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

// formatError converts errors returned by SDK into errors defined in go-storage and go-service-*.
// The original error SHOULD NOT be wrapped.
func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	e, ok := err.(azfile.StorageError)

	if ok {
		switch azfile.StorageErrorCodeType(e.ServiceCode()) {
		case "":
			switch e.Response().StatusCode {
			case fileNotFound:
				return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
			default:
				return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
			}
		case azfile.StorageErrorCodeResourceNotFound:
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		case azfile.StorageErrorCodeInsufficientAccountPermissions:
			return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
		default:
			return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
		}
	}

	return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
}

// getAbsPath will calculate object storage's abs path
func (s *Storage) getAbsPath(path string) string {
	if strings.HasPrefix(path, s.workDir) {
		return strings.TrimPrefix(path, "/")
	}

	prefix := strings.TrimPrefix(s.workDir, "/")
	return prefix + path
}

// getRelPath will get object storage's rel path.
func (s *Storage) getRelPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return strings.TrimPrefix(path, prefix)
}

// getRelativePath will get relative path(fileName or directoryName) based on workDir for DirectoryURL or FileURL.
func (s *Storage) getRelativePath(path string) string {
	relativePath := path
	if strings.HasPrefix(path, s.workDir) {
		relativePath = strings.TrimPrefix(path, s.workDir)
		return strings.TrimPrefix(relativePath, "/")
	}

	return relativePath
}

func (s *Storage) mkDirs(ctx context.Context, path string) (err error) {
	if "." == path {
		return
	}

	subDirs := strings.Split(path, "/")
	i := len(subDirs)
	existedDir := ""
	for i > 0 {
		existedDir = strings.Join(subDirs[0:i], "/")
		_, err = s.client.NewDirectoryURL(existedDir).GetProperties(ctx)
		if err == nil {
			// dir already exists
			break
		} else if !checkError(err, fileNotFound) {
			// Something error other then file not found, return directly.
			return err
		} else {
			// dir does not exist
			i--
			existedDir = ""
		}
	}

	currentDir := existedDir
	for _, v := range subDirs[i:] {
		if currentDir == "" {
			currentDir = v
		} else {
			currentDir += "/" + v
		}

		_, err = s.client.NewDirectoryURL(currentDir).Create(ctx, azfile.Metadata{}, azfile.SMBProperties{})
		if err != nil {
			return err
		}
	}

	return
}

func (s *Storage) newObject(done bool) *types.Object {
	return types.NewObject(s, done)
}

func (s *Storage) formatFileObject(v azfile.FileItem) (o *types.Object, err error) {
	o = s.newObject(true)
	o.ID = v.Name
	o.Path = s.getRelPath(v.Name)
	o.Mode |= types.ModeRead

	if v.Properties.ContentLength != 0 {
		o.SetContentLength(v.Properties.ContentLength)
	}

	return
}

func (s *Storage) formatDirObject(v azfile.DirectoryItem) (o *types.Object, err error) {
	o = s.newObject(true)
	o.ID = v.Name
	o.Path = s.getRelPath(v.Name)
	o.Mode |= types.ModeDir

	return
}

const (
	// File not found error.
	fileNotFound = 404
)

func checkError(err error, expect int) bool {
	e, ok := err.(azfile.StorageError)
	if !ok {
		return false
	}

	return e.Response().StatusCode == expect
}
