package gateway

import (
	"strings"

	"github.com/valyala/fasthttp"
	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol/http"
	"mosn.io/pkg/buffer"
)

type DefaultGatewayResponse struct {
	headers      api.HeaderMap
	dataBuf      buffer.IoBuffer
	trailers     api.HeaderMap
	dataEncoding DataEncodingType
	resultStatus ResponseStatus
}

func NewGatewayResponse(headers api.HeaderMap, dataBuf buffer.IoBuffer, trailers api.HeaderMap) *DefaultGatewayResponse {
	//为兼容客户端解析响应头的逻辑（如 RpcId 等），显式声明禁止转换响应头
	headerImpl := &fasthttp.ResponseHeader{}
	headerImpl.DisableNormalizing()
	httpHeaders := http.ResponseHeader{headerImpl, nil}

	// 做转换，避免入参 headers 不支持 add 方法
	if headers != nil {
		headers.Range(func(key, value string) bool {
			// 不增加默认的 Content-Type，规避多个 Content-Type 响应头引发的乱码问题
			if !strings.EqualFold(key, HeaderContentType) {
				httpHeaders.Add(key, value)
			}
			return true
		})
	}

	if dataBuf == nil {
		dataBuf = buffer.NewIoBuffer(0)
	}

	return &DefaultGatewayResponse{
		headers:      httpHeaders,
		dataBuf:      dataBuf,
		trailers:     trailers,
		resultStatus: PermissionDeny,
	}
}

func (r *DefaultGatewayResponse) Headers() api.HeaderMap {
	return r.headers
}

func (r *DefaultGatewayResponse) GetHeader(headerKey string) string {
	val, _ := r.headers.Get(headerKey)
	return val
}

func (r *DefaultGatewayResponse) SetHeader(headerKey, headerVal string) {
	r.headers.Set(headerKey, headerVal)
}

func (r *DefaultGatewayResponse) AddHeader(headerKey, headerVal string) {
	r.headers.Add(headerKey, headerVal)
}

func (r *DefaultGatewayResponse) DataBuf() buffer.IoBuffer {
	return r.dataBuf
}

func (r *DefaultGatewayResponse) DataBytes() []byte {
	return r.dataBuf.Bytes()
}

func (r *DefaultGatewayResponse) SetDataBytes(data []byte) (n int, err error) {
	r.dataBuf.Reset()
	r.headers.Del("Content-Length")
	r.headers.Del("Content-Encoding")
	return r.dataBuf.Write(data)
}

func (r *DefaultGatewayResponse) Trailers() api.HeaderMap {
	return r.trailers
}

func (r *DefaultGatewayResponse) GetTrailer(trailerKey string) string {
	val, _ := r.trailers.Get(trailerKey)
	return val
}

func (r *DefaultGatewayResponse) SetTrailer(trailerKey, trailerVal string) {
	r.trailers.Set(trailerKey, trailerVal)
}

func (r *DefaultGatewayResponse) DataEncoding() DataEncodingType {
	return r.dataEncoding
}

func (r *DefaultGatewayResponse) SetDataEncoding(encoding DataEncodingType) {
	r.dataEncoding = encoding
}

func (r *DefaultGatewayResponse) ResultStatus() ResponseStatus {
	return r.resultStatus
}

func (r *DefaultGatewayResponse) SetResultStatus(resultStatus ResponseStatus) {
	r.resultStatus = resultStatus
}
