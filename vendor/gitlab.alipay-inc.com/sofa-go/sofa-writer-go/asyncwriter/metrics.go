package asyncwriter

import "sync/atomic"

type Metrics struct {
	commands        *int64
	pendingCommands *int64
	bytes           *int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		commands:        new(int64),
		pendingCommands: new(int64),
		bytes:           new(int64),
	}
}

func (m *Metrics) SetPendingCommands(i *int64) { m.pendingCommands = i }

func (m *Metrics) GetPendingCommands() int64 { return atomic.LoadInt64(m.pendingCommands) }

func (m *Metrics) AddPendingCommands(n int64) { atomic.AddInt64(m.pendingCommands, n) }

func (m *Metrics) AddCommands() { atomic.AddInt64(m.commands, 1) }

func (m *Metrics) SetCommands(i *int64) {
	m.commands = i
}

func (m *Metrics) GetCommands() int64 { return atomic.LoadInt64(m.commands) }

func (m *Metrics) AddBytes(n int64) { atomic.AddInt64(m.bytes, n) }

func (m *Metrics) GetBytes() int64 { return atomic.LoadInt64(m.bytes) }

func (m *Metrics) SetBytes(i *int64) {
	m.bytes = i
}
