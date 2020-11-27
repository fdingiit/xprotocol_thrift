package types

import (
	"mosn.io/mosn/pkg/protocol"
)

type GatewayError interface {
	Code() int
	Error() string
}

type ErrorHandler interface {
	Handle(Context, GatewayError) (httpCode int, headers protocol.CommonHeader, res []byte)
}
