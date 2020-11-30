package sofahessian

import "encoding/binary"

func EncodeBinaryToHessianV1(o *EncodeContext, dst []byte, b []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebinary")
		defer o.tracer.OnTraceStop("encodebinary")
	}

	n := len(b)

	for n > 0x8000 {
		dst = append(dst, "b00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], 0x8000)
		dst = append(dst, b[:0x8000]...)
		b = b[0x8000:]
		n = len(b)
	}

	dst = append(dst, "B00"...)
	binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(n))
	dst = append(dst, b...)

	return dst, nil
}
