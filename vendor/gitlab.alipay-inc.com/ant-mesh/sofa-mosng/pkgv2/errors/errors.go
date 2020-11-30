package errors

import (
	"fmt"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
	"io"
	"net/http"
)


// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(httpCode int) types.GatewayError {
	return &gatewayError{
		code:  httpCode,
		msg:   http.StatusText(httpCode),
	}
}

func NewWithMsg(httpCode int, msg string) types.GatewayError {
	return &gatewayError{
		code: httpCode,
		msg:  msg,
	}
}

func NewWithStack(httpCode int) types.GatewayError {
	return &gatewayError{
		code:  httpCode,
		msg:   http.StatusText(httpCode),
		stack: callers(),
	}
}

func Error(msg string) types.GatewayError {
	return &gatewayError{
		msg:   msg,
		stack: callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) types.GatewayError {
	return &gatewayError{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// gatewayError is an error that has a message and a stack, but no caller.
type gatewayError struct {
	code int
	msg  string
	*stack
}

func (f *gatewayError) Error() string { return f.msg }

func (f *gatewayError) Code() int { return f.code }

func (f *gatewayError) Msg() string { return f.msg }

func (f *gatewayError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}
