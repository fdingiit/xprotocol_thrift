package sofahessian

import (
	"bufio"
	"reflect"
)

func DecodeMapHessian4V2(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	codes, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}
	return decodeMapHessian4V2(o, reader, codes[0])
}

func DecodeMapHessian3V2(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	codes, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}
	return decodeMapHessian3V2(o, reader, codes[0])
}

func decodeMapHessian3V2(o *DecodeContext, reader *bufio.Reader, peek byte) (interface{}, error) {
	if peek == 0x48 {
		return DecodeUntypedMapHessian3V2(o, reader)
	} else if peek == 0x4d {
		return DecodeTypedMapHessian3V2(o, reader)
	}

	return nil, ErrDecodeMalformedMap
}

func decodeMapHessian4V2(o *DecodeContext, reader *bufio.Reader, peek byte) (interface{}, error) {
	if peek == 0x48 {
		return DecodeUntypedMapHessian4V2(o, reader)
	} else if peek == 0x4d {
		return DecodeTypedMapHessian4V2(o, reader)
	}

	return nil, ErrDecodeMalformedMap
}

func DecodeUntypedMapHessian4V2(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	var i interface{}
	err := DecodeUntypedMapToHessian4V2(o, reader, &i)
	return i, err
}

func DecodeUntypedMapToHessian4V2(o *DecodeContext, reader *bufio.Reader, obj *interface{}) error {
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

	m, err := decodeUntypedMapHessian4V2(o, reader)
	if err != nil {
		return err
	}
	*obj = m

	return nil
}

func DecodeTypedMapHessian4V2(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	var i interface{}
	err := DecodeTypedMapToHessian4V2(o, reader, &i)
	return i, err
}

func DecodeTypedMapToHessian4V2(o *DecodeContext, reader *bufio.Reader, obj *interface{}) error {
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

	typ, err := DecodeTypeHessian4V2(o, reader)
	if err != nil {
		return err
	}

	var m map[interface{}]interface{}

	if typ == "" {
		m, err = decodeUntypedMapHessian4V2(o, reader)
		if err != nil {
			return err
		}
		*obj = m
		return nil
	}

	ci, ok := o.loadClassTypeSchema(typ)
	if !ok { // Use JavaMap
		m, err = decodeUntypedMapHessian4V2(o, reader)
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

	structvalue := value.Elem()
	for codes[0] != 0x5a {
		var (
			key interface{}
			val interface{}
		)

		key, err = DecodeHessian4V2(o, reader)
		if err != nil {
			return err
		}

		fieldkey, ok := key.(string)
		if !ok {
			return ErrDecodeTypedMapKeyNotString
		}

		val, err = DecodeHessian4V2(o, reader)
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

func decodeUntypedMapHessian4V2(o *DecodeContext, reader *bufio.Reader) (map[interface{}]interface{}, error) {
	m := make(map[interface{}]interface{}, 4) // Allow config it?

	if err := o.addObjectrefs(m); err != nil {
		return m, err
	}

	codes, err := reader.Peek(1)
	if err != nil {
		return m, err
	}

	for codes[0] != 0x5a {
		var (
			key   interface{}
			value interface{}
		)
		key, err = DecodeHessian4V2(o, reader)
		if err != nil {
			return m, err
		}

		value, err = DecodeHessian4V2(o, reader)
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

	// Discard the peek last byte
	_, err = reader.ReadByte()
	return m, err
}
