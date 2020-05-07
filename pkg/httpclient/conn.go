package httpclient

import (
	"net"
	"time"
)

// Conn is a generic stream-oriented network connection.
type Conn struct {
	net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// Read will read from the conn.
func (c *Conn) Read(buf []byte) (n int, err error) {
	err = c.SetDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return
	}
	defer func() {
		// Clean read timeout so that this will not affect further read
		// It's safe to ignore the returning error: even if it don’t return now, it will return via next read.
		_ = c.SetDeadline(time.Time{})
	}()

	return c.Conn.Read(buf)
}

// Write will write into the conn.
func (c *Conn) Write(buf []byte) (n int, err error) {
	err = c.SetDeadline(time.Now().Add(c.writeTimeout))
	if err != nil {
		return
	}
	defer func() {
		// Clean read timeout so that this will not affect further write
		// It's safe to ignore the returning error: even if it don’t return now, it will return via next write.
		_ = c.SetDeadline(time.Time{})
	}()

	return c.Conn.Write(buf)
}
