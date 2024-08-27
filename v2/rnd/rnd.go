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

// 全局随机数生成器
var rng *rand.Rand

// init 程序启动期间，只设置一次种子
func init() {
	source := rand.NewSource(uint64(time.Now().UnixNano()))
	rng = rand.New(source)
}

// RandNumStrNonSafe 生成一个指定长度的随机数字字符串
func RandNumStrNonSafe(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)

	for i := range result {
		result[i] = digits[rng.Intn(len(digits))]
	}
	return string(result)
}

// UUIDNoDash 生成不带连字符的 UUID
func UUIDNoDash() string {
	u := uuid.New()
	uNoDashes := strings.Replace(u.String(), "-", "", -1)
	return uNoDashes
}

// RandomId 生成一个随机字符串ID
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

// RandomNumberInRange 生成指定范围内的随机数字
func RandomNumberInRange(min, max int) int {
	if min > max {
		panic("min should be less than or equal to max")
	}

	rand.Seed(uint64(time.Now().UnixNano()))

	// 生成范围内的随机数
	return rand.Intn(max-min+1) + min
}
