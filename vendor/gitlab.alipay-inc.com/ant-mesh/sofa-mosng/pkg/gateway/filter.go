package gateway

import (
	"context"
	"runtime/debug"

	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/pkg/buffer"
)

func init() {
	api.RegisterStream("gateway", CreateGatewayFilterFactory)
}

// filterConfigFactory is an implement of api.StreamFilterChainFactory
type GatewayFilterConfigFactory struct {
	config *GatewayFilterConfig
}

type GatewayFilter struct {
	context            context.Context
	pipeline           Pipeline
	downstreamProtocol string
	receiver           api.StreamReceiverFilterHandler
	sender             api.StreamSenderFilterHandler
}

// CreateFilterChain will be invoked in echo request in proxy.NewStreamDetect function if filter has been injected
func (factory *GatewayFilterConfigFactory) CreateFilterChain(context context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewGatewayFilter(context, factory.config)
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	callbacks.AddStreamSenderFilter(filter)
}

// CreateGatewayFilterFactory will be invoked once in mosn init phase
// The filter injection will be skipped if function return is (nil, error)
func CreateGatewayFilterFactory(conf map[string]interface{}) (api.StreamFilterChainFactory, error) {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("create gateway stream filter factory")
	}
	cfg, err := ParseGatewayStreamFilter(conf)
	if err != nil {
		return nil, err
	}
	GetGatewayManager().SetFilterConfig(cfg)
	return &GatewayFilterConfigFactory{cfg}, nil
}

func NewGatewayFilter(ctx context.Context, config *GatewayFilterConfig) *GatewayFilter {
	handlers := config.Handlers
	pl := NewPipeline(handlers)
	filter := &GatewayFilter{
		pipeline:           pl,
		downstreamProtocol: config.DownstreamProtocol,
	}
	return filter
}

func (f *GatewayFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	f.receiver.SetConvert(false)

	//1. init gateway context
	f.context = context.WithValue(ctx, GATEWAY_CONTEXT_NAME, BuildGatewayContext())
	// hack: 传递 mosnctx，用于传递心跳保活协议
	f.context = context.WithValue(f.context, "MOSN_CONTEXT", ctx)

	defer func() {
		if err := recover(); err != nil {
			log.Proxy.Alertf(ctx, "gateway_req_error", "[gateway][%s][filter] Failed to process downstream request, errorMsg=[%+v]\n%s", GetGatewayContext(f.context).UniqueId(), err, string(debug.Stack()))
			f.recover(f.context, err)
		}
	}()
	// 2. decode downstream request
	gateReq, err := f.decodeDownstreamReq(f.context)
	GetGatewayContext(f.context).SetUniqueId(gateReq.RequestId())
	GetGatewayContext(f.context).SetRequest(gateReq)

	if err != nil {
		f.recover(f.context, err)
		return api.StreamFilterStop
	}

	if log.Proxy.GetLogLevel() >= log.DEBUG {
		log.Proxy.Debugf(ctx, "[%s][proxy][gateway] downstream request content, %+v", gateReq.RequestId(), gateReq)
	}

	//3. execute the IN part of gateway handler and handle exception
	if stop, err := f.pipeline.RunInHandlers(f.context); stop {
		if exception, ok := err.(*GatewayError); ok {
			if log.Proxy.GetLogLevel() >= log.INFO {
				log.Proxy.Infof(ctx, "[gateway][%s][filter] gateway pipeline handleIn throw exception %+v", gateReq.RequestId(), exception)
			}
			f.handleException(f.context, exception)
		}
		return api.StreamFilterStop
	}

	f.invokeUpstream(f.context)
	return api.StreamFilterContinue
}

func (f *GatewayFilter) Append(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	//1. init gateway context
	f.context = context.WithValue(ctx, GATEWAY_CONTEXT_NAME, f.context.Value(GATEWAY_CONTEXT_NAME))
	defer func() {
		if err := recover(); err != nil {
			log.Proxy.Alertf(ctx, "gateway_resp_error", "[gateway][%s][filter] Failed to process upstream response, %+v\n%s", GetGatewayContext(f.context).UniqueId(), err, string(debug.Stack()))
			f.recover(f.context, err)
		}
	}()

	//1. encode downstream response
	//设置 upstream 响应结果码，用于网关准确解析响应语义
	GetGatewayContext(f.context).SetAttribute(GatewayAttrUpRespCode, f.sender.RequestInfo().ResponseCode())
	gwResp, err := f.decodeUpstreamResp(f.context)
	GetGatewayContext(f.context).SetResponse(gwResp)

	if log.Proxy.GetLogLevel() >= log.DEBUG {
		log.Proxy.Debugf(ctx, "[gateway][%s][filter] upstream response content, %+v", GetGatewayContext(f.context).UniqueId(), gwResp)
	}

	if err != nil {
		f.recover(f.context, err)
		return api.StreamFilterStop
	}

	//2. execute handler out
	f.pipeline.RunOutHandlers(f.context)

	//3. respond downstream
	f.respondDownstream(f.context, false)
	return api.StreamFilterContinue
}

func (f *GatewayFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.receiver = handler
}

func (f *GatewayFilter) SetSenderFilterHandler(handler api.StreamSenderFilterHandler) {
	f.sender = handler
}

func (f *GatewayFilter) OnDestroy() {
	//noop
}

func (f *GatewayFilter) recover(ctx context.Context, err interface{}) {
	if log.Proxy.GetLogLevel() >= log.ERROR {
		log.Proxy.Errorf(ctx, "[gateway][%s][filter] Recover Gateway process, %+v", GetGatewayContext(ctx).UniqueId(), err)
	}

	var gatewayError *GatewayError
	if exception, ok := err.(*GatewayError); ok {
		gatewayError = exception
	} else {
		gatewayError = NewGatewayError(UnknownError)
	}

	f.handleException(ctx, gatewayError)
}

func (f *GatewayFilter) handleException(ctx context.Context, exception *GatewayError) {
	response := NewGatewayResponse(exception.Headers, exception.DataBuf, exception.Trailers)
	response.SetResultStatus(exception.ResultStatus)
	gwCtx := GetGatewayContext(ctx)
	if gwCtx != nil {
		response.SetDataEncoding(gwCtx.Request().DataEncoding())
		gwCtx.SetResponse(response)
	}

	//2. execute the left part of gateway OUT handler
	f.pipeline.RunOutHandlers(ctx)

	//3. respond downstream
	f.respondDownstream(ctx, true)
}

func (f *GatewayFilter) invokeUpstream(ctx context.Context) {
	gwCtx := GetGatewayContext(ctx)
	if log.Proxy.GetLogLevel() >= log.DEBUG {
		log.Proxy.Debugf(ctx, "[gateway][%s][filter] upstream request content, %+v", gwCtx.UniqueId(), gwCtx.Request())
	}

	headers, bodyBuffer, trailers := f.encodeUpstreamReq(ctx)
	f.receiver.SetRequestHeaders(headers)
	f.receiver.SetRequestData(bodyBuffer)
	f.receiver.SetRequestTrailers(trailers)
}

func (f *GatewayFilter) respondDownstream(ctx context.Context, isHijack bool) {
	if log.Proxy.GetLogLevel() >= log.DEBUG {
		log.Proxy.Debugf(ctx, "[gateway][%s][filter] downstream response content, %+v", GetGatewayContext(ctx).UniqueId(), GetGatewayContext(ctx).Response())
	}

	headers, bodyBuffer, trailers := f.encodeDownstreamResp(ctx)

	if isHijack {
		if bodyBuffer == nil {
			f.receiver.AppendHeaders(headers, true)
		} else {
			f.receiver.AppendHeaders(headers, false)
			f.receiver.AppendData(bodyBuffer, false)
			f.receiver.AppendTrailers(trailers)
		}
	} else {
		f.sender.SetResponseHeaders(headers)
		f.sender.SetResponseData(bodyBuffer)
		f.sender.SetResponseTrailers(trailers)
	}
}

func (f *GatewayFilter) encodeDownstreamResp(ctx context.Context) (api.HeaderMap, buffer.IoBuffer, api.HeaderMap) {
	codec := GetDownstreamCodec(f.downstreamProtocol)
	mapping := GetDownstreamStatusMapping(f.downstreamProtocol)
	return codec.Encode(ctx, mapping)
}

func (f *GatewayFilter) decodeDownstreamReq(ctx context.Context) (GatewayRequest, error) {
	codec := GetDownstreamCodec(f.downstreamProtocol)
	return codec.Decode(ctx, f.receiver.GetRequestHeaders(), f.receiver.GetRequestData(), f.receiver.GetRequestTrailers())
}

func (f *GatewayFilter) encodeUpstreamReq(ctx context.Context) (api.HeaderMap, buffer.IoBuffer, api.HeaderMap) {
	protocol := string(GetGatewayContext(ctx).Service().Upstream().Protocol())
	codec := GetUpstreamCodec(protocol)
	return codec.Encode(ctx)
}

func (f *GatewayFilter) decodeUpstreamResp(ctx context.Context) (GatewayResponse, error) {
	protocol := string(GetGatewayContext(ctx).Service().Upstream().Protocol())
	codec := GetUpstreamCodec(protocol)
	mapping := GetUpstreamStatusMapping(protocol)
	return codec.Decode(ctx, f.sender.GetResponseHeaders(), f.sender.GetResponseData(), f.sender.GetResponseTrailers(), mapping)
}
