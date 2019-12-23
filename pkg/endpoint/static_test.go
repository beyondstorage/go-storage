package endpoint

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPS(t *testing.T) {
	host := uuid.New().String()
	port := 1024
	s := NewHTTPS(host, port)
	assert.Equal(t, ProtocolHTTPS, s.protocol)
	assert.Equal(t, host, s.host)
	assert.Equal(t, port, s.port)
}

func TestNewHTTP(t *testing.T) {
	host := uuid.New().String()
	port := 1024
	s := NewHTTP(host, port)
	assert.Equal(t, ProtocolHTTP, s.protocol)
	assert.Equal(t, host, s.host)
	assert.Equal(t, port, s.port)
}

func TestStatic_Value(t *testing.T) {
	host := uuid.New().String()
	port := 1024
	v := Static{
		protocol: "http",
		host:     host,
		port:     1024,
	}.Value()
	assert.Equal(t, ProtocolHTTP, v.Protocol)
	assert.Equal(t, host, v.Host)
	assert.Equal(t, port, v.Port)
}
