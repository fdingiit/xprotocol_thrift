package sofahessian

import "gitlab.alipay-inc.com/sofa-go/sofa-hessian-go/javaobject"

// RegisterBuiltinJavaClasses registers java classes.
//
// nolint
func RegisterBuiltinJavaClasses() {
	RegisterJavaClass(javaobject.JavaLangStackTraceElements{})
	RegisterJavaClass(javaobject.JavaLangStackTraceElement{})
}

// RegisterSofaRPCJavaClasses registers the java classes.
//
// nolint
func RegisterSofaRPCJavaClasses() {
	RegisterJavaClass(&javaobject.SofaRPCRequest{})
	RegisterJavaClass(&javaobject.SofaRPCResponse{})
	RegisterJavaClass(&javaobject.SofaRPCServerException{})
}
