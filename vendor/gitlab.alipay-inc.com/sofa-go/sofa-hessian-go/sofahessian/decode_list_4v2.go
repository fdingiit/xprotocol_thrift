package sofahessian

import (
	"bufio"
	"fmt"
	"reflect"
)

func DecodeListHessian4V2(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	var i interface{}
	err := DecodeListToHessian4V2(o, reader, &i)
	return i, err
}

func DecodeListToHessian4V2(o *DecodeContext, reader *bufio.Reader, obj *interface{}) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodelist")
		defer o.tracer.OnTraceStop("decodelist")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	var (
		typ    string
		length int32
	)

	switch {
	case c1 == 0x56:
		typ, err = DecodeTypeHessian4V2(o, reader)
		if err != nil {
			return err
		}
		length, err = DecodeInt32Hessian4V2(o, reader)
		if err != nil {
			return err
		}

	case c1 == 0x58:
		length, err = DecodeInt32Hessian4V2(o, reader)
		if err != nil {
			return err
		}

	case c1 >= 0x78 && c1 <= 0x7F:
		length = int32(c1) - 0x78

	case c1 >= 0x70 && c1 <= 0x77:
		typ, err = DecodeTypeHessian4V2(o, reader)
		if err != nil {
			return err
		}

		length = int32(c1) - 0x70
	}

	if length < 0 || int(length) >= o.GetMaxListLength() {
		return ErrDecodeMaxListLengthExceeded
	}

	if typ != "" {
		ci, ok := o.loadClassTypeSchema(typ)
		if ok { // concrete type
			if ci.base.Kind() != reflect.Slice && ci.base.Kind() != reflect.Array {
				return fmt.Errorf("hessian: expect slice/array type but got %s", ci.base.Kind().String())
			}

			value := reflect.MakeSlice(ci.base, int(length), int(length))
			*obj = value.Interface()
			if err = o.addObjectrefs(*obj); err != nil {
				return err
			}

			var list []interface{}

			list, err = decodeBoundedListHessian4V2(o, reader, nil, length)
			if err != nil {
				return err
			}

			if len(list) != int(length) {
				return fmt.Errorf("hessian: expect [%d]T but got [%d]T", length, len(list))
			}

			for i := range list {
				if err = safeSetValueByReflect(value.Index(i), list[i]); err != nil {
					return err
				}
			}

		} else { // java list
			if length > 0 {
				list := make([]interface{}, 0, length)
				jl := &JavaList{class: typ, value: list}
				if err = o.addObjectrefs(jl); err != nil {
					return err
				}
				jl.value, err = decodeBoundedListHessian4V2(o, reader, jl.value, length)
				if err != nil {
					return err
				}
				*obj = jl
			} else {
				*obj = &JavaList{class: typ, value: []interface{}{}}
			}
		}

	} else { // []interface{}
		list := make([]interface{}, length)
		if err := o.addObjectrefs(list); err != nil {
			return err
		}
		for i := 0; i < int(length); i++ {
			obj, err := DecodeHessian4V2(o, reader)
			if err != nil {
				return err
			}
			list[i] = obj
		}
		*obj = list
	}

	return nil
}

func decodeBoundedListHessian4V2(o *DecodeContext, reader *bufio.Reader,
	list []interface{}, length int32) ([]interface{}, error) {
	for i := 0; i < int(length); i++ {
		obj, err := DecodeHessian4V2(o, reader)
		if err != nil {
			return list, err
		}
		list = append(list, obj)
	}
	return list, nil
}
