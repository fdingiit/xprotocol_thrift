package gateway

const (
	HeaderResultCode       = "x-mosn-resultcode"
	HeaderResultStatus     = "x-mosn-result-status"
	HeaderResultTip        = "x-mosn-result-tip"
	HeaderMessage          = "x-mosn-message"
	HeaderAppId            = "x-app-id"
	HeaderDigest           = "digest"
	HeaderTraceId          = "x-mosn-traceid"
	HeaderUpstreamHttpCode = "x-upstream-http-code"

	HeaderContentType = "Content-Type"
	HeaderTimestamp   = "Ts"

	GATEWAY_CONTEXT_NAME = "MOSN_GATEWAY_CONTEXT"

	GatewayAttrUpRespCode = "upstream_resp_code"
)
