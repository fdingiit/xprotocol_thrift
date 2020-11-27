package sofahessian

import "bufio"

func DecodeTypeHessian4V2(o *DecodeContext, reader *bufio.Reader) (string, error) {
	var (
		b   []byte
		err error
	)

	b, err = DecodeTypeToHessian4V2(o, reader, b)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func DecodeTypeToHessian4V2(o *DecodeContext, reader *bufio.Reader, dst []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodetype")
		defer o.tracer.OnTraceStop("decodetype")
	}

	codes, err := reader.Peek(1)
	if err != nil {
		return dst, err
	}

	switch codes[0] {
	case 0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13,
		0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b,
		0x1c, 0x1d, 0x1e, 0x1f,
		0x30, 0x31, 0x32, 0x33,
		0x52, 0x53:
		var typ []byte
		if typ, err = DecodeStringToHessian4V2(o, reader, nil); err != nil {
			return nil, err
		}
		// Copy bytes to string
		if err = o.addTyperefs(string(typ)); err != nil {
			return dst, err
		}
		dst = append(dst, typ...)
		return dst, err

	default:
		var (
			refid int32
			typ   string
		)
		refid, err = DecodeInt32Hessian4V2(o, reader)
		if err != nil {
			return dst, err
		}

		typ, err = o.getTyperefs(int(refid))
		if err != nil {
			return dst, err
		}

		dst = append(dst, typ...)
		return dst, nil
	}
}
