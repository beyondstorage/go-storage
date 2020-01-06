package qingstor

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Xuanwo/storage/types"
	qserror "github.com/yunify/qingstor-sdk-go/v3/request/errors"
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

func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
}

func (s *Storage) getRelPath(path string) string {
	return strings.TrimPrefix(path, s.workDir+"/")
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
