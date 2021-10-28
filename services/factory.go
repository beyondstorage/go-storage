package services

import (
	"errors"
	"fmt"
	"strings"

	"go.beyondstorage.io/v5/types"
)

// Factory is used to initialize a new service or storage.
//
// We will generate a Factory struct which implement this interface for each service.
type Factory interface {
	// FromString fill factory with data parsed from connection string.
	//
	// The connection string will be parsed by the following format:
	//
	//     s3://<credential>@<endpoint>/<name>/<work_dir>?force_path_style
	FromString(conn string) (err error)
	// FromMap fill factory with data parsed from map.
	//
	// The map will be parsed by the following format:
	//
	//   {
	//		"credential": <credential>,
	//		"endpoint": <endpoint>,
	//		"name": <name>,
	//		"work_dir": <work_dir>,
	//		"force_path_style": true
	//   }
	FromMap(m map[string]interface{}) (err error)
	// WithPairs fill factory with data parsed from key-value pairs.
	WithPairs(ps ...types.Pair) (err error)

	// NewServicer will create a new service via already initialized factory.
	//
	// Service should implement `newService() (*Service, error)`
	//
	// It's possible that the factory only support init Storager, but not init Servicer.
	// We will generate an error for it.
	NewServicer() (srv types.Servicer, err error)
	// NewStorager will create a new storage via already initialized factory.
	//
	// Service should implement `newStorage() (*Storage, error)`
	//
	// It's possible that the factory is OK to NewServicer but not OK to NewStorager.
	// We will generate an error for it.
	NewStorager() (sto types.Storager, err error)
}

// factoryRegistry is the registry of all supported services.
var factoryRegistry = make(map[string]Factory)

// RegisterFactory is used to register a new service.
//
// NOTE:
//   - This function is not for public use, it should only be called in service init() function.
//   - This function is not concurrent-safe.
func RegisterFactory(ty string, f Factory) {
	factoryRegistry[ty] = f
}

// NewFactory will create a new factory by service type.
func NewFactory(ty string, ps ...types.Pair) (Factory, error) {
	f, ok := factoryRegistry[ty]
	if !ok {
		return nil, InitError{Op: "new_factory", Type: ty, Err: ErrServiceNotRegistered}
	}

	if len(ps) > 0 {
		err := f.WithPairs(ps...)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

// NewFactoryFromString will create a new factory by service type and connection string.
func NewFactoryFromString(conn string, ps ...types.Pair) (Factory, error) {
	ty, value, err := parseConn(conn)
	if err != nil {
		return nil, InitError{Op: "parse_conn", Type: ty, Err: err, Pairs: ps}
	}

	f, err := NewFactory(ty)
	if err != nil {
		return nil, err
	}

	err = f.FromString(value)
	if err != nil {
		return nil, err
	}

	err = f.WithPairs(ps...)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// NewServicer will initiate a new servicer.
func NewServicer(ty string, ps ...types.Pair) (types.Servicer, error) {
	f, err := NewFactory(ty, ps...)
	if err == nil {
		return f.NewServicer()
	}
	if !errors.Is(err, ErrServiceNotRegistered) {
		return nil, err
	}

	// If the factory is not found, we will try to use the compilable logic.
	fn, ok := servicerFnMap[ty]
	if !ok {
		return nil, InitError{Op: "new_servicer", Type: ty, Err: ErrServiceNotRegistered, Pairs: ps}
	}

	return fn(ps...)
}

// NewServicerFromString will create a new service via connection string.
func NewServicerFromString(conn string, ps ...types.Pair) (types.Servicer, error) {
	f, err := NewFactoryFromString(conn, ps...)
	if err == nil {
		return f.NewServicer()
	}
	if !errors.Is(err, ErrServiceNotRegistered) {
		return nil, err
	}

	// If the factory is not found, we will try to use the compilable logic.
	ty, psc, err := parseConnectionString(conn)
	if err != nil {
		return nil, InitError{Op: "new_servicer", Type: ty, Err: err, Pairs: ps}
	}
	// Append ps after connection string to keep pairs order.
	psc = append(psc, ps...)
	return NewServicer(ty, psc...)
}

// NewStorager will initiate a new storager.
func NewStorager(ty string, ps ...types.Pair) (types.Storager, error) {
	f, err := NewFactory(ty, ps...)
	if err == nil {
		return f.NewStorager()
	}
	if !errors.Is(err, ErrServiceNotRegistered) {
		return nil, err
	}

	// If the factory is not found, we will try to use the compilable logic.
	fn, ok := storagerFnMap[ty]
	if !ok {
		return nil, InitError{Op: "new_storager", Type: ty, Err: ErrServiceNotRegistered, Pairs: ps}
	}

	return fn(ps...)
}

// NewStoragerFromString will create a new storager via connection string.
func NewStoragerFromString(conn string, ps ...types.Pair) (types.Storager, error) {
	f, err := NewFactoryFromString(conn, ps...)
	if err == nil {
		return f.NewStorager()
	}
	if !errors.Is(err, ErrServiceNotRegistered) {
		return nil, err
	}

	// If the factory is not found, we will try to use the compilable logic.
	ty, psc, err := parseConnectionString(conn)
	if err != nil {
		return nil, InitError{Op: "new_storager", Type: ty, Err: err, Pairs: ps}
	}
	// Append ps after connection string to keep pairs order.
	psc = append(psc, ps...)
	return NewStorager(ty, psc...)
}

func parseConn(conn string) (ty, value string, err error) {
	colon := strings.Index(conn, ":")
	if colon == -1 {
		err = fmt.Errorf("%w: %s",
			NewErrorCode("connection string is invalid"), conn)
		return
	}
	return conn[:colon], conn[colon+1:], nil
}
