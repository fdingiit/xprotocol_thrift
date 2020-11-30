package rule

type RpcSignature struct {
	InterfaceName string `json:"interfaceName"`
	MethodName    string `json:"methodName"`
	UniqueId      string `json:"uniqueId"`
	CallerAppName string `json:"callerAppName"`
	TargetAppName string `json:"targetAppName"`
}

const (
	RpcSigKeyInterface = "interfaceName"
	RpcSigKeyMethod    = "methodName"
	RpcSigKeyUniqueId  = "uniqueId"
	RpcSigKeyCaller    = "callerAppName"
	RpcSigKeyTarget    = "targetAppName"
	RpcKeyMark         = "mark"
	RpcKeyTraceId      = "traceId"
)
