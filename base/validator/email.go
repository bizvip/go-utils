package validator

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

// IsValidDomain 全面检查邮箱域名是否有效
func IsValidDomain(domain string) bool {
	// 1. 检查域名格式
	if !isValidDomainFormat(domain) {
		return false
	}

	// 2. 检查DNS记录（优先考虑MX记录，其次是A或AAAA记录）
	hasMX := false
	hasAddr := false

	// 检查MX记录
	mxRecords, _ := net.LookupMX(domain)
	hasMX = len(mxRecords) > 0

	// 如果没有MX记录，检查A/AAAA记录
	if !hasMX {
		ips, _ := net.LookupIP(domain)
		hasAddr = len(ips) > 0
	}

	// 3. 对于邮箱域名，如果至少有一种记录存在，则认为域名有效
	return hasMX || hasAddr
}

// isValidDomainFormat 检查域名格式
func isValidDomainFormat(domain string) bool {
	// 域名长度限制
	if len(domain) > 253 {
		return false
	}

	// 域名格式验证
	const domainRegexPattern = `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`
	re := regexp.MustCompile(domainRegexPattern)
	return re.MatchString(domain)
}

// IsEmailAddrValidWithDomain 增强版：校验邮件地址格式和域名是否有效
func IsEmailAddrValidWithDomain(email string) error {
	if !IsValidEmailFormat(email) {
		return ErrInvalidEmailFormat
	}

	// 提取域名
	domain := email[strings.LastIndex(email, "@")+1:]

	// 使用增强的域名检查
	if !IsValidDomain(domain) {
		return ErrUnresolvableDomain
	}

	return nil
}
