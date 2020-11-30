package sofahessian

import "bufio"

func DecodeBinaryHessian3V2(o *DecodeContext, reader *bufio.Reader) ([]byte, error) {
	p := []byte(nil)
	return DecodeBinaryToHessian3V2(o, reader, p)
}

func DecodeBinaryToHessian3V2(o *DecodeContext, reader *bufio.Reader, dst []byte) ([]byte, error) {
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
		u16    uint16
		length uint16
	)
	for c1 == 0x62 {
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

	if c1 == 0x42 {
		u16, err = readUint16FromReader(reader)
		if err != nil {
			return dst, err
		}
		dst = allocAtLeast(dst, int(u16))
		err = readAtLeastBytesFromReader(reader, int(u16), dst[len(dst)-int(u16):])
		if err != nil {
			return dst, err
		}

	} else if c1 >= 0x20 && c1 <= 0x2f {
		l := c1 - 0x20
		dst = allocAtLeast(dst, int(l))
		err = readAtLeastBytesFromReader(reader, int(l), dst[len(dst)-int(l):])
		if err != nil {
			return dst, err
		}

	} else {
		return dst, ErrDecodeMalformedBinary
	}

	return dst, nil
}
