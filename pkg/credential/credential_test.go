package credential

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	protocol := uuid.New().String()
	args := []string{uuid.New().String(), uuid.New().String()}

	p := &Provider{protocol: protocol, args: args}

	assert.Equal(t, protocol, p.Protocol())
	assert.EqualValues(t, args, p.Value())
}

func TestParse(t *testing.T) {
	cases := []struct {
		name  string
		cfg   string
		value *Provider
		err   error
	}{
		{
			"hmac",
			"hmac:ak:sk",
			&Provider{protocol: ProtocolHmac, args: []string{"ak", "sk"}},
			nil,
		},
		{
			"api key",
			"apikey:key",
			&Provider{protocol: ProtocolAPIKey, args: []string{"key"}},
			nil,
		},
		{
			"file",
			"file:/path/to/file",
			&Provider{protocol: ProtocolFile, args: []string{"/path/to/file"}},
			nil,
		},
		{
			"env",
			"env",
			&Provider{protocol: ProtocolEnv},
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

func TestNewHmac(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		value *Provider
		err   error
	}{
		{
			"normal",
			[]string{"ak", "sk"},
			&Provider{ProtocolHmac, []string{"ak", "sk"}},
			nil,
		},
		{
			"invalid",
			[]string{"ak", "sk", "xxxx"},
			nil,
			ErrInvalidValue,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewHmac(tt.input...)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.err))
			}
			assert.EqualValues(t, tt.value, p)
		})
	}
}

func TestMustNewHmac(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		panic bool
	}{
		{
			"normal",
			[]string{"ak", "sk"},
			false,
		},
		{
			"invalid",
			[]string{"ak", "sk", "xxxx"},
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				assert.Panics(t, func() {
					MustNewHmac(tt.input...)
				})
			} else {
				assert.NotPanics(t, func() {
					MustNewHmac(tt.input...)
				})
			}
		})
	}
}

func TestNewAPIKey(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		value *Provider
		err   error
	}{
		{
			"normal",
			[]string{"key"},
			&Provider{ProtocolAPIKey, []string{"key"}},
			nil,
		},
		{
			"invalid",
			[]string{"ak", "sk", "xxxx"},
			nil,
			ErrInvalidValue,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewAPIKey(tt.input...)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.err))
			}
			assert.EqualValues(t, tt.value, p)
		})
	}
}

func TestMustNewAPIKey(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		panic bool
	}{
		{
			"normal",
			[]string{"key"},
			false,
		},
		{
			"invalid",
			[]string{"ak", "sk", "xxxx"},
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				assert.Panics(t, func() {
					MustNewAPIKey(tt.input...)
				})
			} else {
				assert.NotPanics(t, func() {
					MustNewAPIKey(tt.input...)
				})
			}
		})
	}
}

func TestNewFile(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		value *Provider
		err   error
	}{
		{
			"normal",
			[]string{"/path/to/file"},
			&Provider{ProtocolFile, []string{"/path/to/file"}},
			nil,
		},
		{
			"invalid",
			[]string{"ak", "sk", "xxxx"},
			nil,
			ErrInvalidValue,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewFile(tt.input...)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.err))
			}
			assert.EqualValues(t, tt.value, p)
		})
	}
}

func TestMustNewFile(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		panic bool
	}{
		{
			"normal",
			[]string{"/path/to/file"},
			false,
		},
		{
			"invalid",
			[]string{"ak", "sk", "xxxx"},
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				assert.Panics(t, func() {
					MustNewFile(tt.input...)
				})
			} else {
				assert.NotPanics(t, func() {
					MustNewFile(tt.input...)
				})
			}
		})
	}
}

func TestNewMev(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		value *Provider
		err   error
	}{
		{
			"normal",
			[]string{""},
			&Provider{ProtocolEnv, nil},
			nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewEnv(tt.input...)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.err))
			}
			assert.EqualValues(t, tt.value, p)
		})
	}
}

func TestMustNewEnv(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		panic bool
	}{
		{
			"normal",
			[]string{},
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				assert.Panics(t, func() {
					MustNewEnv(tt.input...)
				})
			} else {
				assert.NotPanics(t, func() {
					MustNewEnv(tt.input...)
				})
			}
		})
	}
}
