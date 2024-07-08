/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package gostr

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// StringToNumber 字符串转换成uint32
func StringToNumber(str string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(str))
	return h.Sum32()
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

func RegexpMatch(text string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return regex.MatchString(text)
}

func UUIDNoDash() string {
	u := uuid.New()
	uNoDashes := strings.Replace(u.String(), "-", "", -1)
	return uNoDashes
}

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

func StrToInt64(intStr string) (int64, error) {
	if strings.Contains(intStr, ".") {
		return 0, fmt.Errorf("this method only accepts number without dots")
	}
	i, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse as int: %v", err)
	}
	return i, nil
}

func Md5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

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

func NanoTimestampStr() string {
	now := time.Now()
	nano := fmt.Sprintf("%06d", now.Nanosecond()/1000)
	return now.Format("20060102150405") + "-" + nano
}

func GetDirNameFromSnowflakeID(snowflakeID int64) string {
	transformedID := snowflakeID >> 10
	dirName := strconv.FormatInt(transformedID%10000, 10)
	return fmt.Sprintf("%04s", dirName)
}

func IsAlphaNum(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func Length(str string) int {
	return len([]rune(str))
}

func ProtoMessageToJson(msg proto.Message) (string, error) {
	marshaller := protojson.MarshalOptions{
		Multiline:     true,
		Indent:        "  ",
		UseProtoNames: true,
	}
	jsonBytes, err := marshaller.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// GenStrBySeed 根据给定的字符串和种子生成一个可重现的新字符串（不建议用到密码）
func GenStrBySeed(input, seed string) string {
	h := hmac.New(sha256.New, []byte(seed))
	h.Write([]byte(input))
	hash := h.Sum(nil)
	encoded := base64.StdEncoding.EncodeToString(hash)
	return encoded
}

func GenSha1(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	r := h.Sum(nil)
	return fmt.Sprintf("%x", r)
}

// Calculator 输入字符串数学表达式，将计算出结果
func Calculator(exp string) (string, error) {
	// 去除空格
	exp = strings.ReplaceAll(exp, " ", "")

	// 定义正则表达式匹配 + - * / ^ sqrt % mod
	re := regexp.MustCompile(`^([\d.]+)([\+\-\*/\^%]|sqrt|mod)([\d.]*)$`)
	matches := re.FindStringSubmatch(exp)

	if len(matches) < 3 {
		return "", fmt.Errorf("无效的表达式")
	}

	// 解析数字和操作符
	num1, err := decimal.NewFromString(matches[1])
	if err != nil {
		return "", fmt.Errorf("无效的数字: %s", matches[1])
	}

	var num2 decimal.Decimal
	if matches[3] != "" {
		num2, err = decimal.NewFromString(matches[3])
		if err != nil {
			return "", fmt.Errorf("无效的数字: %s", matches[3])
		}
	}

	operator := matches[2]

	// 计算结果
	var result decimal.Decimal
	switch operator {
	case "+":
		result = num1.Add(num2)
	case "-":
		result = num1.Sub(num2)
	case "*":
		result = num1.Mul(num2)
	case "/":
		if num2.IsZero() {
			return "", fmt.Errorf("除数不能为零")
		}
		result = num1.Div(num2)
	case "^":
		// decimal 没有直接的幂运算方法，需转换为 float64
		exp, _ := num2.Float64()
		result = num1.Pow(decimal.NewFromFloat(exp))
	case "sqrt":
		if num1.LessThan(decimal.Zero) {
			return "", fmt.Errorf("负数不能开根号")
		}
		// decimal 没有直接的开根号方法，需转换为 float64
		floatVal, _ := num1.Float64()
		result = decimal.NewFromFloat(math.Sqrt(floatVal))
	case "%":
		result = num1.Div(num2).Mul(decimal.NewFromInt(100))
		return result.String() + "%", nil
	case "mod":
		result = num1.Mod(num2)
	default:
		return "", fmt.Errorf("无效的操作符: %s", operator)
	}

	return result.String(), nil
}
