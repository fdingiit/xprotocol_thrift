package sofahessian

import (
	"bufio"
	"reflect"
)

func DecodeMapHessianV1(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	codes, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}
	return decodeMapHessianV1(o, reader, codes[0])
}

func decodeMapHessianV1(o *DecodeContext, reader *bufio.Reader, peek byte) (interface{}, error) {
	if peek == 0x48 {
		return DecodeUntypedMapHessianV1(o, reader)
	} else if peek == 0x4d {
		return DecodeTypedMapHessianV1(o, reader)
	}

	return nil, ErrDecodeMalformedMap
}

func DecodeUntypedMapHessianV1(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	var i interface{}
	err := DecodeUntypedMapToHessianV1(o, reader, &i)
	return i, err
}

func DecodeTypedMapHessianV1(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	var i interface{}
	err := DecodeTypedMapToHessianV1(o, reader, &i)
	return i, err
}

func DecodeUntypedMapToHessianV1(o *DecodeContext, reader *bufio.Reader, obj *interface{}) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodeuntypedmap")
		defer o.tracer.OnTraceStop("decodeuntypedmap")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	if c1 != 0x48 {
		return ErrDecodeMalformedUntypedMap
	}

	m, err := decodeUntypedMapHessianV1(o, reader)
	if err != nil {
		return err
	}
	*obj = m

	return nil
}

func DecodeTypedMapToHessianV1(o *DecodeContext, reader *bufio.Reader, obj *interface{}) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodetypedmap")
		defer o.tracer.OnTraceStop("decodetypedmap")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if c1 != 0x4d {
		return ErrDecodeMalformedTypedMap
	}

	typ, err := DecodeTypeHessianV1(o, reader)
	if err != nil {
		return err
	}

	var m map[interface{}]interface{}

	if typ == "" {
		m, err = decodeUntypedMapHessianV1(o, reader)
		if err != nil {
			return err
		}
		*obj = m
		return nil
	}

	ci, ok := o.loadClassTypeSchema(typ)
	if !ok { // Use JavaMap
		m, err = decodeUntypedMapHessianV1(o, reader)
		if err != nil {
			return err
		}
		*obj = &JavaMap{
			class: typ,
			m:     m,
		}
		return nil
	}

	// Peek byte at first
	codes, err := reader.Peek(1)
	if err != nil {
		return err
	}

	// Concrete type
	value := reflect.New(ci.base)

	if err = o.addObjectrefs(value.Interface()); err != nil {
		return err
	}

	var (
		key interface{}
		val interface{}
	)

	structvalue := value.Elem()
	for codes[0] != 0x7a {
		key, err = DecodeHessianV1(o, reader)
		if err != nil {
			return err
		}

		fieldkey, ok := key.(string)
		if !ok {
			return ErrDecodeTypedMapKeyNotString
		}

		val, err = DecodeHessianV1(o, reader)
		if err != nil {
			return err
		}

		fieldvalue := structvalue.FieldByName(fieldkey)
		if fieldvalue.CanSet() {
			if err = safeSetValueByReflect(fieldvalue, val); err != nil {
				return err
			}
		} else {
			return ErrDecodeTypedMapValueNotAssign
		}

		codes, err = reader.Peek(1)
		if err != nil {
			return err
		}
	}

	// Discard the last byte
	_, err = reader.ReadByte()
	if err != nil {
		return err
	}

	return nil
}

func decodeUntypedMapHessianV1(o *DecodeContext, reader *bufio.Reader) (map[interface{}]interface{}, error) {
	m := make(map[interface{}]interface{}, 4) // Allow config it?

	if err := o.addObjectrefs(m); err != nil {
		return m, err
	}

	codes, err := reader.Peek(1)
	if err != nil {
		return m, err
	}

	var (
		key   interface{}
		value interface{}
	)

	for codes[0] != 0x7a {
		key, err = DecodeHessianV1(o, reader)
		if err != nil {
			return m, err
		}

		value, err = DecodeHessianV1(o, reader)
		if err != nil {
			return m, err
		}

		if !safeSetMap(&m, key, value) {
			// FYI(detailyang): cannot use %+v which maybe infinite recursion because of self-referential data structures
			return m, ErrDecodeMapUnhashable
		}

		codes, err = reader.Peek(1)
		if err != nil {
			return m, err
		}
	}

	_, err = reader.ReadByte()
	if err != nil {
		return m, err
	}

	return m, nil
}
