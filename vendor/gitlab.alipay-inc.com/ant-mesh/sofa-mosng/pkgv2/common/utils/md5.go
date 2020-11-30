package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
)

func Md5Check(a, b interface{}) bool {
	astr := getStr(a)
	bstr := getStr(b)

	if astr == bstr {
		return true
	}

	return false
}

func getStr(a interface{}) string {
	ab, _ := json.Marshal(a)
	as := md5.Sum(ab)
	return fmt.Sprintf("%x", as)
}
