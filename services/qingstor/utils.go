package qingstor

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pengsrc/go-shared/convert"
	qsconfig "github.com/qingstor/qingstor-sdk-go/v4/config"
	iface "github.com/qingstor/qingstor-sdk-go/v4/interface"
	qserror "github.com/qingstor/qingstor-sdk-go/v4/request/errors"
	"github.com/qingstor/qingstor-sdk-go/v4/service"

	"github.com/aos-dev/go-storage/v2"
	"github.com/aos-dev/go-storage/v2/pkg/credential"
	"github.com/aos-dev/go-storage/v2/pkg/headers"
	"github.com/aos-dev/go-storage/v2/pkg/httpclient"
	"github.com/aos-dev/go-storage/v2/services"
	"github.com/aos-dev/go-storage/v2/types"
	"github.com/aos-dev/go-storage/v2/types/info"
	ps "github.com/aos-dev/go-storage/v2/types/pairs"
)

// Service is the qingstor service config.
type Service struct {
	config  *qsconfig.Config
	service iface.Service

	client *http.Client
}

// String implements Service.String
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

	// options for this storager.
	workDir string // workDir dir for all operation.
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

func newServicer(pairs ...*types.Pair) (srv *Service, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Op: services.OpNewServicer, Type: Type, Err: err, Pairs: pairs}
		}
	}()

	opt, err := parsePairServiceNew(pairs)
	if err != nil {
		return nil, err
	}

	srv = &Service{
		client: httpclient.New(opt.HTTPClientOptions),
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	cfg, err := qsconfig.New(cred[0], cred[1])
	if err != nil {
		return nil, err
	}

	// Set config's endpoint
	if opt.HasEndpoint {
		ep := opt.Endpoint.Value()
		cfg.Host = ep.Host
		cfg.Port = ep.Port
		cfg.Protocol = ep.Protocol
	}
	// Set config's http client
	cfg.Connection = srv.client

	srv.config = cfg
	srv.service, _ = service.Init(cfg)
	return
}

// New will create a new qingstor service.
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
	// Handle errors returned by qingstor.
	var e *qserror.QingStorError
	if !errors.As(err, &e) {
		return err
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
		return e
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

func (s *Service) newStorage(pairs ...*types.Pair) (store *Storage, err error) {
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return
	}

	// WorkDir should be an abs path, start and ends with "/"
	if opt.HasWorkDir && !isWorkDirValid(opt.WorkDir) {
		err = ErrInvalidWorkDir
		return
	}
	// set work dir into root path if no work dir passed
	if !opt.HasWorkDir {
		opt.WorkDir = "/"
	}

	if !IsBucketNameValid(opt.Name) {
		err = ErrInvalidBucketName
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

	if opt.HasWorkDir {
		st.workDir = opt.WorkDir
	}
	if opt.HasDisableURICleaning {
		st.config.DisableURICleaning = opt.DisableURICleaning
	}
	return st, nil
}

func (s *Service) detectLocation(name string) (location string, err error) {
	defer func() {
		err = s.formatError("detect_location", err, "")
	}()

	url := fmt.Sprintf("%s://%s.%s:%d", s.config.Protocol, name, s.config.Host, s.config.Port)

	r, err := s.client.Head(url)
	if err != nil {
		return
	}
	if r.StatusCode != http.StatusTemporaryRedirect {
		err = fmt.Errorf("head status is %d instead of %d", r.StatusCode, http.StatusTemporaryRedirect)
		return
	}

	// Example URL: https://bucket.zone.qingstor.com
	location = strings.Split(r.Header.Get(headers.Location), ".")[1]
	return
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

	return &services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}

func (s *Storage) formatFileObject(v *service.KeyType) (o *types.Object, err error) {
	o = &types.Object{
		ID:         *v.Key,
		Name:       s.getRelPath(*v.Key),
		Type:       types.ObjectTypeFile,
		Size:       service.Int64Value(v.Size),
		UpdatedAt:  convertUnixTimestampToTime(service.IntValue(v.Modified)),
		ObjectMeta: info.NewObjectMeta(),
	}

	if v.MimeType != nil {
		o.SetContentType(service.StringValue(v.MimeType))
	}
	if value := service.StringValue(v.StorageClass); value != "" {
		setStorageClass(o.ObjectMeta, value)
	}
	if v.Etag != nil {
		o.SetETag(service.StringValue(v.Etag))
	}
	return o, nil
}

func isObjectDirectory(o *service.KeyType) bool {
	return convert.StringValue(o.MimeType) == "application/x-directory"
}
