/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package gopwd

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// StrToPwd 使用 argon2id 算法生成密码的散列值
func StrToPwd(password string) (string, error) {
	salt, err := GenerateSalt(16)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)

	return fmt.Sprintf("%s$%s", encodedSalt, encodedHash), nil
}

// GenerateSalt 生成指定长度的盐值
func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// IsCorrectPassword 用于比较输入的密码和存储的散列值是否匹配
func IsCorrectPassword(password, hashedPassword string) bool {
	parts := SplitHash(hashedPassword)
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// 使用 subtle.ConstantTimeCompare 比较散列值
	if subtle.ConstantTimeCompare(hash, expectedHash) == 1 {
		return true
	}

	return false
}

// SplitHash 分割存储的散列值和盐值
func SplitHash(hash string) []string {
	return strings.Split(hash, "$")
}
