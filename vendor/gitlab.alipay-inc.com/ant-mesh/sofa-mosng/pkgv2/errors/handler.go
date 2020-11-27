package errors

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
	"mosn.io/mosn/pkg/protocol"
)

func GetDefaultErrorHandler() types.ErrorHandler {
	return &DefaultErrorHandler{}
}

type DefaultErrorHandler struct {
}

func (h *DefaultErrorHandler) Handle(ctx types.Context, err types.GatewayError) (int, protocol.CommonHeader, []byte) {
	headers := protocol.CommonHeader{}
	headers["x-mosng-status"] = string(err.Code())
	headers["x-mosng-status-msg"] = err.Error()
	return err.Code(), headers, nil
}
