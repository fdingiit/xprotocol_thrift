package sofadrm

import (
	"net"

	"github.com/gogo/protobuf/proto"
)

// Transport defines how the underlying transport works.
type Transport interface {
	// Fetch fetches the value from the server.
	Fetch(dataID string, zone string, localVersion int) (value string, version int, err error)
	// Send sends a request to transport and wait a response until success or failure.
	// res nil means no response.
	Send(class string, req proto.Message, res proto.Message) error
	// OnRecv waits a request then invoke fn.
	OnRecv(fn func(err error, class string, req proto.Message, res proto.Message)) error
	// OnRedial waits a redial then invoke fn.
	OnRedial(fn func(conn net.Conn))
}
