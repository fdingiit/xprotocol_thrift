package asyncwriter

import (
	"time"
)

type AsyncWriterOptionSetter interface {
	set(*AsyncWriter)
}

type AsyncWriterOptionSetterFunc func(*AsyncWriter)

func (f AsyncWriterOptionSetterFunc) set(c *AsyncWriter) {
	f(c)
}

func WithAsyncWriterOption(o *Option) AsyncWriterOptionSetterFunc {
	return AsyncWriterOptionSetterFunc(func(c *AsyncWriter) {
		c.option = o
	})
}

func WithAsyncWriterMetrics(m *Metrics) AsyncWriterOptionSetterFunc {
	return AsyncWriterOptionSetterFunc(func(c *AsyncWriter) {
		c.metrics = m
	})
}

// Option configruates the option of write.
type Option struct {
	timeout       time.Duration
	flushInterval time.Duration
	batch         int
	blockwrite    bool
}

// NewOption returns a new Option.
func NewOption() *Option { return &Option{} }

func (o *Option) SetFlushInterval(d time.Duration) *Option {
	o.flushInterval = d
	return o
}

// SetTimeout sets the timeout for write if it's net.Conn
func (o *Option) SetTimeout(d time.Duration) *Option {
	o.timeout = d
	return o
}

// AllowBlockForever indicates caller can blockly write to io.Writer.
func (o *Option) AllowBlockForever() *Option {
	o.blockwrite = true
	return o
}

func (o *Option) SetBatch(b int) *Option {
	o.batch = b
	return o
}
