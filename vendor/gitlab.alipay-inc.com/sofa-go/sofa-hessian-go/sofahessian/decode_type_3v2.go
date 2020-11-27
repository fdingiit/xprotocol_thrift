package sofahessian

import "bufio"

func DecodeTypeHessian3V2(o *DecodeContext, reader *bufio.Reader) (string, error) {
	var (
		b   []byte
		err error
	)

	b, err = DecodeTypeToHessian3V2(o, reader, b)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func DecodeTypeToHessian3V2(o *DecodeContext, reader *bufio.Reader, dst []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodetype")
		defer o.tracer.OnTraceStop("decodetype")
	}

	codes, err := reader.Peek(1)
	if err != nil {
		return dst, err
	}

	var typ string

	switch codes[0] {
	case 0x74:
		// consume one byte
		// nolint
		reader.ReadByte()
		var u16 uint16
		u16, err = readUint16FromReader(reader)
		if err != nil {
			return dst, err
		}

		length := int(u16)
		dst = allocAtLeast(dst, length)
		if err = readAtLeastBytesFromReader(reader, length, dst[len(dst)-length:]); err != nil {
			return dst, err
		}
		typ = string(dst[len(dst)-length:])
		if err = o.addTyperefs(typ); err != nil {
			return dst, err
		}

	case 0x54, 0x75:
		// consume one byte
		// nolint
		reader.ReadByte()

		var u32 uint32
		u32, err = readUint32FromReader(reader)
		if err != nil {
			return dst, err
		}

		typ, err = o.getTyperefs(int(u32))
		if err != nil {
			return dst, err
		}
		dst = append(dst, typ...)

	default:
	}

	return dst, nil
}
