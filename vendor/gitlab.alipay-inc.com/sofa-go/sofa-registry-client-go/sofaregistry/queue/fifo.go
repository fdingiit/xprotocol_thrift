package queue

import (
	"sync"
	"time"
)

type deletingversion struct {
	Version           int64
	DeletionTimestamp time.Time
}

const (
	defaultDeletionKeepDuration = time.Minute
)

type FIFO struct {
	lock     sync.RWMutex
	cond     sync.Cond
	queue    []string
	elements map[string]Element

	deletingVersions     map[string]deletingversion
	deletionKeepDuration time.Duration
}

func NewFIFOWithSize(size int) *FIFO {
	return newFIFOWithSizeAndKeepDuration(size, defaultDeletionKeepDuration)
}
func newFIFOWithSizeAndKeepDuration(size int, keepDuration time.Duration) *FIFO {
	f := &FIFO{
		elements:             make(map[string]Element, size),
		queue:                make([]string, 0, size),
		deletingVersions:     make(map[string]deletingversion, size),
		deletionKeepDuration: keepDuration,
	}
	f.cond.L = &f.lock
	go f.cleanupExpired()
	return f
}

func NewFIFO() *FIFO {
	return NewFIFOWithSize(16)
}

func (f *FIFO) Len() int {
	f.lock.RLock()
	n := len(f.queue)
	f.lock.RUnlock()
	return n
}

func (f *FIFO) cleanupExpired() {
	ticker := time.NewTicker(f.deletionKeepDuration)
	defer ticker.Stop()
	for {
		<-ticker.C
		f.lock.Lock()
		expiredTime := time.Now().Add(-f.deletionKeepDuration)
		for id, lastVersion := range f.deletingVersions {
			if lastVersion.DeletionTimestamp.Before(expiredTime) {
				delete(f.deletingVersions, id)
			}
		}
		f.lock.Unlock()
	}
}

// Push pushes a message to the queue.
func (f *FIFO) Push(i Element) error {
	id := i.ID()
	f.lock.Lock()
	defer f.lock.Unlock()
	if raw, ok := f.elements[id]; !ok {
		if deleting, ok := f.deletingVersions[id]; ok && i.GetVersion() < deleting.Version {
			return nil
		}
		f.queue = append(f.queue, id)
		f.elements[id] = i
	} else {
		zraw, ok := raw.(Element)
		if !ok {
			panic("failed to type assert")
		}

		// always use the newest version
		if i.GetVersion() > zraw.GetVersion() {
			f.elements[id] = i
		}
	}

	f.cond.Broadcast()

	return nil
}

func (f *FIFO) Pop() (Element, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	for {
		for len(f.queue) == 0 {
			f.cond.Wait()
		}
		id := f.queue[0]
		f.queue = f.queue[1:]
		element, ok := f.elements[id]
		if !ok { // stale element
			continue
		}
		if deleting, ok := f.deletingVersions[id]; !ok || element.GetVersion() >= deleting.Version {
			f.deletingVersions[id] = deletingversion{Version: element.GetVersion(), DeletionTimestamp: time.Now()}
		}
		delete(f.elements, id)
		return element, nil
	}
}
