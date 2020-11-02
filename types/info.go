package types

// NewStorageMeta will create a new StorageMeta metadata.
func NewStorageMeta() *StorageMeta {
	return &StorageMeta{
		m: make(map[string]interface{}),
	}
}

// NewStorageStatistic will create a new StorageMeta statistic.
func NewStorageStatistic() *StorageStatistic {
	return &StorageStatistic{
		m: make(map[string]interface{}),
	}
}
