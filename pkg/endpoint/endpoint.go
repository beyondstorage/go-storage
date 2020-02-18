package endpoint

import (
	"fmt"
	"strconv"
	"strings"
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
	s := strings.Split(cfg, ":")
	// TODO: port could be emitted in https/http

	switch s[0] {
	case ProtocolHTTPS:
		port, err := strconv.ParseInt(s[2], 10, 64)
		if err != nil {
			return nil, &Error{"parse", ProtocolHTTPS, s[1:], err}
		}
		return NewHTTPS(s[1], int(port)), nil
	case ProtocolHTTP:
		port, err := strconv.ParseInt(s[2], 10, 64)
		if err != nil {
			return nil, &Error{"parse", ProtocolHTTP, s[1:], err}
		}
		return NewHTTP(s[1], int(port)), nil
	default:
		return nil, &Error{"parse", s[0], nil, ErrUnsupportedProtocol}
	}
}
