package sofahessian

import "bufio"

func DecodeBinaryHessian4V2(o *DecodeContext, reader *bufio.Reader) ([]byte, error) {
	p := []byte(nil)
	return DecodeBinaryToHessian4V2(o, reader, p)
}

func DecodeBinaryToHessian4V2(o *DecodeContext, reader *bufio.Reader, dst []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodebinary")
		defer o.tracer.OnTraceStop("decodebinary")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return dst, err
	}

	if c1 >= 0x20 && c1 <= 0x2f {
		length := int(c1) - 0x20
		dst = allocAtLeast(dst, length)
		err = readAtLeastBytesFromReader(reader, length, dst[len(dst)-length:])
		return dst, err
	}

	var (
		length uint16
		u16    uint16
		c2     byte
	)

	for c1 == 0x41 {
		u16, err = readUint16FromReader(reader)
		if err != nil {
			return dst, err
		}

		length = u16
		dst = allocAtLeast(dst, int(length))
		err = readAtLeastBytesFromReader(reader, int(length), dst[len(dst)-int(length):])
		if err != nil {
			return dst, err
		}

		c1, err = reader.ReadByte()
		if err != nil {
			return dst, err
		}
	}

	switch {
	case c1 == 0x42:
		u16, err = readUint16FromReader(reader)
		if err != nil {
			return dst, err
		}
		length = u16
		dst = allocAtLeast(dst, int(length))
		err = readAtLeastBytesFromReader(reader, int(length), dst[len(dst)-int(length):])
		if err != nil {
			return dst, err
		}
	case c1 >= 0x20 && c1 <= 0x2f:
		length = uint16(c1) - 0x20
		dst = allocAtLeast(dst, int(length))
		err = readAtLeastBytesFromReader(reader, int(length), dst[len(dst)-int(length):])
		if err != nil {
			return dst, err
		}
	case c1 >= 0x34 && c1 <= 0x37:
		c2, err = reader.ReadByte()
		if err != nil {
			return dst, err
		}
		length = uint16(c1-0x34)<<8 + uint16(c2)
		dst = allocAtLeast(dst, int(length))
		err = readAtLeastBytesFromReader(reader, int(length), dst[len(dst)-int(length):])
		if err != nil {
			return dst, err
		}
	default:
		return dst, ErrDecodeMalformedBinary
	}

	return dst, nil
}
