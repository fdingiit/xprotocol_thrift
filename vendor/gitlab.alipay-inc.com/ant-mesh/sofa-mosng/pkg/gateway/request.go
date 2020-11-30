package gateway

import (
	"mosn.io/api"
	"mosn.io/pkg/buffer"
)

type DefaultGatewayRequest struct {
	apiId        string
	headers      api.HeaderMap
	headerMap    map[string]string
	dataBuf      buffer.IoBuffer
	trailers     api.HeaderMap
	dataEncoding DataEncodingType
	attributes   map[string]interface{}
	requestId    string
}

func NewGatewayRequest(headers api.HeaderMap, dataBuf buffer.IoBuffer, trailers api.HeaderMap) *DefaultGatewayRequest {
	headerMap := make(map[string]string, 10)
	headers.Range(func(key, value string) bool {
		headerMap[key] = value
		return true
	})

	if dataBuf == nil {
		dataBuf = buffer.NewIoBuffer(0)
	}

	req := &DefaultGatewayRequest{
		headers:    headers,
		headerMap:  headerMap,
		dataBuf:    dataBuf,
		trailers:   trailers,
		attributes: make(map[string]interface{}, 2),
	}

	if contentType, ok := headers.Get(HeaderContentType); ok {
		req.SetDataEncoding(ParseDataEncoding(contentType))
	}

	return req
}

func (r *DefaultGatewayRequest) ApiId() string {
	return r.apiId
}

func (r *DefaultGatewayRequest) SetApiId(apiId string) {
	r.apiId = apiId
}

func (r *DefaultGatewayRequest) RequestId() string {
	return r.requestId
}

func (r *DefaultGatewayRequest) SetRequestId(requestId string) {
	r.requestId = requestId
}

func (r *DefaultGatewayRequest) Headers() api.HeaderMap {
	return r.headers
}

func (r *DefaultGatewayRequest) HeaderMap() map[string]string {
	return r.headerMap
}

func (r *DefaultGatewayRequest) GetHeader(headerKey string) string {
	if val, ok := r.headerMap[headerKey]; ok {
		return val
	}

	headerVal, _ := r.headers.Get(headerKey)
	return headerVal
}

func (r *DefaultGatewayRequest) SetHeader(headerKey, headerVal string) {
	r.headers.Set(headerKey, headerVal)
	r.headerMap[headerKey] = headerVal
}

func (r *DefaultGatewayRequest) AddHeader(headerKey, headerVal string) {
	r.headers.Add(headerKey, headerVal)
	r.headerMap[headerKey] = headerVal
}

func (r *DefaultGatewayRequest) Trailers() api.HeaderMap {
	return r.trailers
}

func (r *DefaultGatewayRequest) GetTrailer(trailerKey string) string {
	trailerVal, _ := r.trailers.Get(trailerKey)
	return trailerVal
}

func (r *DefaultGatewayRequest) SetTrailer(trailerKey, trailerVal string) {
	r.SetTrailer(trailerKey, trailerVal)
}

func (r *DefaultGatewayRequest) DataBuf() buffer.IoBuffer {
	return r.dataBuf
}

func (r *DefaultGatewayRequest) DataBytes() []byte {
	return r.dataBuf.Bytes()
}

func (r *DefaultGatewayRequest) SetDataBytes(data []byte) (n int, err error) {
	r.dataBuf.Reset()
	r.headers.Del("Content-Length")
	r.headers.Del("Content-Encoding")
	return r.dataBuf.Write(data)
}

func (r *DefaultGatewayRequest) DataEncoding() DataEncodingType {
	return r.dataEncoding
}

func (r *DefaultGatewayRequest) SetDataEncoding(encoding DataEncodingType) {
	r.dataEncoding = encoding
}

func (r *DefaultGatewayRequest) GetAttribute(key string) interface{} {
	return r.attributes[key]
}

func (r *DefaultGatewayRequest) SetAttribute(key string, value interface{}) {
	if r.attributes == nil {
		r.attributes = make(map[string]interface{})
	}
	r.attributes[key] = value
}
