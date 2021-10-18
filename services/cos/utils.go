package cos

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"

	"go.beyondstorage.io/credential"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/httpclient"
	"go.beyondstorage.io/v5/services"
	typ "go.beyondstorage.io/v5/types"
)

// Service is the Tencent oss *Service config.
type Service struct {
	service *cos.Client
	client  *http.Client

	defaultPairs DefaultServicePairs
	features     ServiceFeatures

	typ.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer cos")
}

// Storage is the cos object storage service.
type Storage struct {
	bucket *cos.BucketService
	object *cos.ObjectService

	name     string
	location string
	workDir  string

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	typ.UnimplementedStorager
	typ.UnimplementedDirer
	typ.UnimplementedMultiparter
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager cos {Name: %s, WorkDir: %s}",
		s.name, s.workDir,
	)
}

// New will create both Servicer and Storager.
func New(pairs ...typ.Pair) (_ typ.Servicer, _ typ.Storager, err error) {
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

	srv = &Service{}

	opt, err := parsePairServiceNew(pairs)
	if err != nil {
		return nil, err
	}

	cp, err := credential.Parse(opt.Credential)
	if err != nil {
		return nil, err
	}
	if cp.Protocol() != credential.ProtocolHmac {
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(opt.Credential)}
	}
	ak, sk := cp.Hmac()

	httpClient := httpclient.New(opt.HTTPClientOptions)
	httpClient.Transport = &cos.AuthorizationTransport{
		Transport: httpClient.Transport,
		SecretID:  ak,
		SecretKey: sk,
	}

	srv.client = httpClient
	srv.service = cos.NewClient(nil, srv.client)

	if opt.HasDefaultServicePairs {
		srv.defaultPairs = opt.DefaultServicePairs
	}
	if opt.HasServiceFeatures {
		srv.features = opt.ServiceFeatures
	}
	return
}

// newServicerAndStorager will create a new Tencent oss service.
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

// All available storage classes are listed here.
const (
	// ref: https://cloud.tencent.com/document/product/436/7745
	storageClassHeader = "x-cos-storage-class"

	StorageClassStandard   = "STANDARD"
	StorageClassStandardIA = "STANDARD_IA"
	StorageClassArchive    = "ARCHIVE"
)

// ref: https://www.qcloud.com/document/product/436/7730
func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	// Handle errors returned by cos.
	e, ok := err.(*cos.ErrorResponse)
	if !ok {
		return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
	}

	switch e.Code {
	case "":
		switch e.Response.StatusCode {
		case 404:
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		default:
			return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
		}
	case "NoSuchKey":
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case "AccessDenied":
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

	st = &Storage{}

	url := cos.NewBucketURL(opt.Name, opt.Location, true)
	c := cos.NewClient(&cos.BaseURL{BucketURL: url}, s.client)

	st.bucket = c.Bucket
	st.object = c.Object
	st.name = opt.Name
	st.location = opt.Location

	st.workDir = "/"
	if opt.HasWorkDir {
		st.workDir = opt.WorkDir
	}

	if opt.HasDefaultStoragePairs {
		st.defaultPairs = opt.DefaultStoragePairs
	}

	if opt.HasStorageFeatures {
		st.features = opt.StorageFeatures
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

func (s *Storage) formatFileObject(v cos.Object) (o *typ.Object, err error) {
	o = s.newObject(false)
	o.ID = v.Key
	o.Path = s.getRelPath(v.Key)
	o.Mode |= typ.ModeRead

	o.SetContentLength(int64(v.Size))

	// COS returns different value depends on object upload method or
	// encryption method, so we can't treat this value as content-md5
	//
	// ref: https://cloud.tencent.com/document/product/436/7729
	if v.ETag != "" {
		o.SetEtag(v.ETag)
	}

	// COS uses ISO8601 format: "2019-05-27T11:26:14.000Z" in List
	//
	// ref: https://cloud.tencent.com/document/product/436/7729
	if v.LastModified != "" {
		t, err := time.Parse("2006-01-02T15:04:05.999Z", v.LastModified)
		if err != nil {
			return nil, err
		}
		o.SetLastModified(t)
	}

	var sm ObjectSystemMetadata
	if value := v.StorageClass; value != "" {
		sm.StorageClass = value
	}
	o.SetSystemMetadata(sm)

	return o, nil
}

func (s *Storage) newObject(done bool) *typ.Object {
	return typ.NewObject(s, done)
}

// All available server side algorithm are listed here.
const (
	// ref: https://cloud.tencent.com/document/product/436/7729
	serverSideEncryptionHeader                  = "x-cos-server-side-encryption"
	serverSideEncryptionCosKmsKeyIdHeader       = "x-cos-server-side-encryption-cos-kms-key-id"
	serverSideEncryptionCustomerAlgorithmHeader = "x-cos-server-side-encryption-customer-algorithm"
	serverSideEncryptionCustomerKeyMd5Header    = "x-cos-server-side-encryption-customer-key-MD5"
	serverSideEncryptionContextHeader           = "x-cos-server-side-encryption-context"

	ServerSideEncryptionAes256 = "AES256"
	ServerSideEncryptionCosKms = "cos/kms"
)

func calculateEncryptionHeaders(algo string, key []byte) (algorithm, keyBase64, keyMD5Base64 string, err error) {
	if len(key) != 32 {
		err = ErrServerSideEncryptionCustomerKeyInvalid
		return
	}
	keyBase64 = base64.StdEncoding.EncodeToString(key)
	keyMD5 := md5.Sum(key)
	keyMD5Base64 = base64.StdEncoding.EncodeToString(keyMD5[:])
	return
}

// cos service response error code
//
// ref: https://cloud.tencent.com/document/product/436/7730
const (
	// NoSuchKey the specified key does not exist.
	responseCodeNoSuchKey = "NoSuchKey"
	// NoSuchUpload the specified uploadId dose not exist.
	responseCodeNoSuchUpload = "NoSuchUpload"
)

func checkError(err error, code string) bool {
	if e, ok := err.(*cos.ErrorResponse); ok {
		return strings.Contains(e.Code, code)
	}

	return false
}

// multipartXXX are multipart upload restriction in COS, see more details at:
// https://cloud.tencent.com/document/product/436/7750
const (
	// multipartNumberMaximum is the max part count supported.
	multipartNumberMaximum = 10000
	// multipartSizeMaximum is the maximum size for each part, 5GB.
	multipartSizeMaximum = 5 * 1024 * 1024 * 1024
	// multipartSizeMinimum is the minimum size for each part, 1MB.
	multipartSizeMinimum = 1024 * 1024
)

const (
	// WriteSizeMaximum is the maximum size for write operation, 5GB.
	// ref: https://cloud.tencent.com/document/product/436/7749
	writeSizeMaximum = 5 * 1024 * 1024 * 1024
)
