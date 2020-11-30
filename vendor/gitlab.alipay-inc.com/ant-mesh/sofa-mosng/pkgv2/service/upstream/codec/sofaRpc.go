package codec

import (
	"bytes"
	"strconv"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkg/common/utils"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/constants"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol"
	"mosn.io/mosn/pkg/protocol/http"
	"mosn.io/mosn/pkg/trace"
	"mosn.io/mosn/pkg/trace/sofa"
	"mosn.io/mosn/pkg/trace/sofa/xprotocol"
	mosn "mosn.io/mosn/pkg/types"
	"mosn.io/pkg/buffer"
)

var (
	holderMap map[string]UpstreamReqHeaderHolder
)

type UpstreamReqHeaderHolder struct {
	service string
	vip     string
}

type sofaRpcCodec struct {
	protocol  string
	convertor types.ProtocolConvertor
}

func (src *sofaRpcCodec) Convertor() types.ProtocolConvertor {
	return src.convertor
}

func init() {
	GetCodeFactoryInstance().Register(&sofaRpcCodec{protocol: "SofaRpc", convertor: &sofaRpcProtocolConvertor{}})
	holderMap = make(map[string]UpstreamReqHeaderHolder, 1)
}

func (src *sofaRpcCodec) Protocol() string {
	return src.protocol
}

func (src *sofaRpcCodec) Encode(ctx types.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) {
	headers.Set(constants.MosngServiceHeader, ctx.Router().Service().Name())
	headers.Set(protocol.MosnHeaderPathKey, "/")
}

func (src *sofaRpcCodec) Decode(ctx types.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) {

}

type sofaRpcProtocolConvertor struct {
}

func (*sofaRpcProtocolConvertor) EncodeConv(ctx types.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) (convHeaders api.HeaderMap, convBuf buffer.IoBuffer, convTrailers api.HeaderMap) {
	headerMap := buildUpstreamReqHeader(ctx)

	//set timeout
	timeout := calculateReqTimeout(ctx.Router())
	headerMap[mosn.HeaderGlobalTimeout] = strconv.Itoa(timeout)

	var bufLen int
	if buf == nil {
		bufLen = 0
	} else {
		bufLen = buf.Len()
	}
	boltReq := utils.BuildSofaRequestWithTimeout(headerMap, bufLen, timeout)

	completeEgressSpan(ctx, headerMap[utils.SofaHeaderTargetService], ctx.Router().Conf().Proxy.Method)

	return &boltReq, buf, nil
}

func (*sofaRpcProtocolConvertor) DecodeConv(ctx types.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) (convHeaders api.HeaderMap, convBuf buffer.IoBuffer, convTrailers api.HeaderMap) {
	headers.Set(mosn.HeaderStatus, strconv.Itoa(MappingUpstreamResponseStatus(headers)))
	return headers, buf, trailers
}

func MappingUpstreamResponseStatus(headers api.HeaderMap) int {
	if headers == nil {
		return http.InternalServerError
	}

	if code, ok := headers.Get(mosn.HeaderStatus); ok {
		if codeInt, err := strconv.Atoi(code); err == nil {
			switch codeInt {
			case mosn.CodecExceptionCode:
				return http.BadGateway
			case mosn.UnknownCode:
				return http.InternalServerError
			case mosn.DeserialExceptionCode:
				return http.BadRequest
			default:
				return codeInt
			}
		}
	}

	return http.OK
}

func buildUpstreamReqHeader(context types.Context) map[string]string {
	headerMap := make(map[string]string, 20)
	proxy := context.Router().Conf().Proxy
	headerMap[sofa.TARGET_SERVICE_KEY] = proxy.Interface
	headerMap[sofa.SERVICE_KEY] = proxy.Interface
	headerMap[sofa.TARGET_METHOD] = proxy.Method
	headerMap[sofa.TRACER_ID_KEY] = context.TraceId()
	tracerCtx := GetTracerFromCtx(context)
	for key, value := range tracerCtx {
		headerMap[key] = value
	}
	// tracer
	sofaPenAttrStr := GetSofaPenAttrsStringFromCtx(context)
	headerMap[sofa.SOFA_TRACE_BAGGAGE_DATA] = sofaPenAttrStr

	return headerMap
}

func calculateReqTimeout(r types.Router) int {
	serviceTimeout := r.Conf().Timeout
	if serviceTimeout > 0 {
		return int(serviceTimeout)
	}

	upstreamTimeout := r.Service().Conf().Timeout
	if upstreamTimeout > 0 {
		return int(upstreamTimeout)
	}

	return 3000
}

func GetSofaPenAttrsStringFromCtx(ctx types.Context) string {
	penetrates := GetSofaPenAttrsFromCtx(ctx)
	var sb bytes.Buffer
	for k, v := range penetrates {
		sb.WriteString(escapePercentEqualAnd(k))
		sb.WriteString("=")
		sb.WriteString(escapePercentEqualAnd(v))
		sb.WriteString("&")
	}

	return sb.String()
}

func GetTracerFromCtx(ctx types.Context) map[string]string {
	if tracer, ok := ctx.GetAttribute(constants.AttrTracerCtx).(map[string]string); ok {
		return tracer
	} else {
		t := make(map[string]string, 8)
		ctx.SetAttribute(constants.AttrTracerCtx, t)
		return t
	}
}

func SetValue2Tracer(ctx types.Context, key, value string) {
	tracer := GetTracerFromCtx(ctx)
	tracer[key] = value
}

func GetValueFromTracer(ctx types.Context, key string) string {
	tracer := GetTracerFromCtx(ctx)
	return tracer[key]
}

func SetTracerToCtx(ctx types.Context, tracer map[string]string) {
	ctx.SetAttribute(constants.AttrTracerCtx, tracer)
}

func SetSofaPenAttrToCtx(ctx types.Context, key, value string) {
	attrs := GetSofaPenAttrsFromCtx(ctx)
	attrs[key] = value
}

func GetSofaPenAttrsFromCtx(ctx types.Context) map[string]string {
	if tracer, ok := ctx.GetAttribute(constants.AttrSofaPenAttrs).(map[string]string); ok {
		return tracer
	} else {
		t := make(map[string]string, 4)
		ctx.SetAttribute(constants.AttrSofaPenAttrs, t)
		return t
	}
}

func escapePercentEqualAnd(str string) string {
	var buf bytes.Buffer
	for _, cha := range str {
		switch cha {
		case '%':
			buf.WriteString("%25")
		case '&':
			buf.WriteString("%26")
		case '=':
			buf.WriteString("%3D")
		//case ',':
		//	buf.WriteString("%2C")
		default:
			buf.WriteRune(cha)
		}
	}

	return buf.String()
}

func completeEgressSpan(ctx types.Context, serviceName, method string) {
	span := trace.SpanFromContext(ctx)

	if span != nil {
		span.SetTag(xprotocol.SERVICE_NAME, serviceName)
		span.SetTag(xprotocol.METHOD_NAME, method)
		span.SetTag(xprotocol.PROTOCOL, "SofaRpc")
	}
}
