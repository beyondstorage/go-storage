package azblob

import (
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/types"
)

func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
}

func (s *Storage) getRelPath(path string) string {
	return strings.TrimPrefix(path, s.workDir+"/")
}

// parseStorageClass will parse storageclass.Type into service independent storage class type.
func parseStorageClass(in storageclass.Type) (azblob.AccessTierType, error) {
	switch in {
	case storageclass.Cold:
		return azblob.AccessTierArchive, nil
	case storageclass.Hot:
		return azblob.AccessTierHot, nil
	case storageclass.Warm:
		return azblob.AccessTierCool, nil
	default:
		return azblob.AccessTierType(""), types.ErrStorageClassNotSupported
	}
}

// formatStorageClass will format service independent storage class type into storageclass.Type.
func formatStorageClass(in azblob.AccessTierType) (storageclass.Type, error) {
	switch in {
	case azblob.AccessTierArchive:
		return storageclass.Cold, nil
	case azblob.AccessTierCool:
		return storageclass.Warm, nil
	case azblob.AccessTierHot:
		return storageclass.Hot, nil
	default:
		return "", types.ErrStorageClassNotSupported
	}
}
