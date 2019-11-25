package endpoint

import (
	"fmt"
	"net/url"
	"strconv"
)

// Static is the static endpoint.
type Static struct {
	protocol string
	host     string
	port     int
}

// NewStaticFromParsedURL will create a static endpoint from parsed URL.
func NewStaticFromParsedURL(protocol, host string, port int) Static {
	return Static{
		host:     host,
		port:     port,
		protocol: protocol,
	}
}

// NewStaticFromRawURL will create a static endpoint from raw URL.
func NewStaticFromRawURL(endpoint string) (Static, error) {
	errorMessage := "endpoint NewStaticFromRawURL %s: %w"

	u, err := url.Parse(endpoint)
	if err != nil {
		return Static{}, fmt.Errorf(errorMessage, endpoint, err)
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return Static{}, fmt.Errorf(errorMessage, endpoint, err)
	}
	return Static{
		host:     u.Hostname(),
		port:     port,
		protocol: u.Scheme,
	}, nil
}

// Value implements Provider interface.
func (s Static) Value() Value {
	return Value{
		Protocol: s.protocol,
		Host:     s.host,
		Port:     s.port,
	}
}
