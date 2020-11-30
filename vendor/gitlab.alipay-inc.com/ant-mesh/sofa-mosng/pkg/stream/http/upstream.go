package http

import (
	"context"
	"net/http"
	"strconv"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkg/gateway"
	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol"
	mosnhttp "mosn.io/mosn/pkg/protocol/http"
	"mosn.io/pkg/buffer"
)

type HttpUpstream struct {
	protocol    api.Protocol
	timeout     uint64
	clusterName string
	path        string
}

func init() {
	gateway.RegisterUpstreamParser(string(protocol.HTTP1), parseUpstreamConfig, 1)
	gateway.RegisterUpstreamCodec(string(protocol.HTTP1), &HttpUpstream{}, 1)
}

func (h *HttpUpstream) Protocol() api.Protocol {
	return h.protocol
}

func (h *HttpUpstream) SetProtocol(protocol api.Protocol) {
	h.protocol = protocol
}

func (h *HttpUpstream) Timeout() uint64 {
	return h.timeout
}

func (h *HttpUpstream) SetTimeout(timeout uint64) {
	h.timeout = timeout
}

func (h *HttpUpstream) ClusterName() string {
	return h.clusterName
}

func (h *HttpUpstream) SetClusterName(cluster string) {
	h.clusterName = cluster
}

func (h *HttpUpstream) Path() string {
	return h.path
}

func (h *HttpUpstream) SetPath(path string) {
	h.path = path
}

func parseUpstreamConfig(cfg gateway.UpstreamConfig) gateway.Upstream {
	upt := &HttpUpstream{}
	upt.SetProtocol(cfg.Protocol)
	upt.SetClusterName(cfg.ClusterName)
	upt.SetTimeout(cfg.TimeOut)
	if path, ok := cfg.Config["path"]; ok {
		upt.SetPath(path.(string))
	}
	return upt
}

func (huc *HttpUpstream) Encode(ctx context.Context) (api.HeaderMap, buffer.IoBuffer, api.HeaderMap) {
	gwCtx := gateway.GetGatewayContext(ctx)
	req := gwCtx.Request()

	sc := gwCtx.Service()
	hup := sc.Upstream().(*HttpUpstream)
	req.SetHeader(protocol.MosnHeaderPathKey, hup.Path())

	return req.Headers(), req.DataBuf(), nil
}

func (huc *HttpUpstream) Decode(ctx context.Context, headers api.HeaderMap, dataBuf buffer.IoBuffer, trailers api.HeaderMap, mapping gateway.UpstreamStatusMapping) (gateway.GatewayResponse, error) {
	context := gateway.GetGatewayContext(ctx)

	resp := gateway.NewGatewayResponse(nil, dataBuf, nil)
	resp.SetResultStatus(gateway.BizException)
	resp.SetDataEncoding(context.Request().DataEncoding())
	if respHeaders, ok := headers.(mosnhttp.ResponseHeader); ok {
		if respHeaders.StatusCode() == http.StatusOK {
			resp.SetResultStatus(gateway.ResultSuccess)
		}

		resp.SetHeader(gateway.HeaderUpstreamHttpCode, strconv.Itoa(respHeaders.StatusCode()))
	}
	return resp, nil
}
