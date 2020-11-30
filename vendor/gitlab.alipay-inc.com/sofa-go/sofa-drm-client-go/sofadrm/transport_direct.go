package sofadrm

import (
	"net"
	"sync"

	"github.com/gogo/protobuf/proto"
)

type recvEvent struct {
	err   error
	class string
	req   proto.Message
	res   proto.Message
}

type DirectTransport struct {
	sync.RWMutex
	redialCh      chan net.Conn
	recvCh        chan recvEvent
	sendCallback  func(class string, req proto.Message, res proto.Message) error
	fetchCallback func(dataID string, zone string, localVersion int) (value string, version int, err error)
}

func NewDirectTransport() *DirectTransport {
	return &DirectTransport{
		redialCh: make(chan net.Conn),
		recvCh:   make(chan recvEvent),
	}
}

func (dt *DirectTransport) SetSendCallback(callback func(class string, req proto.Message, res proto.Message) error) {
	dt.Lock()
	dt.sendCallback = callback
	dt.Unlock()
}

func (dt *DirectTransport) SetFetchCallback(callback func(dataID string, zone string,
	localVersion int) (value string, version int, err error)) {
	dt.Lock()
	dt.fetchCallback = callback
	dt.Unlock()
}

func (dt *DirectTransport) Fetch(dataID string, zone string,
	localVersion int) (value string, version int, err error) {
	dt.RLock()
	cb := dt.fetchCallback
	dt.RUnlock()
	return cb(dataID, zone, localVersion)
}

func (dt *DirectTransport) Send(class string, req proto.Message, res proto.Message) error {
	dt.RLock()
	cb := dt.sendCallback
	dt.RUnlock()
	return cb(class, req, res)
}

func (dt *DirectTransport) OnRecv(fn func(err error, class string, req, res proto.Message)) error {
	for e := range dt.recvCh {
		fn(e.err, e.class, e.req, e.res)
	}
	return nil
}

func (dt *DirectTransport) OnRedial(fn func(conn net.Conn)) {
	for conn := range dt.redialCh {
		fn(conn)
	}
}

func (dt *DirectTransport) SendRecvEvent(err error, class string, req, res proto.Message) {
	dt.recvCh <- recvEvent{
		err:   err,
		class: class,
		req:   req,
		res:   res,
	}
}

func (dt *DirectTransport) SendRedialEvent(conn net.Conn) {
	dt.redialCh <- conn
}
