package sofahessian

import (
	"fmt"
	"io"
)

func DecodeBoolHessian4V2(o *DecodeContext, reader io.Reader) (bool, error) {
	var b bool
	err := DecodeBoolToHessian4V2(o, reader, &b)
	return b, err
}

func DecodeBoolToHessian4V2(o *DecodeContext, reader io.Reader, b *bool) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodebool")
		defer o.tracer.OnTraceStop("decodebool")
	}

	var c [1]byte
	n, err := reader.Read(c[:])
	if err != nil {
		return err
	}
	if n < 1 {
		return fmt.Errorf("expect read one byte but got zero")
	}

	switch c[0] {
	case 'T':
		*b = true
	case 'F':
		*b = false
	default:
		return ErrDecodeMalformedBool
	}

	return nil
}
