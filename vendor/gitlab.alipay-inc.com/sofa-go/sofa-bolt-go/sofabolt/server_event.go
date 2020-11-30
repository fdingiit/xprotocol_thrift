package sofabolt

import "net"

//go:generate stringer -type=ServerEvent

type ServerEvent uint16

const (
	ServerTemporaryAcceptEvent    ServerEvent = 0
	ServerWorkerPoolOverflowEvent ServerEvent = 1
	ServerConnErrorEvent          ServerEvent = 2
	ServerConnHijackedEvent       ServerEvent = 3
)

type ServerEventContext struct {
	req   *Request
	res   *Response
	conn  net.Conn
	event ServerEvent
}

func NewServerEventContext(event ServerEvent) *ServerEventContext {
	return &ServerEventContext{event: event}
}

func (s ServerEventContext) GetType() ServerEvent { return s.event }

func (sec *ServerEventContext) SetConn(conn net.Conn) *ServerEventContext {
	sec.conn = conn
	return sec
}

func (sec *ServerEventContext) SetReq(req *Request) *ServerEventContext {
	sec.req = req
	return sec
}

func (sec *ServerEventContext) SetRes(res *Response) *ServerEventContext {
	sec.res = res
	return sec
}

type ServerOnEventHandler func(*Server, error, *ServerEventContext)

var DummyServerOnEventHandler = ServerOnEventHandler(func(*Server, error, *ServerEventContext) {
})
