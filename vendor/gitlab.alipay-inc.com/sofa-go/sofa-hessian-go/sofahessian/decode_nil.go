package sofahessian

import "bufio"

func DecodeNilHessian4V2(o *DecodeContext, reader *bufio.Reader) error {
	return DecodeNilToHessian4V2(o, reader)
}

func DecodeNilToHessian4V2(o *DecodeContext, reader *bufio.Reader) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodenil")
		defer o.tracer.OnTraceStop("decodenil")
	}

	c, err := reader.ReadByte()
	if err != nil {
		return err
	}

	switch c {
	case 'N':
		return nil
	default:
		return ErrDecodeMalformedBool
	}
}

func DecodeNilHessian3V2(o *DecodeContext, reader *bufio.Reader) error {
	return DecodeNilToHessian3V2(o, reader)
}

func DecodeNilToHessian3V2(o *DecodeContext, reader *bufio.Reader) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodenil")
		defer o.tracer.OnTraceStop("decodenil")
	}

	c, err := reader.ReadByte()
	if err != nil {
		return err
	}

	switch c {
	case 'N':
		return nil
	default:
		return ErrDecodeMalformedBool
	}
}

func DecodeNilHessianV1(o *DecodeContext, reader *bufio.Reader) error {
	return DecodeNilToHessianV1(o, reader)
}

func DecodeNilToHessianV1(o *DecodeContext, reader *bufio.Reader) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodenil")
		defer o.tracer.OnTraceStop("decodenil")
	}

	c, err := reader.ReadByte()
	if err != nil {
		return err
	}

	switch c {
	case 'N':
		return nil
	default:
		return ErrDecodeMalformedBool
	}
}
