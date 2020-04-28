package httpclient

import (
	"net/http"
	"time"
)

// Options is the httpclient supported options.
type Options struct {
	// Dialer related options
	DialerConnectTimeout time.Duration

	// Underlying connection related options
	ConnReadTimeout  time.Duration
	ConnWriteTimeout time.Duration

	// HTTP client related options
}

// New will create new http client.
func New(o *Options) *http.Client {
	dialer := NewDialer()

	hc := &http.Client{
		Transport: &http.Transport{
			DialContext: dialer.DialContext,

			// Support http proxy from env.
			Proxy: http.ProxyFromEnvironment,
			// Specify timeout for tls handshake.
			TLSHandshakeTimeout: 10 * time.Second,
			// Specify max idle conns across all hosts.
			MaxIdleConns: 0,
			// Specify max idle conns across per host.
			MaxIdleConnsPerHost: 100,
			// Specify timeout for closing idle (keep-alive) connection.
			IdleConnTimeout: 90 * time.Second,
			// Specify timeout that waiting for server's approve before sending data.
			ExpectContinueTimeout: time.Second,
			// Gzip file should not be auto-decompressed
			DisableCompression: true,
		},
		// http client used in storage don't need to follow redirect, return directly.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		// We will handle timeout by ourselves, and disable http.Client's timeout.
		Timeout: 0,
	}

	if o == nil {
		return hc
	}

	if o.DialerConnectTimeout > 0 {
		dialer.WithConnectTimeout(o.DialerConnectTimeout)
	}
	if o.ConnReadTimeout > 0 {
		dialer.WithReadTimeout(o.ConnReadTimeout)
	}
	if o.ConnWriteTimeout > 0 {
		dialer.WithWriteTimeout(o.ConnWriteTimeout)
	}
	return hc
}
