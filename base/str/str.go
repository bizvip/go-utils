package str

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/goccy/go-json"
	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// 定义错误变量
var (
	ErrInvalidNumberFormat = errors.New("invalid number format, no dots allowed")
	ErrIntParseFailed      = errors.New("failed to parse as int")
	ErrFileOpenFailed      = errors.New("failed to open file")
	ErrHashCalculation     = errors.New("failed to calculate hash")
	ErrRegexCompilation    = errors.New("failed to compile regex pattern")
	ErrJsonMarshalFailed   = errors.New("failed to marshal to JSON")
)

// ToUint32 字符串转换成 uint32
func ToUint32(str string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(str))
	return h.Sum32()
}

// PadCnSpaceChar 使用中文空格为字符串填充
func PadCnSpaceChar(label string, spaces int) string {
	return label + strings.Repeat("\u3000", spaces)
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
func RegexpMatch(txt string, pattern string) (bool, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, ErrRegexCompilation
	}
	return regex.MatchString(txt), nil
}

// ToInt64 将字符串转换成 int64
func ToInt64(intStr string) (int64, error) {
	if strings.Contains(intStr, ".") {
		return 0, ErrInvalidNumberFormat
	}
	i, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0, ErrIntParseFailed
	}
	return i, nil
}

// CalcHash 计算字符串或文件的哈希（支持 MD5、SHA256 和 SHA3-256）
// - input: 要计算哈希的字符串或文件路径
// - useStream: 是否使用流的方式计算哈希
// - hashFunc: 哈希函数 (hash.Hash)
func CalcHash(input string, useStream bool, hashFunc hash.Hash) (string, error) {
	if useStream {
		file, err := os.Open(input) // 使用流的方式计算哈希
		if err != nil {
			return "", ErrFileOpenFailed
		}
		defer func(file *os.File) { _ = file.Close() }(file)

		if _, err := io.Copy(hashFunc, file); err != nil {
			return "", ErrHashCalculation
		}
	} else {
		_, err := hashFunc.Write([]byte(input)) // 使用一次性内存加载的方式计算哈希
		if err != nil {
			return "", ErrHashCalculation
		}
	}

	return hex.EncodeToString(hashFunc.Sum(nil)), nil
}

// Md5 计算字符串或文件的 MD5 哈希
func Md5(input string, useStream bool) (string, error) {
	return CalcHash(input, useStream, md5.New())
}

// Sha256 计算并返回字符串或文件的 SHA256/SHA3-256 哈希值
func Sha256(input string, useStream bool, isSha3 bool) (string, error) {
	var h hash.Hash
	if isSha3 {
		h = sha3.New256()
	} else {
		h = sha256.New()
	}
	return CalcHash(input, useStream, h)
}

// FilterEmptyChar 过滤空字符串
func FilterEmptyChar(str string) string {
	// 一次性移除空格、非打印字符、中文冒号和英文冒号
	replacer := strings.NewReplacer("&nbsp;", "", " ", "", ":", "", "：", "")
	newStr := replacer.Replace(strings.TrimSpace(str))
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, newStr)
}

// GetDirNameFromSnowflakeID 根据 Snowflake ID 生成目录名
func GetDirNameFromSnowflakeID(snowflakeID int64) string {
	transformedID := snowflakeID >> 10
	dirName := strconv.FormatInt(transformedID%10000, 10)
	return fmt.Sprintf("%04s", dirName)
}

// UnicodeLength 计算unicode字符串的字符长度
func UnicodeLength(str string) int {
	return len([]rune(str))
}

// ToPrettyJson 将数据结构转换为格式化的 JSON 字符串
func ToPrettyJson(v interface{}, isProto bool) (string, error) {
	if isProto {
		marshaller := protojson.MarshalOptions{
			Multiline:     true,
			Indent:        "  ",
			UseProtoNames: true,
		}
		jsonBytes, err := marshaller.Marshal(v.(proto.Message))
		if err != nil {
			return "", ErrJsonMarshalFailed
		}
		return string(jsonBytes), nil
	}

	// 普通结构体的 JSON 转换
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", ErrJsonMarshalFailed
	}
	return string(jsonData), nil
}

// GenFixedStrWithSeed 根据给定的字符串和种子生成一个可反复重现的哈希字符串（不适合密码用）
func GenFixedStrWithSeed(input, seed string) string {
	h := hmac.New(sha256.New, []byte(seed))
	h.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GenSha1 计算并返回字符串的 SHA1 哈希值
func GenSha1(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// GenSlug 生成slug
func GenSlug(title string) string {
	var b strings.Builder
	b.Grow(len(title))

	// 直接一趟循环完成：小写 + 过滤 + 合并 -
	prevDash := false
	for _, r := range strings.ToLowerSpecial(unicode.TurkishCase, title) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			// 统一把非 ASCII 字母转成 ASCII（可选，需 transliteration 库）
			if r > 127 {
				// 示例：这里简单跳过；生产可用 github.com/mozillazg/go-unidecode 等做音译
				continue
			}
			b.WriteRune(r)
			prevDash = false
		} else if !prevDash {
			b.WriteByte('-')
			prevDash = true
		}
	}

	slug := strings.Trim(b.String(), "-")

	// 再次截断并去尾/首 -
	if len(slug) > 100 {
		slug = strings.Trim(slug[:100], "-")
	}

	if slug == "" {
		slug = "n-a" // 或者生成 uuid/时间戳等
	}
	return slug
}
