package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func EncryptAES(key, plaintext []byte, ivString string) (string, error) {
	if len(key) != 16 {
		return "", fmt.Errorf("length of secret key should be 16 for 128 bits key size")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	paddedText := append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...)

	// 使用字符串作为 IV
	iv := []byte(ivString)
	if len(iv) != aes.BlockSize {
		return "", fmt.Errorf("length of IV string should be equal to AES block size (16 bytes)")
	}

	ciphertext := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedText)

	finalCiphertext := append(iv, ciphertext...)
	return base64.URLEncoding.EncodeToString(finalCiphertext), nil
}
