package str

import (
	"unicode"
)

// IsAlphaNum 检查字符串是否是字母或数字
func IsAlphaNum(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
