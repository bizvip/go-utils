/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package rnd

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

// RandNumStrNonSafe 生成一个指定长度的随机数字字符串
func RandNumStrNonSafe(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)

	// 创建一个新的随机数生成器实例
	r := rand.New(rand.NewSource(uint64(time.Now().UnixMilli())))

	for i := range result {
		result[i] = digits[r.Intn(len(digits))]
	}
	return string(result)
}

// UUIDNoDash 生成不带连字符的 UUID
func UUIDNoDash() string {
	u := uuid.New()
	uNoDashes := strings.Replace(u.String(), "-", "", -1)
	return uNoDashes
}

// RandomId 生成一个随机 ID
func RandomId() string {
	randomData := make([]byte, 12)
	_, err := rand.Read(randomData)
	if err != nil {
		fmt.Println("Error generating random data:", err)
		return ""
	}

	encoded := base64.StdEncoding.EncodeToString(randomData)
	cleaned := strings.Map(
		func(r rune) rune {
			if ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') || ('0' <= r && r <= '9') {
				return r
			}
			return -1
		}, encoded,
	)

	return cleaned
}
