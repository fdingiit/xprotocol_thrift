// Package coarse100ms export the coarse time.Now() via sleep 100ms.
package coarse100ms

import (
	"sync/atomic"
	"time"
)

var now atomic.Value

func init() {
	t := time.Now()
	now.Store(&t)
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			t := time.Now()
			now.Store(&t)
		}
	}()
}

// Now returns the time.Now() whose's precision is 100ms.
func Now() time.Time {
	tp := now.Load().(*time.Time)
	return *tp
}
