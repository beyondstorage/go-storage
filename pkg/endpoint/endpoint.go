package endpoint

import (
	"fmt"
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
	switch v.Protocol {
	case "http":
		if v.Port == 80 {
			return fmt.Sprintf("http://%s", v.Host)
		}
		return fmt.Sprintf("http://%s:%d", v.Host, v.Port)
	// If user not input protocol, we will set default to "https"
	case "", "https":
		if v.Port == 443 {
			return fmt.Sprintf("https://%s", v.Host)
		}
		return fmt.Sprintf("https://%s:%d", v.Host, v.Port)
	default:
		panic(fmt.Errorf("invalid protocol: %s", v.Protocol))
	}
}
