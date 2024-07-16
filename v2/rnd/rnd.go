/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package rnd

import (
	"math/rand"
	"time"
)

// RandNumStrNonSafe 生成一个指定长度的随机数字字符串
func RandNumStrNonSafe(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)

	// 创建一个新的随机数生成器实例
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range result {
		result[i] = digits[r.Intn(len(digits))]
	}
	return string(result)
}
