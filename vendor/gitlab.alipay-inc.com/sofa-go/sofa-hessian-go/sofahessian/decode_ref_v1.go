package sofahessian

import (
	"bufio"
)

func DecodeRefHessianV1(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decoderef")
		defer o.tracer.OnTraceStop("decoderef")
	}

	var i interface{}
	err := DecodeRefToHessianV1(o, reader, &i)
	return i, err
}

func DecodeRefToHessianV1(o *DecodeContext, reader *bufio.Reader, obj *interface{}) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	var refid uint32
	if c1 == 0x4a {
		c1, err = reader.ReadByte()
		if err != nil {
			return err
		}
		refid = uint32(c1)

	} else if c1 == 0x4b {
		var u16 uint16
		u16, err = readUint16FromReader(reader)
		if err != nil {
			return err
		}
		refid = uint32(u16)

	} else if c1 == 0x52 {
		refid, err = readUint32FromReader(reader)
		if err != nil {
			return err
		}

	} else {
		return ErrDecodeMalformedReference
	}

	i, err := o.getObjectrefs(int(refid))
	if err != nil {
		return err
	}
	*obj = i

	return nil
}
