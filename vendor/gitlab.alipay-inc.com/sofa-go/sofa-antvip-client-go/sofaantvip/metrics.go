package sofaantvip

import (
	"sync/atomic"
	"time"
)

type Metrics struct {
	CloudSyncer  CloudSyncerMetric
	alipaysyncer AlipaySyncerMetric
}

func NewMetrics() *Metrics { return &Metrics{} }

func (m *Metrics) GetAlipaySyncerMetric() *AlipaySyncerMetric { return &m.alipaysyncer }

func (m *Metrics) GetCloudSyncerMetric() *CloudSyncerMetric { return &m.CloudSyncer }

type AlipaySyncerMetric struct {
	success  int64
	failure  int64
	lastsync int64
}

func (m *AlipaySyncerMetric) LoadSuccess() int64    { return atomic.LoadInt64(&m.success) }
func (m *AlipaySyncerMetric) LoadFailure() int64    { return atomic.LoadInt64(&m.failure) }
func (m *AlipaySyncerMetric) LoadLastSyncAt() int64 { return atomic.LoadInt64(&m.lastsync) }

func (m *AlipaySyncerMetric) addSuccess() {
	atomic.AddInt64(&m.success, 1)
	atomic.StoreInt64(&m.lastsync, time.Now().Unix())
}

func (m *AlipaySyncerMetric) addFailure() {
	atomic.AddInt64(&m.failure, 1)
	atomic.StoreInt64(&m.lastsync, time.Now().Unix())
}

type CloudSyncerMetric struct {
	success  int64
	failure  int64
	lastsync int64
}

func (m *CloudSyncerMetric) LoadSuccess() int64    { return atomic.LoadInt64(&m.success) }
func (m *CloudSyncerMetric) LoadFailure() int64    { return atomic.LoadInt64(&m.failure) }
func (m *CloudSyncerMetric) LoadLastSyncAt() int64 { return atomic.LoadInt64(&m.lastsync) }

func (m *CloudSyncerMetric) addSuccess() {
	atomic.AddInt64(&m.success, 1)
	atomic.StoreInt64(&m.lastsync, time.Now().Unix())
}

func (m *CloudSyncerMetric) addFailure() {
	atomic.AddInt64(&m.failure, 1)
	atomic.StoreInt64(&m.lastsync, time.Now().Unix())
}
