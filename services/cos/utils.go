package cos

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aos-dev/go-storage/v2/pkg/httpclient"
	"github.com/aos-dev/go-storage/v2/types/info"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/aos-dev/go-storage/v2"
	"github.com/aos-dev/go-storage/v2/pkg/credential"
	"github.com/aos-dev/go-storage/v2/services"
	"github.com/aos-dev/go-storage/v2/types"
	ps "github.com/aos-dev/go-storage/v2/types/pairs"
)

// Service is the Tencent oss *Service config.
type Service struct {
	service *cos.Client
	client  *http.Client
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
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager cos {Name: %s, WorkDir: %s}",
		s.name, s.workDir,
	)
}

// New will create both Servicer and Storager.
func New(pairs ...*types.Pair) (_ storage.Servicer, _ storage.Storager, err error) {
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

	srv = &Service{}

	opt, err := parsePairServiceNew(pairs)
	if err != nil {
		return nil, err
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	httpClient := httpclient.New(opt.HTTPClientOptions)
	httpClient.Transport = &cos.AuthorizationTransport{
		Transport: httpClient.Transport,
		SecretID:  cred[0],
		SecretKey: cred[1],
	}

	srv.client = httpClient
	srv.service = cos.NewClient(nil, srv.client)
	return
}

// newServicerAndStorager will create a new Tencent oss service.
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
	// Handle errors returned by cos.
	e, ok := err.(*cos.ErrorResponse)
	if !ok {
		return err
	}

	switch e.Code {
	case "":
		switch e.Response.StatusCode {
		case 404:
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		default:
			return err
		}
	case "NoSuchKey":
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case "AccessDenied":
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

func (s *Storage) formatFileObject(v cos.Object) (o *types.Object, err error) {
	o = &types.Object{
		ID:         v.Key,
		Name:       s.getRelPath(v.Key),
		Type:       types.ObjectTypeFile,
		Size:       int64(v.Size),
		ObjectMeta: info.NewObjectMeta(),
	}

	// COS returns different value depends on object upload method or
	// encryption method, so we can't treat this value as content-md5
	//
	// ref: https://cloud.tencent.com/document/product/436/7729
	if v.ETag != "" {
		o.SetETag(v.ETag)
	}

	// COS uses ISO8601 format: "2019-05-27T11:26:14.000Z" in List
	//
	// ref: https://cloud.tencent.com/document/product/436/7729
	if v.LastModified != "" {
		t, err := time.Parse("2006-01-02T15:04:05.999Z", v.LastModified)
		if err != nil {
			return nil, err
		}
		o.UpdatedAt = t
	}

	if value := v.StorageClass; value != "" {
		setStorageClass(o.ObjectMeta, value)
	}

	return o, nil
}
