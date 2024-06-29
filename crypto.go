/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package goutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

type CryptoUtils struct{}

func NewCryptoUtils() *CryptoUtils {
	return &CryptoUtils{}
}

func (u *CryptoUtils) Encrypt(text string, pass string) (string, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	key := pbkdf2.Key([]byte(pass), salt, 10000, 32, sha3.New512)

	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	ciphertext = append(salt, append(nonce, ciphertext...)...)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (u *CryptoUtils) Decrypt(cipherText string, pass string) (string, error) {
	enc, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(enc) < 28 {
		return "", fmt.Errorf("text too short")
	}

	salt := enc[:16]
	nonce := enc[16:28]
	cipherBytes := enc[28:]

	key := pbkdf2.Key([]byte(pass), salt, 10000, 32, sha3.New512)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesGCM.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
