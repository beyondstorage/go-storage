package qingstor

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/types"
	"github.com/yunify/qingstor-sdk-go/v3/config"
	qserror "github.com/yunify/qingstor-sdk-go/v3/request/errors"
	"github.com/yunify/qingstor-sdk-go/v3/service"
)

// New will create a new qingstor service.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	const errorMessage = "qingstor New: %w"

	srv := &Service{
		noRedirectClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, nil, fmt.Errorf(errorMessage, credential.ErrUnsupportedProtocol)
	}
	cfg, err := config.New(cred[0], cred[1])
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}
	if opt.HasEndpoint {
		ep := opt.Endpoint.Value()
		cfg.Host = ep.Host
		cfg.Port = ep.Port
		cfg.Protocol = ep.Protocol
	}

	srv.config = cfg
	srv.service, _ = service.Init(cfg)

	store, err := srv.newStorage(pairs...)
	if err != nil && errors.Is(err, types.ErrPairRequired) {
		return srv, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}
	return srv, store, nil
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

func handleQingStorError(err error) error {
	if err == nil {
		panic("error must not be nil")
	}

	var e *qserror.QingStorError
	e, ok := err.(*qserror.QingStorError)
	if !ok {
		return fmt.Errorf("%w: %v", types.ErrUnhandledError, err)
	}

	if e.Code == "" {
		switch e.StatusCode {
		case 404:
			return fmt.Errorf("%w: %v", types.ErrObjectNotExist, err)
		default:
			return fmt.Errorf("%w: %v", types.ErrUnhandledError, err)
		}
	}

	switch e.Code {
	case "permission_denied":
		return fmt.Errorf("%w: %v", types.ErrPermissionDenied, err)
	case "object_not_exists":
		return fmt.Errorf("%w: %v", types.ErrObjectNotExist, err)
	case "invalid_access_key_id":
		return fmt.Errorf("%w: %v", types.ErrConfigIncorrect, err)
	default:
		return fmt.Errorf("%w: %v", types.ErrUnhandledError, err)
	}
}

func convertUnixTimestampToTime(v int) time.Time {
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(int64(v), 0)
}

const (
	storageClassStandard   = "STANDARD"
	storageClassStandardIA = "STANDARD_IA"
)

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (string, error) {
	switch in {
	case storageclass.Hot:
		return storageClassStandard, nil
	case storageclass.Warm:
		return storageClassStandardIA, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in string) (storageclass.Type, error) {
	switch in {
	case storageClassStandard, "":
		return storageclass.Hot, nil
	case storageClassStandardIA:
		return storageclass.Warm, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}
