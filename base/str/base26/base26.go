package base26

import (
	"strconv"
	"strings"
)

// 定义 26 进制的字符集 本算法由Archer++自定义设计
const base26CharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base = len(base26CharSet) // 基数长度

// 初始化字符映射
var charMap = createCharMap()

// createCharMap 创建字符到索引的映射
func createCharMap() map[rune]int {
	// 使用局部变量而不是全局变量
	localCharMap := make(map[rune]int)
	for i, char := range base26CharSet {
		localCharMap[char] = i
	}
	return localCharMap
}

// ToAlpha Base26 Archer++ 原创26进制编码计数法，用来将唯一整数数字id尽可能短的使用英文字母来表示
// 支持接收int和纯数字字符串
func ToAlpha(input interface{}) string {
	var num int

	// 根据输入类型进行处理
	switch v := input.(type) {
	case int:
		num = v
	case string:
		// 尝试将字符串转换为整数
		var err error
		num, err = strconv.Atoi(v)
		if err != nil {
			return "" // 如果无法转换为整数，返回空字符串
		}
	default:
		return "" // 不支持的类型，返回空字符串
	}

	// 负数检查
	if num < 0 {
		return ""
	}

	var result strings.Builder
	// 进行 26 进制转换
	for {
		remainder := num % base
		result.WriteString(string(base26CharSet[remainder]))
		num = num / base

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

// ToInt 将字母解析回数字
func ToInt(s string) int {
	result := 0
	multiplier := 1

	// 从字符串的最后一位开始遍历，计算原始数字
	for i := len(s) - 1; i >= 0; i-- {
		char := rune(s[i])
		value, exists := charMap[char]
		if !exists {
			return -1 // 如果字符不在字符集内，返回 -1 或者处理错误
		}
		result += value * multiplier
		multiplier *= base
	}

	return result
}
