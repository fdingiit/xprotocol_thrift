package javaobject

type JavaUtilConcurrentAtomicLong struct {
	Value int64 `hessian:"value"`
}

func (j JavaUtilConcurrentAtomicLong) GetJavaClassName() string {
	return "java.util.concurrent.atomic.AtomicLong"
}

type JavaStringArray []string

func (j JavaStringArray) GetJavaClassName() string {
	return "[string"
}

// JavaLangInteger represents java.lang.Integer.
type JavaLangInteger int32

func (j JavaLangInteger) GetJavaClassName() string { return "java.lang.Integer" }

// JavaLangStackTraceElements represents []java.lang.StackTrace.
type JavaLangStackTraceElements []JavaLangStackTraceElement

func (j JavaLangStackTraceElements) GetJavaClassName() string {
	return "[java.lang.StackTraceElement"
}

// JavaLangStackTraceElement represents java.lang.StackTrace.
type JavaLangStackTraceElement struct {
	DeclaringClass string `hessian:"declaringClass"`
	MethodName     string `hessian:"methodName"`
	FileName       string `hessian:"fileName"`
	LineNumber     int32  `hessian:"lineNumber"`
}

func (j JavaLangStackTraceElement) GetJavaClassName() string {
	return "java.lang.StackTraceElement"
}
