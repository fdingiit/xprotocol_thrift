package sofadrm

import "sync/atomic"

type Metrics struct {
	redial     uint64
	serverPush uint64
	hearbeat   uint64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) addServerPush() {
	atomic.AddUint64(&m.serverPush, 1)
}

func (m *Metrics) GetServerPush() uint64 {
	return atomic.LoadUint64(&m.serverPush)
}

func (m *Metrics) addRedial() {
	atomic.AddUint64(&m.redial, 1)
}

func (m *Metrics) GetRedial() uint64 {
	return atomic.LoadUint64(&m.redial)
}

func (m *Metrics) addHeartbeat() {
	atomic.AddUint64(&m.hearbeat, 1)
}

func (m *Metrics) GetHeartbeat() uint64 {
	return atomic.LoadUint64(&m.hearbeat)
}
