package s3

import (
	"fmt"
	"strings"

	"github.com/Xuanwo/storage/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
