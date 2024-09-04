/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

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
