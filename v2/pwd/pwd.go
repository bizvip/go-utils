/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package pwd

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/argon2"
)

// 定义常见错误消息
var (
	ErrInvalidSecPwdLength      = errors.New("password must be 6 digits")
	ErrInvalidSecPwdConsecutive = errors.New("password must not contain three consecutive identical digits")
	ErrInvalidSecPwdSequential  = errors.New("password must not contain three or more sequential digits")
)

// NumberPwdValidate 验证指定长度的纯数字字符串密码
func NumberPwdValidate(secPwd string, length int) error {
	// 构建正则表达式以匹配指定长度的数字字符串
	regexPattern := fmt.Sprintf(`^[0-9]{%d}$`, length)
	matched, err := regexp.MatchString(regexPattern, secPwd)
	if err != nil || !matched {
		return ErrInvalidSecPwdLength
	}

	// 检查是否包含三个或更多连续相同的数字
	for i := 0; i < len(secPwd)-2; i++ {
		if secPwd[i] == secPwd[i+1] && secPwd[i] == secPwd[i+2] {
			return ErrInvalidSecPwdConsecutive
		}
	}

	// 检查是否为三个或更多递增或递减的顺子
	for i := 0; i < len(secPwd)-2; i++ {
		if secPwd[i+1] == secPwd[i]+1 && secPwd[i+2] == secPwd[i]+2 {
			return ErrInvalidSecPwdSequential
		}
		if secPwd[i+1] == secPwd[i]-1 && secPwd[i+2] == secPwd[i]-2 {
			return ErrInvalidSecPwdSequential
		}
	}

	return nil
}

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
