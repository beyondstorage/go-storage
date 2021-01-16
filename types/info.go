package types

// NewStorageMeta will create a new StorageMeta.
func NewStorageMeta() *StorageMeta {
	return &StorageMeta{
		m: make(map[string]interface{}),
	}
}
