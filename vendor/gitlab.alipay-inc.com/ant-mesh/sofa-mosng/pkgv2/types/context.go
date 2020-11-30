package types

import (
	ctx "context"

	"github.com/valyala/fasthttp"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/constants"

	"mosn.io/mosn/pkg/protocol"
	"mosn.io/mosn/pkg/trace"
	"mosn.io/pkg/buffer"

	"mosn.io/mosn/pkg/protocol/http"
	"mosn.io/mosn/pkg/types"
	mosnt "mosn.io/mosn/pkg/types"
)

// todo
type Context interface {
	ctx.Context
	InitRequest(headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap)
	Request() Request
	InitResponse(headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap)
	Response() Response
	Router() Router
	SetRouter(router Router)
	GetAttribute(key AttributeKey) interface{}
	SetAttribute(key AttributeKey, value interface{})
	DownStreamProtocol() mosnt.ProtocolName
	UpstreamProtocol() mosnt.ProtocolName
	SetDownStreamProtocol(mosnt.ProtocolName)
	SetUpstreamProtocol(mosnt.ProtocolName)
	TraceId() string
	SetTraceId(string)
}

type context struct {
	ctx.Context
	request      Request
	response     Response
	attributes   map[AttributeKey]interface{}
	router       Router
	downProtocol mosnt.ProtocolName
	upProtocol   mosnt.ProtocolName
	traceId      string
}

func (c *context) TraceId() string {
	return c.traceId
}

func (c *context) SetTraceId(traceId string) {
	c.traceId = traceId
}

func (c *context) GetAttribute(key AttributeKey) interface{} {
	return c.attributes[key]
}

func (c *context) DownStreamProtocol() mosnt.ProtocolName {
	if c.downProtocol == "" {
		return protocol.HTTP1
	}
	return c.downProtocol
}
func (c *context) UpstreamProtocol() mosnt.ProtocolName {
	if c.upProtocol == "" {
		return protocol.HTTP1
	}
	return c.upProtocol
}

func (c *context) SetDownStreamProtocol(p mosnt.ProtocolName) {
	c.downProtocol = p
}
func (c *context) SetUpstreamProtocol(p mosnt.ProtocolName) {
	c.upProtocol = p
}

func (c *context) SetAttribute(key AttributeKey, value interface{}) {
	// todo mux?
	c.attributes[key] = value
}

func (c *context) InitRequest(headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) {
	req := &request{
		baseReqRes: baseReqRes{
			headers:  headers,
			buf:      buf,
			trailers: trailers,
			// todo
			dataEncoding: APPLICATION_JSON,
		},
	}
	c.request = req
}

func (c *context) Request() Request {
	return c.request
}

func (c *context) InitResponse(headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) {
	if c.downProtocol == protocol.HTTP1 {
		if c.upProtocol != protocol.HTTP1 {
			// 后端协议不是 http，忽略 response header
			c.response = NewHttpResponse(nil, buf, trailers)
		} else {
			c.response = NewHttpResponse(headers, buf, trailers)
		}
	} else {
		c.response = NewResponse(headers, buf, trailers)
	}
}

func (c *context) Response() Response {
	return c.response
}

func (c *context) Router() Router {
	return c.router
}

func (c *context) SetRouter(router Router) {
	c.router = router
}

func BuildContext(parentCtx ctx.Context) Context {
	traceID := trace.IdGen().GenerateTraceId()
	c := ctx.WithValue(parentCtx, constants.ContextKeyTraceId, traceID)
	return &context{
		Context:    c,
		attributes: make(map[AttributeKey]interface{}),
		traceId:    traceID,
		response:   NewHttpResponse(nil, nil, nil),
	}
}

type AttributeKey string

type DataEncodingType string

const (
	APPLICATION_JSON DataEncodingType = "application/json"
	PROTOBUF         DataEncodingType = "application/protobuf"
)

type BaseReqRes interface {
	GetHeaders() types.HeaderMap

	SetHeaders(headers types.HeaderMap)

	GetHeader(key string) string

	SetHeader(key, value string)

	DelHeader(key string)

	GetDataBuf() types.IoBuffer

	SetDataBuf(buf types.IoBuffer)

	GetDataBytes() []byte

	SetDataBytes(data []byte) (n int, err error)

	GetTrailers() types.HeaderMap

	SetTrailers(trailers types.HeaderMap)

	GetTrailer(trailerKey string) string

	SetTrailer(trailerKey, trailerVal string)

	GetAttribute(key AttributeKey) interface{}

	SetAttribute(key AttributeKey, value interface{})

	DataEncoding() DataEncodingType
	//
	SetDataEncoding(encoding DataEncodingType)

	AppendHeaders(headers types.HeaderMap)
}

type baseReqRes struct {
	headers      types.HeaderMap
	attributes   map[AttributeKey]interface{}
	buf          types.IoBuffer
	trailers     types.HeaderMap
	dataEncoding DataEncodingType
}

func (c *baseReqRes) SetDataEncoding(dataEncoding DataEncodingType) {
	c.dataEncoding = dataEncoding
}

func (c *baseReqRes) DataEncoding() DataEncodingType {
	return c.dataEncoding
}

func (c *baseReqRes) DelHeader(key string) {
	c.headers.Del(key)
}

func (c *baseReqRes) GetHeaders() types.HeaderMap {
	return c.headers
}

func (c *baseReqRes) SetHeaders(headers types.HeaderMap) {
	c.headers = headers
}

func (c *baseReqRes) GetHeader(key string) string {
	if val, ok := c.headers.Get(key); ok {
		return val
	}
	return ""
}

func (c *baseReqRes) AppendHeaders(headers types.HeaderMap) {
	if headers != nil {
		headers.Range(func(key, value string) bool {
			c.SetHeader(key, value)
			return true
		})
	}
}

func (c *baseReqRes) SetHeader(key, value string) {
	if c.headers == nil {
		c.headers = &http.RequestHeader{
			RequestHeader:     &fasthttp.RequestHeader{},
			EmptyValueHeaders: map[string]bool{},
		}
	}
	c.headers.Set(key, value)
}

func (c *baseReqRes) GetDataBuf() types.IoBuffer {
	return c.buf
}

func (c *baseReqRes) SetDataBuf(buf types.IoBuffer) {
	c.buf = buf
}

func (c *baseReqRes) GetDataBytes() []byte {
	if c.buf == nil {
		return nil
	}
	return c.buf.Bytes()
}

func (c *baseReqRes) SetDataBytes(data []byte) (n int, err error) {
	if data == nil {
		return
	}

	if c.buf == nil {
		c.buf = buffer.NewIoBuffer(len(data))
	}
	c.buf.Reset()
	return c.buf.Write(data)
}

func (c *baseReqRes) GetTrailers() types.HeaderMap {
	return c.trailers
}

func (c *baseReqRes) SetTrailers(trailers types.HeaderMap) {
	c.trailers = trailers
}

func (c *baseReqRes) GetTrailer(key string) string {
	if c.trailers == nil {
		return ""
	}
	if val, ok := c.trailers.Get(key); ok {
		return val
	}
	return ""
}

func (c *baseReqRes) SetTrailer(key, val string) {
	if c.trailers == nil {
		c.trailers = &http.RequestHeader{
			RequestHeader:     &fasthttp.RequestHeader{},
			EmptyValueHeaders: map[string]bool{},
		}
	}
	c.trailers.Set(key, val)
}

func (c *baseReqRes) GetAttribute(key AttributeKey) interface{} {
	return c.attributes[key]
}

func (c *baseReqRes) SetAttribute(key AttributeKey, value interface{}) {
	c.attributes[key] = value
}

type Request interface {
	BaseReqRes
	GetQueryString() map[string]string
}

type request struct {
	baseReqRes
	queryString map[string]string
}

func (c *request) SetHeader(key, value string) {
	if c.headers == nil {
		c.headers = &http.RequestHeader{
			RequestHeader:     &fasthttp.RequestHeader{},
			EmptyValueHeaders: map[string]bool{},
		}
	}
	c.headers.Set(key, value)
}

func (c *request) GetQueryString() map[string]string {
	if c.queryString == nil {
		if qs, ok := c.headers.Get(protocol.MosnHeaderQueryStringKey); ok {
			c.queryString = http.ParseQueryString(qs)
		} else {
			return map[string]string{}
		}
	}

	return c.queryString
}

type Response interface {
	BaseReqRes

	StatusCode() int

	SetStatusCode(httpCode int)
}

type ResponseStatus string

type response struct {
	baseReqRes
	httpCode int
}

func (r *response) StatusCode() int {
	return r.httpCode
}

func (r *response) SetStatusCode(httpCode int) {
	r.httpCode = httpCode
}

func (c *response) SetHeader(key, value string) {
	if c.headers == nil {
		c.headers = &http.ResponseHeader{
			ResponseHeader:    &fasthttp.ResponseHeader{},
			EmptyValueHeaders: map[string]bool{},
		}
	}
	c.headers.Set(key, value)
}

func (c *response) SetTrailer(key, val string) {
	if c.trailers == nil {
		c.trailers = &http.ResponseHeader{
			ResponseHeader:    &fasthttp.ResponseHeader{},
			EmptyValueHeaders: map[string]bool{},
		}
	}
	c.trailers.Set(key, val)
}

func NewHttpResponse(headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) Response {
	headerImpl := &fasthttp.ResponseHeader{}
	headerImpl.DisableNormalizing()
	httpHeaders := http.ResponseHeader{ResponseHeader: headerImpl, EmptyValueHeaders: map[string]bool{}}

	if headers != nil {
		headers.Range(func(key, value string) bool {
			httpHeaders.Add(key, value)
			return true
		})
	}

	resp := &response{
		baseReqRes: baseReqRes{
			headers:  httpHeaders,
			buf:      buf,
			trailers: trailers,
		},
	}

	return resp
}

func NewResponse(headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) Response {
	resp := &response{
		baseReqRes: baseReqRes{
			headers:  headers,
			buf:      buf,
			trailers: trailers,
		},
	}

	return resp
}
