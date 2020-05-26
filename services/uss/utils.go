package uss

import (
	"fmt"
	"strings"

	"github.com/Xuanwo/storage/types/info"
	"github.com/upyun/go-sdk/upyun"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/httpclient"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// Storage is the uss service.
type Storage struct {
	bucket *upyun.UpYun

	name    string
	workDir string
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager uss {Name: %s, WorkDir: %s}",
		s.name, s.workDir)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...*types.Pair) (storage.Storager, error) {
	return newStorager(pairs...)
}

func newStorager(pairs ...*types.Pair) (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Op: services.OpNewStorager, Type: Type, Err: err, Pairs: pairs}
		}
	}()

	store = &Storage{}

	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	cfg := &upyun.UpYunConfig{
		Bucket:   opt.Name,
		Operator: cred[0],
		Password: cred[1],
	}
	store.bucket = upyun.NewUpYun(cfg)
	// Set http client
	store.bucket.SetHTTPClient(httpclient.New(opt.HTTPClientOptions))
	store.name = opt.Name
	store.workDir = "/"
	if opt.HasWorkDir {
		store.workDir = opt.WorkDir
	}
	return
}

// ref: https://help.upyun.com/knowledge-base/errno/
func formatError(err error) error {
	fn := func(s string) bool {
		return strings.Contains(err.Error(), `"code": `+s)
	}

	switch {
	case !fn(""):
		// If body is empty
		switch {
		case strings.Contains(err.Error(), "404"):
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		default:
			return err
		}
	case fn("40400001"):
		// 40400001:	file or directory not found
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case fn("40100017"), fn("40100019"), fn("40300011"):
		// 40100017: user need permission
		// 40100019: account forbidden
		// 40300011: has no permission to delete
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
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

func (s *Storage) formatFileObject(v *upyun.FileInfo) (o *types.Object, err error) {
	o = &types.Object{
		ID:         v.Name,
		Name:       s.getRelPath(v.Name),
		Type:       types.ObjectTypeFile,
		Size:       v.Size,
		UpdatedAt:  v.Time,
		ObjectMeta: info.NewObjectMeta(),
	}

	if v.ETag != "" {
		o.SetETag(v.ETag)
	}
	if v.ContentType != "" {
		o.SetContentType(v.ContentType)
	}

	return o, nil
}
