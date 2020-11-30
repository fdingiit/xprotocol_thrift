package forwarder

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// domainForwarder is in charge of sending Transactions to Backend over HTTP
// and retrying them if needed. One domainForwarder is created per HTTP backend.
type domainForwarder struct {
	domain              string
	workerNum           int
	chanSize            int
	flushInterval       time.Duration
	logger              Logger
	highPrioChan        chan Transaction
	lowPrioChan         chan Transaction
	requeueChan         chan Transaction
	workers             []*Worker
	retryQueue          []Transaction
	stopRetry           chan bool
	state               uint32
	isRetrying          int32
	transactionDropped  int32
	transactionRetried  int32
	transactionRequeued int32
	m                   sync.RWMutex
}

func newDomainForwarder(domain string, workerNum, chanSize int, flushInterval time.Duration, logger Logger) *domainForwarder {
	f := &domainForwarder{
		domain:        domain,
		workerNum:     workerNum,
		chanSize:      chanSize,
		flushInterval: flushInterval,
		logger:        logger,
		state:         Stopped,
	}
	return f
}

func (f *domainForwarder) init() {
	f.highPrioChan = make(chan Transaction, f.chanSize)
	f.lowPrioChan = make(chan Transaction, f.chanSize)
	f.requeueChan = make(chan Transaction, f.chanSize)
	f.workers = []*Worker{}
	f.retryQueue = []Transaction{}
	f.stopRetry = make(chan bool)
}

// Start starts a domainForwarder.
func (f *domainForwarder) Start() error {
	// Lock so we can't stop a Forwarder while is starting
	f.m.Lock()
	defer f.m.Unlock()

	if f.state == Started {
		return fmt.Errorf("the forwarder is already started")
	}

	// reset internal state to purge transactions from past starts
	f.init()

	for i := 0; i < f.workerNum; i++ {
		w := NewWorker(f.highPrioChan, f.lowPrioChan, f.requeueChan, f.logger)
		w.Start()
		f.workers = append(f.workers, w)
	}

	go f.HandleFailedTransactions()

	f.logger.Infof("domainForwarder started")
	f.state = Started
	return nil
}

func (f *domainForwarder) HandleFailedTransactions() {
	ticker := time.NewTicker(f.flushInterval)
	for {
		select {
		case <-ticker.C:
			f.retryTransactions()
		case t := <-f.requeueChan:
			f.requeueTransaction(t)
		case <-f.stopRetry:
			ticker.Stop()
			return
		}
	}
}

type byCreateTime []Transaction

func (v byCreateTime) Len() int           { return len(v) }
func (v byCreateTime) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byCreateTime) Less(i, j int) bool { return v[i].GetCreatedAt().After(v[j].GetCreatedAt()) }

func (f *domainForwarder) retryTransactions() {
	// In case it takes more that flushInterval to sort and retry transactions, we skip a retry
	if !atomic.CompareAndSwapInt32(&f.isRetrying, 0, 1) {
		f.logger.Errorf("the forwarder is still retrying Transactions: this should never happen and you might lower the flush interval")
		return
	}
	defer atomic.StoreInt32(&f.isRetrying, 0)

	var dropped, retried int32
	sort.Sort(byCreateTime(f.retryQueue))
	for _, t := range f.retryQueue {
		select {
		case f.lowPrioChan <- t:
			retried++
		default:
			dropped++
		}
	}

	f.retryQueue = []Transaction{}
	atomic.AddInt32(&f.transactionDropped, dropped)
	atomic.AddInt32(&f.transactionRetried, retried)
	if dropped > 0 {
		f.logger.Errorf("dropped %d transactions in this retry attempt, because the workers are too busy", dropped)
	}
}

func (f *domainForwarder) requeueTransaction(t Transaction) {
	f.retryQueue = append(f.retryQueue, t)
	atomic.AddInt32(&f.transactionRequeued, 1)
}

// Stop stops a domainForwarder, all transactions not yet flushed will be lost.
func (f *domainForwarder) Stop() {
	// Lock so we can't start a Forwarder while is stopping
	f.m.Lock()
	defer f.m.Unlock()

	if f.state == Stopped {
		f.logger.Warnf("the forwarder is already stopped")
		return
	}

	f.stopRetry <- true
	for _, w := range f.workers {
		w.Stop()
	}
	f.workers = []*Worker{}
	f.retryQueue = []Transaction{}
	close(f.highPrioChan)
	close(f.lowPrioChan)
	close(f.requeueChan)
	f.logger.Infof("domainForwarder stopped")
	f.state = Stopped
}

func (f *domainForwarder) SendHTTPTransaction(t Transaction) error {
	select {
	case f.highPrioChan <- t:
	default:
		return fmt.Errorf("the forwarder input queue for %s is full: dropping transaction", f.domain)
	}
	return nil
}
