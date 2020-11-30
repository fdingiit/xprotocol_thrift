package sofahessian

import (
	"encoding/binary"
	"math"
)

// EncodeFloat64ToHessianV1 encodes float64 to dst with hessian3 v2 protocol.
func EncodeFloat64ToHessianV1(o *EncodeContext, dst []byte, d float64) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodefloat64")
		defer o.tracer.OnTraceStop("encodefloat64")
	}

	dst = append(dst, "D00000000"...)
	binary.BigEndian.PutUint64(dst[len(dst)-8:], math.Float64bits(d))

	return dst, nil
}
