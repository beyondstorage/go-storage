package dropbox

import (
	"fmt"
	"strings"

	ps "github.com/Xuanwo/storage/types/pairs"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
)

// New will create a new client.
func New(pairs ...*types.Pair) (_ storage.Servicer, _ storage.Storager, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Type: Type, Err: err, Pairs: pairs}
		}
	}()

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, nil, err
	}

	cfg := dropbox.Config{}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	switch credProtocol {
	case credential.ProtocolAPIKey:
		cfg.Token = cred[0]
	default:
		return nil, nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}

	store := &Storage{
		client: files.New(cfg),

		workDir: "/",
	}

	if opt.HasWorkDir {
		store.workDir = opt.WorkDir
	}
	return nil, store, nil
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
