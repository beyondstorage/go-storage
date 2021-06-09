package endpoint

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// ProtocolHTTPS is the https credential protocol.
	//
	// Deprecated: Moved to github.com/beyondstorage/go-endpoint
	ProtocolHTTPS = "https"
	// ProtocolHTTP is the http credential protocol.
	//
	// Deprecated: Moved to github.com/beyondstorage/go-endpoint
	ProtocolHTTP = "http"
)

// Value is the required info to connect a service.
//
// Deprecated: Moved to github.com/beyondstorage/go-endpoint
type Value struct {
	Protocol string
	Host     string
	Port     int
}

// String will compose all info into a valid URL.
func (v Value) String() string {
	if defaultPort[v.Protocol] == v.Port {
		return fmt.Sprintf("%s://%s", v.Protocol, v.Host)
	}
	return fmt.Sprintf("%s://%s:%d", v.Protocol, v.Host, v.Port)
}

// Parse will parse config string to create a endpoint Provider.
//
// Deprecated: Moved to github.com/beyondstorage/go-endpoint
func Parse(cfg string) (p Value, err error) {
	s := strings.Split(cfg, ":")
	if len(s) < 2 {
		return Value{}, &Error{"parse", ErrInvalidValue, s[0], nil}
	}

	defer func() {
		if err != nil {
			err = &Error{"parse", err, s[0], s[1:]}
		}
	}()

	var port int
	if len(s) >= 3 {
		xport, err := strconv.ParseInt(s[2], 10, 64)
		if err != nil {
			return Value{}, err
		}
		port = int(xport)
	}

	switch s[0] {
	case ProtocolHTTPS:
		if port == 0 {
			port = defaultPort[ProtocolHTTPS]
		}
		return NewHTTPS(s[1], port), nil
	case ProtocolHTTP:
		if port == 0 {
			port = defaultPort[ProtocolHTTP]
		}
		return NewHTTP(s[1], port), nil
	default:
		return Value{}, ErrUnsupportedProtocol
	}
}

// NewHTTPS will create a static endpoint from parsed URL.
//
// Deprecated: Moved to github.com/beyondstorage/go-endpoint
func NewHTTPS(host string, port int) Value {
	return Value{
		Protocol: ProtocolHTTPS,
		Host:     host,
		Port:     port,
	}
}

// NewHTTP will create a static endpoint from parsed URL.
//
// Deprecated: Moved to github.com/beyondstorage/go-endpoint
func NewHTTP(host string, port int) Value {
	return Value{
		Protocol: ProtocolHTTP,
		Host:     host,
		Port:     port,
	}
}

var defaultPort = map[string]int{
	ProtocolHTTP:  80,
	ProtocolHTTPS: 443,
}
