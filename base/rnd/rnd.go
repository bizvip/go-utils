package rnd

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

// 全局随机数生成器
var rng *rand.Rand

// init 程序启动期间，只设置一次种子
func init() {
	source := rand.NewSource(uint64(time.Now().UnixNano()))
	rng = rand.New(source)
}

// RandNumStrNonSafe 生成一个指定长度的随机数字字符串
func RandNumStrNonSafe(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)

	for i := range result {
		result[i] = digits[rng.Intn(len(digits))]
	}
	return string(result)
}

// UUID 可选无横线的UUID
func UUID(isNoDash bool) string {
	u := uuid.New()
	if isNoDash {
		return strings.Replace(u.String(), "-", "", -1)
	} else {
		return u.String()
	}
}

// GenRandomAlphaNumeric 生成一个只有大小写字母和数字的随机字符串
func GenRandomAlphaNumeric() string {
	randomData := make([]byte, 12)
	_, err := rng.Read(randomData)
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

// GenNumberInRange 生成指定范围内的随机数字
func GenNumberInRange(min, max int) int {
	if min > max {
		panic("min should be less than or equal to max")
	}
	return rng.Intn(max-min+1) + min
}
