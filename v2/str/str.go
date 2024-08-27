/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.
 * Author ORCID: https://orcid.org/0009-0003-8150-367X
 ******************************************************************************/

package str

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/goccy/go-json"
	"golang.org/x/crypto/sha3"
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

// Md5 计算字符串的 MD5
func Md5(str string) string {
	h := md5.Sum([]byte(str))
	return hex.EncodeToString(h[:])
}

// Sha256 计算并返回字符串的 SHA(2)/SHA3-256 哈希值
func Sha256(s string, isSha3 bool) (string, error) {
	var h hash.Hash
	if isSha3 {
		h = sha3.New256()
	} else {
		h = sha256.New()
	}
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", err
	}
	hashBytes := h.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
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
