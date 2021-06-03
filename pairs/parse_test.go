package pairs

import (
	"fmt"
	"math"
	"testing"

	"github.com/beyondstorage/go-storage/v4/services"
	. "github.com/beyondstorage/go-storage/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestNotRegistered(t *testing.T) {
	_, err := ParseMap("not_a_type", nil)
	assert.ErrorIs(t, err, services.ErrServiceNotRegistered)
}

func TestParseGlobalPair(t *testing.T) {
	cases := []struct {
		name string
		k    string
		v    string
		pair Pair
		err  error
	}{
		// {"not parseable type", "interceptor", "", Pair{}, ErrPairTypeNotParsable},
		{"unknown pair", "not_a_pair", "", Pair{}, ErrPairNotRegistered},
		{"string", "name", "abc", WithName("abc"), nil},
		{"int64", "size", "114514", WithSize(114514), nil},
		{"int64 not a number", "size", "NAN", Pair{}, ErrPairValueInvalid},
		{"int64 out of range", "size", fmt.Sprint(math.MaxInt64) + "0", Pair{}, ErrPairValueInvalid},
		{"list mode number", "list_mode", fmt.Sprint(uint8(ListModeBlock | ListModePrefix)), WithListMode(ListModeBlock | ListModePrefix), nil},
		// {"list mode string", "list_mode", "  block  |prefix", WithListMode(ListModeBlock | ListModePrefix), nil},
		// {"list mode empty", "list_mode", " ", WithListMode(0), nil},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := globalPairMap.parse(tt.k, tt.v)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorIs(t, err, tt.err)
			}
			assert.Equal(t, tt.pair, p)
		})
	}
}
