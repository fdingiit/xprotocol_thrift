package sofahessian

// EncodeRef4V2 encodes refid to dst.
//
// ref ::= x51 int
func EncodeRef4V2(o *EncodeContext, dst []byte, refid uint32) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encoderef")
		defer o.tracer.OnTraceStop("encoderef")
	}

	dst = append(dst, "Q"...)
	return EncodeInt32ToHessian4V2(o, dst, int32(refid))
}

func encodeObjectrefToHessian4V2(o *EncodeContext, dst []byte, obj interface{}) ([]byte, int, error) {
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
		dst, err = EncodeRef4V2(o, dst, uint32(ref))
		return dst, ref, err
	}
	_, err = o.addObjectrefs(obj)

	return dst, -1, err
}

func encodeClassrefToHessian4V2(o *EncodeContext, dst []byte, obj interface{}) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodeoclassref")
		defer o.tracer.OnTraceStop("encodeoclassref")
	}

	ref, err := o.getClassrefs(obj)
	if err != nil {
		return dst, err
	}

	if ref >= 0 {
		return EncodeRef4V2(o, dst, uint32(ref))
	}
	_, _, err = o.addClassrefs(obj)
	return dst, err
}

// encodeTyperefToHessian4V2 encodes type to dst.
// type ::= string
//      ::= int
//  Figure 30
// A map or list includes a type attribute indicating the type name of the map or list for object-oriented languages.
//
// Each type is added to the type map for future reference.
func encodeTyperefToHessian4V2(o *EncodeContext, dst []byte, typ string) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodetyperef")
		defer o.tracer.OnTraceStop("encodetyperef")
	}

	if typ == "" {
		return dst, nil
	}

	v, ok, err := o.getTyperefs(typ)
	if err != nil {
		return dst, err
	}

	if !ok {
		if err = o.addTyperefs(typ); err != nil {
			return dst, err
		}
		return EncodeStringToHessian4V2(o, dst, typ)
	}

	return EncodeInt32ToHessian4V2(o, dst, int32(v))
}
