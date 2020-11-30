package sofahessian

// EncodeNilToHessian4V2 encodes nil to dst.
//
// null ::= N
func EncodeNilToHessian4V2(o *EncodeContext, dst []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodenil")
		defer o.tracer.OnTraceStop("encodenil")
	}
	dst = append(dst, 'N')
	return dst, nil
}

// EncodeNilToHessian3V2 encoddes nil to dst with hessian3 v2 protocol.
func EncodeNilToHessian3V2(e *EncodeContext, dst []byte) ([]byte, error) {
	return EncodeNilToHessian4V2(e, dst)
}

// EncodeNilToHessianV1 encoddes nil to dst with hessian3 v2 protocol.
func EncodeNilToHessianV1(e *EncodeContext, dst []byte) ([]byte, error) {
	return EncodeNilToHessian4V2(e, dst)
}
