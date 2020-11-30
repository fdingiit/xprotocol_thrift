package gateway

type ResponseStatus string

const (
	// success
	ResultSuccess ResponseStatus = "ResultSuccess"

	PermissionDeny    ResponseStatus = "PermissionDeny"
	InvokeExceedLimit ResponseStatus = "InvokeExceedLimit"
	HumanCheckDeny    ResponseStatus = "HumanCheckDeny"
	AclCheckFail      ResponseStatus = "AclCheckFail"
	DeviceIdCheckFail ResponseStatus ="DeviceIdCheckFail"

	SessionStatus ResponseStatus = "SessionStatus"

	ServiceMissed     ResponseStatus = "ServiceMissed"
	RequestDataMissed ResponseStatus = "RequestDataMissed"
	ValueInvalid      ResponseStatus = "ValueInvalid"
	EncryptionError   ResponseStatus = "EncryptionError"

	RequestTimeOut        ResponseStatus = "RequestTimeOut"
	RemoteAccessException ResponseStatus = "RemoteAccessException"
	CreateProxyError      ResponseStatus = "CreateProxyError"

	UnknownError ResponseStatus = "UnknownError"

	ServiceNotFound     ResponseStatus = "ServiceNotFound"
	MethodNotFound      ResponseStatus = "MethodNotFound"
	IllegalAccess       ResponseStatus = "IllegalAccess"
	DataParserException ResponseStatus = "DataParserException"
	IllegalArgument     ResponseStatus = "IllegalArgument"
	BizException        ResponseStatus = "BizException"

	// sing check result
	SignKeyNotFound    ResponseStatus = "SignKeyNotFound"
	SignParamMissing   ResponseStatus = "SignParamMissing"
	SignVerifyFailed   ResponseStatus = "SignVerifyFailed"
	SignTimeStampError ResponseStatus = "SignTimeStampError"

	ResponseDataNotModified ResponseStatus = "ResponseDataNotModified"
	CORSOptions             ResponseStatus = "CORSOptions"
)
