package sofahessian

import "encoding/binary"

func encodeObjectrefToHessianV1(o *EncodeContext, dst []byte, obj interface{}) ([]byte, int, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodeoobjectref")
		defer o.tracer.OnTraceStop("encodeoobjectref")
	}

	if o.disableObjectrefs {
		return dst, -1, nil
	}

	ref, err := o.getObjectrefs(obj)
	if err != nil {
		return dst, -1, err
	}

	if ref >= 0 {
		dst, err = EncodeRefHessianV1(o, dst, uint32(ref))
		return dst, ref, err
	}
	_, err = o.addObjectrefs(obj)

	return dst, -1, err
}

// EncodeRefHessianV1 encodes refid to dst.
func EncodeRefHessianV1(o *EncodeContext, dst []byte, refid uint32) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encoderef")
		defer o.tracer.OnTraceStop("encoderef")
	}
	dst = append(dst, "R0000"...)
	binary.BigEndian.PutUint32(dst[len(dst)-4:], refid)

	return dst, nil
}
