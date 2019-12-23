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

// Provider can provide all authenticate needed info.
type Provider interface {
	Value() Value
}

// Value is the credential value.
type Value struct {
	AccessKey string
	SecretKey string
}

// Parse will parse config string to create a credential Provider.
func Parse(cfg string) (Provider, error) {
	errorMessage := "parse credential config [%s]: <%w>"

	s := strings.Split(cfg, ":")
	if len(s) < 2 {
		return nil, fmt.Errorf(errorMessage, cfg, ErrInvalidConfig)
	}

	switch s[0] {
	case ProtocolStatic:
		return NewStatic(s[1], s[2]), nil
	default:
		return nil, fmt.Errorf(errorMessage, cfg, ErrUnsupportedProtocol)
	}
}
