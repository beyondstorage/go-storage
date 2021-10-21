package azblob

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"

	"go.beyondstorage.io/credential"
	"go.beyondstorage.io/endpoint"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/httpclient"
	"go.beyondstorage.io/v5/services"
	typ "go.beyondstorage.io/v5/types"
)

// Service is the azblob config.
type Service struct {
	service azblob.ServiceURL

	defaultPairs DefaultServicePairs
	features     ServiceFeatures

	typ.UnimplementedServicer
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

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	typ.UnimplementedStorager
	typ.UnimplementedAppender
	typ.UnimplementedDirer
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager azblob {Name: %s, WorkDir: %s}",
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
func newServicer(pairs ...typ.Pair) (srv *Service, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_servicer", Type: Type, Err: formatError(err), Pairs: pairs}
		}
	}()

	srv = &Service{}

	opt, err := parsePairServiceNew(pairs)
	if err != nil {
		return nil, err
	}

	ep, err := endpoint.Parse(opt.Endpoint)
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
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(opt.Endpoint)}
	}

	primaryURL, _ := url.Parse(uri)

	cred, err := credential.Parse(opt.Credential)
	if err != nil {
		return nil, err
	}
	if cred.Protocol() != credential.ProtocolHmac {
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(opt.Credential)}
	}

	credValue, err := azblob.NewSharedKeyCredential(cred.Hmac())
	if err != nil {
		return nil, err
	}

	httpClient := httpclient.New(opt.HTTPClientOptions)

	p := azblob.NewPipeline(credValue, azblob.PipelineOptions{
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

	if opt.HasDefaultServicePairs {
		srv.defaultPairs = opt.DefaultServicePairs
	}
	if opt.HasServiceFeatures {
		srv.features = opt.ServiceFeatures
	}
	return srv, nil
}

func newServicerAndStorager(pairs ...typ.Pair) (srv *Service, store *Storage, err error) {
	srv, err = newServicer(pairs...)
	if err != nil {
		return
	}

	store, err = srv.newStorage(pairs...)
	if err != nil {
		err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
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
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	// Handle errors returned by azblob.
	e, ok := err.(azblob.StorageError)
	if !ok {
		return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
	}

	switch azblob.StorageErrorCodeType(e.ServiceCode()) {
	case "":
		switch e.Response().StatusCode {
		case 404:
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		default:
			return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
		}
	case azblob.StorageErrorCodeBlobNotFound:
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case azblob.StorageErrorCodeInsufficientAccountPermissions:
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

	bucket := s.service.NewContainerURL(opt.Name)

	st = &Storage{
		bucket: bucket,

		name:    opt.Name,
		workDir: "/",
	}

	if opt.HasDefaultStoragePairs {
		st.defaultPairs = opt.DefaultStoragePairs
	}

	if opt.HasStorageFeatures {
		st.features = opt.StorageFeatures
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

func (s *Storage) formatFileObject(v azblob.BlobItemInternal) (o *typ.Object, err error) {
	o = s.newObject(false)
	o.ID = v.Name
	o.Path = s.getRelPath(v.Name)
	o.Mode |= typ.ModeRead

	o.SetLastModified(v.Properties.LastModified)
	o.SetEtag(string(v.Properties.Etag))

	if v.Properties.ContentLength != nil {
		o.SetContentLength(*v.Properties.ContentLength)
	}
	if v.Properties.ContentType != nil {
		o.SetContentType(*v.Properties.ContentType)
	}
	if len(v.Properties.ContentMD5) > 0 {
		o.SetContentMd5(base64.StdEncoding.EncodeToString(v.Properties.ContentMD5))
	}

	var sm ObjectSystemMetadata
	if value := v.Properties.AccessTier; value != "" {
		sm.AccessTier = string(value)
	}
	if v.Properties.CustomerProvidedKeySha256 != nil {
		sm.EncryptionKeySha256 = *v.Properties.CustomerProvidedKeySha256
	}
	if v.Properties.EncryptionScope != nil {
		sm.EncryptionScope = *v.Properties.EncryptionScope
	}
	if v.Properties.ServerEncrypted != nil {
		sm.ServerEncrypted = *v.Properties.ServerEncrypted
	}
	o.SetSystemMetadata(sm)

	return o, nil
}

func (s *Storage) newObject(done bool) *typ.Object {
	return typ.NewObject(s, done)
}

func calculateEncryptionHeaders(key []byte, scope string) (cpk azblob.ClientProvidedKeyOptions, err error) {
	if len(key) != 32 {
		err = ErrEncryptionKeyInvalid
		return
	}
	keyBase64 := base64.StdEncoding.EncodeToString(key)
	keySha256 := sha256.Sum256(key)
	keySha256Base64 := base64.StdEncoding.EncodeToString(keySha256[:])
	cpk = azblob.ClientProvidedKeyOptions{
		EncryptionKey:       &keyBase64,
		EncryptionKeySha256: &keySha256Base64,
		EncryptionAlgorithm: "AES256",
		EncryptionScope:     &scope,
	}
	return
}

const (
	// AppendBlobIfMaxSizeLessThanOrEqual ensures that the AppendBlock operation succeeds only if the append blob's size is less than or equal to a value.
	// ref: https://docs.microsoft.com/rest/api/storageservices/append-block.
	AppendBlobIfMaxSizeLessThanOrEqual = 4 * 1024 * 1024 * 50000
	// AppendSizeMaximum is the max append size in per append operation.
	// ref: https://docs.microsoft.com/rest/api/storageservices/append-block.
	AppendSizeMaximum = 4 * 1024 * 1024
	// AppendNumberMaximum is the max append numbers in append operation.
	// ref: https://docs.microsoft.com/rest/api/storageservices/append-block.
	AppendNumberMaximum = 50000

	// WriteSizeMaximum is the maximum size for write operation.
	// ref: https://docs.microsoft.com/en-us/rest/api/storageservices/put-blob
	WriteSizeMaximum = 5000 * 1024 * 1024
)

func checkError(err error, expect azblob.ServiceCodeType) bool {
	e, ok := err.(azblob.StorageError)
	if !ok {
		return false
	}

	return e.ServiceCode() == expect
}
