/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package num

import (
	"errors"
	"regexp"
)

// 定义常见错误消息
var (
	ErrInvalidSecPwdLength      = errors.New("security password must be 6 digits")
	ErrInvalidSecPwdConsecutive = errors.New("security password must not contain three consecutive numbers")
)

// ValidateSecPwd 验证交易密码
func ValidateSecPwd(secPwd string) error {
	matched, err := regexp.MatchString(`^[0-9]{6}$`, secPwd)
	if err != nil || !matched {
		return ErrInvalidSecPwdLength
	}
	// 禁止三个连续数字
	for i := 0; i < len(secPwd)-2; i++ {
		if secPwd[i+1] == secPwd[i]+1 && secPwd[i+2] == secPwd[i]+2 {
			return ErrInvalidSecPwdConsecutive
		}
	}
	return nil
}
