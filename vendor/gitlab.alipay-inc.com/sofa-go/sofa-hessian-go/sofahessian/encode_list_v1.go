package sofahessian

import (
	"encoding/binary"
	"reflect"
)

func EncodeListToHessianV1(o *EncodeContext, dst []byte, obj interface{}) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodelist")
		defer o.tracer.OnTraceStop("encodelist")
	}

	if obj == nil {
		return EncodeNilToHessianV1(o, dst)
	}

	value := reflect.ValueOf(obj)
	classname := getInterfaceName(obj)
	return encodeListToHessianV1(o, dst, value, classname)
}

func EncodeListBeginToHessianV1(o *EncodeContext, dst []byte, length int, typ string) ([]byte, bool, error) {
	refid, _, err := o.getTyperefs(typ)
	if err != nil {
		return dst, false, err
	}

	if refid >= 0 {
		dst = append(dst, 'R')
		dst, err = EncodeInt32RefToHessianV1(o, dst, int32(refid))
		if err != nil {
			return dst, false, err
		}
		return dst, false, err
	}

	dst = append(dst, 'V')
	dst, err = encodeTyperefToHessian3V2(o, dst, typ)
	if err != nil {
		return dst, false, err
	}

	dst = append(dst, "l0000"...)
	binary.BigEndian.PutUint32(dst[len(dst)-4:], uint32(length))

	return dst, true, nil
}

func EncodeListEndToHessianV1(o *EncodeContext, dst []byte, end bool) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodelistend")
		defer o.tracer.OnTraceStop("encodelistend")
	}

	if end {
		dst = append(dst, 'z')
	}
	return dst, nil
}

func encodeListToHessianV1(o *EncodeContext, dst []byte, slice reflect.Value, typ string) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodelistbegin")
		defer o.tracer.OnTraceStop("encodelistbegin")
	}

	// Unwrap the pointer if we can
	slice = reflect.Indirect(slice)

	if slice.Kind() != reflect.Slice &&
		slice.Kind() != reflect.Array {
		return dst, ErrEncodeNotSliceType
	}

	var (
		err   error
		refid int
	)

	if !o.disableObjectrefs {
		if slice.Kind() == reflect.Slice {
			// []interface{} cannot be hashed so we use address instead.
			dst, refid, err = encodeObjectrefToHessianV1(o, dst, slice.Pointer())
		} else { // Array
			dst, refid, err = encodeObjectrefToHessianV1(o, dst, slice.Interface())
		}

		if err != nil {
			return dst, err
		}

		if refid >= 0 {
			return dst, nil
		}
	}

	var (
		end    bool
		length = slice.Len()
	)

	dst, end, err = EncodeListBeginToHessianV1(o, dst, length, typ)
	if err != nil {
		return dst, err
	}

	for i := 0; i < length; i++ {
		if slice.Index(i).CanInterface() {
			if dst, err = EncodeToHessianV1(o, dst, slice.Index(i).Interface()); err != nil {
				return dst, err
			}
		} else {
			return dst, ErrEncodeSliceElemCannotBeInterfaced
		}
	}

	return EncodeListEndToHessianV1(o, dst, end)
}
