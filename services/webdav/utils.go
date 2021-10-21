package webdav

import (
	"go.beyondstorage.io/v5/types"
)

// Storage is the example client.
type Storage struct {
	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	panic("implement me")
}

// NewStorager will create Storager only.
func NewStorager(pairs ...types.Pair) (types.Storager, error) {
	panic("implement me")
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	panic("implement me")
}
