/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package goutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
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
	key := pbkdf2.Key([]byte(pass), salt, 10000, 32, sha512.New)

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

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	ciphertext = append(salt, ciphertext...) // 将 salt 添加到密文开头
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (u *CryptoUtils) Decrypt(text string, pass string) (string, error) {
	enc, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	if len(enc) < 16 {
		return "", fmt.Errorf("text too short")
	}
	salt := enc[:16]
	enc = enc[16:]

	key := pbkdf2.Key([]byte(pass), salt, 10000, 32, sha512.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(enc) < nonceSize {
		return "", fmt.Errorf("text too short")
	}
	nonce, cipherBytes := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
