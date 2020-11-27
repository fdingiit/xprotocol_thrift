package javaobject

type SofaRPCRequest struct {
	TargetAppName           string                 `hessian:"targetAppName"`
	TargetServiceUniqueName string                 `hessian:"targetServiceUniqueName"`
	MethodName              string                 `hessian:"methodName"`
	MethodArgSigs           []string               `hessian:"methodArgSigs"`
	RequestProps            map[string]interface{} `hessian:"requestProps"`
}

func (s *SofaRPCRequest) GetJavaClassName() string {
	return "com.alipay.sofa.rpc.core.request.SofaRequest"
}

type SofaRPCResponse struct {
	IsError       bool              `hessian:"isError"`
	ErrorMsg      string            `hessian:"errorMsg"`
	AppResponse   interface{}       `hessian:"appResponse"`
	ResponseProps map[string]string `hessian:"responseProps"`
}

func (s *SofaRPCResponse) GetJavaClassName() string {
	return "com.alipay.sofa.rpc.core.response.SofaResponse"
}

type SofaRPCServerException struct {
	DetailMessage        string                     `hessian:"detailMessage"`
	Cause                interface{}                `hessian:"cause"`
	StackTrace           JavaLangStackTraceElements `hessian:"stackTrace"`
	SuppressedExceptions interface{}                `hessian:"suppressedExceptions"`
}

func (s *SofaRPCServerException) GetJavaClassName() string {
	return "com.alipay.remoting.rpc.exception.RpcServerException"
}
