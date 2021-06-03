package pairs

import (
	"fmt"
	"reflect"

	"github.com/beyondstorage/go-storage/v4/services"
	. "github.com/beyondstorage/go-storage/v4/types"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

var (
	// ErrPairNotRegistered means the pair is not registered.
	ErrPairNotRegistered = services.NewErrorCode("pair not registered")
	// ErrPairValueInvalid means the pair's value is invalid.
	ErrPairValueInvalid = services.NewErrorCode("pair value invalid")
)

// Parse a slice of Pairs from a map.
func ParseMap(ty string, m map[string]string) (ps []Pair, err error) {
	servicePairMap, ok := servicePairMaps[ty]
	if !ok {
		err = fmt.Errorf("%w: %v", services.ErrServiceNotRegistered, ty)
		return
	}

	for k, v := range m {
		var pair Pair
		if _, ok := globalPairMap[k]; ok {
			pair, err = globalPairMap.parse(k, v)
		} else {
			pair, err = servicePairMap.parse(k, v)
			pair.Key = ty + "_" + pair.Key
		}
		if err != nil {
			return nil, err
		}
		ps = append(ps, pair)
	}
	return
}

// Parse a Pair.
func Parse(ty string, k string, v string) (pair Pair, err error) {
	servicePairMap, ok := servicePairMaps[ty]
	if !ok {
		err = fmt.Errorf("%w: %v", services.ErrServiceNotRegistered, ty)
		return
	}

	if _, ok := globalPairMap[k]; ok {
		pair, err = globalPairMap.parse(k, v)
	} else {
		pair, err = servicePairMap.parse(k, v)
		pair.Key = ty + "_" + pair.Key
	}
	if err != nil {
		pair = Pair{}
	}
	return
}

func (m PairMap) parse(k string, v string) (pair Pair, err error) {
	pairInfo, ok := m[k]
	if !ok {
		err = fmt.Errorf("%w: %v", ErrPairNotRegistered, k)
		return Pair{}, err
	}

	value := pairInfo.New()
	err = decoder.Decode(value, map[string][]string{"value": {v}})
	if err != nil {
		return Pair{}, fmt.Errorf("%w: %v", ErrPairValueInvalid, err)
	}
	return Pair{Key: k, Value: reflect.ValueOf(value).Elem().FieldByName("Value").Interface()}, nil
}
