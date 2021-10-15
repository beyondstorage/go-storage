package endpoint

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// ProtocolHTTP is the http endpoint protocol.
	ProtocolHTTP = "http"
	// ProtocolHTTPS is the https endpoint protocol.
	ProtocolHTTPS = "https"
	// ProtocolFile is the file endpoint protocol
	ProtocolFile = "file"
	// ProtocolTCP is the tcp endpoint protocol
	ProtocolTCP = "tcp"
)

// Parse will parse config string to create a endpoint Endpoint.
func Parse(cfg string) (p Endpoint, err error) {
	s := strings.Split(cfg, ":")

	//delete headmost '//'
	if len(s) > 1 && strings.HasPrefix(s[1], "//") {
		s[1] = s[1][2:]
	}
	switch s[0] {
	case ProtocolHTTP:
		host, port, err := parseHostPort(s[1:])
		if err != nil || strings.HasPrefix(host, "/") {
			return Endpoint{}, &Error{"parse", ErrInvalidValue, s[0], s[1:]}
		}
		if port == 0 {
			port = 80
		}
		return NewHTTP(host, port), nil
	case ProtocolHTTPS:
		host, port, err := parseHostPort(s[1:])
		if err != nil || strings.HasPrefix(host, "/") {
			return Endpoint{}, &Error{"parse", ErrInvalidValue, s[0], s[1:]}
		}
		if port == 0 {
			port = 443
		}
		return NewHTTPS(host, port), nil
	case ProtocolFile:
		// Handle file paths that contains ":" (often happens on windows platform)
		//
		// See issue: https://github.com/beyondstorage/go-endpoint/issues/3
		path := strings.Join(s[1:], ":")
		return NewFile(path), nil
	case ProtocolTCP:
		//See issue: https://github.com/beyondstorage/go-endpoint/issues/7
		host, port, err := parseHostPort(s[1:])
		if err != nil || strings.HasPrefix(host, "/") {
			return Endpoint{}, &Error{"parse", ErrInvalidValue, s[0], s[1:]}
		}
		return NewTCP(host, port), nil
	default:
		return Endpoint{}, &Error{"parse", ErrUnsupportedProtocol, s[0], nil}
	}
}

type hostPort struct {
	host string
	port int
}

func (hp hostPort) String() string {
	return fmt.Sprintf("%s:%d", hp.host, hp.port)
}

func parseHostPort(s []string) (host string, port int, err error) {
	if len(s) == 1 {
		return s[0], 0, nil
	}
	v, err := strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		return "", 0, err
	}
	return s[0], int(v), nil
}

type Endpoint struct {
	protocol string
	args     interface{}
}

func NewHTTP(host string, port int) Endpoint {
	return Endpoint{
		protocol: ProtocolHTTP,
		args:     hostPort{host, port},
	}
}

func NewHTTPS(host string, port int) Endpoint {
	return Endpoint{
		protocol: ProtocolHTTPS,
		args:     hostPort{host, port},
	}
}

func NewFile(path string) Endpoint {
	return Endpoint{
		protocol: ProtocolFile,
		args:     path,
	}
}

func NewTCP(host string, port int) Endpoint {
	return Endpoint{
		protocol: ProtocolTCP,
		args:     hostPort{host, port},
	}
}

func (p Endpoint) Protocol() string {
	return p.protocol
}

func (p Endpoint) String() string {
	if p.args == nil {
		return p.protocol
	}
	return fmt.Sprintf("%s:%s", p.protocol, p.args)
}

func (p Endpoint) HTTP() (url, host string, port int) {
	if p.protocol != ProtocolHTTP {
		panic(Error{
			Op:       "http",
			Err:      ErrInvalidValue,
			Protocol: p.protocol,
			Values:   p.args,
		})
	}

	hp := p.args.(hostPort)
	if hp.port == 80 {
		return fmt.Sprintf("%s://%s", p.protocol, hp.host), hp.host, 80
	}
	return fmt.Sprintf("%s://%s:%d", p.protocol, hp.host, hp.port), hp.host, hp.port
}

func (p Endpoint) HTTPS() (url, host string, port int) {
	if p.protocol != ProtocolHTTPS {
		panic(Error{
			Op:       "https",
			Err:      ErrInvalidValue,
			Protocol: p.protocol,
			Values:   p.args,
		})
	}

	hp := p.args.(hostPort)
	if hp.port == 443 {
		return fmt.Sprintf("%s://%s", p.protocol, hp.host), hp.host, 443
	}
	return fmt.Sprintf("%s://%s:%d", p.protocol, hp.host, hp.port), hp.host, hp.port
}

func (p Endpoint) File() (path string) {
	if p.protocol != ProtocolFile {
		panic(Error{
			Op:       "file",
			Err:      ErrInvalidValue,
			Protocol: p.protocol,
			Values:   p.args,
		})
	}

	return p.args.(string)
}

func (p Endpoint) TCP() (addr, host string, port int) {
	if p.protocol != ProtocolTCP {
		panic(Error{
			Op:       "tcp",
			Err:      ErrInvalidValue,
			Protocol: p.protocol,
			Values:   p.args,
		})
	}
	hp := p.args.(hostPort)
	return fmt.Sprintf("%s:%d", hp.host, hp.port), hp.host, hp.port
}
