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

	"go.beyondstorage.io/credential"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/httpclient"
	"go.beyondstorage.io/v5/services"
	typ "go.beyondstorage.io/v5/types"
)

// Service is the gcs config.
type Service struct {
	service   *gs.Client
	projectID string

	defaultPairs DefaultServicePairs
	features     ServiceFeatures

	typ.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer gcs")
}

// Storage is the gcs service client.
type Storage struct {
	bucket *gs.BucketHandle

	name    string
	workDir string

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	typ.UnimplementedStorager
	typ.UnimplementedDirer
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
	return newServicerAndStorager(pairs...)
}

// NewServicer will create Servicer only.
func NewServicer(pairs ...typ.Pair) (typ.Servicer, error) {
	return newServicer(pairs...)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...typ.Pair) (typ.Storager, error) {
	_, store, err := newServicerAndStorager(pairs...)
	return store, err
}

func newServicer(pairs ...typ.Pair) (srv *Service, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_servicer", Type: Type, Err: formatError(err), Pairs: pairs}
		}
	}()

	ctx := context.Background()
	srv = &Service{}

	opt, err := parsePairServiceNew(pairs)
	if err != nil {
		return nil, err
	}

	hc := httpclient.New(opt.HTTPClientOptions)

	var creds *google.Credentials

	cp, err := credential.Parse(opt.Credential)
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
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(opt.Credential)}
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
	srv.projectID = opt.ProjectID

	if opt.HasDefaultServicePairs {
		srv.defaultPairs = opt.DefaultServicePairs
	}
	if opt.HasServiceFeatures {
		srv.features = opt.ServiceFeatures
	}
	return
}

// New will create a new aliyun oss service.
func newServicerAndStorager(pairs ...typ.Pair) (srv *Service, store *Storage, err error) {
	srv, err = newServicer(pairs...)
	if err != nil {
		return
	}

	store, err = srv.newStorage(pairs...)
	if err != nil {
		err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
		return nil, nil, err
	}
	return srv, store, nil
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
func (s *Service) newStorage(pairs ...typ.Pair) (st *Storage, err error) {
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return nil, err
	}

	bucket := s.service.Bucket(opt.Name)

	store := &Storage{
		bucket: bucket,
		name:   opt.Name,

		workDir: "/",
	}

	if opt.HasDefaultStoragePairs {
		store.defaultPairs = opt.DefaultStoragePairs
	}
	if opt.HasStorageFeatures {
		store.features = opt.StorageFeatures
	}
	if opt.HasWorkDir {
		store.workDir = opt.WorkDir
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
