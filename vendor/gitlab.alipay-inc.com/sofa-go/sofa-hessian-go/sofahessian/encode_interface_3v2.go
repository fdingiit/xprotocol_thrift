package sofahessian

import (
	"errors"
	"reflect"
	"time"
)

// EncodeValueToHessian3v2 encodes reflect.Value to dst.
func EncodeValueToHessian3v2(o *EncodeContext, dst []byte, value reflect.Value) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encodevalue")
		defer o.tracer.OnTraceStop("encodevalue")
	}

	if value.Kind() == reflect.Invalid {
		return dst, ErrEncodeCannotInvalidValue
	} else if value.Kind() == reflect.Ptr && value.IsNil() {
		return EncodeNilToHessian3V2(o, dst)
	}

	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		if value.CanInterface() {
			return EncodeListToHessian3V2(o, dst, value.Interface())
		}
		return dst, ErrEncodeSliceCannotBeInterfaced

	case reflect.Map:
		if value.CanInterface() {
			return EncodeMapToHessian3V2(o, dst, value.Interface())
		}
		return dst, ErrEncodeMapCannotBeInterfaced

	case reflect.Struct:
		if value.CanInterface() {
			return EncodeObjectToHessian3V2(o, dst, value.Interface())
		}
		return dst, ErrEncodeStructCannotBeInterfaced

	case reflect.Ptr:
		// **T => *T
		indir := value.Elem()
		if indir.Kind() == reflect.Struct {
			return EncodeObjectToHessian3V2(o, dst, value.Interface())
		}

		if value.CanInterface() {
			return EncodeToHessian3V2(o, dst, indir.Interface())
		}

		return dst, ErrEncodePtrCannotBeInterfaced

	case reflect.Bool:
		return EncodeBoolToHessian3V2(o, dst, value.Bool())
	case reflect.Int:
		return EncodeInt64ToHessian3V2(o, dst, value.Int())
	case reflect.Int8:
		return EncodeInt32ToHessian3V2(o, dst, int32(value.Int()))
	case reflect.Int16:
		return EncodeInt32ToHessian3V2(o, dst, int32(value.Int()))
	case reflect.Int32:
		return EncodeInt32ToHessian3V2(o, dst, int32(value.Int()))
	case reflect.Int64:
		return EncodeInt64ToHessian3V2(o, dst, value.Int())
	case reflect.Uint:
		return EncodeInt64ToHessian3V2(o, dst, int64(value.Uint()))
	case reflect.Uint8:
		return EncodeInt32ToHessian3V2(o, dst, int32(value.Uint()))
	case reflect.Uint16:
		return EncodeInt32ToHessian3V2(o, dst, int32(value.Uint()))
	case reflect.Uint32:
		return EncodeInt32ToHessian3V2(o, dst, int32(value.Uint()))
	case reflect.Uint64:
		return EncodeInt64ToHessian3V2(o, dst, int64(value.Uint()))
	case reflect.Float32:
		return EncodeFloat64ToHessian3V2(o, dst, value.Float())
	case reflect.Float64:
		return EncodeFloat64ToHessian3V2(o, dst, value.Float())
	case reflect.String:
		return EncodeStringToHessian3V2(o, dst, value.String())
	case reflect.Interface:
		return EncodeToHessian3V2(o, dst, value.Elem())
	case reflect.Uintptr:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.UnsafePointer:
		fallthrough
	default:
		return dst, errors.New("hessian: cannot encode type " + value.Kind().String())
	}
}

// EncodeHessian3V2 encodes the interface to dst.
func EncodeHessian3V2(o *EncodeContext, value interface{}) ([]byte, error) {
	return EncodeToHessian3V2(o, nil, value)
}

// EncodeToHessian3V2 encodes the interface to dst.
func EncodeToHessian3V2(o *EncodeContext, dst []byte, value interface{}) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("encode")
		defer o.tracer.OnTraceStop("encode")
	}

	if o.maxdepth > 0 {
		o.addDepth()
		if o.depth > o.maxdepth {
			return nil, ErrEncodeMaxDepthExceeded
		}
		defer o.subDepth()
	}

	switch v := value.(type) {
	// Fast path without reflection
	case HessianEncoder:
		return v.HessianEncode(o, dst)
	case *[]byte:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeBinaryToHessian3V2(o, dst, *v)
	case []byte:
		return EncodeBinaryToHessian3V2(o, dst, v)

	case *string:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeStringToHessian3V2(o, dst, *v)
	case string:
		return EncodeStringToHessian3V2(o, dst, v)

	case int:
		return EncodeInt64ToHessian3V2(o, dst, int64(v))
	case *int:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeInt64ToHessian3V2(o, dst, int64(*v))
	case uint:
		return EncodeInt64ToHessian3V2(o, dst, int64(v))
	case *uint:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeInt64ToHessian3V2(o, dst, int64(*v))
	case uint8:
		return EncodeInt32ToHessian3V2(o, dst, int32(v))
	case int8:
		return EncodeInt32ToHessian3V2(o, dst, int32(v))
	case uint16:
		return EncodeInt32ToHessian3V2(o, dst, int32(v))
	case int16:
		return EncodeInt32ToHessian3V2(o, dst, int32(v))
	case uint32:
		return EncodeInt32ToHessian3V2(o, dst, int32(v))
	case int32:
		return EncodeInt32ToHessian3V2(o, dst, v)
	case uint64:
		return EncodeInt64ToHessian3V2(o, dst, int64(v))
	case int64:
		return EncodeInt64ToHessian3V2(o, dst, v)

	case *uint8:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeInt32ToHessian3V2(o, dst, int32(*v))
	case *int8:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeInt32ToHessian3V2(o, dst, int32(*v))
	case *uint16:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeInt32ToHessian3V2(o, dst, int32(*v))
	case *int16:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeInt32ToHessian3V2(o, dst, int32(*v))
	case *uint32:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeInt64ToHessian3V2(o, dst, int64(*v))
	case *int32:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeInt32ToHessian3V2(o, dst, *v)
	case *uint64:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeInt64ToHessian3V2(o, dst, int64(*v))
	case *int64:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeInt64ToHessian3V2(o, dst, *v)

	case *float32:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeFloat64ToHessian3V2(o, dst, float64(*v))
	case float32:
		return EncodeFloat64ToHessian3V2(o, dst, float64(v))
	case *float64:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeFloat64ToHessian3V2(o, dst, *v)
	case float64:
		return EncodeFloat64ToHessian3V2(o, dst, v)

	case *bool:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeBoolToHessian3V2(o, dst, *v)
	case bool:
		return EncodeBoolToHessian3V2(o, dst, v)

	case nil:
		return EncodeNilToHessian3V2(o, dst)

	case *time.Time:
		if v == nil {
			return EncodeNilToHessian3V2(o, dst)
		}
		return EncodeDateToHessian3V2(o, dst, *v)
	case time.Time:
		return EncodeDateToHessian3V2(o, dst, v)

	default:
		return EncodeValueToHessian3v2(o, dst, reflect.ValueOf(v))
	}
}
