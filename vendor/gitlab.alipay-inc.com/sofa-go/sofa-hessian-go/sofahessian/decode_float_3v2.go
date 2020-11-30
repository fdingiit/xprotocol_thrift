package sofahessian

import (
	"bufio"
	"math"
)

func DecodeFloat64Hessian3V2(o *DecodeContext, reader *bufio.Reader) (float64, error) {
	var i float64
	err := DecodeFloat64ToHessian3V2(o, reader, &i)
	return i, err
}

func DecodeFloat64ToHessian3V2(o *DecodeContext, reader *bufio.Reader, i *float64) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	switch c1 {
	case 0x44:
		u64, err := readUint64FromReader(reader)
		if err != nil {
			return err
		}
		*i = math.Float64frombits(u64)
		return nil

	case 0x67:
		*i = 0.0

	case 0x68:
		*i = 1.0

	case 0x69:
		u8, err := reader.ReadByte()
		if err != nil {
			return err
		}
		*i = float64(u8)

	case 0x6a:
		u16, err := readUint16FromReader(reader)
		if err != nil {
			return err
		}
		*i = float64(u16)

	case 0x6b:
		u32, err := readUint32FromReader(reader)
		if err != nil {
			return err
		}
		*i = float64(math.Float32frombits(u32))
		return nil

	default:
		return ErrDecodeMalformedDouble
	}

	return nil
}
