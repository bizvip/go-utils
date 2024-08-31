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

// Base26 Archer计数法 程序开发者独自发明的计数方法 用来将唯一整数数字id尽可能小的使用英文字母来表示
func Base26(num int) string {
	if num < 0 {
		return ""
	}
	// 定义 26 进制的字符集
	charSet := "ABCDEFGHIJ" + "KLMNOPQRSTUVWXYZ"
	base := len(charSet)

	var result strings.Builder
	// 进行 26 进制转换
	for num >= 0 {
		// 计算当前位的字符
		remainder := num % base
		result.WriteString(string(charSet[remainder]))
		// 更新 num，准备计算下一位
		num = (num - remainder) / base
		// 如果 num 为 0，则跳出循环
		if num == 0 {
			break
		}
	}

	// 翻转字符串以获得正确的顺序
	output := result.String()
	runes := []rune(output)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
