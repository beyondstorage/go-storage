package endpoint

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrInvalidConfig will return if config string invalid.
	ErrInvalidConfig = errors.New("invalid config")
	// ErrUnsupportedProtocol will return if protocol is unsupported.
	ErrUnsupportedProtocol = errors.New("unsupported protocol")
)

// Provider will return all info needed to connect a service.
type Provider interface {
	Value() Value
}

// Value is the required info to connect a service.
type Value struct {
	Protocol string
	Host     string
	Port     int
}

// String will compose all info into a valid URL.
func (v Value) String() string {
	return fmt.Sprintf("%s://%s:%d", v.Protocol, v.Host, v.Port)
}

// Parse will parse config string to create a endpoint Provider.
func Parse(cfg string) (Provider, error) {
	errorMessage := "parse credential config [%s]: <%w>"

	s := strings.Split(cfg, ":")
	if len(s) < 2 {
		return nil, fmt.Errorf(errorMessage, cfg, ErrInvalidConfig)
	}

	switch s[0] {
	case ProtocolHTTPS:
		port, err := strconv.ParseInt(s[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf(errorMessage, cfg, err)
		}
		return NewHTTPS(s[1], int(port)), nil
	case ProtocolHTTP:
		port, err := strconv.ParseInt(s[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf(errorMessage, cfg, err)
		}
		return NewHTTP(s[1], int(port)), nil
	default:
		return nil, fmt.Errorf(errorMessage, cfg, ErrUnsupportedProtocol)
	}
}
