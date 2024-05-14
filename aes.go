package utils

import (
	"bytes"
)

// PKCS7Padding adds padding to the original text to fit the block size
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	latest := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, latest...)
}
