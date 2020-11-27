package sofahessian

import (
	"bufio"
	"math"
)

func DecodeFloat64ToHessianV1(o *DecodeContext, reader *bufio.Reader, i *float64) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	switch c1 {
	case 'D':
		u64, err := readUint64FromReader(reader)
		if err != nil {
			return err
		}
		*i = math.Float64frombits(u64)
		return nil

	default:
		return ErrDecodeMalformedDouble
	}
}

func DecodeFloat64HessianV1(o *DecodeContext, reader *bufio.Reader) (float64, error) {
	var i float64
	err := DecodeFloat64ToHessianV1(o, reader, &i)
	return i, err
}
