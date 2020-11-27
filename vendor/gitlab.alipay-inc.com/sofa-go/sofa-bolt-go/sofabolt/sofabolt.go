package sofabolt

import (
	"errors"
)

var (
	ErrBufferNotEnough       = errors.New("sofabolt: buffer not enough")
	ErrMalformedProto        = errors.New("sofabolt: malformed proto")
	ErrMalformedType         = errors.New("sofabolt: malformed type")
	ErrServerHandler         = errors.New("sofabolt: server handler cannot be nil")
	ErrServerNotARequest     = errors.New("sofabolt: server received a response")
	ErrClientExpectResponse  = errors.New("sofabolt: receive a request")
	ErrClientTimeout         = errors.New("sofabolt: client do timeout")
	ErrClientNotARequest     = errors.New("sofabolt: client send a response")
	ErrClientWasClosed       = errors.New("sofabolt: client was closed")
	ErrClientTooManyRequests = errors.New("sofabolt: client too many requests")
	ErrClientServerTimeout   = errors.New("sofabolt: clientserver do timeout")
	ErrClientDisableRedial   = errors.New("sofabolt: disable redial")
	ErrClientNilConnection   = errors.New("sofabolt: client connection is nil")
)
