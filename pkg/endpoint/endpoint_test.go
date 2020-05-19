package endpoint

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue_String(t *testing.T) {
	t.Run("standard port", func(t *testing.T) {
		v := &Value{
			Protocol: "http",
			Host:     "example.com",
			Port:     80,
		}

		assert.Equal(t, "http://example.com", v.String())
	})
	t.Run("non-standard port", func(t *testing.T) {
		v := &Value{
			Protocol: "http",
			Host:     "example.com",
			Port:     8080,
		}

		assert.Equal(t, "http://example.com:8080", v.String())
	})
}

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
			ErrInvalidValue,
		},
		{
			"normal http",
			"http:example.com:80",
			NewHTTP("example.com", 80),
			nil,
		},
		{
			"normal http without port",
			"http:example.com",
			NewHTTP("example.com", 80),
			nil,
		},
		{
			"wrong port number in http",
			"http:example.com:xxx",
			nil,
			strconv.ErrSyntax,
		},
		{
			"normal https",
			"https:example.com:443",
			NewHTTPS("example.com", 443),
			nil,
		},
		{
			"normal https without port",
			"https:example.com",
			NewHTTPS("example.com", 443),
			nil,
		},
		{
			"wrong port number in https",
			"https:example.com:xxx",
			nil,
			strconv.ErrSyntax,
		},
		{
			"not supported protocol",
			"notsupported:abc.com",
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
				e := &strconv.NumError{}
				ok := errors.As(err, &e)
				if ok {
					assert.True(t, errors.Is(e.Err, tt.err))
				} else {
					assert.True(t, errors.Is(err, tt.err))
				}
				// FIXME: strconv.NumError doesn't implement Unwrap for now, waiting for go1.14
				// assert.True(t, errors.Is(err, tt.err))
			}
			assert.EqualValues(t, tt.value, p)
		})
	}
}
