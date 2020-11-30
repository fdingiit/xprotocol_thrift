package sofahessian

import "encoding/binary"

func EncodeBinaryToHessian3V2(o *EncodeContext, dst []byte, b []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebinary")
		defer o.tracer.OnTraceStop("encodebinary")
	}

	n := len(b)
	if n < 16 {
		dst = append(dst, uint8(n)+0x20)
		dst = append(dst, b...)
		return dst, nil
	}

	for n > 0x8000 {
		dst = append(dst, "b00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], 0x8000)
		dst = append(dst, b[:0x8000]...)
		b = b[0x8000:]
		n = len(b)
	}

	if n < 16 {
		dst = append(dst, uint8(n+0x20))
	} else {
		dst = append(dst, "B00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(n))
	}
	dst = append(dst, b...)

	return dst, nil
}
