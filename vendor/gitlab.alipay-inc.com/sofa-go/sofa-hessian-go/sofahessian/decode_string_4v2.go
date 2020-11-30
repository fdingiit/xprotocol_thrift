package sofahessian

import "bufio"

func DecodeStringHessian4V2(o *DecodeContext, reader *bufio.Reader) (string, error) {
	var (
		b   []byte
		err error
	)

	b, err = DecodeStringToHessian4V2(o, reader, b)
	if err != nil {
		return "", err
	}
	return string(b), err
}

// DecodeStringToHessian4V2 decodes dst to string.
func DecodeStringToHessian4V2(o *DecodeContext, reader *bufio.Reader, s []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodestring")
		defer o.tracer.OnTraceStop("decodestring")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return s, err
	}

	length, err := decodeStringLengthToHessian4V2(o, reader, c1)
	if err != nil {
		return s, err
	}

	lastChunk := true
	if c1 == 'R' {
		lastChunk = false
	}

	for i := uint16(0); ; {
		if i == length {
			if lastChunk {
				return s, nil
			}

			c1, err := reader.ReadByte()
			if err != nil {
				return s, err
			}

			if c1 != 'R' {
				lastChunk = true
			}

			sublength, err := decodeStringLengthToHessian4V2(o, reader, c1)
			if err != nil {
				return s, err
			}
			length += sublength

			continue
		}

		// Read UTF8 codepoint
		r, size, err := reader.ReadRune()
		if err != nil {
			return s, err
		}

		// Convert the UTF8 codepoint to bytes
		s = appendRune(s, uint32(r), size)

		i++
	}
}

func decodeStringLengthToHessian4V2(o *DecodeContext, reader *bufio.Reader, peek byte) (uint16, error) {
	switch {
	case 0x00 <= peek && peek <= 0x1f:
		return uint16(peek - 0x00), nil

	case 0x30 <= peek && peek <= 0x33:
		c2, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		return (uint16(peek)-0x30)<<8 + uint16(c2), nil
	case peek == 'R' || peek == 'S':
		u16, err := readUint16FromReader(reader)
		if err != nil {
			return 0, err
		}

		return u16, nil
	}

	return 0, ErrDecodeMalformedString
}
