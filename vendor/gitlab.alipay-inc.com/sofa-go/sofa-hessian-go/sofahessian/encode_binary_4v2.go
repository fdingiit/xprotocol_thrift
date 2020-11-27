package sofahessian

import "encoding/binary"

// EncodeBinaryToHessian4V2 encodes binary to dst.
//
// binary0000 ::= x41 b1 b0 <binary-data> binary # non-final chunk
// 000000000000::= 'B' b1 b0 <binary-data>00000000# final chunk
// 000000000000::= [x20-x2f] <binary-data>00000000# binary data of
// 00000000000000000000000000000000000000000000   #  length 0-15
// 000000000000::= [x34-x37] <binary-data>00000000# binary data of
// 00000000000000000000000000000000000000000000   #  length 0-1023
//
// Binary data is encoded in chunks. The octet x42 ('B') encodes the final chunk and
// x62 ('b') represents any non-final chunk. Each chunk has a 16-bit length value.
//
// len = 256 * b1 + b0
func EncodeBinaryToHessian4V2(o *EncodeContext, dst []byte, b []byte) ([]byte, error) {
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

	for n > 4093 {
		dst = append(dst, "A00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], 4093)
		dst = append(dst, b[:4093]...)
		b = b[4093:]
		n = len(b)
	}

	if n < 16 {
		dst = append(dst, uint8(n+0x20))
	} else if n < 1024 {
		dst = append(dst, uint8(n>>8+0x34), uint8(n))
	} else {
		dst = append(dst, "B00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(n))
	}
	dst = append(dst, b...)

	return dst, nil
}
