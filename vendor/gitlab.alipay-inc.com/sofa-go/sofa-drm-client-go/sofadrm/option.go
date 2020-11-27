package sofadrm

import (
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

// ClientOptionSetter configures a client.
type ClientOptionSetter interface {
	set(*Client)
}

type ClientOptionSetterFunc func(*Client)

func (f ClientOptionSetterFunc) set(c *Client) {
	f(c)
}

func WithClientLogger(logger sofalogger.Logger) ClientOptionSetterFunc {
	return ClientOptionSetterFunc(func(c *Client) {
		c.logger = logger
	})
}

func WithClientConfig(config *Config) ClientOptionSetterFunc {
	return ClientOptionSetterFunc(func(c *Client) {
		c.config = config
	})
}

func WithClientTransport(transport Transport) ClientOptionSetterFunc {
	return ClientOptionSetterFunc(func(c *Client) {
		c.transport = transport
	})
}

func WithClientMetrics(m *Metrics) ClientOptionSetterFunc {
	return ClientOptionSetterFunc(func(c *Client) {
		c.metrics = m
	})
}
