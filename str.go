package goutils

import (
	"crypto/md5"
	"crypto/rand"
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

type S struct{}

func Str() *S {
	return &S{}
}
func (s *S) PadCnSpaceChar(label string, spaces int) string {
	for i := 0; i < spaces; i++ {
		label += string('\u3000')
	}
	return label
}

// UniqueStrings 返回一个新的切片，其中包含原切片中的唯一字符串
func (s *S) UniqueStrings(input []string) []string {
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

func (s *S) RegexpMatch(text string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logs.Logger().Error("Regex compile error:", err)
		return false
	}
	return regex.MatchString(text)
}

func (s *S) UUIDNoDash() string {
	u := uuid.New()
	uNoDashes := strings.Replace(u.String(), "-", "", -1)
	return uNoDashes
}

func (s *S) RandomId() string {
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

func (s *S) StrToInt64(intStr string) (int64, error) {
	if strings.Contains(intStr, ".") {
		return 0, fmt.Errorf("this method only accepts number without dots")
	}
	i, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse as int: %v", err)
	}
	return i, nil
}

func (s *S) Md5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func (s *S) FilterEmptyChar(str string) string {
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

func (s *S) NanoTimestampStr() string {
	now := time.Now()
	nano := fmt.Sprintf("%06d", now.Nanosecond()/1000)
	return now.Format("20060102150405") + "-" + nano
}

func (s *S) GetDirNameFromSnowflakeID(snowflakeID int64) string {
	transformedID := snowflakeID >> 10
	dirName := strconv.FormatInt(transformedID%10000, 10)
	return fmt.Sprintf("%04s", dirName)
}

func (s *S) IsAlphaNum(str string) bool {
	//	isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(str)
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func (s *S) Length(str string) int {
	return len([]rune(str))
}

func (s *S) ProtoMessageToJson(msg proto.Message) (string, error) {
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
