package validator

import (
	"errors"
	"regexp"
	"unicode"
)

// 校验是否 MD5 字符串的错误定义
var (
	ErrInvalidMD5Length  = errors.New("invalid length for MD5 hash")
	ErrInvalidMD5Pattern = errors.New("value does not match MD5 hash pattern")
)

// IsMd5 验证字符串是否是有效的 MD5 值
func IsMd5(input string) error {
	if len(input) != 32 {
		return ErrInvalidMD5Length
	}
	match, err := regexp.MatchString(`^[a-fA-F0-9]{32}$`, input)
	if err != nil {
		return err
	}
	if !match {
		return ErrInvalidMD5Pattern
	}
	return nil
}

// IsAlphaNum 检查字符串是否是字母或数字
func IsAlphaNum(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsLengthBetween 检查字符串长度是否在指定范围内（包括边界值）
func IsLengthBetween(str string, min, max int) bool {
	length := len(str)
	return length >= min && length <= max
}
