package endpoint

import (
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name  string
		cfg   string
		value Endpoint
		err   error
	}{
		{
			"invalid string",
			"abcx",
			Endpoint{},
			ErrUnsupportedProtocol,
		},
		{
			"normal http",
			"http:example.com:80",
			Endpoint{ProtocolHTTP, hostPort{"example.com", 80}},
			nil,
		},
		{
			"normal http with //",
			"http://example.com:80",
			Endpoint{ProtocolHTTP, hostPort{"example.com", 80}},
			nil,
		},
		{
			"normal http with multi /",
			"http://////example.com:80",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"normal http without port",
			"http:example.com",
			Endpoint{ProtocolHTTP, hostPort{"example.com", 80}},
			nil,
		},
		{
			"normal http without port, with //",
			"http://example.com",
			Endpoint{ProtocolHTTP, hostPort{"example.com", 80}},
			nil,
		},
		{
			"normal http without port, with multi /",
			"http://///example.com",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"wrong port number in http",
			"http:example.com:xxx",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"wrong port number in http, with //",
			"http://example.com:xxx",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"wrong port number in http, with multi /",
			"http://///example.com:xxx",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"normal https",
			"https:example.com:443",
			Endpoint{ProtocolHTTPS, hostPort{"example.com", 443}},
			nil,
		},
		{
			"normal https with //",
			"https://example.com:443",
			Endpoint{ProtocolHTTPS, hostPort{"example.com", 443}},
			nil,
		},
		{
			"normal https with multi /",
			"https://///example.com:443",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"normal https without port",
			"https:example.com",
			Endpoint{ProtocolHTTPS, hostPort{"example.com", 443}},
			nil,
		},
		{
			"normal https without port with //",
			"https://example.com",
			Endpoint{ProtocolHTTPS, hostPort{"example.com", 443}},
			nil,
		},
		{
			"normal https without port with multi /",
			"https://///example.com",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"wrong port number in https",
			"https:example.com:xxx",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"wrong port number in https with //",
			"https://example.com:xxx",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"wrong port number in https with multi /",
			"https://///example.com:xxx",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"not supported protocol",
			"notsupported:abc.com",
			Endpoint{},
			ErrUnsupportedProtocol,
		},
		{
			"not supported protocol with //",
			"notsupported://abc.com",
			Endpoint{},
			ErrUnsupportedProtocol,
		},
		{
			"normal file",
			"file:/root/data",
			Endpoint{ProtocolFile, "/root/data"},
			nil,
		},
		{
			"normal file with multi /",
			"file:///root/data",
			Endpoint{ProtocolFile, "/root/data"},
			nil,
		},
		{
			"files contains `:`",
			"file:C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\TestStorage_Stat286526883\\001\\199446694",
			Endpoint{ProtocolFile, "C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\TestStorage_Stat286526883\\001\\199446694"},
			nil,
		},
		{
			"files contains `:` with muti /",
			"file:///C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\TestStorage_Stat286526883\\001\\199446694",
			Endpoint{ProtocolFile, "/C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\TestStorage_Stat286526883\\001\\199446694"},
			nil,
		},
		{
			"normal tcp",
			"tcp:127.0.0.1:8000",
			Endpoint{ProtocolTCP, hostPort{"127.0.0.1", 8000}},
			nil,
		},
		{
			"normal tcp with //",
			"tcp://127.0.0.1:8000",
			Endpoint{ProtocolTCP, hostPort{"127.0.0.1", 8000}},
			nil,
		},
		{
			"normal tcp with multi /",
			"tcp://///127.0.0.1:8000",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"wrong port number in tcp",
			"tcp:127.0.0.1:xxx",
			Endpoint{},
			ErrInvalidValue,
		},
		{
			"wrong port number in tcp with //",
			"tcp://127.0.0.1:xxx",
			Endpoint{},
			ErrInvalidValue,
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

func TestNewFile(t *testing.T) {
	assert.Equal(t, Endpoint{ProtocolFile, "/example"}, NewFile("/example"))
}

func TestNewHTTP(t *testing.T) {
	assert.Equal(t,
		Endpoint{ProtocolHTTP, hostPort{"example.com", 8080}},
		NewHTTP("example.com", 8080),
	)
}

func TestNewHTTPS(t *testing.T) {
	assert.Equal(t,
		Endpoint{ProtocolHTTPS, hostPort{"example.com", 4433}},
		NewHTTPS("example.com", 4433),
	)
}

func TestNewTCP(t *testing.T) {
	assert.Equal(t,
		Endpoint{ProtocolTCP, hostPort{"127.0.0.1", 8000}},
		NewTCP("127.0.0.1", 8000),
	)
}

func TestEndpoint_Protocol(t *testing.T) {
	ep := NewFile("/test")

	assert.Equal(t, ProtocolFile, ep.Protocol())
}

func TestEndpoint_Protocol2(t *testing.T) {
	ep := NewTCP("127.0.0.1", 8000)

	assert.Equal(t, ProtocolTCP, ep.Protocol())
}

func TestEndpoint_String(t *testing.T) {
	cases := []struct {
		name     string
		value    Endpoint
		expected string
	}{
		{
			"file",
			Endpoint{ProtocolFile, "/test"},
			"file:/test",
		},
		{
			"http without port",
			Endpoint{ProtocolHTTP, hostPort{"example.com", 80}},
			"http:example.com:80",
		},
		{
			"http with port",
			Endpoint{ProtocolHTTP, hostPort{"example.com", 8080}},
			"http:example.com:8080",
		},
		{
			"https without port",
			Endpoint{ProtocolHTTPS, hostPort{"example.com", 443}},
			"https:example.com:443",
		},
		{
			"https with port",
			Endpoint{ProtocolHTTPS, hostPort{"example.com", 4433}},
			"https:example.com:4433",
		},
		{
			"tcp with port",
			Endpoint{ProtocolTCP, hostPort{"127.0.0.1", 8000}},
			"tcp:127.0.0.1:8000",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.value.String())
		})
	}
}

func TestEndpoint(t *testing.T) {
	p := NewFile("/test")

	assert.Panics(t, func() {
		p.HTTP()
	})
	assert.Panics(t, func() {
		p.HTTPS()
	})

	assert.Panics(t, func() {
		p.TCP()
	})

	assert.Equal(t, "/test", p.File())
}

func TestEndpoint_HTTP(t *testing.T) {
	p := NewHTTP("example.com", 80)

	url, host, port := p.HTTP()
	assert.Equal(t, "http://example.com", url)
	assert.Equal(t, "example.com", host)
	assert.Equal(t, 80, port)

	p = NewHTTP("example.com", 8080)
	url, host, port = p.HTTP()
	assert.Equal(t, "http://example.com:8080", url)
	assert.Equal(t, "example.com", host)
	assert.Equal(t, 8080, port)
}

func TestEndpoint_HTTPS(t *testing.T) {
	p := NewHTTPS("example.com", 443)

	url, host, port := p.HTTPS()
	assert.Equal(t, "https://example.com", url)
	assert.Equal(t, "example.com", host)
	assert.Equal(t, 443, port)

	p = NewHTTPS("example.com", 4433)
	url, host, port = p.HTTPS()
	assert.Equal(t, "https://example.com:4433", url)
	assert.Equal(t, "example.com", host)
	assert.Equal(t, 4433, port)
}

func TestEndpoint_TCP(t *testing.T) {
	p := NewTCP("127.0.0.1", 8000)

	addr, host, port := p.TCP()
	assert.Equal(t, "127.0.0.1:8000", addr)
	assert.Equal(t, "127.0.0.1", host)
	assert.Equal(t, 8000, port)
}

func ExampleParse() {
	ep, err := Parse("http:example.com")
	if err != nil {
		log.Fatal(err)
	}

	switch ep.Protocol() {
	case ProtocolHTTP:
		url, host, port := ep.HTTP()
		log.Println("url: ", url)
		log.Println("host: ", host)
		log.Println("port: ", port)
	case ProtocolHTTPS:
		url, host, port := ep.HTTPS()
		log.Println("url: ", url)
		log.Println("host: ", host)
		log.Println("port: ", port)
	case ProtocolFile:
		path := ep.File()
		log.Println("path: ", path)
	case ProtocolTCP:
		addr, host, port := ep.TCP()
		log.Println("addr:", addr)
		log.Println("host:", host)
		log.Println("port", port)
	default:
		panic("unsupported protocol")
	}
}
