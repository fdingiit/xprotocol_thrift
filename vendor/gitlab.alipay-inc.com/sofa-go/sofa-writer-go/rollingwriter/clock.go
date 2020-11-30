package rollingwriter

import (
	"sync"
	"time"
)

type Clocker interface {
	Now() time.Time
}

var (
	_ Clocker = (*WallClocker)(nil)
	_ Clocker = (*FakeClocker)(nil)
)

type WallClocker struct {
}

func (wc WallClocker) Now() time.Time { return time.Now() }

type FakeClocker struct {
	sync.RWMutex
	now time.Time
}

func (f *FakeClocker) SetNow(n time.Time) {
	f.Lock()
	f.now = n
	f.Unlock()
}

func (f *FakeClocker) Now() time.Time {
	f.RLock()
	defer f.RUnlock()
	return f.now
}
