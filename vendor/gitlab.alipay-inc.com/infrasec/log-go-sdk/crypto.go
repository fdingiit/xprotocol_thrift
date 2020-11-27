package sls

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

var ivspec = []byte("0000000000000000")

func Encrypt(input, key string) (string, error) {
	if len(input) == 0 {
		return "", fmt.Errorf("failed to enctrypt, plain content empty")
	}

	key, err := formatKey(key)
	if err != nil {
		return "", err
	}

	//first base64 encode
	base64encodeStr := string(base64.StdEncoding.EncodeToString([]byte(input)))
	//aes encrypt
	aesEncodeStr, err := aesEncodeStr(base64encodeStr, key)
	if err != nil {
		return "", err
	}

	//second base64 encode
	base64encodeStr = string(base64.StdEncoding.EncodeToString([]byte(aesEncodeStr)))

	return base64encodeStr, nil
}

func Decrypt(input, key string) (string, error) {
	if len(input) == 0 {
		return "", fmt.Errorf("failed to dectrypt, cipher content empty")
	}

	key, err := formatKey(key)
	if err != nil {
		return "", err
	}

	//first base64 decode
	base64decodeStr, err := base64DecodeStr(input)
	if err != nil {
		return "", err
	}

	//aes decrypt
	aesDecodeStr, err := aesDecodeStr(base64decodeStr, key)
	if err != nil {
		return "", err
	}

	//second base64 decode
	base64decodeStr, err = base64DecodeStr(aesDecodeStr)
	if err != nil {
		return "", err
	}

	return base64decodeStr, nil
}

func formatKey(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key is none")
	}
	key = key + "0000000000000000"
	return key[:16], nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

//aes encrypt
func aesEncodeStr(src, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(src) == 0 {
		return "", fmt.Errorf("aes encrypt, plain content empty")
	}

	ecb := cipher.NewCBCEncrypter(block, ivspec)
	content := []byte(src)
	content = pkcs5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	return hex.EncodeToString(crypted), nil
}

//aes decrypt
func aesDecodeStr(crypt, key string) (string, error) {
	crypted, err := hex.DecodeString(strings.ToLower(crypt))
	if err != nil {
		return "", err
	}

	if len(crypted) == 0 {
		return "", fmt.Errorf("aes decrypt, cipher content empty")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ecb := cipher.NewCBCDecrypter(block, ivspec)
	decrypted := make([]byte, len(crypted))
	ecb.CryptBlocks(decrypted, crypted)

	// fix panic
	if len(decrypted) <= 0 || int(decrypted[len(decrypted)-1]) >= len(decrypted) {
		return "", fmt.Errorf("decrypted error : decrypted out of range")
	}

	return string(pkcs5Trimming(decrypted)), nil
}

// base64 decode
func base64DecodeStr(src string) (string, error) {
	a, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	return string(a), nil
}
