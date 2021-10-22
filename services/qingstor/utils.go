package qingstor

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pengsrc/go-shared/convert"
	qsconfig "github.com/qingstor/qingstor-sdk-go/v4/config"
	iface "github.com/qingstor/qingstor-sdk-go/v4/interface"
	qserror "github.com/qingstor/qingstor-sdk-go/v4/request/errors"
	"github.com/qingstor/qingstor-sdk-go/v4/service"

	"go.beyondstorage.io/credential"
	"go.beyondstorage.io/endpoint"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/headers"
	"go.beyondstorage.io/v5/pkg/httpclient"
	"go.beyondstorage.io/v5/services"
	typ "go.beyondstorage.io/v5/types"
)

// Service is the qingstor service config.
type Service struct {
	config  *qsconfig.Config
	service iface.Service

	client *http.Client

	defaultPairs DefaultServicePairs
	features     ServiceFeatures

	typ.UnimplementedServicer
}

// String implements Service.String.
func (s *Service) String() string {
	if s.config == nil {
		return fmt.Sprintf("Servicer qingstor")
	}
	return fmt.Sprintf("Servicer qingstor {AccessKey: %s}", s.config.AccessKeyID)
}

// Storage is the qingstor object storage client.
type Storage struct {
	bucket     iface.Bucket
	config     *qsconfig.Config
	properties *service.Properties

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	// options for this storager.
	workDir string // workDir dir for all operation.

	typ.UnimplementedStorager
	typ.UnimplementedCopier
	typ.UnimplementedFetcher
	typ.UnimplementedMover
	typ.UnimplementedMultiparter
	typ.UnimplementedReacher
	typ.UnimplementedAppender
	typ.UnimplementedDirer
	typ.UnimplementedLinker
	typ.UnimplementedStorageHTTPSigner
}

// String implements Storager.String
func (s *Storage) String() string {
	// qingstor work dir should start and end with "/"
	return fmt.Sprintf(
		"Storager qingstor {Name: %s, Location: %s, WorkDir: %s}",
		*s.properties.BucketName, *s.properties.Zone, s.workDir,
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

	opt, err := parsePairServiceNew(pairs)
	if err != nil {
		return nil, err
	}

	srv = &Service{
		client: httpclient.New(opt.HTTPClientOptions),
	}

	var cfg *qsconfig.Config

	// Set config's credential.
	cp, err := credential.Parse(opt.Credential)
	if err != nil {
		return nil, err
	}
	switch cp.Protocol() {
	case credential.ProtocolHmac:
		cfg, err = qsconfig.New(cp.Hmac())
		if err != nil {
			return nil, err
		}
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(opt.Credential)}
	}

	// Set config's endpoint
	if opt.HasEndpoint {
		ep, err := endpoint.Parse(opt.Endpoint)
		if err != nil {
			return nil, err
		}

		switch ep.Protocol() {
		case endpoint.ProtocolHTTPS:
			_, cfg.Host, cfg.Port = ep.HTTPS()
		case endpoint.ProtocolHTTP:
			_, cfg.Host, cfg.Port = ep.HTTP()
		default:
			return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(opt.Endpoint)}
		}

		cfg.Protocol = ep.Protocol()
	}
	// Set config's http client
	cfg.Connection = srv.client

	srv.config = cfg
	srv.service, _ = service.Init(cfg)

	if opt.HasDefaultServicePairs {
		srv.defaultPairs = opt.DefaultServicePairs
	}
	if opt.HasServiceFeatures {
		srv.features = opt.ServiceFeatures
	}
	return
}

// New will create a new qingstor service.
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

// multipartXXX are multipart upload restriction in QingStor, see more detail at:
// https://docs.qingcloud.com/qingstor/api/object/multipart/index.html#%E5%88%86%E6%AE%B5%E4%B8%8A%E4%BC%A0%E9%99%90%E5%88%B6
const (
	// multipartNumberMinimum is the min part count supported
	multipartNumberMinimum = 0
	// multipartNumberMaximum is the max part count supported
	multipartNumberMaximum = 10000
	// multipartSizeMaximum is the maximum size for each part, 5GB
	multipartSizeMaximum = 5 * 1024 * 1024 * 1024
	// multipartSizeMinimum is the minimum size for each part, except the last part, 4MB
	multipartSizeMinimum = 4 * 1024 * 1024
)

const (
	// writeSizeMaximum is the maximum size for write operation, 5GB.
	// ref: https://docs.qingcloud.com/qingstor/#object
	writeSizeMaximum = 5 * 1024 * 1024 * 1024
	// copySizeMaximum is the maximum size for copy operation, 5GB.
	// ref: https://docs.qingcloud.com/qingstor/api/object/copy
	copySizeMaximum = 5 * 1024 * 1024 * 1024
	// appendSizeMaximum is the maximum append size for per append operation, 5GB.
	// ref: https://docs.qingcloud.com/qingstor/api/object/append
	appendSizeMaximum = 5 * 1024 * 1024 * 1024
	// appendSizeMaximum is the total maximum size for an append object, 5TB.
	// ref: https://docs.qingcloud.com/qingstor/api/object/append
	appendTotalSizeMaximum = 50 * 1024 * 1024 * 1024 * 1024
)

// bucketNameRegexp is the bucket name regexp, which indicates:
// 1. length: 6-63;
// 2. contains lowercase letters, digits and strikethrough;
// 3. starts and ends with letter or digit.
var bucketNameRegexp = regexp.MustCompile(`^[a-z\d][a-z-\d]{4,61}[a-z\d]$`)

// IsBucketNameValid will check whether given string is a valid bucket name.
func IsBucketNameValid(s string) bool {
	return bucketNameRegexp.MatchString(s)
}

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	// Handle errors returned by qingstor.
	var e *qserror.QingStorError
	if !errors.As(err, &e) {
		return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
	}

	switch e.Code {
	case "":
		// code=="" means this response doesn't have body.
		switch e.StatusCode {
		case 404:
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, e)
		default:
			return e
		}
	case "permission_denied":
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, e)
	case "object_not_exists":
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, e)
	default:
		return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
	}
}

func convertUnixTimestampToTime(v int) time.Time {
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(int64(v), 0)
}

// All available storage classes are listed here.
const (
	StorageClassStandard   = "STANDARD"
	StorageClassStandardIA = "STANDARD_IA"
)

func (s *Service) newStorage(pairs ...typ.Pair) (store *Storage, err error) {
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return
	}

	// WorkDir should be an abs path, start and ends with "/"
	if opt.HasWorkDir && !isWorkDirValid(opt.WorkDir) {
		err = ErrWorkDirInvalid
		return
	}
	// set work dir into root path if no work dir passed
	if !opt.HasWorkDir {
		opt.WorkDir = "/"
	}

	if !IsBucketNameValid(opt.Name) {
		err = ErrBucketNameInvalid
		return
	}

	// Detect location automatically
	if !opt.HasLocation {
		opt.Location, err = s.detectLocation(opt.Name)
		if err != nil {
			return
		}
	}

	bucket, err := s.service.Bucket(opt.Name, opt.Location)
	if err != nil {
		return
	}

	st := &Storage{
		bucket:     bucket,
		config:     bucket.Config,
		properties: bucket.Properties,

		workDir: "/",
	}

	if opt.HasDisableURICleaning {
		st.config.DisableURICleaning = opt.DisableURICleaning
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

func (s *Service) detectLocation(name string) (location string, err error) {
	defer func() {
		err = s.formatError("detect_location", err, "")
	}()

	u := fmt.Sprintf("%s://%s:%d/%s", s.config.Protocol, s.config.Host, s.config.Port, name)

	r, err := s.client.Head(u)
	if err != nil {
		return
	}
	if r.StatusCode != http.StatusMovedPermanently {
		err = fmt.Errorf("%w: head status is %d instead of %d", services.ErrUnexpected, r.StatusCode, http.StatusMovedPermanently)
		return
	}

	// Example URL: https://zone.qingstor.com/bucket
	locationUrl, err := url.Parse(r.Header.Get(headers.Location))
	if err != nil {
		return
	}
	location = strings.Split(locationUrl.Host, ".")[0]
	return
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

// isWorkDirValid check qingstor work dir
// work dir must start with only one "/" (abs path), and end with only one "/" (a dir).
// If work dir is the root path, set it to "/".
func isWorkDirValid(wd string) bool {
	return strings.HasPrefix(wd, "/") && // must start with "/"
		strings.HasSuffix(wd, "/") && // must end with "/"
		!strings.HasPrefix(wd, "//") && // not start with more than one "/"
		!strings.HasSuffix(wd, "//") // not end with more than one "/"
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

func (s *Storage) newObject(done bool) *typ.Object {
	return typ.NewObject(s, done)
}

func (s *Storage) formatFileObject(v *service.KeyType) (o *typ.Object, err error) {
	o = s.newObject(false)
	o.ID = *v.Key
	o.Path = s.getRelPath(*v.Key)
	// If you have enabled virtual link, you will not get the accurate object type.
	// If you want to get the exact object mode, please use `stat`
	o.Mode |= typ.ModeRead

	o.SetContentLength(service.Int64Value(v.Size))
	o.SetLastModified(convertUnixTimestampToTime(service.IntValue(v.Modified)))

	if v.MimeType != nil {
		o.SetContentType(service.StringValue(v.MimeType))
	}
	if v.Etag != nil {
		o.SetEtag(service.StringValue(v.Etag))
	}

	var sm ObjectSystemMetadata
	if value := service.StringValue(v.StorageClass); value != "" {
		sm.StorageClass = value
	}
	o.SetSystemMetadata(sm)

	return o, nil
}

func isObjectDirectory(o *service.KeyType) bool {
	return convert.StringValue(o.MimeType) == "application/x-directory"
}

// All available SSE customer algorithms are listed here.
const (
	SseCustomerAlgorithmAes256 = "AES256"
)

func calculateEncryptionHeaders(algo string, key []byte) (algorithm, keyBase64, keyMD5Base64 *string, err error) {
	if len(key) != 32 {
		err = ErrEncryptionCustomerKeyInvalid
		return
	}
	kB64 := base64.StdEncoding.EncodeToString(key)
	kMD5 := md5.Sum(key)
	kMD5B64 := base64.StdEncoding.EncodeToString(kMD5[:])
	return &algo, &kB64, &kMD5B64, nil
}

func (s *Storage) formatGetObjectInput(opt pairStorageRead) (input *service.GetObjectInput, err error) {
	input = &service.GetObjectInput{}
	if opt.HasEncryptionCustomerAlgorithm {
		input.XQSEncryptionCustomerAlgorithm, input.XQSEncryptionCustomerKey, input.XQSEncryptionCustomerKeyMD5, err = calculateEncryptionHeaders(opt.EncryptionCustomerAlgorithm, opt.EncryptionCustomerKey)
		if err != nil {
			return
		}
	}

	if opt.HasOffset || opt.HasSize {
		rs := headers.FormatRange(opt.Offset, opt.Size)
		input.Range = &rs
	}

	return
}

func (s *Storage) formatPutObjectInput(size int64, opt pairStorageWrite) (input *service.PutObjectInput, err error) {
	input = &service.PutObjectInput{
		ContentLength: &size,
	}
	if opt.HasContentMd5 {
		input.ContentMD5 = service.String(opt.ContentMd5)
	}
	if opt.HasStorageClass {
		input.XQSStorageClass = service.String(opt.StorageClass)
	}
	if opt.HasEncryptionCustomerAlgorithm {
		input.XQSEncryptionCustomerAlgorithm, input.XQSEncryptionCustomerKey, input.XQSEncryptionCustomerKeyMD5, err = calculateEncryptionHeaders(opt.EncryptionCustomerAlgorithm, opt.EncryptionCustomerKey)
		if err != nil {
			return
		}
	}

	return
}
