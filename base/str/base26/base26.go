package base26

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

// 定义 26 进制的字符集
const base26CharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base = len(base26CharSet) // 基数长度

// 编译时初始化字符映射
var charMap = map[rune]int{
	'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4, 'F': 5, 'G': 6, 'H': 7,
	'I': 8, 'J': 9, 'K': 10, 'L': 11, 'M': 12, 'N': 13, 'O': 14, 'P': 15,
	'Q': 16, 'R': 17, 'S': 18, 'T': 19, 'U': 20, 'V': 21, 'W': 22, 'X': 23,
	'Y': 24, 'Z': 25,
}

// uint64ToBase26 是一个内部辅助函数，将 uint64 转换为无符号的 base26 字符串。
func uint64ToBase26(num uint64) string {
	if num == 0 {
		return string(base26CharSet[0]) // 0 对应 "A"
	}

	var resultBuilder strings.Builder
	tempNum := num
	for {
		remainder := tempNum % uint64(base)
		resultBuilder.WriteString(string(base26CharSet[remainder]))
		tempNum = tempNum / uint64(base)
		if tempNum == 0 {
			break
		}
	}

	alphaNumStr := resultBuilder.String()
	runes := []rune(alphaNumStr)
	slices.Reverse(runes) // 使用更高效的 slices.Reverse
	return string(runes)
}

// Uint64ToAlpha 将 uint64 转换为 base26 字母表示。
func Uint64ToAlpha(input uint64) (string, error) {
	// uint64 本身不会是负数，直接转换
	return uint64ToBase26(input), nil
}

// Int64ToAlpha 将 int64 转换为 base26 字母表示。
// 负数会带负号，0 会转换为 "A"。
func Int64ToAlpha(input int64) (string, error) {
	if input == 0 {
		return uint64ToBase26(0), nil // "A"
	}

	isNegative := false
	var absVal uint64

	if input < 0 {
		isNegative = true
		if input == math.MinInt64 { // math.MinInt64 的绝对值是 2^63
			absVal = uint64(1) << 63
		} else {
			absVal = uint64(-input)
		}
	} else {
		absVal = uint64(input)
	}

	alphaStr := uint64ToBase26(absVal)

	if isNegative {
		return "-" + alphaStr, nil
	}
	return alphaStr, nil
}

// StrNumToAlpha 将字符串形式的十进制数字转换为 base26 字母表示。
// 会检查字符串是否为有效数字，负数会带负号。
func StrNumToAlpha(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input string is empty")
	}

	// 特殊处理字符串 "-0" 以便输出 "-A"
	if input == "-0" {
		return "-" + uint64ToBase26(0), nil
	}

	// 尝试按 int64 解析。这能正确处理正负号和大多数常见整数。
	valInt, errInt := strconv.ParseInt(input, 10, 64)
	if errInt == nil {
		// 解析为 int64 成功
		return Int64ToAlpha(valInt)
	}

	// 如果作为 int64 解析失败，检查是否是因为数字太大（但仍可能是有效的 uint64）
	// 只有当原始字符串不以"-"开头且错误是 strconv.ErrRange 时，才尝试 ParseUint
	if !strings.HasPrefix(input, "-") {
		numErr, ok := errInt.(*strconv.NumError)
		if ok && numErr.Err == strconv.ErrRange {
			valUint, errUint := strconv.ParseUint(input, 10, 64)
			if errUint == nil {
				// 解析为 uint64 成功 (通常是正数且大于 math.MaxInt64)
				return Uint64ToAlpha(valUint)
			}
			// 如果 ParseUint 也失败 (例如，超出 uint64 范围)
			return "", fmt.Errorf("input string '%s' is out of uint64 range: %w", input, errUint)
		}
	}

	// 对于其他所有解析错误 (例如，非数字字符)
	return "", fmt.Errorf("input string '%s' is not a valid number: %w", input, errInt)
}

// ToNum 将（可能带负号的）base26 字母表示的字符串解析回十进制数字字符串。
func ToNum(alphaStr string) (string, error) {
	if alphaStr == "" {
		return "", fmt.Errorf("input alpha string is empty")
	}

	isNegative := false
	sourceStr := alphaStr

	if strings.HasPrefix(sourceStr, "-") {
		isNegative = true
		sourceStr = sourceStr[1:]
		if sourceStr == "" { // 输入仅为 "-"
			return "", fmt.Errorf("invalid input alpha string: '-'")
		}
	}

	var valAbs uint64 = 0
	var multiplier uint64 = 1

	for i := len(sourceStr) - 1; i >= 0; i-- {
		char := rune(sourceStr[i])
		charValue, exists := charMap[char]
		if !exists {
			return "", fmt.Errorf("invalid character '%c' in alpha string '%s'", char, alphaStr)
		}

		term := uint64(charValue) * multiplier
		// 检查乘法溢出: 如果 multiplier 或 charValue 非零，但 term/multiplier != charValue
		if multiplier != 0 && uint64(charValue) != 0 && term/multiplier != uint64(charValue) {
			return "", fmt.Errorf("numeric overflow during conversion (term calculation) for alpha string '%s'", alphaStr)
		}
		// 检查加法溢出: 如果 valAbs 很大，加上 term 后会小于 valAbs
		if valAbs > math.MaxUint64-term {
			return "", fmt.Errorf("numeric overflow during conversion (sum calculation) for alpha string '%s'", alphaStr)
		}
		valAbs += term

		if i > 0 { // 准备下一个乘数，并检查其是否会溢出
			if multiplier > math.MaxUint64/uint64(base) {
				return "", fmt.Errorf("numeric overflow during conversion (multiplier calculation) for alpha string '%s'", alphaStr)
			}
			multiplier *= uint64(base)
		}
	}

	resultNumStr := strconv.FormatUint(valAbs, 10)

	if isNegative {
		if valAbs == 0 { // 输入为 "-A" (代表 -0), 标准输出为 "0"
			return "0", nil
		}
		return "-" + resultNumStr, nil
	}
	return resultNumStr, nil
}

// IsValidBase26 检查字符串是否是有效的 base26 格式
func IsValidBase26(s string) bool {
	if s == "" {
		return false
	}

	// 处理负号
	if strings.HasPrefix(s, "-") {
		s = s[1:]
		if s == "" { // 仅有负号
			return false
		}
	}

	// 检查每个字符是否都是 A-Z
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}

	return true
}
