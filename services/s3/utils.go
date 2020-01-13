package s3

import (
	"fmt"
	"strings"

	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handleS3Error(err error) error {
	if err == nil {
		panic("error must not be nil")
	}

	var e awserr.Error
	e, ok := err.(awserr.Error)
	if !ok {
		return fmt.Errorf("%w: %v", types.ErrUnhandledError, err)
	}

	switch e.Code() {
	default:
		return fmt.Errorf("%w: %v", types.ErrUnhandledError, err)
	}
}

func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
}

func (s *Storage) getRelPath(path string) string {
	return strings.TrimPrefix(path, s.workDir+"/")
}

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (string, error) {
	switch in {
	case storageclass.Hot:
		return s3.ObjectStorageClassStandard, nil
	case storageclass.Warm:
		return s3.ObjectStorageClassStandardIa, nil
	case storageclass.Cold:
		return s3.ObjectStorageClassGlacier, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in string) (storageclass.Type, error) {
	switch in {
	case s3.ObjectStorageClassStandard:
		return storageclass.Hot, nil
	case s3.ObjectStorageClassStandardIa:
		return storageclass.Warm, nil
	case s3.ObjectStorageClassGlacier:
		return storageclass.Cold, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}

}
