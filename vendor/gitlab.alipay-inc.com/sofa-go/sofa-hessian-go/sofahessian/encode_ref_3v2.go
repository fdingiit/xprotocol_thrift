package sofahessian

import "encoding/binary"

// EncodeRefHessian3V2 encodes refid to dst.
func EncodeRefHessian3V2(o *EncodeContext, dst []byte, refid uint32) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encoderef")
		defer o.tracer.OnTraceStop("encoderef")
	}

	if refid < 0x100 {
		dst = append(dst, 'J')
		dst = append(dst, uint8(refid))

	} else if refid < 0x10000 {
		dst = append(dst, "K00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(refid))

	} else {
		dst = append(dst, "R0000"...)
		binary.BigEndian.PutUint32(dst[len(dst)-4:], refid)
	}

	return dst, nil
}

func encodeTyperefToHessian3V2(o *EncodeContext, dst []byte, typ string) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodetyperef")
		defer o.tracer.OnTraceStop("encodetyperef")
	}

	if typ == "" {
		return dst, nil
	}

	refid, ok, err := o.getTyperefs(typ)
	if err != nil {
		return dst, err
	}

	if !ok {
		if err = o.addTyperefs(typ); err != nil {
			return dst, err
		}
		dst = append(dst, "t00"...)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(len(typ)))
		dst = append(dst, typ...)
		return dst, nil
	}

	return EncodeInt32ToHessian3V2(o, dst, int32(refid))
}

func encodeObjectrefToHessian3V2(o *EncodeContext, dst []byte, obj interface{}) ([]byte, int, error) {
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
		dst, err = EncodeRefHessian3V2(o, dst, uint32(ref))
		return dst, ref, err
	}
	_, err = o.addObjectrefs(obj)

	return dst, -1, err
}
