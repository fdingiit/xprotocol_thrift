package sofahessian

import (
	"encoding/binary"
	"reflect"
	"sort"
)

func EncodeMapBeginToHessianV1(o *EncodeContext, dst []byte, typ string) ([]byte, error) {
	dst = append(dst, 'M')
	if typ != "" {
		dst = append(dst, 0x74, 0, 0)
		binary.BigEndian.PutUint16(dst[len(dst)-2:], uint16(len(typ)))
		dst = append(dst, typ...)
	}

	return dst, nil
}

func EncodeMapEndToHessianV1(o *EncodeContext, dst []byte) ([]byte, error) {
	dst = append(dst, 0x7a)
	return dst, nil
}

// EncodeMapToHessianV1 encodes map to dst.
func EncodeMapToHessianV1(o *EncodeContext, dst []byte, obj interface{}) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodemap")
		defer o.tracer.OnTraceStop("encodemap")
	}

	if obj == nil {
		return EncodeNilToHessianV1(o, dst)
	}

	// Allow *map to reduce recursive encodeto call
	t := reflect.TypeOf(obj)
	if kind := t.Kind(); kind != reflect.Map {
		if kind == reflect.Ptr {
			if t.Elem().Kind() != reflect.Map {
				return dst, ErrEncodeNotMapType
			}
		} else {
			return dst, ErrEncodeNotMapType
		}
	}

	v := reflect.ValueOf(obj)

	var (
		err   error
		refid int
	)

	if !o.disableObjectrefs {
		// Map cannot be hashed, use pointer instead.
		dst, refid, err = encodeObjectrefToHessianV1(o, dst, v.Pointer())
		if err != nil {
			return dst, err
		}

		if refid >= 0 {
			return dst, nil
		}
	}

	classname := getInterfaceName(obj)
	dst, err = EncodeMapBeginToHessianV1(o, dst, classname)
	if err != nil {
		return dst, err
	}

	// Unwrap the pointer if can
	v = reflect.Indirect(v)

	// Map in golang is unordered but other languages maybe or not maybe unordered.
	keys := v.MapKeys()
	if o.less == nil { // Fast path
		for i := range keys {
			key := keys[i]
			if key.CanInterface() { // Fast path
				if dst, err = EncodeToHessianV1(o, dst, key.Interface()); err != nil {
					return dst, err
				}
			} else {
				if dst, err = EncodeValueToHessianV1(o, dst, key); err != nil {
					return dst, err
				}
			}

			value := v.MapIndex(key)
			if value.CanInterface() { // Fast path
				if dst, err = EncodeToHessianV1(o, dst, value.Interface()); err != nil {
					return dst, err
				}
			} else {
				if dst, err = EncodeValueToHessianV1(o, dst, value); err != nil {
					return dst, err
				}
			}
		}

	} else {
		keys := keys
		sorted := make([]reflect.Value, 0, len(keys))
		for i := range keys {
			sorted = append(sorted, keys[i])
		}

		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].CanInterface() && sorted[j].CanInterface() {
				ii := sorted[i]
				keyi := ii.Interface()
				valuei := v.MapIndex(ii)
				ji := sorted[j]
				keyj := ji.Interface()
				valuej := v.MapIndex(ji)
				return o.less(keyi, keyj, valuei, valuej)
			}
			return false
		})
		for i := 0; i < len(sorted); i++ {
			key := sorted[i]
			if key.CanInterface() { // Fast path
				if dst, err = EncodeToHessianV1(o, dst, key.Interface()); err != nil {
					return dst, err
				}
			} else {
				if dst, err = EncodeValueToHessianV1(o, dst, key); err != nil {
					return dst, err
				}
			}

			value := v.MapIndex(key)
			if value.CanInterface() { // Fast path
				if dst, err = EncodeToHessianV1(o, dst, value.Interface()); err != nil {
					return dst, err
				}
			} else {
				if dst, err = EncodeValueToHessianV1(o, dst, value); err != nil {
					return dst, err
				}
			}
		}
	}

	return EncodeMapEndToHessianV1(o, dst)
}
