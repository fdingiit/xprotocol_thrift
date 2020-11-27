package sofahessian

import (
	"bufio"
	"time"
)

func DecodeDateHessian4V2(o *DecodeContext, reader *bufio.Reader) (time.Time, error) {
	var t time.Time
	err := DecodeDateToHessian4V2(o, reader, &t)
	return t, err
}

func DecodeDateToHessian4V2(o *DecodeContext, reader *bufio.Reader, t *time.Time) error {
	c1, err := reader.ReadByte()
	if err != nil {
		return err
	}

	if c1 == 0x4a {
		u64, err := readUint64FromReader(reader)
		if err != nil {
			return err
		}

		if u64/1000/3600/24/365 >= 2262-1970 { // try to save the time when year after 2262
			var mt time.Time
			*t = mt.Round(time.Millisecond * time.Duration(u64))
		} else {
			*t = time.Unix(int64(u64/1000), int64(u64)%1000*10e5)
		}

		return nil
	}

	if c1 == 0x4b {
		i32, err := readUint32FromReader(reader)
		if err != nil {
			return err
		}

		*t = time.Unix(int64(i32*60), 0)
		return nil
	}

	return ErrDecodeMalformedDate
}
