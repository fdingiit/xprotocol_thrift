package sofahessian

// EncodeBoolToHessian4V2 encodes bool to dst.
// The octet 'F' represents false and the octet T represents true.
// boolean ::= T
//         ::= F
func EncodeBoolToHessian4V2(o *EncodeContext, dst []byte, b bool) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebool")
		defer o.tracer.OnTraceStop("encodebool")
	}

	if b {
		dst = append(dst, 'T')
	} else {
		dst = append(dst, 'F')
	}
	return dst, nil
}

func EncodeBoolToHessian3V2(o *EncodeContext, dst []byte, b bool) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebool")
		defer o.tracer.OnTraceStop("encodebool")
	}

	if b {
		dst = append(dst, 'T')
	} else {
		dst = append(dst, 'F')
	}
	return dst, nil
}

func EncodeBoolToHessianV1(o *EncodeContext, dst []byte, b bool) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebool")
		defer o.tracer.OnTraceStop("encodebool")
	}

	if b {
		dst = append(dst, 'T')
	} else {
		dst = append(dst, 'F')
	}
	return dst, nil
}
