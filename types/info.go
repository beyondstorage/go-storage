package types

// StorageMeta is the static metadata for StorageMeta.
type StorageMeta storageMeta

// NewStorageMeta will create a new StorageMeta metadata.
func NewStorageMeta() StorageMeta {
	return StorageMeta{
		m: make(map[string]interface{}),
	}
}

// StorageStatistic is the statistic metadata for StorageMeta.
type StorageStatistic storageStatistic

// NewStorageStatistic will create a new StorageMeta statistic.
func NewStorageStatistic() StorageStatistic {
	return StorageStatistic{
		m: make(map[string]interface{}),
	}
}
