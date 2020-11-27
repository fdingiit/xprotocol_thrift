package gateway

import (
	"context"
	"strings"

	"mosn.io/api"
	"mosn.io/pkg/buffer"
)

type ApiStatus string

const (
	UNKNOWN ApiStatus = "unknown"
	OPENED  ApiStatus = "opened"
	CLOSED  ApiStatus = "closed"
)

type DataEncodingType string

const (
	URLENCODED DataEncodingType = "application/x-www-form-urlencoded"
	JSON       DataEncodingType = "application/json"
	PROTOBUF   DataEncodingType = "application/protobuf"
	STREAM     DataEncodingType = "application/octet-stream"
	EXT        DataEncodingType = "application/rpc-ext"
)

type HandleStatus string

const (
	HandleStatusContinue      HandleStatus = "Continue"
	HandleStatusStopAndReturn HandleStatus = "StopAndReturn"
)

func GetDataEncodingName(encodingType DataEncodingType) string {
	switch encodingType {
	case URLENCODED:
		return "URLENCODED"
	case JSON:
		return "JSON"
	case PROTOBUF:
		return "PROTOBUF"
	case STREAM:
		return "STREAM"
	case EXT:
		return "EXT"
	default:
		return "JSON"
	}
}

func GetDataEncodingValue(encodingType DataEncodingType) int32 {
	switch encodingType {
	case URLENCODED:
		return 0
	case JSON:
		return 1
	case PROTOBUF:
		return 2
	case STREAM:
		return 3
	case EXT:
		return 99
	default:
		return 1
	}
}

func ParseDataEncoding(contentType string) DataEncodingType {
	if len(contentType) < 1 {
		return ""
	}

	if index := strings.Index(contentType, ";"); index > 0 {
		contentType = contentType[0:index]
	}

	switch contentType {
	case string(URLENCODED):
		return URLENCODED
	case string(JSON):
		return JSON
	case string(PROTOBUF):
		return PROTOBUF
	case string(STREAM):
		return STREAM
	case string(EXT):
		return EXT
	default:
		return ""
	}
	return ""
}

type Service interface {
	ServiceKey() string

	SetServiceKey(serviceKey string)

	Status() ApiStatus

	SetStatus(status ApiStatus)

	Timeout() uint64

	SetTimeout(timeout uint64)

	GetAttribute(key string) interface{}

	SetAttribute(key string, value interface{})

	Upstream() Upstream

	SetUpstream(upstream Upstream)
}

type Upstream interface {
	Protocol() api.Protocol

	SetProtocol(protocol api.Protocol)

	Timeout() uint64

	SetTimeout(timeout uint64)

	ClusterName() string

	SetClusterName(cluster string)
}

type GatewayRequest interface {
	ApiId() string

	SetApiId(apiId string)

	RequestId() string

	SetRequestId(uniqueId string)

	Headers() api.HeaderMap

	HeaderMap() map[string]string

	GetHeader(headerKey string) string

	SetHeader(headerKey, headerVal string)

	AddHeader(headerKey, headerVal string)

	DataBuf() buffer.IoBuffer

	DataBytes() []byte

	SetDataBytes(data []byte) (n int, err error)

	Trailers() api.HeaderMap

	GetTrailer(trailerKey string) string

	SetTrailer(trailerKey, trailerVal string)

	DataEncoding() DataEncodingType

	SetDataEncoding(encoding DataEncodingType)

	GetAttribute(key string) interface{}

	SetAttribute(key string, value interface{})
}

type GatewayResponse interface {
	Headers() api.HeaderMap

	GetHeader(headerKey string) string

	SetHeader(headerKey, headerVal string)

	AddHeader(headerKey, headerVal string)

	DataBuf() buffer.IoBuffer

	DataBytes() []byte

	SetDataBytes(data []byte) (n int, err error)

	Trailers() api.HeaderMap

	GetTrailer(trailerKey string) string

	SetTrailer(trailerKey, trailerVal string)

	DataEncoding() DataEncodingType

	SetDataEncoding(encoding DataEncodingType)

	ResultStatus() ResponseStatus

	SetResultStatus(resultStatus ResponseStatus)
}

type Pipeline interface {
	RunInHandlers(ctx context.Context) (bool, error)

	RunOutHandlers(ctx context.Context) (bool, error)

	//AddFirst(handler Handler) Pipeline
	//
	//AddLast(handler Handler) Pipeline
	//
	//AddBefore(baseHandlerName string, handler Handler) Pipeline
	//
	//AddAfter(baseHandlerName string, handler Handler) Pipeline
	//
	//AddListFirst(handlers ...Handler) Pipeline
	//
	//AddListLast(handlers ...Handler) Pipeline
	//
	//Copy() Pipeline
}

type Handler interface {
	Name() string

	HandleIn(ctx context.Context) (HandleStatus, error)

	HandleOut(ctx context.Context) (HandleStatus, error)
}

type ConfigListener interface {
	Update(key string, value interface{})
}

type GatewayManager interface {
	AddOrUpdateGateway(gateway GatewayConfig) bool

	AddOrUpdateService(service ServiceConfig) bool

	AddOrUpdateApp(app AppConfig) bool

	AddOrUpdateCluster(cluster ClusterConfig) bool

	AddOrUpdateGatewayRule(rule GatewayRule) bool

	GetService(id string) Service

	GetApp(appName string) *AppConfig

	GetCluster(clusterName string) *ClusterConfig

	GetGatewayRule(ruleName string) interface{}

	RemoveService(service ServiceConfig) bool

	RemoveApp(app AppConfig) bool

	RemoveCluster(cluster ClusterConfig) bool

	RemoveGatewayRule(rule GatewayRule) bool

	SetFilterConfig(conf *GatewayFilterConfig)

	GetFilterConfig() *GatewayFilterConfig

	Dump()
}

type ConfigFilter interface {
	OnAddOrUpdateService(service Service)

	OnAddOrUpdateApp(app AppConfig)

	OnAddOrUpdateCluster(cluster ClusterConfig)

	OnAddOrUpdateGatewayRule(rule GatewayRule)
}

type DownstreamCodec interface {
	Decode(ctx context.Context, headers api.HeaderMap, dataBuf buffer.IoBuffer, trailers api.HeaderMap) (GatewayRequest, error)

	Encode(ctx context.Context, mapping DownstreamStatusMapping) (api.HeaderMap, buffer.IoBuffer, api.HeaderMap)
}

type UpstreamCodec interface {
	Encode(ctx context.Context) (api.HeaderMap, buffer.IoBuffer, api.HeaderMap)

	Decode(ctx context.Context, headers api.HeaderMap, dataBuf buffer.IoBuffer, trailers api.HeaderMap, mapping UpstreamStatusMapping) (GatewayResponse, error)
}

// ResponseMapping maps the gateway status to downstream status
type DownstreamStatusMapping func(ctx context.Context, status ResponseStatus) map[string]string

// ResponseMapping maps the upstream status to gateway status
type UpstreamStatusMapping func(ctx context.Context, headers api.HeaderMap) ResponseStatus

type UpstreamParser func(cfg UpstreamConfig) Upstream

type ServiceRuleParser func(cfg interface{}) (interface{}, error)
