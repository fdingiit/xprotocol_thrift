package convertor_helper

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"time"

	error2 "gitlab.alipay-inc.com/octopus/tbase-go-client/model/error"
)

// Int32ToString converts a int32 to string, work faster than fmt.Sprintf
func Int32ToString(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

// byteArrayToInt64 is a specialized version of strconv.byteArrayToInt64 that parses a base-10 encoded signed integer from a []byte.
//
// This can be used to avoid allocating a string, since strconv.byteArrayToInt64 only takes a string.
func ByteArrayToInt64(b []byte) (int64, error) {
	if len(b) == 0 {
		return 0, error2.NewTBaseClientInternalError("empty slice given to parseInt")
	}

	var neg bool
	if b[0] == '-' || b[0] == '+' {
		neg = b[0] == '-'
		b = b[1:]
	}

	n, err := byteArrayToUInt64(b)
	if err != nil {
		return 0, err
	}

	if neg {
		return -int64(n), nil
	}

	return int64(n), nil
}

// byte array to uint64
func byteArrayToUInt64(b []byte) (uint64, error) {
	if len(b) == 0 {
		return 0, errors.New("empty slice given to parseUint")
	}

	var n uint64

	for i, c := range b {
		if c < '0' || c > '9' {
			return 0, error2.NewTBaseClientInternalError(fmt.Sprintf("invalid character %c at position %d in parseUint", c, i))
		}

		n *= 10
		n += uint64(c - '0')
	}

	return n, nil
}

// BufferedBytesDelim reads a line from br and checks that the line ends with \r\n, returning the line without \r\n.
func BufferedBytesDelim(br *bufio.Reader) ([]byte, error) {
	b, err := br.ReadSlice('\n')
	if err != nil {
		return nil, err
	} else if len(b) < 2 || b[len(b)-2] != '\r' {
		return nil, fmt.Errorf("malformed resp %q", b)
	}
	return b[:len(b)-2], err
}

// Shuffle the string slice
func Shuffle(vals []string) []string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(vals); n > 0; n-- {
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
	}
	return vals
}

// get abs of int
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// get abs of int64
func AbsInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
