package asyncwriteconn

import (
	"net"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-writer-go/asyncwriter"
)

type OptionSetter interface {
	set(*Conn)
}

type OptionSetterFunc func(*Conn)

func (f OptionSetterFunc) set(c *Conn) {
	f(c)
}

func WithOption(o *Option) OptionSetterFunc {
	return OptionSetterFunc(func(c *Conn) {
		c.option = o
	})
}

func WithMetrics(m *Metrics) OptionSetterFunc {
	return OptionSetterFunc(func(c *Conn) {
		c.metrics = m
	})
}

type Conn struct {
	option  *Option
	metrics *Metrics
	writer  *asyncwriter.AsyncWriter
	conn    net.Conn
}

type Metrics = asyncwriter.Metrics
type Option = asyncwriter.Option

func NewMetrics() *Metrics {
	return asyncwriter.NewMetrics()
}

func NewOption() *Option { return asyncwriter.NewOption() }

func New(conn net.Conn, options ...OptionSetterFunc) (*Conn, error) {
	c := &Conn{conn: conn}

	for i := range options {
		options[i].set(c)
	}

	if err := c.polyfill(); err != nil {
		return nil, err
	}

	aw, err := asyncwriter.New(conn,
		asyncwriter.WithAsyncWriterOption(c.option),
		asyncwriter.WithAsyncWriterMetrics(c.metrics),
	)
	if err != nil {
		return nil, err
	}
	c.writer = aw

	return c, nil
}

// nolint
func (c *Conn) polyfill() error {
	if c.option == nil {
		c.option = asyncwriter.NewOption()
	}
	if c.metrics == nil {
		c.metrics = asyncwriter.NewMetrics()
	}
	return nil
}

func (c *Conn) Read(b []byte) (n int, err error) {
	return c.conn.Read(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	return c.writer.Write(b)
}

func (c *Conn) Close() error {
	// discard the error
	_ = c.writer.Close()
	return c.conn.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
