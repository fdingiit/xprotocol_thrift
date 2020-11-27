package sofaregistry

import "sync/atomic"

type Metrics struct {
	successTask          uint64
	failureTask          uint64
	redial               uint64
	serverPush           uint64
	publisherRegister    uint64
	publisherUnregister  uint64
	subscriberRegister   uint64
	subscriberUnregister uint64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) addPublisherRegister() {
	atomic.AddUint64(&m.publisherRegister, 1)
}

func (m *Metrics) GetPublisherRegister() uint64 { return atomic.LoadUint64(&m.publisherRegister) }

func (m *Metrics) addPublisherUnregister() {
	atomic.AddUint64(&m.publisherUnregister, 1)
}

func (m *Metrics) GetPublisherUnregister() uint64 { return atomic.LoadUint64(&m.publisherUnregister) }

func (m *Metrics) addSubscriberRegister() {
	atomic.AddUint64(&m.subscriberRegister, 1)
}
func (m *Metrics) GetSubscriberRegister() uint64 { return atomic.LoadUint64(&m.subscriberRegister) }

func (m *Metrics) addSubscriberUnregister() {
	atomic.AddUint64(&m.subscriberUnregister, 1)
}

func (m *Metrics) GetSubscriberUnregister() uint64 { return atomic.LoadUint64(&m.subscriberUnregister) }

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

func (m *Metrics) addSuccessTask() {
	atomic.AddUint64(&m.successTask, 1)
}

func (m *Metrics) addFailureTask() {
	atomic.AddUint64(&m.failureTask, 1)
}

func (m *Metrics) GetSuccessTask() uint64 {
	return atomic.LoadUint64(&m.successTask)
}

func (m *Metrics) GetFailureTask() uint64 {
	return atomic.LoadUint64(&m.failureTask)
}
