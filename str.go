/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package goutils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/bizvip/go-utils/logs"
)

type StrUtils struct{}

func NewStrUtils() *StrUtils {
	return &StrUtils{}
}

// PadCnSpaceChar 使用中文空格为字符串填充
func (s *StrUtils) PadCnSpaceChar(label string, spaces int) string {
	for i := 0; i < spaces; i++ {
		label += string('\u3000')
	}
	return label
}

// UniqueStrings 返回一个新的切片，其中包含原切片中的唯一字符串
func (s *StrUtils) UniqueStrings(input []string) []string {
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

func (s *StrUtils) RegexpMatch(text string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logs.Logger().Error("Regex compile error:", err)
		return false
	}
	return regex.MatchString(text)
}

func (s *StrUtils) UUIDNoDash() string {
	u := uuid.New()
	uNoDashes := strings.Replace(u.String(), "-", "", -1)
	return uNoDashes
}

func (s *StrUtils) RandomId() string {
	randomData := make([]byte, 12)
	_, err := rand.Read(randomData)
	if err != nil {
		fmt.Println("Error generating random data:", err)
		return ""
	}

	encoded := base64.StdEncoding.EncodeToString(randomData)
	cleaned := strings.Map(func(r rune) rune {
		if ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') || ('0' <= r && r <= '9') {
			return r
		}
		return -1
	}, encoded)

	return cleaned
}

func (s *StrUtils) StrToInt64(intStr string) (int64, error) {
	if strings.Contains(intStr, ".") {
		return 0, fmt.Errorf("this method only accepts number without dots")
	}
	i, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse as int: %v", err)
	}
	return i, nil
}

func (s *StrUtils) Md5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func (s *StrUtils) FilterEmptyChar(str string) string {
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

func (s *StrUtils) NanoTimestampStr() string {
	now := time.Now()
	nano := fmt.Sprintf("%06d", now.Nanosecond()/1000)
	return now.Format("20060102150405") + "-" + nano
}

func (s *StrUtils) GetDirNameFromSnowflakeID(snowflakeID int64) string {
	transformedID := snowflakeID >> 10
	dirName := strconv.FormatInt(transformedID%10000, 10)
	return fmt.Sprintf("%04s", dirName)
}

func (s *StrUtils) IsAlphaNum(str string) bool {
	//	isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(str)
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func (s *StrUtils) Length(str string) int {
	return len([]rune(str))
}

func (s *StrUtils) ProtoMessageToJson(msg proto.Message) (string, error) {
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
func (s *StrUtils) GenStrBySeed(input, seed string) string {
	h := hmac.New(sha256.New, []byte(seed))
	h.Write([]byte(input))
	hash := h.Sum(nil)
	encoded := base64.StdEncoding.EncodeToString(hash)
	return encoded
}

func (s *StrUtils) GenSha1(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	r := h.Sum(nil)
	return fmt.Sprintf("%x", r)
}
