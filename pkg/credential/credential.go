package credential

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidConfig will return if config string invalid.
	ErrInvalidConfig = errors.New("invalid config")
	// ErrUnsupportedProtocol will return if protocol is unsupported.
	ErrUnsupportedProtocol = errors.New("unsupported protocol")
)

const (
	// ProtocolHmac will hold access key and secret key credential.
	//
	// HMAC means hash-based message authentication code, it may be inaccurate to represent credential
	// protocol ak/sk(access key + secret key with hmac), but it's simple and no confuse with other
	// protocol, so just keep this.
	ProtocolHmac = "hmac"
	// ProtocolAPIKey will hold api key credential.
	ProtocolAPIKey = "apikey"
	// ProtocolFile will hold file credential.
	ProtocolFile = "file"
	// ProtocolEnv will represent credential from env.
	ProtocolEnv = "env"
)

// Value is the credential value.
type Provider struct {
	protocol string
	args     []string
}

func (p *Provider) Protocol() string {
	return p.protocol
}

func (p *Provider) Value() []string {
	return p.args
}

// Parse will parse config string to create a credential Provider.
func Parse(cfg string) (*Provider, error) {
	errorMessage := "parse credential config [%s]: %w"

	s := strings.Split(cfg, ":")
	if len(s) < 2 {
		return nil, fmt.Errorf(errorMessage, cfg, ErrInvalidConfig)
	}

	switch s[0] {
	case ProtocolHmac:
		return NewHmac(s[1:]...)
	case ProtocolAPIKey:
		return NewAPIKey(s[1:]...)
	case ProtocolFile:
		return NewFile(s[1:]...)
	default:
		return nil, fmt.Errorf(errorMessage, cfg, ErrUnsupportedProtocol)
	}
}

func NewHmac(value ...string) (*Provider, error) {
	errorMessage := "parse hmac credential [%s]: %w"

	if len(value) != 2 {
		return nil, fmt.Errorf(errorMessage, value, ErrInvalidConfig)
	}
	return &Provider{ProtocolHmac, []string{value[0], value[1]}}, nil
}

func MustNewHmac(value ...string) *Provider {
	p, err := NewHmac(value...)
	if err != nil {
		panic(err)
	}
	return p
}

func NewAPIKey(value ...string) (*Provider, error) {
	errorMessage := "parse apikey credential [%s]: %w"

	if len(value) != 1 {
		return nil, fmt.Errorf(errorMessage, value, ErrInvalidConfig)
	}
	return &Provider{ProtocolAPIKey, []string{value[0]}}, nil
}

func MustNewAPIKey(value ...string) *Provider {
	p, err := NewAPIKey(value...)
	if err != nil {
		panic(err)
	}
	return p
}

func NewFile(value ...string) (*Provider, error) {
	errorMessage := "parse file credential [%s]: %w"

	if len(value) != 1 {
		return nil, fmt.Errorf(errorMessage, value, ErrInvalidConfig)
	}
	return &Provider{ProtocolFile, []string{value[0]}}, nil
}

func MustNewFile(value ...string) *Provider {
	p, err := NewFile(value...)
	if err != nil {
		panic(err)
	}
	return p
}
