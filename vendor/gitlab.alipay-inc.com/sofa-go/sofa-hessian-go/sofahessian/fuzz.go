// +build fuzz

package sofahessian

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"reflect"
)

func Fuzz(data []byte) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(hex.EncodeToString(data))
			panic(r)
		}
	}()
	dctx := NewDecodeContext()
	ectx := NewEncodeContext()
	br := bufio.NewReader(bytes.NewReader(data))
	v, err := DecodeHessian4V2(dctx, br)
	if err != nil {
		return 0
	}

	if fv, ok := v.(float64); ok && math.IsNaN(fv) {
		return 0
	}

	dst, err := EncodeHessian4V2(ectx, v)
	if err != nil {
		fmt.Println(hex.EncodeToString(data))
		panic(err)
	}

	vv, err := DecodeHessian4V2(dctx, bufio.NewReader(bytes.NewReader(dst)))
	if err != nil {
		if err == ErrDecodeMapUnhashable {
			return 0
		}

		fmt.Println(hex.EncodeToString(data))
		panic(err)
	}

	if mv, ok := v.(map[interface{}]interface{}); ok {
		if ov, ok := vv.(map[interface{}]interface{}); ok {
			_ = mv
			_ = ov
			return 0
		}
		log.Fatal("expect map[interface{}]interface{}")
	}

	if !reflect.DeepEqual(v, vv) {
		fmt.Println(hex.EncodeToString(data))
		fmt.Println(hex.EncodeToString(dst))
		panic("not equal")
	}

	return 1
}
