package sofabolt

import (
	"net"
	"time"
)

// clientOptionSetter configures a client.
type clientOptionSetter interface {
	set(*Client)
}

type clientOptionSetterFunc func(*Client)

func (f clientOptionSetterFunc) set(c *Client) {
	f(c)
}

func WithClientMetrics(cm *ClientMetrics) clientOptionSetterFunc {
	return clientOptionSetterFunc(func(c *Client) {
		c.metrics = cm
	})
}

func WithClientDisableAutoIncrementRequestID(b bool) clientOptionSetterFunc {
	return clientOptionSetterFunc(func(c *Client) {
		c.options.disableAutoIncrementRequestID = b
	})
}

func WithClientTimeout(readtimeout,
	writetimeout,
	idletimeout time.Duration,
	flushInterval time.Duration) clientOptionSetterFunc {
	return clientOptionSetterFunc(func(c *Client) {
		c.options.readTimeout = readtimeout
		c.options.writeTimeout = writetimeout
		c.options.idleTimeout = idletimeout
	})
}

func WithClientConn(conn net.Conn) clientOptionSetterFunc {
	return clientOptionSetterFunc(func(c *Client) {
		c.conn = conn
	})
}

func WithClientHeartbeat(heartbeatinterval, heartbeattimeout time.Duration,
	heartbeatprobes int, onheartbeat func(success bool)) clientOptionSetterFunc {
	return clientOptionSetterFunc(func(c *Client) {
		c.options.heartbeatTimeout = heartbeattimeout
		c.options.heartbeatInterval = heartbeatinterval
		c.options.heartbeatProbes = heartbeatprobes
		c.options.onHeartbeat = onheartbeat
	})
}

func WithClientMaxPendingCommands(m int) clientOptionSetterFunc {
	return clientOptionSetterFunc(func(c *Client) {
		c.options.maxPendingCommands = m
	})
}

func WithClientRedial(dialer Dialer) clientOptionSetterFunc {
	return clientOptionSetterFunc(func(c *Client) {
		c.options.dialer = dialer
	})
}

func WithClientHandler(handler Handler) clientOptionSetterFunc {
	return clientOptionSetterFunc(func(c *Client) {
		c.options.handler = handler
	})
}
