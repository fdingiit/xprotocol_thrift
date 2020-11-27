package sofahessian

import (
	"encoding/binary"
	"time"
)

// EncodeDateToHessian4V2 encodes data to dst.
// date ::= x4a b7 b6 b5 b4 b3 b2 b1 b0
//      ::= x4b b4 b3 b2 b1 b0
// Date represented by a 64-bit long of milliseconds since Jan 1 1970 00:00H, UTC.
func EncodeDateToHessian4V2(o *EncodeContext, dst []byte, t time.Time) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodedate")
		defer o.tracer.OnTraceStop("encodedate")
	}

	ts := t.UnixNano() / 1000 / 1000

	if ts%60000 == 0 {
		minutes := ts / 60000
		if (minutes>>31) == 0 || (minutes>>31) == -1 {
			dst = append(dst, "K0000"...)
			binary.BigEndian.PutUint32(dst[len(dst)-4:], uint32(minutes))
			return dst, nil
		}
	}

	dst = append(dst, "J00000000"...)
	binary.BigEndian.PutUint64(dst[len(dst)-8:], uint64(ts))
	return dst, nil
}

func EncodeDateToHessian3V2(o *EncodeContext, dst []byte, t time.Time) ([]byte, error) {
	ts := t.UnixNano() / 1000 / 1000
	dst = append(dst, "d00000000"...)
	binary.BigEndian.PutUint64(dst[len(dst)-8:], uint64(ts))
	return dst, nil
}

func EncodeDateToHessianV1(o *EncodeContext, dst []byte, t time.Time) ([]byte, error) {
	ts := t.UnixNano() / 1000 / 1000
	dst = append(dst, "d00000000"...)
	binary.BigEndian.PutUint64(dst[len(dst)-8:], uint64(ts))
	return dst, nil
}
