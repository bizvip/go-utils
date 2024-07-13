/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.
 * Author ORCID: https://orcid.org/0009-0003-8150-367X
 ******************************************************************************/

package str

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

// ToPrintJSON 将结构体转换为适合打印的格式化的 JSON 字符串
func ToPrintJSON(v interface{}) (string, error) {
	// 使用 json.MarshalIndent 进行格式化 JSON 序列化
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// ToUint32 字符串转换成 uint32
func ToUint32(str string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(str))
	return h.Sum32()
}

// ToInt64 字符串数字转换成 int64
func ToInt64(str string) (int64, error) {
	// 使用 strconv.ParseInt 将字符串转换为 int64 类型
	number, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number format: %v", err)
	}
	return number, nil
}

// PadCnSpaceChar 使用中文空格为字符串填充
func PadCnSpaceChar(label string, spaces int) string {
	for i := 0; i < spaces; i++ {
		label += string('\u3000')
	}
	return label
}

// UniqueStrings 返回一个新的切片，其中包含原切片中的唯一字符串
func UniqueStrings(input []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// RegexpMatch 使用正则表达式匹配字符串
func RegexpMatch(text string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return regex.MatchString(text)
}

// UUIDNoDash 生成不带连字符的 UUID
func UUIDNoDash() string {
	u := uuid.New()
	uNoDashes := strings.Replace(u.String(), "-", "", -1)
	return uNoDashes
}

// RandNumStr 生成一个指定长度的随机数字字符串
func RandNumStr(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			panic(err) // rand.Int 应该不会失败，但如果失败，直接 panic
		}
		result[i] = digits[num.Int64()]
	}
	return string(result)
}

// RandomId 生成一个随机 ID
func RandomId() string {
	randomData := make([]byte, 12)
	_, err := rand.Read(randomData)
	if err != nil {
		fmt.Println("Error generating random data:", err)
		return ""
	}

	encoded := base64.StdEncoding.EncodeToString(randomData)
	cleaned := strings.Map(
		func(r rune) rune {
			if ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') || ('0' <= r && r <= '9') {
				return r
			}
			return -1
		}, encoded,
	)

	return cleaned
}

// Md5 计算字符串的 MD5
func Md5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

// FilterEmptyChar 过滤空字符串
func FilterEmptyChar(str string) string {
	newStr := strings.ReplaceAll(strings.TrimSpace(str), "&nbsp;", "")
	newStr = strings.ReplaceAll(newStr, " ", "")
	newStr = strings.Map(
		func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, newStr,
	)
	newStr = strings.ReplaceAll(newStr, ":", "")
	newStr = strings.ReplaceAll(newStr, "：", "")
	return newStr
}

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
