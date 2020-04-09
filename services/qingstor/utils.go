package qingstor

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/yunify/qingstor-sdk-go/v3/config"
	qserror "github.com/yunify/qingstor-sdk-go/v3/request/errors"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

// New will create a new qingstor service.
func New(pairs ...*types.Pair) (_ storage.Servicer, _ storage.Storager, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Type: Type, Err: err, Pairs: pairs}
		}
	}()

	srv := &Service{
		noRedirectClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}

	opt, err := parseServicePairNew(pairs...)
	if err != nil {
		return nil, nil, err
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, nil, services.NewPairUnsupportedError(ps.WithCredential(opt.Credential))
	}
	cfg, err := config.New(cred[0], cred[1])
	if err != nil {
		return nil, nil, err
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
	if err != nil {
		if e := services.NewPairRequiredError(); errors.As(err, &e) {
			return srv, nil, nil
		}
		return nil, nil, err
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
		return "", services.NewPairUnsupportedError(ps.WithStorageClass(in))
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in string) storageclass.Type {
	switch in {
	case storageClassStandard, "":
		return storageclass.Hot
	case storageClassStandardIA:
		return storageclass.Warm
	default:
		return ""
	}
}
