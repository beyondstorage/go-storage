package gcs

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	gs "cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"github.com/beyondstorage/go-storage/credential"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/services"
	typ "github.com/beyondstorage/go-storage/v5/types"
)

// Service is the gcs config.
type Service struct {
	f         Factory
	service   *gs.Client
	projectID string

	defaultPairs typ.DefaultServicePairs
	features     typ.ServiceFeatures

	typ.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer gcs")
}

// Storage is the gcs service client.
type Storage struct {
	f      Factory
	bucket *gs.BucketHandle

	name    string
	workDir string

	defaultPairs typ.DefaultStoragePairs
	features     typ.StorageFeatures

	typ.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager gcs {Name: %s, WorkDir: %s}",
		s.name, s.workDir,
	)
}

// New will create both Servicer and Storager.
func New(pairs ...typ.Pair) (typ.Servicer, typ.Storager, error) {
	f := Factory{}
	err := f.WithPairs(pairs...)
	if err != nil {
		return nil, nil, err
	}
	srv, err := f.NewServicer()
	if err != nil {
		return nil, nil, err
	}
	sto, err := f.NewStorager()
	if err != nil {
		return nil, nil, err
	}
	return srv, sto, nil
}

// NewServicer will create Servicer only.
func NewServicer(pairs ...typ.Pair) (typ.Servicer, error) {
	f := Factory{}
	err := f.WithPairs(pairs...)
	if err != nil {
		return nil, err
	}
	return f.NewServicer()
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

func (f *Factory) newService() (srv *Service, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_servicer", Type: Type, Err: formatError(err)}
		}
	}()

	ctx := context.Background()
	srv = &Service{
		f:        *f,
		features: f.serviceFeatures(),
	}

	hc := &http.Client{}

	var creds *google.Credentials

	cp, err := credential.Parse(f.Credential)
	if err != nil {
		return nil, err
	}
	switch cp.Protocol() {
	case credential.ProtocolFile:
		credJSON, err := os.ReadFile(cp.File())
		if err != nil {
			return nil, err
		}
		creds, err = google.CredentialsFromJSON(ctx, credJSON, gs.ScopeFullControl)
		if err != nil {
			return nil, err
		}
	case credential.ProtocolBase64:
		credJSON, err := base64.StdEncoding.DecodeString(cp.Base64())
		if err != nil {
			return nil, err
		}
		creds, err = google.CredentialsFromJSON(ctx, credJSON, gs.ScopeFullControl)
		if err != nil {
			return nil, err
		}
	case credential.ProtocolEnv:
		// Google provide DefaultCredentials support via env.
		// It will read credentials via:
		// - file path in GOOGLE_APPLICATION_CREDENTIALS
		// - Well known files on different platforms
		//   - On unix platform: `~/.config/gcloud/application_default_credentials.json`
		//   - On windows platform: `$APPDATA/gcloud/application_default_credentials.json`
		// - Metadata server in Google App Engine or Google Compute Engine
		creds, err = google.FindDefaultCredentials(ctx, gs.ScopeFullControl)
		if err != nil {
			return nil, err
		}
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(f.Credential)}
	}

	ot := &oauth2.Transport{
		Source: creds.TokenSource,
		Base:   hc.Transport,
	}
	hc.Transport = ot

	client, err := gs.NewClient(ctx, option.WithHTTPClient(hc))
	if err != nil {
		return nil, err
	}

	srv.service = client
	srv.projectID = f.ProjectID

	return
}

// All available storage classes are listed here.
const (
	StorageClassStandard = "STANDARD"
	StorageClassNearLine = "NEARLINE"
	StorageClassColdLine = "COLDLINE"
	StorageClassArchive  = "ARCHIVE"
)

// ref: https://cloud.google.com/storage/docs/json_api/v1/status-codes
func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	// gcs sdk could return explicit error, we should handle them.
	if errors.Is(err, gs.ErrObjectNotExist) {
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	}

	e, ok := err.(*googleapi.Error)
	if !ok {
		return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
	}

	switch e.Code {
	case http.StatusNotFound:
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case http.StatusForbidden:
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
	}
}

// newStorage will create a new client.
func (f *Factory) newStorage() (st *Storage, err error) {
	s, err := f.newService()
	if err != nil {
		return nil, err
	}

	bucket := s.service.Bucket(f.Name)

	store := &Storage{
		f:        *f,
		features: f.storageFeatures(),
		bucket:   bucket,
		name:     f.Name,

		workDir: "/",
	}

	if f.WorkDir != "" {
		store.workDir = f.WorkDir
	}
	return store, nil
}

func (s *Service) formatError(op string, err error, name string) error {
	if err == nil {
		return nil
	}

	return services.ServiceError{
		Op:       op,
		Err:      formatError(err),
		Servicer: s,
		Name:     name,
	}
}

// getAbsPath will calculate object storage's abs path
func (s *Storage) getAbsPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return prefix + path
}

// getRelPath will get object storage's rel path.
func (s *Storage) getRelPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return strings.TrimPrefix(path, prefix)
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

func (s *Storage) formatFileObject(v *gs.ObjectAttrs) (o *typ.Object, err error) {
	o = s.newObject(true)
	o.ID = v.Name
	o.Path = s.getRelPath(v.Name)
	o.Mode |= typ.ModeRead

	o.SetContentLength(v.Size)
	o.SetLastModified(v.Updated)

	if v.ContentType != "" {
		o.SetContentType(v.ContentType)
	}
	if v.Etag != "" {
		o.SetEtag(v.Etag)
	}
	if len(v.MD5) > 0 {
		o.SetContentMd5(base64.StdEncoding.EncodeToString(v.MD5))
	}

	var sm ObjectSystemMetadata
	if value := v.StorageClass; value != "" {
		sm.StorageClass = value
	}
	if value := v.CustomerKeySHA256; value != "" {
		sm.EncryptionKeySha256 = value
	}
	o.SetSystemMetadata(sm)

	return
}

func (s *Storage) newObject(done bool) *typ.Object {
	return typ.NewObject(s, done)
}
