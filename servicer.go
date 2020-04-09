package storage

import (
	"context"

	"github.com/Xuanwo/storage/types"
)

/*
Servicer can maintain multipart storage services.

Implementer can choose to implement this interface or not.
*/
type Servicer interface {
	// String will implement Stringer.
	String() string

	// List will list all storager instances under this service.
	List(pairs ...*types.Pair) (err error)
	// ListWithContext will list all storager instances under this service.
	ListWithContext(ctx context.Context, pairs ...*types.Pair) (err error)
	// Get will get a valid storager instance for service.
	Get(name string, pairs ...*types.Pair) (Storager, error)
	// GetWithContext will get a valid storager instance for service.
	GetWithContext(ctx context.Context, name string, pairs ...*types.Pair) (Storager, error)
	// Create will create a new storager instance.
	Create(name string, pairs ...*types.Pair) (Storager, error)
	// CreateWithContext will create a new storager instance.
	CreateWithContext(ctx context.Context, name string, pairs ...*types.Pair) (Storager, error)
	// Delete will delete a storager instance.
	Delete(name string, pairs ...*types.Pair) (err error)
	// DeleteWithContext will delete a storager instance.
	DeleteWithContext(ctx context.Context, name string, pairs ...*types.Pair) (err error)
}
