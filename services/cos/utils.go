package cos

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
	// ref: https://cloud.tencent.com/document/product/436/7745
	storageClassHeader = "x-cos-storage-class"

	storageClassStandard   = "STANDARD"
	storageClassStandardIA = "STANDARD_IA"
	storageClassArchive    = "ARCHIVE"
)

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (string, error) {
	switch in {
	case storageclass.Cold:
		return storageClassArchive, nil
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
	case storageClassArchive:
		return storageclass.Cold, nil
	case storageClassStandardIA:
		return storageclass.Warm, nil
	// cos only return storage class while not standard, we should handle empty string
	case storageClassStandard, "":
		return storageclass.Hot, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}
