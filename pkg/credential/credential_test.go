package credential

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name  string
		cfg   string
		value Provider
		err   error
	}{
		{
			"invalid string",
			"abcx",
			nil,
			ErrInvalidConfig,
		},
		{
			"normal static",
			"static:ak:sk",
			NewStatic("ak", "sk"),
			nil,
		},
		{
			"not supported protocol",
			"notsupported:ak:sk",
			nil,
			ErrUnsupportedProtocol,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Parse(tt.cfg)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.err))
			}
			assert.EqualValues(t, tt.value, p)
		})
	}
}
