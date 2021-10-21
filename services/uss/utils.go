package uss

import (
	"fmt"
	"strings"

	"github.com/upyun/go-sdk/v3/upyun"

	"go.beyondstorage.io/credential"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/httpclient"
	"go.beyondstorage.io/v5/services"
	typ "go.beyondstorage.io/v5/types"
)

// Storage is the uss service.
type Storage struct {
	bucket *upyun.UpYun

	name    string
	workDir string

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	typ.UnimplementedStorager
	typ.UnimplementedDirer
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager uss {Name: %s, WorkDir: %s}",
		s.name, s.workDir)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...typ.Pair) (typ.Storager, error) {
	return newStorager(pairs...)
}

func newStorager(pairs ...typ.Pair) (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
		}
	}()

	store = &Storage{}

	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return
	}

	cp, err := credential.Parse(opt.Credential)
	if err != nil {
		return nil, err
	}
	if cp.Protocol() != credential.ProtocolHmac {
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(opt.Credential)}
	}

	operator, password := cp.Hmac()
	cfg := &upyun.UpYunConfig{
		Bucket:   opt.Name,
		Operator: operator,
		Password: password,
	}
	store.bucket = upyun.NewUpYun(cfg)
	// Set http client
	store.bucket.SetHTTPClient(httpclient.New(opt.HTTPClientOptions))
	store.name = opt.Name
	store.workDir = "/"

	if opt.HasDefaultStoragePairs {
		store.defaultPairs = opt.DefaultStoragePairs
	}
	if opt.HasStorageFeatures {
		store.features = opt.StorageFeatures
	}

	if opt.HasWorkDir {
		store.workDir = opt.WorkDir
	}
	return
}

// ref: https://help.upyun.com/knowledge-base/errno/
func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	if ae, ok := err.(*upyun.Error); ok {
		switch ae.Code {
		case responseCodeFileOrDirectoryNotFound, responseCodeNotFoundMarkAsDeleted, responseCodeNotFoundBlockDeleted:
			// responseCodeFileOrDirectoryNotFound: 40400001, file or directory not found
			// responseCodeNotFoundMarkAsDeleted:   40401004, not found, mark as deleted
			// responseCodeNotFoundBlockDeleted:    40401005: not found, block deleted
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		case responseCodeUserNeedPermission, responseCodeAccountForbidden, responseCodeHasNoPermissionToDelete:
			// responseCodeUserNeedPermission:      40100017, user need permission
			// responseCodeAccountForbidden:        40100019, account forbidden
			// responseCodeHasNoPermissionToDelete: 40300011, has no permission to delete
			return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
		default:
			return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
		}
	}

	return fmt.Errorf("%w, %v", services.ErrUnexpected, err)
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

func (s *Storage) formatFileObject(v *upyun.FileInfo) (o *typ.Object, err error) {
	o = s.newObject(true)
	o.ID = v.Name
	o.Path = s.getRelPath(v.Name)

	o.SetContentLength(v.Size)
	o.SetLastModified(v.Time)
	// v.Meta means all the k-v in header with key which has prefix `x-upyun-meta-`
	// so we consider it as user's metadata
	// see more details at: https://github.com/upyun/go-sdk/blob/master/upyun/fileinfo.go#L39
	o.SetUserMetadata(v.Meta)

	if v.MD5 != "" {
		o.SetContentMd5(v.MD5)
	}
	if v.ContentType != "" {
		o.SetContentType(v.ContentType)
	}
	if v.IsDir {
		o.Mode |= typ.ModeDir
	} else {
		o.Mode |= typ.ModeRead
	}

	return o, nil
}

func (s *Storage) newObject(stated bool) *typ.Object {
	return typ.NewObject(s, stated)
}

// uss service response error code
//
// ref: http://docs.upyun.com/api/errno/
const (
	// responseCodeFileOrDirectoryNotFound file or directory not found
	responseCodeFileOrDirectoryNotFound = 40400001
	// responseCodeNotFoundMarkAsDeleted not found, mark as deleted
	responseCodeNotFoundMarkAsDeleted = 40401004
	// responseCodeNotFoundBlockDeleted not found, block deleted
	responseCodeNotFoundBlockDeleted = 40401005
	// responseCodeFolderAlreadyExist folder already exists
	responseCodeFolderAlreadyExist = 40600002
	// responseCodeUserNeedPermission user need permission
	responseCodeUserNeedPermission = 40100017
	// responseCodeAccountForbidden account forbidden
	responseCodeAccountForbidden = 40100019
	// responseCodeHasNoPermissionToDelete has no permission to delete
	responseCodeHasNoPermissionToDelete = 40300011
)

func checkErrorCode(err error, code int) bool {
	if ae, ok := err.(*upyun.Error); ok {
		return ae.Code == code
	}

	return false
}
