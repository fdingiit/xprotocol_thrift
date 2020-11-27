package sofaregistry

import (
	"net"

	"github.com/gogo/protobuf/proto"
)

// Transport defines how the underlying transport works.
type Transport interface {
	// Send sends a request to transport and wait a response until success or failure.
	Send(class string, req proto.Message, res proto.Message) error
	// OnRecv waits a request then invoke fn.
	OnRecv(fn func(err error, req proto.Message)) error
	// OnRedial waits a redial then invoke fn.
	OnRedial(fn func(conn net.Conn))
	// Close closes the underlying connection then redial if necessary.
	Close() error
}
