package validator

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// ValidationError 验证错误类型
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("field %s: %s", e.Field, e.Message)
	}
	return e.Message
}

// ValidationRule 验证规则接口
type ValidationRule[T any] interface {
	Validate(value T) error
}

// Validator 泛型验证器
type Validator[T any] struct {
	rules []ValidationRule[T]
}

// NewValidator 创建新的验证器
func NewValidator[T any](rules ...ValidationRule[T]) *Validator[T] {
	return &Validator[T]{rules: rules}
}

// AddRule 添加验证规则
func (v *Validator[T]) AddRule(rule ValidationRule[T]) {
	v.rules = append(v.rules, rule)
}

// Validate 执行验证
func (v *Validator[T]) Validate(value T) error {
	for _, rule := range v.rules {
		if err := rule.Validate(value); err != nil {
			return err
		}
	}
	return nil
}

// 字符串验证规则

// StringLengthRule 字符串长度验证
type StringLengthRule struct {
	Min, Max int
	Field    string
}

func (r StringLengthRule) Validate(value string) error {
	length := len([]rune(value))
	if length < r.Min {
		return ValidationError{r.Field, fmt.Sprintf("length must be at least %d characters", r.Min)}
	}
	if r.Max > 0 && length > r.Max {
		return ValidationError{r.Field, fmt.Sprintf("length must not exceed %d characters", r.Max)}
	}
	return nil
}

// RegexRule 正则表达式验证
type RegexRule struct {
	Pattern *regexp.Regexp
	Message string
	Field   string
}

func NewRegexRule(pattern, message, field string) *RegexRule {
	return &RegexRule{
		Pattern: regexp.MustCompile(pattern),
		Message: message,
		Field:   field,
	}
}

func (r RegexRule) Validate(value string) error {
	if !r.Pattern.MatchString(value) {
		return ValidationError{r.Field, r.Message}
	}
	return nil
}

// RequiredRule 必填验证
type RequiredRule[T comparable] struct {
	Field string
}

func (r RequiredRule[T]) Validate(value T) error {
	var zero T
	if value == zero {
		return ValidationError{r.Field, "is required"}
	}
	return nil
}

// InRule 枚举值验证
type InRule[T comparable] struct {
	ValidValues []T
	Field       string
}

func (r InRule[T]) Validate(value T) error {
	if !slices.Contains(r.ValidValues, value) {
		return ValidationError{r.Field, fmt.Sprintf("must be one of: %v", r.ValidValues)}
	}
	return nil
}

// RangeRule 数值范围验证
type RangeRule[T ~int | ~int32 | ~int64 | ~float32 | ~float64] struct {
	Min, Max T
	Field    string
}

func (r RangeRule[T]) Validate(value T) error {
	if value < r.Min || value > r.Max {
		return ValidationError{r.Field, fmt.Sprintf("must be between %v and %v", r.Min, r.Max)}
	}
	return nil
}

// 常用验证器

// ValidateEmail 验证邮箱格式
func ValidateEmail(email, field string) error {
	rule := NewRegexRule(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		"invalid email format",
		field,
	)
	return rule.Validate(email)
}

// ValidatePhone 验证手机号（中国）
func ValidatePhone(phone, field string) error {
	rule := NewRegexRule(
		`^1[3-9]\d{9}$`,
		"invalid phone number format",
		field,
	)
	return rule.Validate(phone)
}

// ValidateIDCard 验证身份证号（中国）
func ValidateIDCard(idCard, field string) error {
	// 18位身份证验证
	if len(idCard) != 18 {
		return ValidationError{field, "ID card must be 18 digits"}
	}

	rule := NewRegexRule(
		`^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`,
		"invalid ID card format",
		field,
	)

	if err := rule.Validate(idCard); err != nil {
		return err
	}

	// 校验码验证
	return validateIDCardChecksum(idCard, field)
}

// validateIDCardChecksum 验证身份证校验码
func validateIDCardChecksum(idCard, field string) error {
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checksums := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

	sum := 0
	for i := 0; i < 17; i++ {
		digit, err := strconv.Atoi(string(idCard[i]))
		if err != nil {
			return ValidationError{field, "invalid ID card format"}
		}
		sum += digit * weights[i]
	}

	expectedChecksum := checksums[sum%11]
	actualChecksum := strings.ToUpper(string(idCard[17]))

	if expectedChecksum != actualChecksum {
		return ValidationError{field, "invalid ID card checksum"}
	}

	return nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password, field string) error {
	validators := NewValidator[string](
		StringLengthRule{Min: 8, Max: 64, Field: field},
		*NewRegexRule(`[a-z]`, "must contain at least one lowercase letter", field),
		*NewRegexRule(`[A-Z]`, "must contain at least one uppercase letter", field),
		*NewRegexRule(`\d`, "must contain at least one digit", field),
		*NewRegexRule(`[!@#$%^&*(),.?":{}|<>]`, "must contain at least one special character", field),
	)

	return validators.Validate(password)
}
