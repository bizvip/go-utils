package base26

import (
	"strings"
)

// 定义 26 进制的字符集
const base64CharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base = len(base64CharSet) // 基数长度

// 初始化字符映射
var charMap = createCharMap()

// createCharMap 创建字符到索引的映射
func createCharMap() map[rune]int {
	charMap := make(map[rune]int)
	for i, char := range base64CharSet {
		charMap[char] = i
	}
	return charMap
}

// Base26 Archer计数法 程序开发者独自发明的计数方法 用来将唯一整数数字id尽可能小的使用英文字母来表示
func Base26(num int) string {
	if num < 0 {
		return ""
	}

	var result strings.Builder
	// 进行 26 进制转换
	for num >= 0 {
		// 计算当前位的字符
		remainder := num % base
		result.WriteString(string(base64CharSet[remainder]))
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

// DeBase26 将字母解析回数字
func DeBase26(s string) int {
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
