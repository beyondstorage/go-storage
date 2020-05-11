/*
Package info intend to provide all available info from storage/object.
*/
package info

// StorageMeta is the static metadata for StorageMeta.
type StorageMeta struct {
	Name    string
	WorkDir string

	m map[string]interface{}
}

// NewStorageMeta will create a new StorageMeta metadata.
func NewStorageMeta() StorageMeta {
	return StorageMeta{
		m: make(map[string]interface{}),
	}
}

// StorageStatistic is the statistic metadata for StorageMeta.
type StorageStatistic struct {
	m map[string]interface{}
}

// NewStorageStatistic will create a new StorageMeta statistic.
func NewStorageStatistic() StorageStatistic {
	return StorageStatistic{
		m: make(map[string]interface{}),
	}
}

// ObjectMeta is the metadata for ObjectMeta.
type ObjectMeta struct {
	m map[string]interface{}
}

// NewObjectMeta will create a new ObjectMeta metadata.
func NewObjectMeta() ObjectMeta {
	return ObjectMeta{
		m: make(map[string]interface{}),
	}
}

// Get will get meta from object meta.
func (m ObjectMeta) Get(key string) (interface{}, bool) {
	v, ok := m.m[key]
	if !ok {
		return nil, false
	}
	return v, true
}

// Set will get meta from object meta.
func (m ObjectMeta) Set(key string, value interface{}) ObjectMeta {
	m.m[key] = value
	return m
}
