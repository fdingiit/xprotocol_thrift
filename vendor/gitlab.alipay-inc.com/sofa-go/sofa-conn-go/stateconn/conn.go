package stateconn

import (
	"net"
	"sync/atomic"
	"time"
)

type StateGetter interface {
	GetState() (time.Time, State)
	GetConn() net.Conn
}

type StateConn struct {
	state int64
	net.Conn
}

func (s *StateConn) Write(p []byte) (n int, err error) {
	s.SetState(StateActive)
	return s.Conn.Write(p)
}

func (s *StateConn) Read(p []byte) (n int, err error) {
	s.SetState(StateActive)
	return s.Conn.Read(p)
}

func (s *StateConn) Close() error {
	// close once
	if _, state := s.GetState(); state != StateClosed {
		s.SetState(StateClosed)
		return s.Conn.Close()
	}
	return nil
}

func (s *StateConn) GetConn() net.Conn {
	return s.Conn
}

func (s *StateConn) SetState(state State) {
	atomic.StoreInt64(&s.state, packState(state))
}

func (s *StateConn) GetState() (time.Time, State) {
	i64 := atomic.LoadInt64(&s.state)
	return time.Unix(i64>>8, 0), State(i64 & 0xFF)
}
