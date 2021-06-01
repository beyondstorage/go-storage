package pairs

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/beyondstorage/go-storage/v4/services"
	. "github.com/beyondstorage/go-storage/v4/types"
)

var (
	// ErrPairTypeNotParsable means the pair's type is not parseable.
	ErrPairTypeNotParsable = services.NewErrorCode("type not parseable")
	// ErrPairNotRegistered means the pair is not registered.
	ErrPairNotRegistered = services.NewErrorCode("pair not registered")
	// ErrPairValueInvalid means the pair's value is invalid.
	ErrPairValueInvalid = services.NewErrorCode("pair value invalid")
)

// parseFunc is a function that can parse a string into a given type.
type parseFunc func(string) (interface{}, error)

var parseMap map[reflect.Type]parseFunc

// Parse a slice of Pairs from a map.
func ParseMap(ty string, m map[string]string) (ps []Pair, err error) {
	servicePairMap, ok := servicePairMaps[ty]
	if !ok {
		err = fmt.Errorf("%w: %v", services.ErrServiceNotRegistered, ty)
		return
	}

	for k, v := range m {
		pair := Pair{Key: k}
		if _, ok := globalPairMap[k]; ok {
			pair.Value, err = globalPairMap.parse(k, v)
		} else {
			pair.Value, err = servicePairMap.parse(k, v)
		}
		if err != nil {
			return
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
		return globalPairMap.parse(k, v)
	} else {
		return servicePairMap.parse(k, v)
	}
}

func (m PairMap) parse(k string, v string) (pair Pair, err error) {
	vType, ok := m[k]
	if !ok {
		err = fmt.Errorf("%w: %v", ErrPairNotRegistered, k)
		return Pair{}, err
	}
	parseFn, ok := parseMap[vType]
	if !ok {
		err = fmt.Errorf("%w: %v, %v", ErrPairTypeNotParsable, k, vType)
		return Pair{}, err
	}
	pair.Key = k
	pair.Value, err = parseFn(v)
	if err != nil {
		err = fmt.Errorf("%w: %v, %v, %v: %v", ErrPairValueInvalid, k, vType, v, err)
		return Pair{}, err
	}
	return
}

var (
	stringType = reflect.TypeOf("")
	boolType   = reflect.TypeOf(true)

	int64Type = reflect.TypeOf(int64(0))
	int32Type = reflect.TypeOf(int32(0))
	int16Type = reflect.TypeOf(int16(0))
	int8Type  = reflect.TypeOf(int8(0))
	intType   = reflect.TypeOf(int(0))

	uint64Type = reflect.TypeOf(uint64(0))
	uint32Type = reflect.TypeOf(uint32(0))
	uint16Type = reflect.TypeOf(uint16(0))
	uint8Type  = reflect.TypeOf(uint8(0))
	uintType   = reflect.TypeOf(uint(0))

	listModeType = reflect.TypeOf(ListModeBlock)
)

func parseString(s string) (interface{}, error) {
	return s, nil
}

func parseBool(s string) (interface{}, error) {
	return strconv.ParseBool(s)
}

// Support base prefix: 2 for "0b", 8 for "0" or "0o", 16 for "0x", and 10 otherwise
func mkParseInt(bitSize int) func(s string) (interface{}, error) {
	return func(s string) (interface{}, error) {
		var v interface{}
		v64, err := strconv.ParseInt(s, 0, bitSize)
		if err != nil {
			return nil, err
		}
		switch bitSize {
		case 64:
			v = int64(v64)
		case 32:
			v = int32(v64)
		case 16:
			v = int16(v64)
		case 8:
			v = int8(v64)
		case 0:
			v = int(v64)
		}
		return v, err
	}
}

func mkParseUint(bitSize int) func(s string) (interface{}, error) {
	return func(s string) (interface{}, error) {
		var v interface{}
		v64, err := strconv.ParseUint(s, 0, bitSize)
		if err != nil {
			return nil, err
		}
		switch bitSize {
		case 64:
			v = uint64(v64)
		case 32:
			v = uint32(v64)
		case 16:
			v = uint16(v64)
		case 8:
			v = uint8(v64)
		case 0:
			v = uint(v64)
		}
		return v, err
	}
}

func parseListMode(s string) (interface{}, error) {
	v64, err := strconv.ParseUint(s, 0, 8)
	if err == nil || errors.Is(err, strconv.ErrRange) {
		return ListMode(v64), nil
	}
	if strings.TrimSpace(s) == "" {
		return ListMode(0), nil
	}
	modes := strings.Split(s, "|")
	var l ListMode
	for _, mode := range modes {
		mode = strings.TrimSpace(mode)
		switch mode {
		case ListModeBlock.String():
			l |= ListModeBlock
		case ListModeDir.String():
			l |= ListModeDir
		case ListModePart.String():
			l |= ListModePart
		case ListModePrefix.String():
			l |= ListModePrefix
		default:
			return nil, errors.New("invalid mode " + mode)
		}
	}
	return l, nil
}

func init() {
	parseMap = map[reflect.Type]parseFunc{
		stringType: parseString,
		boolType:   parseBool,

		int64Type: mkParseInt(64),
		int32Type: mkParseInt(32),
		int16Type: mkParseInt(16),
		int8Type:  mkParseInt(8),
		intType:   mkParseInt(0),

		uint64Type: mkParseUint(64),
		uint32Type: mkParseUint(32),
		uint16Type: mkParseUint(16),
		uint8Type:  mkParseUint(8),
		uintType:   mkParseUint(0),

		listModeType: parseListMode,
	}
}
