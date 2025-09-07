package rnd

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/rs/zerolog/log"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

// RandNumStr 生成一个指定长度的随机数字字符串（加密安全）
func RandNumStr(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			log.Error().Err(err).Msg("Error generating random number")
			return ""
		}
		result[i] = digits[n.Int64()]
	}
	return string(result)
}

// RandNumStrNonSafe 生成一个指定长度的随机数字字符串（兼容性保留，使用加密安全版本）
func RandNumStrNonSafe(length int) string {
	return RandNumStr(length)
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

// GenRandomAlphaNumeric 生成一个只有大小写字母和数字的随机字符串（加密安全）
func GenRandomAlphaNumeric() string {
	randomData := make([]byte, 12)
	_, err := rand.Read(randomData)
	if err != nil {
		log.Error().Err(err).Msg("Error generating random data")
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

// GenNumberInRange 生成指定范围内的随机数字（加密安全）
func GenNumberInRange(min, max int) int {
	if min > max {
		panic("min should be less than or equal to max")
	}

	rangeSize := max - min + 1
	n, err := rand.Int(rand.Reader, big.NewInt(int64(rangeSize)))
	if err != nil {
		log.Error().Err(err).Msg("Error generating random number in range")
		return min
	}

	return int(n.Int64()) + min
}
