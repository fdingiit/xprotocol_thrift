package sofaantvip

import (
	"sync"
)

type domainListener struct {
	sync.RWMutex
	name      string
	domain    *VipDomain
	listeners []Listener
}

func newDomainListener(name string, ln ...Listener) *domainListener {
	dl := &domainListener{
		name: name,
	}
	dl.listeners = append(dl.listeners, ln...)
	return dl
}

func (l *domainListener) getName() string {
	l.RLock()
	defer l.RUnlock()
	return l.name
}

func (l *domainListener) getDomain() *VipDomain {
	l.RLock()
	defer l.RUnlock()
	return l.domain
}

func (l *domainListener) setDomain(d *VipDomain) {
	l.Lock()
	defer l.Unlock()
	l.domain = d
}

func (l *domainListener) addListener(ln Listener) {
	l.Lock()
	defer l.Unlock()
	l.listeners = append(l.listeners, ln)
}

func (l *domainListener) getListeners() []Listener {
	l.RLock()
	defer l.RUnlock()
	return l.listeners
}
