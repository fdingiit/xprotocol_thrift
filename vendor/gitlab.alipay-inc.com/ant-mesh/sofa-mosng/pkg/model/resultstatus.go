package model

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkg/gateway"
	"mosn.io/mosn/pkg/types"
	"strconv"
	"mosn.io/mosn/pkg/protocol"
	"net/http"
	"context"
)

var respStatusMap = make(map[gateway.ResponseStatus]*ResultStatus)
var resultCodeMap = make(map[int]*ResultStatus)

var (
	ResultSuccess           = ResultStatus{ResultCode: 1000, HttpCode: http.StatusOK, Memo: "ok", Tips: "ok", RespStatus: gateway.ResultSuccess}
	PermissionDeny          = ResultStatus{ResultCode: 1001, HttpCode: http.StatusForbidden, Memo: "Permission Deny", Tips: "Permission Deny", RespStatus: gateway.PermissionDeny}
	InvokeExceedLimit       = ResultStatus{ResultCode: 1002, HttpCode: http.StatusTooManyRequests, Memo: "Invoke Exceed Limit", Tips: "顾客太多，客官请稍候", RespStatus: gateway.InvokeExceedLimit}
	HumanCheckDeny          = ResultStatus{ResultCode: 1004, HttpCode: http.StatusForbidden, Memo: "Human Check Deny", Tips: "抱歉，请求不合法，请稍后再试", RespStatus: gateway.HumanCheckDeny}
	AclCheckFail            = ResultStatus{ResultCode: 1005, HttpCode: http.StatusUnauthorized, Memo: "Acl Check Failed", Tips: "抱歉，请求不合法，请稍后再试", RespStatus: gateway.AclCheckFail}
	DeviceIdCheckFail       = ResultStatus{ResultCode: 1006, HttpCode: http.StatusUnauthorized, Memo: "DeviceId Check Failed", Tips: "抱歉，请求不合法，请稍后再试", RespStatus: gateway.DeviceIdCheckFail}
	SessionStatus           = ResultStatus{ResultCode: 2000, HttpCode: http.StatusUnauthorized, Memo: "Illegal Session Status", Tips: "登录超时，请重新登录", RespStatus: gateway.SessionStatus}
	ServiceMissed           = ResultStatus{ResultCode: 3000, HttpCode: http.StatusNotFound, Memo: "Service Missed", Tips: "抱歉，暂时无法操作，请稍后再试", RespStatus: gateway.ServiceMissed}
	RequestDataMissed       = ResultStatus{ResultCode: 3001, HttpCode: http.StatusBadRequest, Memo: "Request Data Missed", Tips: "系统繁忙，请稍后再试", RespStatus: gateway.RequestDataMissed}
	ValueInvalid            = ResultStatus{ResultCode: 3002, HttpCode: http.StatusBadRequest, Memo: "Request Data Invalid", Tips: "抱歉，暂时无法操作，请稍后再试", RespStatus: gateway.ValueInvalid}
	EncryptionError         = ResultStatus{ResultCode: 3003, HttpCode: http.StatusBadRequest, Memo: "Encryption Error", Tips: "抱歉，暂时无法操作，请稍后再试", RespStatus: gateway.EncryptionError}
	RequestTimeOut          = ResultStatus{ResultCode: 4001, HttpCode: http.StatusGatewayTimeout, Memo: "Request Timeout", Tips: "请求超时，请稍后再试", RespStatus: gateway.RequestTimeOut}
	RemoteAccessException   = ResultStatus{ResultCode: 4002, HttpCode: http.StatusBadGateway, Memo: "Remote Access Exception", Tips: "系统繁忙，请稍后再试", RespStatus: gateway.RemoteAccessException}
	CreateProxyError        = ResultStatus{ResultCode: 4003, HttpCode: http.StatusServiceUnavailable, Memo: "Create Proxy Error", Tips: "系统繁忙，请稍后再试", RespStatus: gateway.CreateProxyError}
	UnknownError            = ResultStatus{ResultCode: 5000, HttpCode: http.StatusInternalServerError, Memo: "Unknown Error", Tips: "抱歉，暂时无法操作，请稍后再试", RespStatus: gateway.UnknownError}
	ServiceNotFound         = ResultStatus{ResultCode: 6000, HttpCode: http.StatusNotImplemented, Memo: "RPC-ServiceNotFound", Tips: "抱歉，暂时无法操作，请稍后再试", RespStatus: gateway.ServiceNotFound}
	MethodNotFound          = ResultStatus{ResultCode: 6001, HttpCode: http.StatusNotImplemented, Memo: "RPC-MethodNotFound", Tips: "抱歉，暂时无法操作，请稍后再试", RespStatus: gateway.MethodNotFound}
	IllegalAccess           = ResultStatus{ResultCode: 6003, HttpCode: http.StatusForbidden, Memo: "RPC-IllegalAccess", Tips: "抱歉，暂时无法操作，请稍后再试", RespStatus: gateway.IllegalAccess}
	DataParserException     = ResultStatus{ResultCode: 6004, HttpCode: http.StatusUnsupportedMediaType, Memo: "RPC-DataParserException", Tips: "抱歉，暂时无法操作，请稍后再试", RespStatus: gateway.DataParserException}
	IllegalArgument         = ResultStatus{ResultCode: 6005, HttpCode: http.StatusBadRequest, Memo: "RPC-IllegalArgument", Tips: "抱歉，暂时无法操作，请稍后再试", RespStatus: gateway.IllegalArgument}
	BizException            = ResultStatus{ResultCode: 6666, HttpCode: http.StatusBadGateway, Memo: "RPC-BizException", Tips: "抱歉，暂时无法操作，请稍后再试", RespStatus: gateway.BizException}
	SignKeyNotFound         = ResultStatus{ResultCode: 7000, HttpCode: http.StatusForbidden, Memo: "Sign Key Not Found", Tips: "验签失败", RespStatus: gateway.SignKeyNotFound}
	SignParamMissing        = ResultStatus{ResultCode: 7001, HttpCode: http.StatusForbidden, Memo: "Sign Param Missed", Tips: "抱歉，请求参数不合法", RespStatus: gateway.SignParamMissing}
	SignVerifyFailed        = ResultStatus{ResultCode: 7002, HttpCode: http.StatusForbidden, Memo: "Sign Verify Failed", Tips: "验签失败", RespStatus: gateway.SignVerifyFailed}
	SignTimeStampError      = ResultStatus{ResultCode: 7003, HttpCode: http.StatusForbidden, Memo: "Sign Timestamp Error", Tips: "手机时间异常，请到系统时间设置，将其设为最新", RespStatus: gateway.SignTimeStampError}
	ResponseDataNotModified = ResultStatus{ResultCode: 8001, HttpCode: http.StatusNotModified, Memo: "etag", Tips: "etag", RespStatus: gateway.ResponseDataNotModified}
	CORSOptions             = ResultStatus{ResultCode: 8002, HttpCode: http.StatusOK, Memo: "CORS preflight", Tips: "跨域预检请求", RespStatus: gateway.CORSOptions}
)

type ResultStatus struct {
	ResultCode int
	HttpCode   int
	Memo       string
	Tips       string
	RespStatus gateway.ResponseStatus
}

func init() {
	gateway.RegisterDownstreamStatusMapping(string(protocol.HTTP1), MappingResponseStatus, 1)
	ResultSuccess.AddResultStatus()
	PermissionDeny.AddResultStatus()
	InvokeExceedLimit.AddResultStatus()
	HumanCheckDeny.AddResultStatus()
	AclCheckFail.AddResultStatus()
	DeviceIdCheckFail.AddResultStatus()
	SessionStatus.AddResultStatus()
	ServiceMissed.AddResultStatus()
	RequestDataMissed.AddResultStatus()
	ValueInvalid.AddResultStatus()
	EncryptionError.AddResultStatus()
	RequestTimeOut.AddResultStatus()
	RemoteAccessException.AddResultStatus()
	CreateProxyError.AddResultStatus()
	UnknownError.AddResultStatus()
	ServiceNotFound.AddResultStatus()
	MethodNotFound.AddResultStatus()
	IllegalAccess.AddResultStatus()
	DataParserException.AddResultStatus()
	IllegalAccess.AddResultStatus()
	DataParserException.AddResultStatus()
	IllegalArgument.AddResultStatus()
	BizException.AddResultStatus()
	SignKeyNotFound.AddResultStatus()
	SignParamMissing.AddResultStatus()
	SignVerifyFailed.AddResultStatus()
	SignTimeStampError.AddResultStatus()
	ResponseDataNotModified.AddResultStatus()
	CORSOptions.AddResultStatus()
}

func (r *ResultStatus) AddResultStatus() {
	resultCodeMap[r.ResultCode] = r
	respStatusMap[r.RespStatus] = r
}

func GetByRespStatus(s gateway.ResponseStatus) *ResultStatus {
	return respStatusMap[s]
}

func GetByResultCode(resultCode int) *ResultStatus {
	return resultCodeMap[resultCode]
}

func MappingResponseStatus(ctx context.Context, s gateway.ResponseStatus) map[string]string {
	mapped := GetByRespStatus(s)
	resultMap := make(map[string]string)
	if mapped == nil {
		return resultMap
	}

	resultMap[types.HeaderStatus] = strconv.Itoa(mapped.HttpCode)
	resultMap[gateway.HeaderResultCode] = strconv.Itoa(mapped.ResultCode)
	resultMap[gateway.HeaderMessage] = mapped.Memo
	resultMap[gateway.HeaderResultTip] = mapped.Tips
	return resultMap
}
