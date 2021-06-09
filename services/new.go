package services

import (
	"sync"

	"github.com/beyondstorage/go-storage/v4/types"
)

type (
	// NewServicerFunc is a function that can initiate a new servicer.
	NewServicerFunc func(ps ...types.Pair) (types.Servicer, error)
	// NewStoragerFunc is a function that can initiate a new storager.
	NewStoragerFunc func(ps ...types.Pair) (types.Storager, error)
)

var (
	serviceFnMap map[string]NewServicerFunc
	serviceLock  sync.Mutex

	storagerFnMap map[string]NewStoragerFunc
	storagerLock  sync.Mutex
)

// RegisterServicer will register a servicer.
func RegisterServicer(ty string, fn NewServicerFunc) {
	serviceLock.Lock()
	defer serviceLock.Unlock()

	serviceFnMap[ty] = fn
}

// NewServicer will initiate a new servicer.
func NewServicer(ty string, ps ...types.Pair) (types.Servicer, error) {
	serviceLock.Lock()
	defer serviceLock.Unlock()

	fn, ok := serviceFnMap[ty]
	if !ok {
		return nil, InitError{Op: "new_servicer", Type: ty, Err: ErrServiceNotRegistered, Pairs: ps}
	}

	return fn(ps...)
}

// RegisterStorager will register a storager.
func RegisterStorager(ty string, fn NewStoragerFunc) {
	storagerLock.Lock()
	defer storagerLock.Unlock()

	storagerFnMap[ty] = fn
}

// NewStorager will initiate a new storager.
func NewStorager(ty string, ps ...types.Pair) (types.Storager, error) {
	storagerLock.Lock()
	defer storagerLock.Unlock()

	fn, ok := storagerFnMap[ty]
	if !ok {
		return nil, InitError{Op: "new_storager", Type: ty, Err: ErrServiceNotRegistered, Pairs: ps}
	}

	return fn(ps...)
}

func init() {
	serviceFnMap = make(map[string]NewServicerFunc)
	storagerFnMap = make(map[string]NewStoragerFunc)
}

// PairMap is a map of the types of pairs.
type PairMap map[string]string

var (
	globalPairMap   PairMap
	servicePairMaps map[string]PairMap
)

func RegisterGlobalPairMap(m PairMap) {
	globalPairMap = m
}

// RegisterServicePairMap will register a service's pair map.
func RegisterServicePairMap(ty string, m PairMap) {
	for pair := range m {
		if _, ok := globalPairMap[pair]; ok {
			panic("service pair name is duplicate with global pair: " + pair)
		}
	}
	servicePairMaps[ty] = m
}

func NewServicerFromMap(ty string, m map[string]string) (types.Servicer, error) {
	ps, err := parseMap(ty, m)
	if err != nil {
		return nil, InitError{Op: "new_servicer", Type: ty, Err: err, Pairs: ps}
	}
	return NewServicer(ty, ps...)
}

func NewStoragerFromMap(ty string, m map[string]string) (types.Storager, error) {
	ps, err := parseMap(ty, m)
	if err != nil {
		return nil, InitError{Op: "new_storager", Type: ty, Err: err, Pairs: ps}
	}
	return NewStorager(ty, ps...)
}

func init() {
	servicePairMaps = make(map[string]PairMap)
}
