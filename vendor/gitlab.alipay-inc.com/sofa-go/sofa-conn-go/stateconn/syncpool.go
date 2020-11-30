package stateconn

import (
	"net"
	"sync"
)

var (
	statePool = sync.Pool{
		New: func() interface{} {
			return &StateConn{}
		},
	}
)

func AcquireConn(conn net.Conn) *StateConn {
	sc, ok := statePool.Get().(*StateConn)
	if !ok {
		sc = &StateConn{}
	}
	sc.SetState(StateNew)
	sc.Conn = conn
	return sc
}

func ReleaseConn(sc *StateConn) {
	statePool.Put(sc)
}
