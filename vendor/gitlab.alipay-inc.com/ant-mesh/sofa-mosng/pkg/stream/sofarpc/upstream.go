package sofarpc

import (
	"context"
	"strconv"

	"github.com/golang/protobuf/proto"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkg/common/utils"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkg/gateway"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkg/model"
	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/buffer"
)

const (
	SofaMosngServiceClass   = "com.alipay.mobilegw.adapterservice.mosng.MosngAdapterService:1.0"
	SofaMosngServiceMethod  = "mosngMobileService"
	SofaVipAddrPostfix      = "-pool.cz10a.alipay.net:12200"
	SofaTraceContextAttrKey = "rpc_trace_context"
	SofaTargetLocalAttrKey  = "sofa_target_local"
)

type SofaRpcUpstream struct {
	protocol    api.Protocol
	timeout     uint64
	clusterName string
	targetUrl   string
}

func init() {
	gateway.RegisterUpstreamParser(string(protocol.SofaRPC), parseConfig, 1)
	gateway.RegisterUpstreamCodec(string(protocol.SofaRPC), &SofaRpcUpstream{}, 1)
	gateway.RegisterUpstreamStatusMapping(string(protocol.SofaRPC), MappingUpstreamResponseStatus, 1)
}

func (u *SofaRpcUpstream) Protocol() api.Protocol {
	return u.protocol
}

func (u *SofaRpcUpstream) SetProtocol(protocol api.Protocol) {
	u.protocol = protocol
}

func (u *SofaRpcUpstream) Timeout() uint64 {
	return u.timeout
}

func (u *SofaRpcUpstream) SetTimeout(timeout uint64) {
	u.timeout = timeout
}

func (u *SofaRpcUpstream) ClusterName() string {
	return u.clusterName
}

func (u *SofaRpcUpstream) SetClusterName(cluster string) {
	u.clusterName = cluster
}

func (u *SofaRpcUpstream) TargetUrl() string {
	return u.targetUrl
}

func (u *SofaRpcUpstream) SetTargetUrl(targetUrl string) {
	u.targetUrl = targetUrl
}

func parseConfig(cfg gateway.UpstreamConfig) gateway.Upstream {
	upt := &SofaRpcUpstream{}
	upt.SetProtocol(cfg.Protocol)
	upt.SetClusterName(cfg.ClusterName)
	upt.SetTimeout(cfg.TimeOut)
	if targetUrl, ok := cfg.Config["target_url"]; ok {
		upt.SetTargetUrl(targetUrl.(string))
	}
	return upt
}

func (huc *SofaRpcUpstream) Encode(context context.Context) (api.HeaderMap, buffer.IoBuffer, api.HeaderMap) {
	var headerPbList []*model.RequestDefaultMapEntry
	gwCtx := gateway.GetGatewayContext(context)
	request := gwCtx.Request()
	request.Headers().Range(func(key, value string) bool {
		item := &model.RequestDefaultMapEntry{
			Key:   proto.String(key),
			Value: proto.String(value),
		}
		headerPbList = append(headerPbList, item)
		return true
	})

	upstreamReq := &model.MobileRpcRequestPB{
		UniqueId:            proto.String(request.RequestId()),
		OperationType:       proto.String(request.ApiId()),
		RequestDataBytes:    request.DataBytes(),
		RequestDataEncoding: proto.Int32(gateway.GetDataEncodingValue(request.DataEncoding())),
		Headers:             headerPbList,
	}

	content, err := proto.Marshal(upstreamReq)
	if err != nil {
		return nil, nil, nil
	}

	headerMap := buildRequestHeader(gwCtx)

	//set timeout
	timeout := calculateReqTimeout(gwCtx.Service())
	headerMap[types.HeaderGlobalTimeout] = strconv.Itoa(timeout)
	boltReq := utils.BuildSofaRequestWithTimeout(headerMap, len(content), timeout)

	return &boltReq, buffer.NewIoBufferBytes(content), nil
}

func (huc *SofaRpcUpstream) Decode(ctx context.Context, headers api.HeaderMap, dataBuf buffer.IoBuffer, trailers api.HeaderMap, mapping gateway.UpstreamStatusMapping) (gateway.GatewayResponse, error) {
	resp := gateway.NewGatewayResponse(nil, dataBuf, trailers)

	if mapping != nil {
		status := mapping(ctx, headers)
		if status != "" {
			return resp, gateway.NewGatewayError(status)
		}
	}

	upstreamResp := &model.MobileRpcResponsePB{}
	if dataBuf == nil {
		return resp, gateway.NewGatewayError(gateway.BizException)
	}

	err := proto.Unmarshal(dataBuf.Bytes(), upstreamResp)
	if err != nil {
		return resp, gateway.NewGatewayError(gateway.DataParserException)
	}

	pbHeaders := upstreamResp.Headers
	for _, h := range pbHeaders {
		resp.SetHeader(h.GetKey(), h.GetValue())
	}

	resp.SetDataBytes(upstreamResp.ResponseDataBytes)
	resp.SetDataEncoding(gateway.GetGatewayContext(ctx).Request().DataEncoding())
	resultStatus := upstreamResp.GetResultStatus()
	status := model.GetByResultCode(int(resultStatus))
	resp.SetResultStatus(status.RespStatus)
	return resp, nil
}

func buildRequestHeader(context *gateway.GatewayContext) map[string]string {
	headerMap := make(map[string]string)
	cluster := context.Service().Upstream().ClusterName()
	headerMap[utils.SofaHeaderTargetService] = SofaMosngServiceClass + ":" + cluster
	headerMap[utils.SofaHeaderService] = SofaMosngServiceClass + ":" + cluster
	headerMap[utils.SofaHeaderMethodName] = SofaMosngServiceMethod
	headerMap[utils.SofaHeaderVipAddr] = cluster + SofaVipAddrPostfix

	if traceCtx, ok := context.GetAttribute(SofaTraceContextAttrKey).(map[string]string); ok {
		for k, v := range traceCtx {
			headerMap[k] = v
		}
	}

	// set local invoke flag
	if target, ok := context.GetAttribute(SofaTargetLocalAttrKey).(string); ok {
		headerMap[utils.SofaHeaderTargetLocal] = target
	}
	return headerMap
}

func calculateReqTimeout(svc gateway.Service) int {
	serviceTimeout := svc.Timeout()
	upstreamTimeout := svc.Upstream().Timeout()

	if serviceTimeout > 0 {
		return int(serviceTimeout)
	}

	if upstreamTimeout > 0 {
		return int(upstreamTimeout)
	}

	return 3000
}

func MappingUpstreamResponseStatus(ctx context.Context, headers api.HeaderMap) gateway.ResponseStatus {
	var status gateway.ResponseStatus
	if headers == nil {
		return status
	}

	if code, ok := headers.Get(types.HeaderStatus); ok {
		if codeInt, err := strconv.Atoi(code); err == nil {
			switch codeInt {
			case types.RouterUnavailableCode:
				status = gateway.ServiceNotFound
			case types.NoHealthUpstreamCode:
				status = gateway.RemoteAccessException
			case types.UpstreamOverFlowCode:
				status = gateway.RemoteAccessException
			case types.TimeoutExceptionCode:
				status = gateway.RequestTimeOut
			default:
				return status
			}
		}
	}

	return status
}
