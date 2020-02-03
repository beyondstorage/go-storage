package dropbox

import (
	"fmt"
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/types"
)

// New will create a new client.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	const errorMessage = "dropbox New: %w"

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	cfg := dropbox.Config{}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	switch credProtocol {
	case credential.ProtocolAPIKey:
		cfg.Token = cred[0]
	default:
		return nil, nil, fmt.Errorf(errorMessage, credential.ErrUnsupportedProtocol)
	}

	store := &Storage{
		client: files.New(cfg),
	}

	return nil, store, nil
}

func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
}
