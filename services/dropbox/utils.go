package dropbox

import (
	"fmt"
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/httpclient"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// Storage is the dropbox client.
type Storage struct {
	client files.Client

	workDir string
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager dropbox {WorkDir: %s}",
		s.workDir,
	)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...*types.Pair) (storage.Storager, error) {
	return newStorager(pairs...)
}

// New will create a new client.
func newStorager(pairs ...*types.Pair) (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Op: services.OpNewStorager, Type: Type, Err: err, Pairs: pairs}
		}
	}()

	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return
	}

	cfg := dropbox.Config{
		Client: httpclient.New(opt.HTTPClientOptions),
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	switch credProtocol {
	case credential.ProtocolAPIKey:
		cfg.Token = cred[0]
	default:
		return nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	store = &Storage{
		client: files.New(cfg),

		workDir: "/",
	}

	if opt.HasWorkDir {
		store.workDir = opt.WorkDir
	}
	return
}

// ref: https://www.dropbox.com/developers/documentation/http/documentation
//
// FIXME: I don't know how to handle dropbox's API error correctly, please give me some help.
func formatError(err error) error {
	fn := func(errorSummary, s string) bool {
		return strings.HasPrefix(errorSummary, s)
	}

	switch e := err.(type) {
	case files.DownloadAPIError:
		if fn(e.ErrorSummary, "not_found") {
			err = fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		}
	case auth.AccessAPIError:
		err = fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	}
	return err
}
func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
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
