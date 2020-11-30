package sofahessian

import (
	"bufio"
)

func DecodeStringHessian3V2(o *DecodeContext, reader *bufio.Reader) (string, error) {
	var (
		b   []byte
		err error
	)

	b, err = DecodeStringToHessian3V2(o, reader, b)
	if err != nil {
		return "", err
	}
	return string(b), err
}

// DecodeStringToHessian3V2 decodes dst to string.
func DecodeStringToHessian3V2(o *DecodeContext, reader *bufio.Reader, s []byte) ([]byte, error) {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodestring")
		defer o.tracer.OnTraceStop("decodestring")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return s, err
	}

	for c1 == 's' {
		s, err = readLenAndUTF8StringFromReader(reader, s)
		if err != nil {
			return s, err
		}

		c1, err = reader.ReadByte()
		if err != nil {
			return s, err
		}
	}

	if c1 >= 0x00 && c1 <= 0x1F {
		s, err = readUTF8StringFromReader(reader, s, int(c1))
	} else if c1 == 0x53 {
		s, err = readLenAndUTF8StringFromReader(reader, s)
	} else {
		return s, ErrDecodeMalformedString
	}

	return s, err
}
