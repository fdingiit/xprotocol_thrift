package sofahessian

import (
	"bufio"
	"math"
)

func DecodeFloat64Hessian4V2(o *DecodeContext, reader *bufio.Reader) (float64, error) {
	var i float64
	err := DecodeFloat64ToHessian4V2(o, reader, &i)
	return i, err
}

func DecodeFloat64ToHessian4V2(o *DecodeContext, reader *bufio.Reader, i *float64) error {
	if o.tracer != nil {
		o.tracer.OnTraceStart("decodefloat64")
		defer o.tracer.OnTraceStop("decodefloat64")
	}

	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	if c1 == 0x44 {
		u64, err := readUint64FromReader(reader)
		if err != nil {
			return err
		}
		*i = math.Float64frombits(u64)
		return nil
	}

	switch c1 {
	case 0x5b:
		*i = 0.0
		return nil

	case 0x5c:
		*i = 1.0
		return nil

	case 0x5d:
		c2, err := reader.ReadByte()
		if err != nil {
			return err
		}
		*i = float64(int8(c2))
		return nil

	case 0x5e:
		i16, err := readInt16FromReader(reader)
		if err != nil {
			return err
		}
		*i = float64(i16)
		return nil

	case 0x5f:
		i32, err := readInt32FromReader(reader)
		if err != nil {
			return err
		}
		*i = float64(i32) * 0.001
		return nil
	}

	return ErrDecodeMalformedDouble
}
