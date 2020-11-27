package sofadrm

import (
	"sync"
	"sync/atomic"
)

//go:generate syncmap -pkg sofadrm -o mutlisteners_generated.go -name MutListenerMap map[string]*mutlistener

// mutlistener means always mutable
type mutlistener struct {
	sync.RWMutex
	ln         Listener
	registered uint32
}

func buildMutListener(ln Listener) *mutlistener {
	return &mutlistener{
		ln: ln,
	}
}

func (l *mutlistener) OverWrite(ln Listener) {
	l.Lock()
	l.ln = ln
	l.Unlock()
}

func (l *mutlistener) IsRegistered() bool {
	return atomic.LoadUint32(&l.registered) == 1
}

func (l *mutlistener) MarkRegistered() {
	atomic.StoreUint32(&l.registered, 1)
}

func (l *mutlistener) MarkUnRegistered() {
	atomic.StoreUint32(&l.registered, 0)
}

func (l *mutlistener) OnDRMPush(dataID string, value string) {
	l.RLock()
	if l.ln != nil {
		l.ln.OnDRMPush(dataID, value)
	}
	l.RUnlock()
}
