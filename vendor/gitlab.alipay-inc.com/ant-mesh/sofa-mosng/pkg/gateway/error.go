package gateway

import (
	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol"
	"mosn.io/pkg/buffer"
)

type GatewayError struct {
	ResultStatus ResponseStatus
	Headers      api.HeaderMap
	DataBuf      buffer.IoBuffer
	Trailers     api.HeaderMap
	DataEncoding DataEncodingType
}

func NewGatewayError(resultstatus ResponseStatus) *GatewayError {
	return &GatewayError{
		ResultStatus: resultstatus,
		Headers:      protocol.CommonHeader(make(map[string]string)),
		DataBuf:      buffer.NewIoBuffer(0),
		Trailers:     protocol.CommonHeader(make(map[string]string)),
		DataEncoding: JSON,
	}
}

func NewGatewayErrorByResult(status ResponseStatus, headers api.HeaderMap, dataBuf buffer.IoBuffer, trailers api.HeaderMap) *GatewayError {
	return &GatewayError{
		ResultStatus: status,
		Headers:      headers,
		DataBuf:      dataBuf,
		Trailers:     trailers,
		DataEncoding: JSON,
	}
}

func NewGatewayErrorByResultAndEncoding(status ResponseStatus, headers api.HeaderMap, dataBuf buffer.IoBuffer, trailers api.HeaderMap, dataEncoding DataEncodingType) *GatewayError {
	return &GatewayError{
		ResultStatus: status,
		Headers:      headers,
		DataBuf:      dataBuf,
		Trailers:     trailers,
		DataEncoding: dataEncoding,
	}
}

func (e *GatewayError) Error() string {
	return string(e.ResultStatus)
}
