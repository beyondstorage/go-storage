package endpoint

const (
	// ProtocolHTTPS is the https credential protocol.
	ProtocolHTTPS = "https"
	// ProtocolHTTP is the http credential protocol.
	ProtocolHTTP = "http"
)

// Static is the static endpoint.
type Static struct {
	protocol string
	host     string
	port     int
}

// Value implements Provider interface.
func (s Static) Value() Value {
	return Value{
		Protocol: s.protocol,
		Host:     s.host,
		Port:     s.port,
	}
}

// NewHTTPS will create a static endpoint from parsed URL.
func NewHTTPS(host string, port int) Static {
	return Static{
		protocol: ProtocolHTTPS,
		host:     host,
		port:     port,
	}
}

// NewHTTP will create a static endpoint from parsed URL.
func NewHTTP(host string, port int) Static {
	return Static{
		protocol: ProtocolHTTP,
		host:     host,
		port:     port,
	}
}
