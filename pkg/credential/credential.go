package credential

import (
	"strings"
)

const (
	// ProtocolHmac will hold access key and secret key credential.
	//
	// HMAC means hash-based message authentication code, it may be inaccurate to represent credential
	// protocol ak/sk(access key + secret key with hmac), but it's simple and no confuse with other
	// protocol, so just keep this.
	//
	// value = [Access Key, Secret Key]
	ProtocolHmac = "hmac"
	// ProtocolAPIKey will hold api key credential.
	//
	// value = [API Key]
	ProtocolAPIKey = "apikey"
	// ProtocolFile will hold file credential.
	//
	// value = [File Path], service decide how to use this file
	ProtocolFile = "file"
	// ProtocolEnv will represent credential from env.
	//
	// value = [], service retrieves credential value from env.
	ProtocolEnv = "env"
	// ProtocolBase64 will represents credential binary data in base64
	//
	// Storage service like gcs will take token files as input, we provide base64 protocol so that user
	// can pass token binary data directly.
	ProtocolBase64 = "base64"
)

// Provider will provide credential protocol and values.
type Provider struct {
	protocol string
	args     []string
}

// Protocol provides current credential's protocol.
func (p *Provider) Protocol() string {
	return p.protocol
}

// Value provides current credential's value in string array.
func (p *Provider) Value() []string {
	return p.args
}

// Parse will parse config string to create a credential Provider.
func Parse(cfg string) (*Provider, error) {
	s := strings.Split(cfg, ":")

	switch s[0] {
	case ProtocolHmac:
		return NewHmac(s[1:]...)
	case ProtocolAPIKey:
		return NewAPIKey(s[1:]...)
	case ProtocolFile:
		return NewFile(s[1:]...)
	case ProtocolEnv:
		return NewEnv()
	case ProtocolBase64:
		return NewBase64(s[1:]...)
	default:
		return nil, &Error{"parse", ErrUnsupportedProtocol, s[0], nil}
	}
}

// NewHmac create a hmac provider.
func NewHmac(value ...string) (*Provider, error) {
	if len(value) != 2 {
		return nil, &Error{"new", ErrInvalidValue, ProtocolHmac, value}
	}
	return &Provider{ProtocolHmac, []string{value[0], value[1]}}, nil
}

// MustNewHmac make sure Provider must be created if no panic happened.
func MustNewHmac(value ...string) *Provider {
	p, err := NewHmac(value...)
	if err != nil {
		panic(err)
	}
	return p
}

// NewAPIKey create a api key provider.
func NewAPIKey(value ...string) (*Provider, error) {
	if len(value) != 1 {
		return nil, &Error{"new", ErrInvalidValue, ProtocolAPIKey, value}
	}
	return &Provider{ProtocolAPIKey, []string{value[0]}}, nil
}

// MustNewAPIKey make sure Provider must be created if no panic happened.
func MustNewAPIKey(value ...string) *Provider {
	p, err := NewAPIKey(value...)
	if err != nil {
		panic(err)
	}
	return p
}

// NewFile create a file provider.
func NewFile(value ...string) (*Provider, error) {
	if len(value) != 1 {
		return nil, &Error{"new", ErrInvalidValue, ProtocolFile, value}
	}
	return &Provider{ProtocolFile, []string{value[0]}}, nil
}

// MustNewFile make sure Provider must be created if no panic happened.
func MustNewFile(value ...string) *Provider {
	p, err := NewFile(value...)
	if err != nil {
		panic(err)
	}
	return p
}

// NewEnv create a env provider.
func NewEnv(_ ...string) (*Provider, error) {
	return &Provider{ProtocolEnv, nil}, nil
}

// MustNewEnv make sure Provider must be created if no panic happened.
func MustNewEnv(value ...string) *Provider {
	p, _ := NewEnv(value...)
	return p
}

// NewBase64 create a base64 provider.
func NewBase64(value ...string) (*Provider, error) {
	if len(value) != 1 {
		return nil, &Error{"new", ErrInvalidValue, ProtocolFile, value}
	}
	return &Provider{ProtocolBase64, value}, nil
}

// MustNewBase64 make sure Provider must be created if no panic happened.
func MustNewBase64(value ...string) *Provider {
	p, err := NewBase64(value...)
	if err != nil {
		panic(err)
	}
	return p
}
