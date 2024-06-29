/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package goutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/argon2"
)

type PwdUtils struct{}

// GenerateSalt 生成一个随机盐值
func (p *PwdUtils) GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(salt), nil
}

// HashPassword 使用 Argon2id 对密码进行加盐哈希
func (p *PwdUtils) HashPassword(password, salt string) (string, error) {
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)
	return hashEncoded, nil
}

// Encrypt 使用 AES 加密数据
func (p *PwdUtils) Encrypt(plaintext, key string) (string, error) {
	keyHash := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.RawStdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 使用 AES 解密数据
func (p *PwdUtils) Decrypt(ciphertext, key string) (string, error) {
	keyHash := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertextBytes, err := base64.RawStdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertextBytes := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// VerifyPassword 验证密码
func (p *PwdUtils) VerifyPassword(password, salt, hash string) (bool, error) {
	expectedHash, err := p.HashPassword(password, salt)
	if err != nil {
		return false, err
	}
	return hash == expectedHash, nil
}
