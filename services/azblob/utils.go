package azblob

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/aos-dev/go-storage/v2/types/info"

	"github.com/aos-dev/go-storage/v2"
	"github.com/aos-dev/go-storage/v2/pkg/credential"
	"github.com/aos-dev/go-storage/v2/pkg/httpclient"
	"github.com/aos-dev/go-storage/v2/services"
	"github.com/aos-dev/go-storage/v2/types"
	ps "github.com/aos-dev/go-storage/v2/types/pairs"
)

// Service is the azblob config.
type Service struct {
	service azblob.ServiceURL
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer azblob")
}

// Storage is the azblob service client.
type Storage struct {
	bucket azblob.ContainerURL

	name    string
	workDir string
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager azblob {Name: %s, WorkDir: %s}",
		s.name, s.workDir,
	)
}

// New will create both Servicer and Storager.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	return newServicerAndStorager(pairs...)
}

// NewServicer will create Servicer only.
func NewServicer(pairs ...*types.Pair) (storage.Servicer, error) {
	return newServicer(pairs...)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...*types.Pair) (storage.Storager, error) {
	_, store, err := newServicerAndStorager(pairs...)
	return store, err
}

// newServicer will create a azure blob servicer
//
// azblob use different URL to represent different sub services.
// - ServiceURL's          methods perform operations on a storage account.
//   - ContainerURL's     methods perform operations on an account's container.
//      - BlockBlobURL's  methods perform operations on a container's block blob.
//      - AppendBlobURL's methods perform operations on a container's append blob.
//      - PageBlobURL's   methods perform operations on a container's page blob.
//      - BlobURL's       methods perform operations on a container's blob regardless of the blob's type.
//
// Our Service will store a ServiceURL for operation.
func newServicer(pairs ...*types.Pair) (srv *Service, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Op: services.OpNewServicer, Type: Type, Err: err, Pairs: pairs}
		}
	}()

	srv = &Service{}

	opt, err := parsePairServiceNew(pairs)
	if err != nil {
		return nil, err
	}
	primaryURL, _ := url.Parse(opt.Endpoint.Value().String())

	credProtocol, credValue := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	cred, err := azblob.NewSharedKeyCredential(credValue[0], credValue[1])
	if err != nil {
		return nil, err
	}

	httpClient := httpclient.New(opt.HTTPClientOptions)

	p := azblob.NewPipeline(cred, azblob.PipelineOptions{
		HTTPSender: pipeline.FactoryFunc(func(next pipeline.Policy, po *pipeline.PolicyOptions) pipeline.PolicyFunc {
			return func(ctx context.Context, request pipeline.Request) (pipeline.Response, error) {
				r, err := httpClient.Do(request.WithContext(ctx))
				if err != nil {
					err = pipeline.NewError(err, "HTTP request failed")
				}
				return pipeline.NewHTTPResponse(r), err
			}
		}),
		// We don't need sdk level retry and we will handle read timeout by ourselves.
		Retry: azblob.RetryOptions{
			// Use a fixed back-off retry policy.
			Policy: 1,
			// A value of 1 means 1 try and no retries.
			MaxTries: 1,
			// Set a long enough timeout to adopt our timeout control.
			// This value could be adjusted to context deadline if request context has a deadline set.
			TryTimeout: 720 * time.Hour,
		},
	})
	srv.service = azblob.NewServiceURL(*primaryURL, p)

	return srv, nil
}

func newServicerAndStorager(pairs ...*types.Pair) (srv *Service, store *Storage, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Op: services.OpNewStorager, Type: Type, Err: err, Pairs: pairs}
		}
	}()

	srv, err = newServicer(pairs...)
	if err != nil {
		return
	}

	store, err = srv.newStorage(pairs...)
	if err != nil {
		return
	}
	return
}

// StorageClass is the storage class used in storage lib.
type StorageClass azblob.AccessTierType

// All available storage classes are listed here.
const (
	StorageClassArchive = azblob.AccessTierArchive
	StorageClassCool    = azblob.AccessTierCool
	StorageClassHot     = azblob.AccessTierHot
	StorageClassNone    = azblob.AccessTierNone
)

// ref: https://docs.microsoft.com/en-us/rest/api/storageservices/status-and-error-codes2
func formatError(err error) error {
	// Handle errors returned by azblob.
	e, ok := err.(azblob.StorageError)
	if !ok {
		return err
	}

	switch azblob.StorageErrorCodeType(e.ServiceCode()) {
	case "":
		switch e.Response().StatusCode {
		case 404:
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		default:
			return err
		}
	case azblob.StorageErrorCodeBlobNotFound:
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case azblob.StorageErrorCodeInsufficientAccountPermissions:
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
	}
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...*types.Pair) (st *Storage, err error) {
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return nil, err
	}

	bucket := s.service.NewContainerURL(opt.Name)

	st = &Storage{
		bucket: bucket,

		name:    opt.Name,
		workDir: "/",
	}

	if opt.HasWorkDir {
		st.workDir = opt.WorkDir
	}
	return st, nil
}

func (s *Service) formatError(op string, err error, name string) error {
	if err == nil {
		return nil
	}

	return &services.ServiceError{
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

	return &services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}

func (s *Storage) formatFileObject(v azblob.BlobItem) (o *types.Object, err error) {
	o = &types.Object{
		ID:         v.Name,
		Name:       s.getRelPath(v.Name),
		Type:       types.ObjectTypeFile,
		UpdatedAt:  v.Properties.LastModified,
		ObjectMeta: info.NewObjectMeta(),
	}

	o.SetETag(string(v.Properties.Etag))

	if v.Properties.ContentLength != nil {
		o.Size = *v.Properties.ContentLength
	}
	if v.Properties.ContentType != nil && *v.Properties.ContentType != "" {
		o.SetContentType(*v.Properties.ContentType)
	}
	if len(v.Properties.ContentMD5) > 0 {
		o.SetContentMD5(base64.StdEncoding.EncodeToString(v.Properties.ContentMD5))
	}
	if value := v.Properties.AccessTier; value != "" {
		setStorageClass(o.ObjectMeta, StorageClass(value))
	}

	return o, nil
}

func (s *Storage) formatDirObject(v azblob.BlobPrefix) (o *types.Object) {
	o = &types.Object{
		ID:         v.Name,
		Name:       s.getRelPath(v.Name),
		Type:       types.ObjectTypeDir,
		Size:       0,
		ObjectMeta: info.NewObjectMeta(),
	}

	return o
}
