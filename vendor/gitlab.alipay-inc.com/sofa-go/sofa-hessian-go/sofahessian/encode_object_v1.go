package sofahessian

import (
	"encoding/binary"
	"reflect"
)

func EncodeObjectToHessianV1(o *EncodeContext, dst []byte, obj interface{}) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodeobject")
		defer o.tracer.OnTraceStop("encodeobject")
	}

	if obj == nil {
		return EncodeNilToHessianV1(o, dst)
	}

	var (
		t     = reflect.TypeOf(obj)
		err   error
		refid int
	)

	dst, refid, err = encodeObjectrefToHessianV1(o, dst, obj)
	if err != nil {
		return dst, err
	}
	if refid >= 0 {
		return dst, nil
	}

	v := decAllocReflectValue(reflect.ValueOf(obj))
	t = decAllocReflectType(t)
	dst, err = encodeObjectDefinitionHessianV1(o, dst, obj, t, v)
	if err != nil {
		return dst, err
	}

	// Write object field
	for i := 0; i < v.NumField(); i++ {
		vfield := v.Field(i)
		if !vfield.CanInterface() {
			continue
		}
		tfield := t.Field(i)
		name := tfield.Name
		if hname := tfield.Tag.Get("hessian"); hname != "" {
			if hname == "-" {
				continue
			}
			name = hname
		}

		dst, err = EncodeStringToHessianV1(o, dst, name)
		if err != nil {
			return dst, err
		}
		dst, err = EncodeToHessianV1(o, dst, vfield.Interface())
		if err != nil {
			return dst, err
		}
	}
	dst, err = EncodeMapEndToHessianV1(o, dst)
	if err != nil {
		return dst, err
	}

	return dst, nil
}

func EncodeClassDefinitionToHessianV1(o *EncodeContext, dst []byte, classname string, fields []string) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodeobjectdefinition")
		defer o.tracer.OnTraceStop("encodeobjectdefinition")
	}

	dst, ref, err := encodeObjectBeginToHessianV1(o, dst, classname)
	if err != nil {
		return dst, err
	}

	if ref == -1 {
		if dst, err = EncodeInt32ToHessianV1(o, dst, int32(len(fields))); err != nil {
			return dst, err
		}

		for i := range fields {
			if dst, err = EncodeStringToHessianV1(o, dst, fields[i]); err != nil {
				return dst, err
			}
		}

		// Set the ref
		dst, _, err = encodeObjectBeginToHessianV1(o, dst, classname)
		if err != nil {
			return dst, err
		}
	}

	return dst, nil
}

func encodeObjectDefinitionHessianV1(o *EncodeContext, dst []byte,
	obj interface{}, t reflect.Type, v reflect.Value) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodebjectdefinition")
		defer o.tracer.OnTraceStop("encodeobjectdefinition")
	}

	var (
		classname = getInterfaceName(obj)
	)

	dst = append(dst, 0x4d, 't', '0', '0')
	binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(len(classname)))
	dst = append(dst, classname...)

	return dst, nil
}

func encodeObjectBeginToHessianV1(o *EncodeContext, dst []byte, typ string) ([]byte, int, error) {
	refid, referenced, err := o.addClassrefs(typ)
	if err != nil {
		return dst, -1, err
	}

	if referenced {
		if refid >= 0 {
			dst = append(dst, 0x6f)
			dst, err = EncodeInt32ToHessianV1(o, dst, int32(refid))
			return dst, refid, err
		}
	}

	// class definition
	dst = append(dst, 0x4F)
	dst, err = EncodeInt32ToHessianV1(o, dst, int32(len(typ)))
	if err != nil {
		return dst, -1, err
	}
	dst = append(dst, typ...)

	return dst, -1, err
}
