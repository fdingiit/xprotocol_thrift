package errorconn

import (
	"net"
	"time"
)

type Conn struct {
	err error
}

func New(err error) net.Conn {
	return &Conn{err: err}
}

func (dc *Conn) Read(b []byte) (n int, err error) {
	return 0, dc.err
}

func (dc *Conn) Write(b []byte) (n int, err error) {
	return 0, dc.err
}

func (dc *Conn) Close() error {
	return dc.err
}

func (dc *Conn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (dc *Conn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (dc *Conn) SetDeadline(t time.Time) error {
	return dc.err
}

func (dc *Conn) SetReadDeadline(t time.Time) error {
	return dc.err
}

func (dc *Conn) SetWriteDeadline(t time.Time) error {
	return dc.err
}
