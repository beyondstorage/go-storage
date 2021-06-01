package pairs

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/beyondstorage/go-storage/v4/services"
	. "github.com/beyondstorage/go-storage/v4/types"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

func TestNotRegistered(t *testing.T) {
	_, err := ParseMap("not_a_type", nil)
	assert.ErrorIs(t, err, services.ErrServiceNotRegistered)
}

func TestParseInt(t *testing.T) {
	cases := []struct {
		name string
		s    string
		typ  reflect.Type
		v    interface{}
		err  error
	}{
		{"not a number", "NAN", intType, nil, strconv.ErrSyntax},
		{"int", "123", intType, 123, nil},
		{"int base 2", "0b10", intType, 2, nil},
		{"int base 8", "010", intType, 8, nil},
		{"int base 16", "0x10", intType, 16, nil},
		{"int8", "123", int8Type, int8(123), nil},
		{"int8 out of range", "12345", int8Type, nil, strconv.ErrRange},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			value, err := parseMap[tt.typ](tt.s)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorIs(t, err, tt.err)
			}
			assert.Equal(t, tt.v, value)
		})
	}
}

// Parse the string representation of a value should get an identical result.
func TestParseIdentical(t *testing.T) {
	f := fuzz.New()
	N := 100
	for typ, parser := range parseMap {
		t.Run(typ.String(), func(t *testing.T) {
			for i := 0; i < N; i++ {
				vPtr := reflect.New(typ)
				f.Fuzz(vPtr.Interface())
				v := vPtr.Elem().Interface()
				value, err := parser(fmt.Sprint(v))
				b := assert.Nil(t, err)
				// We check they have same type and string representation,
				// since in special cases like ListMode, two different value have same string representation.
				b = b && assert.Equal(t, reflect.TypeOf(v), reflect.TypeOf(value))
				b = b && assert.Equal(t, fmt.Sprint(v), fmt.Sprint(value))
				if !b {
					break
				}
			}
		})
	}
}

func TestParseGlobalPair(t *testing.T) {
	cases := []struct {
		name string
		k    string
		v    string
		pair Pair
		err  error
	}{
		{"not parseable type", "interceptor", "", Pair{}, ErrPairTypeNotParsable},
		{"unknown pair", "not_a_pair", "", Pair{}, ErrPairNotRegistered},
		{"string", "name", "abc", WithName("abc"), nil},
		{"int64", "size", "114514", WithSize(114514), nil},
		{"int64 not a number", "size", "NAN", Pair{}, ErrPairValueInvalid},
		{"int64 out of range", "size", fmt.Sprint(math.MaxInt64) + "0", Pair{}, ErrPairValueInvalid},
		{"list mode number", "list_mode", fmt.Sprint(uint8(ListModeBlock | ListModePrefix)), WithListMode(ListModeBlock | ListModePrefix), nil},
		{"list mode string", "list_mode", "  block  |prefix", WithListMode(ListModeBlock | ListModePrefix), nil},
		{"list mode empty", "list_mode", " ", WithListMode(0), nil},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := globalPairMap.parse(tt.k, tt.v)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.err))
			}
			assert.Equal(t, tt.pair, p)
		})
	}
}
