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

	p := Provider{protocol: protocol, args: args}

	assert.Equal(t, protocol, p.Protocol())
	assert.EqualValues(t, args, p.Value())
}

func TestParse(t *testing.T) {
	cases := []struct {
		name  string
		cfg   string
		value Provider
		err   error
	}{
		{
			"hmac",
			"hmac:ak:sk",
			Provider{protocol: ProtocolHmac, args: []string{"ak", "sk"}},
			nil,
		},
		{
			"api key",
			"apikey:key",
			Provider{protocol: ProtocolAPIKey, args: []string{"key"}},
			nil,
		},
		{
			"file",
			"file:/path/to/file",
			Provider{protocol: ProtocolFile, args: []string{"/path/to/file"}},
			nil,
		},
		{
			"env",
			"env",
			Provider{protocol: ProtocolEnv},
			nil,
		},
		{
			"base64",
			"base64:aGVsbG8sd29ybGQhCg==",
			Provider{protocol: ProtocolBase64, args: []string{"aGVsbG8sd29ybGQhCg=="}},
			nil,
		},
		{
			"not supported protocol",
			"notsupported:ak:sk",
			Provider{},
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
