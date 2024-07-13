/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package valid

import (
	"errors"
	"net"
	"regexp"
	"strings"
)

// 定义错误变量
var (
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrUnresolvableDomain = errors.New("domain cannot be resolved")
)

// IsValidEmailFormat 校验邮件地址格式是否合法
func IsValidEmailFormat(email string) bool {
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegexPattern)
	return re.MatchString(email)
}

// IsDomainResolvable 校验域名是否可解析
func IsDomainResolvable(domain string) bool {
	_, err := net.LookupMX(domain)
	return err == nil
}

// EmailAddr 校验邮件地址格式和域名是否可解析
func EmailAddr(email string) error {
	if !IsValidEmailFormat(email) {
		return ErrInvalidEmailFormat
	}
	// 提取域名
	domain := email[strings.LastIndex(email, "@")+1:]
	if !IsDomainResolvable(domain) {
		return ErrUnresolvableDomain
	}
	return nil
}
