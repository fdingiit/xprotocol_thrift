package http

import (
	"context"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkg/gateway"
	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol"
	"mosn.io/mosn/pkg/trace"
	"mosn.io/pkg/buffer"
)

const (
	HeaderOperationType = "Operation-Type"
)

type HttpDownstream struct {
	protocol api.Protocol
	appName  string
}

func init() {
	gateway.RegisterDownstreamCodec(string(protocol.HTTP1), &HttpDownstream{}, 1)
}

func (hdc *HttpDownstream) Protocol() api.Protocol {
	return hdc.protocol
}

func (hdc *HttpDownstream) SetProtocol(protocol api.Protocol) {
	hdc.protocol = protocol
}

func (hdc *HttpDownstream) AppName() string {
	return hdc.appName
}

func (hdc *HttpDownstream) SetAppName(app string) {
	hdc.appName = app
}

func (hdc *HttpDownstream) Decode(ctx context.Context, headers api.HeaderMap, dataBuf buffer.IoBuffer, trailers api.HeaderMap) (gateway.GatewayRequest, error) {
	request := gateway.NewGatewayRequest(headers, dataBuf, nil)

	var apiId string
	if operationType, ok := headers.Get(HeaderOperationType); ok {
		apiId = operationType
	} else if path, ok := headers.Get(protocol.MosnHeaderPathKey); ok {
		apiId = path
	}

	traceId := trace.IdGen().GenerateTraceId()
	request.SetRequestId(traceId)
	request.SetApiId(apiId)
	return request, nil
}

func (hdc *HttpDownstream) Encode(ctx context.Context, mapping gateway.DownstreamStatusMapping) (api.HeaderMap, buffer.IoBuffer, api.HeaderMap) {
	resp := gateway.GetGatewayContext(ctx).Response()

	if mapping != nil {
		statusMap := mapping(ctx, resp.ResultStatus())
		for k, v := range statusMap {
			resp.SetHeader(k, v)
		}
	} else {
		resp.SetHeader(gateway.HeaderResultStatus, string(resp.ResultStatus()))
	}

	return resp.Headers(), resp.DataBuf(), nil
}
