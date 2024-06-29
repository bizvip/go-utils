/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package goutils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

type AesUtils struct{}

func NewAesUtils() *AesUtils { return &AesUtils{} }

// PKCS7Padding adds padding to the original text to fit the block size
func (a *AesUtils) PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	latest := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, latest...)
}

func (a *AesUtils) EncryptAES(key, plaintext []byte, ivString string) (string, error) {
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
