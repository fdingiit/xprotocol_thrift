package fastsimplemap

import (
	"encoding/hex"
	"fmt"
)

func Fuzz(data []byte) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(hex.EncodeToString(data))
			panic(r)
		}
	}()

	var (
		om FastSimpleMap
		nm FastSimpleMap
	)

	err := om.Decode(data)
	if err != nil {
		return 0
	}

	n := om.GetEncodeSize()
	e := make([]byte, n)
	om.Encode(e)

	if err := nm.Decode(e); err != nil {
		panic("failed to decode")
	}

	if !nm.Equal(&om) {
		panic("failed to equal")
	}

	return 1
}
