package oss

import (
	"strings"

	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/types"
)

func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
}

func (s *Storage) getRelPath(path string) string {
	return strings.TrimPrefix(path, s.workDir+"/")
}

const (
	// ref: https://www.alibabacloud.com/help/doc-detail/31984.htm
	storageClassHeader = "x-oss-storage-class"

	// ref: https://www.alibabacloud.com/help/doc-detail/51374.htm
	storageClassStandard = "STANDARD"
	storageClassIA       = "IA"
	storageClassArchive  = "Archive"
)

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (string, error) {
	switch in {
	case storageclass.Hot:
		return storageClassStandard, nil
	case storageclass.Warm:
		return storageClassIA, nil
	case storageclass.Cold:
		return storageClassArchive, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in string) (storageclass.Type, error) {
	switch in {
	case storageClassStandard:
		return storageclass.Hot, nil
	case storageClassIA:
		return storageclass.Warm, nil
	case storageClassArchive:
		return storageclass.Cold, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}
