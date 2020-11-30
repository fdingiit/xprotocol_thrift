package sofahessian

import (
	"bufio"
	"errors"
	"reflect"
)

func DecodeObjectToHessian3V2(o *DecodeContext, reader *bufio.Reader, obj interface{}) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodeotobject")
		defer o.tracer.OnTraceStop("decodetoobject")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	var refid int32

	if c1 == 0x4f {
		err = decodeObjectDefinitionHessian3V2(o, reader)
		if err != nil {
			return err
		}
		return DecodeObjectToHessian3V2(o, reader, obj)

	} else if c1 == 0x6f {
		refid, err = DecodeInt32Hessian3V2(o, reader)
		if err != nil {
			return err
		}

	} else {
		return ErrDecodeMalformedObject
	}

	cd, err := o.getClassrefs(int(refid))
	if err != nil {
		return err
	}

	name := getInterfaceName(obj)
	if name != cd.class {
		return ErrDecodeUnmatchedObject
	}

	structvalue := decAllocReflectValue(reflect.ValueOf(obj))
	if err := o.addObjectrefs(obj); err != nil {
		return err
	}
	rt := decAllocReflectType(reflect.TypeOf(obj))

	for i := range cd.fields {
		field := cd.fields[i]
		fi, ok := lookupReflectField(rt, field)
		if !ok {
			if o.disallowMissingField {
				return errors.New("hessian: malformed class field (not found) " + field)
			}
			_, err := DecodeHessian3V2(o, reader)
			if err != nil {
				return err
			}

		} else {
			key := structvalue.Field(fi)
			if !key.CanSet() {
				return errors.New("hessian: malformed class field (unassignable) " + field)
			}

			value, err := DecodeHessian3V2(o, reader)
			if err != nil {
				return err
			}

			err = safeSetValueByReflect(key, value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DecodeObjectHessian3V2(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodeobject")
		defer o.tracer.OnTraceStop("decodeobject")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	var refid int32

	if c1 == 0x4f {
		err = decodeObjectDefinitionHessian3V2(o, reader)
		if err != nil {
			return nil, err
		}
		return DecodeObjectHessian3V2(o, reader)

	} else if c1 == 0x6f {
		refid, err = DecodeInt32Hessian3V2(o, reader)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, ErrDecodeMalformedObject
	}

	cd, err := o.getClassrefs(int(refid))
	if err != nil {
		return nil, err
	}

	ci, ok := o.loadClassTypeSchema(cd.class)
	if !ok { // Generic java object
		jo := &JavaObject{
			class:  cd.class,
			names:  make([]string, 0, len(cd.fields)),
			values: make([]interface{}, 0, len(cd.fields)),
		}

		if err := o.addObjectrefs(jo); err != nil {
			return nil, err
		}

		for i := range cd.fields {
			fieldname := cd.fields[i]
			if len(fieldname) == 0 {
				return nil, ErrDecodeObjectFieldCannotBeNull
			}

			fieldvalue, err := DecodeHessian3V2(o, reader)
			if err != nil {
				return nil, err
			}
			jo.names = append(jo.names, fieldname)
			jo.values = append(jo.values, fieldvalue)
		}

		return jo, nil
	}

	// Concrete type
	value := reflect.New(ci.base)
	structvalue := value.Elem()

	if err := o.addObjectrefs(value.Interface()); err != nil {
		return nil, err
	}

	for i := range cd.fields {
		field := cd.fields[i]
		fi, ok := lookupReflectField(ci.base, field)
		if !ok {
			if o.disallowMissingField {
				return nil, errors.New("hessian: malformed class field (not found) " + field)
			}
			_, err := DecodeHessian3V2(o, reader)
			if err != nil {
				return nil, err
			}
		} else {
			key := structvalue.Field(fi)
			if !key.CanSet() {
				return nil, errors.New("hessian: malformed class field (unassignable) " + field)
			}

			value, err := DecodeHessian3V2(o, reader)
			if err != nil {
				return nil, err
			}

			err = safeSetValueByReflect(key, value)
			if err != nil {
				return nil, err
			}
		}
	}

	return value.Interface(), nil
}

func decodeObjectDefinitionHessian3V2(o *DecodeContext, reader *bufio.Reader) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodeobjectdefinition")
		defer o.tracer.OnTraceStop("decodeobjectdefinition")
	}

	u32, err := DecodeInt32Hessian3V2(o, reader)
	if err != nil {
		return err
	}

	// TODO(detailyang): config it to avoid DDOS
	typ := string(make([]byte, u32))
	err = readAtLeastBytesFromReader(reader, int(u32), s2b(typ))
	if err != nil {
		return err
	}

	// TODO(detailyang): cache class definition
	fieldslen, err := DecodeInt32Hessian3V2(o, reader)
	if err != nil {
		return err
	}

	if fieldslen < 0 || int(fieldslen) > o.GetMaxObjectFields() {
		return ErrDecodeMaxObjectFieldsExceeded
	}

	cd := ClassDefinition{
		class:  typ,
		fields: make([]string, 0, fieldslen),
	}

	for i := 0; i < int(fieldslen); i++ {
		field, err := DecodeStringHessian3V2(o, reader)
		if err != nil {
			return err
		}
		cd.fields = append(cd.fields, field)
	}

	if err := o.addClassrefs(cd); err != nil {
		return err
	}

	return nil
}
