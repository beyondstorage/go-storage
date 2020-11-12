package types

// NewStorageMeta will create a new StorageMeta.
func NewStorageMeta() *StorageMeta {
	return &StorageMeta{
		m: make(map[string]interface{}),
	}
}

// NewStorageStatistic will create a new StorageStatistic.
func NewStorageStatistic() *StorageStatistic {
	return &StorageStatistic{
		m: make(map[string]interface{}),
	}
}
