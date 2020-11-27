// +build fuzz

package sofabolt

import (
	"bytes"
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

	var cmd Command
	_, err := cmd.Read(NewReadOption(), bytes.NewReader(data))
	if err != nil {
		return 0
	}

	d, err := cmd.Write(NewWriteOption(), nil)
	if err != nil {
		return 0
	}

	if !bytes.Equal(data, d) {
		panic("not equal")
	}

	return 1
}
