package sofaregistry

import (
	"net"
	"sync"

	"github.com/gogo/protobuf/proto"
)

var _ Transport = (*DirectTransport)(nil)

type recvEvent struct {
	err error
	req proto.Message
}

type DirectTransport struct {
	sync.RWMutex
	redialCh chan net.Conn
	recvCh   chan recvEvent
	callback func(class string, req proto.Message, res proto.Message) error
}

func NewDirectTransport() *DirectTransport {
	return &DirectTransport{
		redialCh: make(chan net.Conn),
		recvCh:   make(chan recvEvent),
	}
}

func (dt *DirectTransport) SetCallback(callback func(class string, req proto.Message, res proto.Message) error) {
	dt.Lock()
	dt.callback = callback
	dt.Unlock()
}

func (dt *DirectTransport) Send(class string, req proto.Message, res proto.Message) error {
	dt.RLock()
	cb := dt.callback
	dt.RUnlock()
	return cb(class, req, res)
}

func (dt *DirectTransport) OnRecv(fn func(err error, req proto.Message)) error {
	for e := range dt.recvCh {
		fn(e.err, e.req)
	}
	return nil
}

func (dt *DirectTransport) OnRedial(fn func(conn net.Conn)) {
	for conn := range dt.redialCh {
		fn(conn)
	}
}

func (dt *DirectTransport) SendRecvEvent(err error, req proto.Message) {
	dt.recvCh <- recvEvent{
		err: err,
		req: req,
	}
}

func (dt *DirectTransport) SendRedialEvent(conn net.Conn) {
	dt.redialCh <- conn
}

func (dt *DirectTransport) Close() error {
	return nil
}
