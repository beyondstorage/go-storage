package kodo

import (
	"strings"
	"time"

	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/types"
)

func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
}

func (s *Storage) getRelPath(path string) string {
	return strings.TrimPrefix(path, s.workDir+"/")
}

func convertUnixTimestampToTime(v int64) time.Time {
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(v, 0)
}

const (
	// ref: https://developer.qiniu.com/kodo/api/3710/chtype
	storageClassStandard   = 0
	storageClassStandardIA = 1
	storageClassArchive    = 2
)

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (int, error) {
	switch in {
	case storageclass.Hot:
		return storageClassStandard, nil
	case storageclass.Warm:
		return storageClassStandardIA, nil
	case storageclass.Cold:
		return storageClassArchive, nil
	default:
		return 0, types.ErrStorageClassNotSupported
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in int) (storageclass.Type, error) {
	switch in {
	case 0:
		return storageclass.Hot, nil
	case 1:
		return storageclass.Warm, nil
	case 2:
		return storageclass.Cold, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}
