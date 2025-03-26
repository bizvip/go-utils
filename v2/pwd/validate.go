package pwd

import (
	"errors"
	"fmt"
	"regexp"
)

// 定义常见错误消息
var (
	ErrInvalidSecPwdLength      = errors.New("password must be 6 digits")
	ErrInvalidSecPwdConsecutive = errors.New("password must not contain three consecutive identical digits")
	ErrInvalidSecPwdSequential  = errors.New("password must not contain three or more sequential digits")
	ErrCiphertextTooShort       = errors.New("ciphertext too short")
	ErrInvalidSHA256Length      = errors.New("SHA256 hash must be 64 characters long")
	ErrInvalidSHA256Format      = errors.New("SHA256 hash must be a valid hexadecimal string")
)

// ValidateSixNumberAsPwd 验证指定长度的纯数字字符串密码
func ValidateSixNumberAsPwd(secPwd string, length int) error {
	// 构建正则表达式以匹配指定长度的数字字符串
	regexPattern := fmt.Sprintf(`^[0-9]{%d}$`, length)
	matched, err := regexp.MatchString(regexPattern, secPwd)
	if err != nil || !matched {
		return ErrInvalidSecPwdLength
	}

	// 检查是否包含三个或更多连续相同的数字
	for i := 0; i < len(secPwd)-2; i++ {
		if secPwd[i] == secPwd[i+1] && secPwd[i] == secPwd[i+2] {
			return ErrInvalidSecPwdConsecutive
		}
	}

	// 检查是否为三个或更多递增或递减的顺子
	for i := 0; i < len(secPwd)-2; i++ {
		if secPwd[i+1] == secPwd[i]+1 && secPwd[i+2] == secPwd[i]+2 {
			return ErrInvalidSecPwdSequential
		}
		if secPwd[i+1] == secPwd[i]-1 && secPwd[i+2] == secPwd[i]-2 {
			return ErrInvalidSecPwdSequential
		}
	}

	return nil
}

// ValidateSHA256 验证字符串是否为 64 位长度的 SHA-256 哈希值
func ValidateSHA256(hash string) error {
	// 检查长度是否为 64
	if len(hash) != 64 {
		return ErrInvalidSHA256Length
	}

	// 使用正则表达式检查是否为十六进制字符
	matched, err := regexp.MatchString(`^[a-fA-F0-9]{64}$`, hash)
	if err != nil || !matched {
		return ErrInvalidSHA256Format
	}

	return nil
}
