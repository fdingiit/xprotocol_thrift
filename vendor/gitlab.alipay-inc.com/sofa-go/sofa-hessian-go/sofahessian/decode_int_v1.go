package sofahessian

import (
	"bufio"
)

func DecodeInt32HessianV1(o *DecodeContext, reader *bufio.Reader) (int32, error) {
	var i int32
	err := DecodeInt32ToHessianV1(o, reader, &i)
	return i, err
}

func DecodeInt64HessianV1(o *DecodeContext, reader *bufio.Reader) (int64, error) {
	var i int64
	err := DecodeInt64ToHessianV1(o, reader, &i)
	return i, err
}

func DecodeInt32ToHessianV1(o *DecodeContext, reader *bufio.Reader, i *int32) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if c1 == 'I' {
		ix, err := readInt32FromReader(reader)
		*i = ix
		return err
	}

	return ErrDecodeCannotDecodeInt32
}

func DecodeInt64ToHessianV1(o *DecodeContext, reader *bufio.Reader, i *int64) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if c1 == 'L' {
		ix, err := readUint64FromReader(reader)
		if err != nil {
			return err
		}
		*i = int64(ix)
		return err
	}

	return ErrDecodeCannotDecodeInt32
}
