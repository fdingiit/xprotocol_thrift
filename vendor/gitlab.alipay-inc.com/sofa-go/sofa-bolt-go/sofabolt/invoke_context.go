package sofabolt

import (
	"sync"
	"time"
)

type InvokeContext struct {
	// nolint
	noCopy   noCopy
	timeout  time.Duration
	created  time.Time
	req      *Request
	res      *Response
	ireslock sync.Mutex
	ires     Response
	errCh    chan error
	doneCh   chan struct{}
	callback ClientCallbacker
}

func (i *InvokeContext) GetDeadline() time.Time {
	return i.created.Add(i.timeout)
}

func (i *InvokeContext) GetTimeout() time.Duration { return i.timeout }
func (i *InvokeContext) GetCreated() time.Time     { return i.created }
func (i *InvokeContext) GetRequest() *Request      { return i.req }

func (i *InvokeContext) AssignResponse(res *Response) {
	i.ireslock.Lock()
	i.ires.Reset()
	res.CopyTo(&i.ires)
	i.ireslock.Unlock()
}

func (i *InvokeContext) Invoke(err error, res *Response) {
	if i.callback != nil {
		i.res = res
		i.callback.Invoke(err, i)

	} else {
		i.AssignResponse(res)
		// Notify the sender
		i.errCh <- err
	}
}
