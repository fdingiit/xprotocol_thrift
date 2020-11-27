package forwarder

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Stopped uint32 = iota // Stopped represent the state of an unstarted Forwarder.
	Started               // Started represent the state of an started Forwarder.
)

const (
	WorkerNumber  = 4
	ChanSize      = 100
	FlushInterval = 5 * time.Second
)

const (
	useragentHTTPHeaderKey = "User-Agent"
)

var (
	forwardEndpoint = "/private_api/v2/metrics"
)

type Forwarder interface {
	Start() error
	Stop()
	Forward(payload []byte, extra http.Header) error
}

type DefaultForwarder struct {
	domainForwarders map[string]*domainForwarder
	workerNum        int // number of concurrent HTTP request made by the DefaultForwarder
	logger           Logger
	state            uint32
	m                sync.RWMutex
}

// NewDefaultForwarder returns a new DefaultForwarder.
func NewDefaultForwarder(domain string, workerNum, chanSize int, flushInterval time.Duration, logger Logger) *DefaultForwarder {
	f := &DefaultForwarder{
		domainForwarders: make(map[string]*domainForwarder),
		workerNum:        workerNum,
		state:            Stopped,
		logger:           logger,
	}
	f.domainForwarders[domain] = newDomainForwarder(domain, workerNum, chanSize, flushInterval, logger)
	return f
}

// Start initialize and runs the forwarder.
func (f *DefaultForwarder) Start() error {
	// Lock so we can't stop a Forwarder while is starting
	f.m.Lock()
	defer f.m.Unlock()

	if f.state == Started {
		return fmt.Errorf("the forwarder is already started")
	}

	for _, df := range f.domainForwarders {
		_ = df.Start()
	}

	f.state = Started
	return nil
}

// Stop all the component of a forwarder and free resources.
func (f *DefaultForwarder) Stop() {
	f.m.Lock()
	defer f.m.Unlock()

	if f.state == Stopped {
		f.logger.Warnf("the forwarder is already stopped")
		return
	}

	for _, df := range f.domainForwarders {
		df.Stop()
	}
	f.domainForwarders = map[string]*domainForwarder{}
	f.state = Stopped
}

func (f *DefaultForwarder) Forward(payload []byte, extra http.Header) error {
	ts := f.createHTTPTransactions(forwardEndpoint, payload, extra)
	return f.sendHTTPTransactions(ts)
}

func (f *DefaultForwarder) createHTTPTransactions(endpoint string, payload []byte, extra http.Header) []*HTTPTransaction {
	var ts []*HTTPTransaction
	for domain := range f.domainForwarders {
		header := make(http.Header)
		for key := range extra {
			header.Set(key, extra.Get(key))
		}
		t := NewHTTPTransaction(domain, endpoint, header, payload, f.logger)
		ts = append(ts, t)
	}
	return ts
}

func (f *DefaultForwarder) sendHTTPTransactions(ts []*HTTPTransaction) error {
	if atomic.LoadUint32(&f.state) == Stopped {
		return fmt.Errorf("the forwarder is not started")
	}
	for _, t := range ts {
		if err := f.domainForwarders[t.Domain].SendHTTPTransaction(t); err != nil {
			f.logger.Error(err.Error())
		}
	}
	return nil
}
