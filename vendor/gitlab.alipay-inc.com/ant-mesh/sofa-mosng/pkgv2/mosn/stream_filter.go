package mosn

import (
	"context"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"mosn.io/mosn/pkg/protocol"
	mt "mosn.io/mosn/pkg/types"
	"strconv"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/errors"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/router"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
	mosn "mosn.io/api"

	"net/http"
)

func init() {
	mosn.RegisterStream("mosng", CreateGwStreamFilterFactory)
}

type StreamFilter struct {
	context       types.Context
	receiver      mosn.StreamReceiverFilterHandler
	sender        mosn.StreamSenderFilterHandler
	err           types.GatewayError
	receiveStatus bool
	senderStatus  bool
}

type StreamFilterFactory struct{}

func (sff *StreamFilterFactory) CreateFilterChain(context context.Context, callbacks mosn.StreamFilterChainFactoryCallbacks) {
	f := NewGatewayFilter()
	callbacks.AddStreamReceiverFilter(f, mosn.BeforeRoute)
	callbacks.AddStreamSenderFilter(f)
}

func CreateGwStreamFilterFactory(conf map[string]interface{}) (mosn.StreamFilterChainFactory, error) {
	return &StreamFilterFactory{}, nil
}

func (f *StreamFilter) OnReceive(ctx context.Context, headers mosn.HeaderMap, buf mt.IoBuffer, trailers mosn.HeaderMap) mosn.StreamFilterStatus {
	defer func() {
		if err := recover(); err != nil {
			log.ProxyLogger().Errorf("[%s][StreamFilter][OnReceive][recover] invoke OnReceive error: %v", f.context.TraceId(), errors.Errorf("%v", err))
			f.recover(err, false)
		}
	}()

	if f.receiveStatus {
		return mosn.StreamFilterContinue
	}

	f.receiveStatus = true

	// build ctx
	f.context = types.BuildContext(ctx)
	f.context.InitRequest(headers, buf, trailers)

	name := ctx.Value(mt.ContextKeyListenerName)
	if name == nil {
		panic(errors.NewWithMsg(http.StatusBadGateway, "no listener name for this request"))
	}

	if server := config.StoreInstance().Get(api.GATEWAY, name.(string)); server != nil {
		gateway := server.(*v1.Gateway)

		f.context.SetUpstreamProtocol(mosn.Protocol(gateway.UpstreamProtocol))
		f.context.SetDownStreamProtocol(mosn.Protocol(gateway.DownstreamProtocol))
	} else {
		panic(errors.NewWithMsg(http.StatusBadGateway, "no listener instance for this request"))
	}

	// match router
	r := router.GetRouterManager().Match(name.(string), headers)
	if r == nil {
		panic(errors.NewWithMsg(http.StatusBadRequest, "no router match for this request"))
	}

	f.context.SetRouter(r)

	// run in bound filters
	if r.Pipeline() != nil {
		if err := r.Pipeline().DoInBound(f.context); err != nil {
			log.ProxyLogger().Errorf("[%s][StreamFilter][OnReceive][DoInBound] invoke inbound pipeline err: %v", f.context.TraceId(), err)
			f.recover(err, false)
			return mosn.StreamFilterStop
		}
	} else {
		log.ProxyLogger().Infof("[%s][StreamFilter][OnReceive][DoInBound] no filter found: %v", f.context.TraceId(), r)
	}

	// set service for proxy
	if err := r.Service().Upstream().Invoke(f.context); err != nil {
		log.ProxyLogger().Errorf("[%s][StreamFilter][OnReceive][Invoke] invoke service err: %v", f.context.TraceId(), err)
		f.recover(err, false)
		return mosn.StreamFilterStop
	}

	// end receive
	f.endReceive()

	return mosn.StreamFilterContinue
}

func (f *StreamFilter) Append(ctx context.Context, headers mosn.HeaderMap, buf mt.IoBuffer, trailers mosn.HeaderMap) mosn.StreamFilterStatus {
	defer func() {
		if err := recover(); err != nil {
			log.ProxyLogger().Errorf("[%s][StreamFilter][Append][recover] invoke Append error: %v", f.context.TraceId(), errors.Errorf("%v", err))
			f.recover(err, true)
		}
	}()

	if f.senderStatus {
		return mosn.StreamFilterContinue
	}

	f.senderStatus = true

	if f.context == nil {
		panic(errors.NewWithMsg(http.StatusBadRequest, "no context found"))
	}

	// init response
	f.context.InitResponse(headers, buf, trailers)

	if f.err != nil {
		f.endWithErr()
		return mosn.StreamFilterStop
	}

	if r := f.context.Router(); r != nil {
		// upstream receive response
		if err := r.Service().Upstream().Receive(f.context); err != nil {
			log.ProxyLogger().Errorf("[%s][StreamFilter][Append][Receive] receive service err: %v", f.context.TraceId(), err)
			f.recover(err, true)
			return mosn.StreamFilterStop
		}

		// run out bound filters
		if r.Pipeline() != nil {
			if err := r.Pipeline().DoOutBound(f.context); err != nil {
				log.ProxyLogger().Errorf("[%s][StreamFilter][Append][DoOutBound] invoke outbound pipeline err: %v", f.context.TraceId(), err)
				f.recover(err, true)
				return mosn.StreamFilterStop
			}
		}

		f.endAppend()

		return mosn.StreamFilterContinue
	}

	log.ProxyLogger().Errorf("[%s][StreamFilter][Append] no router found when invoke Append", f.context.TraceId())
	return mosn.StreamFilterStop
}

func (f *StreamFilter) SetReceiveFilterHandler(handler mosn.StreamReceiverFilterHandler) {
	f.receiver = handler
}

func (f *StreamFilter) SetSenderFilterHandler(handler mosn.StreamSenderFilterHandler) {
	f.sender = handler
}

func (f *StreamFilter) OnDestroy() {
	// noop
}

func (f *StreamFilter) Context() types.Context {
	return f.context
}

func (f *StreamFilter) endReceive() {
	f.receiver.SetRequestHeaders(f.context.Request().GetHeaders())
	f.receiver.SetRequestData(f.context.Request().GetDataBuf())
	f.receiver.SetRequestTrailers(f.context.Request().GetTrailers())
}

func (f *StreamFilter) endAppend() {
	f.sender.SetResponseHeaders(f.context.Response().GetHeaders())
	f.sender.SetResponseData(f.context.Response().GetDataBuf())
	f.sender.SetResponseTrailers(f.context.Response().GetTrailers())
}

func (f *StreamFilter) recover(err interface{}, hasProxy bool) {
	if gError, ok := err.(types.GatewayError); ok {
		f.onError(gError, hasProxy)
	} else {
		f.onError(errors.New(http.StatusInternalServerError), hasProxy)
	}
}

func (f *StreamFilter) onError(err types.GatewayError, hasProxy bool) {
	f.err = err

	if hasProxy {
		f.endWithErr()
	}
}

func (f *StreamFilter) endWithErr() {
	defer func() {
		if err := recover(); err != nil {
			log.ProxyLogger().Errorf("[%s][StreamFilter][onError] cause error again after recover: %v", f.context.TraceId(), err)
		}
	}()

	var (
		httpCode = 502
		headers  protocol.CommonHeader
		bytes    []byte
	)

	if c := f.context; c == nil {
		f.context = types.BuildContext(nil)
		f.context.InitResponse(nil, nil, nil)
	}

	if r := f.context.Router(); r != nil {
		if err := r.Pipeline().DoOutBound(f.context); err != nil {
			log.ProxyLogger().Errorf("[%s][StreamFilter][onError][DoOutBound] invoke outbound pipeline err again: %v", f.context.TraceId(), err)
		}
		httpCode, headers, bytes = r.Pipeline().HandleErr(f.context, f.err)
	}

	res := f.context.Response()
	res.AppendHeaders(headers)
	res.SetHeader(mt.HeaderStatus, strconv.Itoa(httpCode))
	res.SetDataBytes(bytes)

	f.endAppend()
}

func NewGatewayFilter() *StreamFilter {
	return &StreamFilter{
		receiveStatus: false,
		senderStatus:  false,
	}
}
