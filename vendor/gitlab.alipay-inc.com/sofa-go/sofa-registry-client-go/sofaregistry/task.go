package sofaregistry

import (
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"
)

type task struct {
	req      proto.Message
	res      proto.Message
	class    string
	id       string
	done     uint32
	version  int64
	failures uint32
}

func (t *task) tryDelay() {
	if failures := t.getFailures(); failures > 0 {
		// delay max 1s if always failed
		if failures > 1000 {
			failures = 1000
		}
		time.Sleep(time.Duration(failures) * time.Millisecond)
	}
}

func (t *task) getFailures() uint32 { return atomic.LoadUint32(&t.failures) }

func (t *task) addFailures() {
	atomic.AddUint32(&t.failures, 10)
}

func (t *task) Done() {
	atomic.StoreUint32(&t.done, 1)
}

func (t *task) IsDone() bool { return atomic.LoadUint32(&t.done) == 1 }

func (t *task) ID() string {
	return t.id
}

func (t *task) GetVersion() int64 { return t.version }
