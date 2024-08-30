/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package id

import (
	"fmt"
	"strings"
)

// ToAlpha 转换数字ID为默认大写字母表示，并补足长度
func ToAlpha(id int, length int, customMap map[int]string) string {
	idStr := fmt.Sprintf("%d", id)
	var result strings.Builder

	// 使用提供的自定义映射，如果没有提供则使用默认映射
	numToCharMap := customMap
	if numToCharMap == nil {
		numToCharMap = map[int]string{
			0: "C",
			1: "B",
			2: "A",
			3: "D",
			4: "F",
			5: "E",
			6: "G",
			7: "J",
			8: "I",
			9: "H",
		}
	}

	// 将每个数字映射到对应的大写字母
	for _, digit := range idStr {
		num := int(digit - '0')
		if char, exists := numToCharMap[num]; exists {
			result.WriteString(char)
		}
	}

	// 补足到指定的最小长度
	for result.Len() < length {
		result.WriteString("X")
	}

	return result.String()
}
