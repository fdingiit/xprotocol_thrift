package sofadrm

import (
	"sync"
	"sync/atomic"
)

var dummyListenerFunc ListenerFunc = func(dataID string, value string) {}

//go:generate syncmap -pkg sofadrm -o listeners_generated.go -name ListenerMap map[string]*listener

type Listener interface {
	OnDRMPush(dataID string, value string)
}

type ListenerFunc func(dataID string, value string)

func (s ListenerFunc) OnDRMPush(dataID string, value string) {
	s(dataID, value)
}

type listener struct {
	sync.RWMutex
	ln         []Listener
	registered uint32
}

func buildListener(ln Listener) *listener {
	m := &listener{}
	m.Add(ln)
	return m
}

func (l *listener) Size() int {
	l.RLock()
	n := len(l.ln)
	l.RUnlock()
	return n
}

func (l *listener) Del(ln Listener) {
	l.Lock()
	for i := range l.ln {
		if l.ln[i] == ln {
			l.ln = append(l.ln[:i], l.ln[i+1:]...)
			break
		}
	}
	l.Unlock()
}

func (l *listener) Add(ln Listener) {
	if ln == nil { // ignore nil listener
		return
	}
	l.Lock()
	l.ln = append(l.ln, ln)
	l.Unlock()
}

func (l *listener) IsRegistered() bool {
	return atomic.LoadUint32(&l.registered) == 1
}

func (l *listener) MarkRegistered() {
	atomic.StoreUint32(&l.registered, 1)
}

func (l *listener) MarkUnRegistered() {
	atomic.StoreUint32(&l.registered, 0)
}

func (l *listener) OnDRMPush(dataID string, value string) {
	l.RLock()
	for i := range l.ln {
		l.ln[i].OnDRMPush(dataID, value)
	}
	l.RUnlock()
}
