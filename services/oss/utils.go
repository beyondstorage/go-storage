package oss

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/beyondstorage/go-storage/credential"
	"github.com/beyondstorage/go-storage/endpoint"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/services"
	typ "github.com/beyondstorage/go-storage/v5/types"
)

// Service is the aliyun oss *Service config.
type Service struct {
	f       Factory
	service *oss.Client

	defaultPairs DefaultServicePairs
	features     ServiceFeatures

	typ.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer oss")
}

// Storage is the aliyun object storage service.
type Storage struct {
	f Factory

	bucket *oss.Bucket

	name    string
	workDir string

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	typ.UnimplementedStorager
	typ.UnimplementedAppender
	typ.UnimplementedMultiparter
	typ.UnimplementedDirer
	typ.UnimplementedLinker
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager oss {Name: %s, WorkDir: %s}",
		s.bucket.BucketName, s.workDir,
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

	srv = &Service{
		f:        *f,
		features: f.serviceFeatures(),
	}

	cp, err := credential.Parse(f.Credential)
	if err != nil {
		return nil, err
	}
	if cp.Protocol() != credential.ProtocolHmac {
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(f.Credential)}
	}
	ak, sk := cp.Hmac()

	ep, err := endpoint.Parse(f.Endpoint)
	if err != nil {
		return nil, err
	}

	var url string
	switch ep.Protocol() {
	case endpoint.ProtocolHTTP:
		url, _, _ = ep.HTTP()
	case endpoint.ProtocolHTTPS:
		url, _, _ = ep.HTTPS()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(f.Endpoint)}
	}

	var copts []oss.ClientOption
	copts = append(copts, oss.HTTPClient(&http.Client{}))

	srv.service, err = oss.New(url, ak, sk, copts...)
	if err != nil {
		return nil, err
	}
	return
}

// All available storage classes are listed here.
const (
	// ref: https://www.alibabacloud.com/help/doc-detail/31984.htm
	storageClassHeader = "x-oss-storage-class"

	// ref: https://www.alibabacloud.com/help/doc-detail/51374.htm
	StorageClassStandard = "STANDARD"
	StorageClassIA       = "IA"
	StorageClassArchive  = "Archive"
)

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	switch e := err.(type) {
	case oss.ServiceError:
		switch e.Code {
		case "":
			switch e.StatusCode {
			case 404:
				return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
			default:
				return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
			}
		case "NoSuchKey":
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		case "AccessDenied":
			return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
		}
	case oss.UnexpectedStatusCodeError:
		switch e.Got() {
		case 404:
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		case 403:
			return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
		}
	}

	return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
}

// newStorage will create a new client.
func (f *Factory) newStorage() (st *Storage, err error) {
	s, err := f.newService()
	if err != nil {
		return nil, err
	}

	bucket, err := s.service.Bucket(f.Name)
	if err != nil {
		return nil, err
	}

	store := &Storage{
		f:        *f,
		features: f.storageFeatures(),
		bucket:   bucket,

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

func (s *Storage) formatFileObject(v oss.ObjectProperties) (o *typ.Object, err error) {
	o = s.newObject(false)
	o.ID = v.Key
	o.Path = s.getRelPath(v.Key)
	if v.Type == "Symlink" {
		o.Mode |= typ.ModeLink
	} else {
		o.Mode |= typ.ModeRead
	}

	o.SetContentLength(v.Size)
	o.SetLastModified(v.LastModified)

	// OSS advise us don't use Etag as Content-MD5.
	//
	// ref: https://help.aliyun.com/document_detail/31965.html
	if v.ETag != "" {
		o.SetEtag(v.ETag)
	}

	var sm ObjectSystemMetadata
	if value := v.Type; value != "" {
		sm.StorageClass = value
	}
	o.SetSystemMetadata(sm)

	return
}

func (s *Storage) newObject(done bool) *typ.Object {
	return typ.NewObject(s, done)
}

// All available encryption algorithms are listed here.
const (
	serverSideEncryptionHeader      = "x-oss-server-side-encryption"
	serverSideEncryptionKeyIdHeader = "x-oss-server-side-encryption-key-id"

	ServerSideEncryptionAES256 = "AES256"
	ServerSideEncryptionKMS    = "KMS"
	ServerSideEncryptionSM4    = "SM4"

	ServerSideDataEncryptionSM4 = "SM4"
)

// OSS response error code.
//
// ref: https://error-center.alibabacloud.com/status/product/Oss
const (
	// responseCodeNoSuchUpload will be returned while the specified upload does not exist.
	responseCodeNoSuchUpload = "NoSuchUpload"
)

func checkError(err error, code string) bool {
	e, ok := err.(oss.ServiceError)
	if !ok {
		return false
	}

	return e.Code == code
}

// multipartXXX are multipart upload restriction in OSS, see more details at:
// https://help.aliyun.com/document_detail/31993.html
const (
	// multipartNumberMaximum is the max part count supported.
	multipartNumberMaximum = 10000
	// multipartSizeMaximum is the maximum size for each part, 5GB.
	multipartSizeMaximum = 5 * 1024 * 1024 * 1024
	// multipartSizeMinimum is the minimum size for each part, 100KB.
	multipartSizeMinimum = 100 * 1024
)

const (
	// writeSizeMaximum is the maximum size for each object with a single PUT operation, 5GB.
	// ref: https://help.aliyun.com/document_detail/31978.html#title-gkg-amg-aes
	writeSizeMaximum = 5 * 1024 * 1024 * 1024
	// appendSizeMaximum is the total maximum size for an append object, 5GB.
	// ref: https://help.aliyun.com/document_detail/31981.html?spm=a2c4g.11186623.6.1684.479a3ea7S8dRgB#title-22f-5c3-0sv
	appendTotalSizeMaximum = 5 * 1024 * 1024 * 1024
)
