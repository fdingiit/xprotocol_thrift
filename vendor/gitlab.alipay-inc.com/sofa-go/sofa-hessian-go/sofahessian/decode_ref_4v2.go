package sofahessian

import "bufio"

func DecodeRefHessian4V2(o *DecodeContext, reader *bufio.Reader) (interface{}, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decoderef")
		defer o.tracer.OnTraceStop("decoderef")
	}

	var i interface{}
	err := DecodeRefToHessian4V2(o, reader, &i)
	return i, err
}

func DecodeRefToHessian4V2(o *DecodeContext, reader *bufio.Reader, obj *interface{}) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	if c1 != 0x51 {
		return ErrDecodeMalformedReference
	}

	refid, err := DecodeInt32Hessian4V2(o, reader)
	if err != nil {
		return err
	}

	i, err := o.getObjectrefs(int(refid))
	if err != nil {
		return err
	}
	*obj = i

	return nil
}
