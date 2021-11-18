package gdrive

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	. "github.com/dgraph-io/ristretto"

	"go.beyondstorage.io/credential"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

const (
	numCounters = 1e7     // number of keys to track frequency of (10M).
	maxCost     = 1 << 30 // maximum cost of cache (1GB).
	bufferItems = 64      // number of keys per Get buffer.
	cost        = 1
	expireTime  = 100
)

// Service is the gdrive config.
// It is not usable, only for generate code
type Service struct {
	f Factory

	defaultPairs types.DefaultServicePairs
	features     types.ServiceFeatures

	types.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer gdrive")
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

	name         string
	workDir      string
	service      *drive.Service
	cache        *Cache
	defaultPairs types.DefaultStoragePairs
	features     types.StorageFeatures

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager gdrive {Name: %s, WorkDir: %s}",
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

func (f *Factory) newStorage() (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err)}
		}
	}()

	store = &Storage{
		f:        *f,
		features: f.storageFeatures(),
		name:     f.Name,
		workDir:  "/",
	}

	// Init cache for storager
	ch, err := initCache(numCounters, maxCost, bufferItems)
	if err != nil {
		return nil, err
	}
	store.cache = ch

	if f.WorkDir != "" {
		store.workDir = f.WorkDir
	}

	ctx := context.Background()

	// Google drive only support authorized by Oauth2
	// Ref:https://developers.google.com/drive/api/v3/about-auth
	hc := &http.Client{}

	var credJSON []byte

	cp, err := credential.Parse(f.Credential)
	if err != nil {
		return nil, err
	}
	switch cp.Protocol() {
	case credential.ProtocolFile:
		credJSON, err = os.ReadFile(cp.File())
		if err != nil {
			return nil, err
		}
	case credential.ProtocolBase64:
		credJSON, err = base64.StdEncoding.DecodeString(cp.Base64())
		if err != nil {
			return nil, err
		}
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(f.Credential)}
	}

	// Loading token source from binary data.
	// DriveScope means full control of gdrive
	creds, err := google.CredentialsFromJSON(ctx, credJSON, drive.DriveScope)
	if err != nil {
		return nil, err
	}
	ot := &oauth2.Transport{
		Source: creds.TokenSource,
		Base:   hc.Transport,
	}
	hc.Transport = ot

	store.service, err = drive.NewService(ctx, option.WithHTTPClient(hc))
	if err != nil {
		return nil, err
	}

	return store, nil
}

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	e, ok := err.(*googleapi.Error)
	if !ok {
		return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
	}

	// FIXME: find better way to deal with errors.
	//Ref: https://developers.google.com/drive/api/v3/handle-errors
	switch e.Code {
	case 400:
		return fmt.Errorf("%w: %v", services.ErrCapabilityInsufficient, err)
	case 401:
		return fmt.Errorf("%w: %v", credential.ErrInvalidValue, err)
	case 403:
		return fmt.Errorf("%w: %v", services.ErrRequestThrottled, err)
	case 404:
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case 500:
		return fmt.Errorf("%w: %v", services.ErrServiceInternal, err)
	default:
		return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
	}
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

// getRelativePath will get relative path(fileName or directoryName) based on workDir
func (s *Storage) getRelativePath(path string, name string) string {
	relativePath := path
	prefix := strings.TrimPrefix(s.workDir, "/")
	if strings.HasPrefix(path, prefix) {
		relativePath = strings.TrimPrefix(path, prefix)
	}

	return strings.TrimPrefix(relativePath, "/") + "/" + name
}

// getAbsPath will calculate object storage's abs path
func (s *Storage) getAbsPath(path string) string {
	if strings.HasPrefix(path, s.workDir) {
		return strings.TrimPrefix(path, "/")
	} else if strings.TrimPrefix(s.workDir, "/") == strings.Trim(path, "/") {
		return strings.TrimPrefix(s.workDir, "/")
	} else if path == "" {
		return strings.TrimPrefix(s.workDir, "/")
	} else {
		prefix := strings.TrimPrefix(s.workDir, "/")
		if !strings.HasPrefix(path, prefix) {
			return prefix + "/" + path
		} else {
			return strings.TrimPrefix(path, "/")
		}
	}
}

// getRelPath will get object storage's rel path.
func (s *Storage) getRelPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/") + "/"
	return strings.TrimPrefix(path, prefix)
}

// getFileName will get a file's name without path
func (s *Storage) getFileName(path string) string {
	if strings.Contains(path, "/") {
		tmp := strings.Split(path, "/")
		return tmp[len(tmp)-1]
	} else {
		return path
	}
}

func initCache(nc int64, mc int64, bi int64) (cache *Cache, err error) {
	config := &Config{
		NumCounters: nc,
		MaxCost:     mc,
		BufferItems: bi,
	}

	cache, err = NewCache(config)

	if err != nil {
		return nil, err
	}
	return cache, nil
}

func (s *Storage) setCache(path string, fileId string) {
	s.cache.SetWithTTL(path, fileId, cost, expireTime)
}

func (s *Storage) getCache(path string) (string, bool) {
	id, found := s.cache.Get(path)
	if found {
		return id.(string), true
	}
	return "", false
}
