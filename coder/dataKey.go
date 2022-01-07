package coder

import (
	"crypto/sha1"
	_ "embed"
)

var (
	//go:embed salt.bin
	saltBin []byte
)

func DataKey(key, iv []byte, otpSecret string) []byte {
	newKey := make([]byte, len(key))
	copy(newKey, key)
	keySum := saltBin
	keySum = append(keySum, key...)
	keySum = append(keySum, iv...)
	keySum = append(keySum, []byte(otpSecret)...)
	hasher := sha1.New()
	hasher.Write(keySum)
	sha1Key := hasher.Sum(nil)
	copy(newKey, sha1Key)
	return newKey
}
