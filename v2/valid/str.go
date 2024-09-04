/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package valid

import (
	"errors"
	"regexp"
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
