package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var digitAndLetters = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = digitAndLetters[rand.Intn(len(digitAndLetters))]
	}
	return string(b)
}
