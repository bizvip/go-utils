/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package base62

import (
	"fmt"
	"math/big"
	"strings"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// SHA256ToBase62 将SHA256转换为Base62编码
func SHA256ToBase62(sha256Hash string) (string, error) {
	// 1. 验证输入
	if len(sha256Hash) != 64 {
		return "", fmt.Errorf("invalid SHA256 hash length: %d", len(sha256Hash))
	}

	// 2. 将16进制转换为big.Int
	num := new(big.Int)
	num.SetString(sha256Hash, 16)

	// 3. 转换为Base62
	var encoded strings.Builder
	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	// 当num不为0时继续转换
	for num.Cmp(zero) > 0 {
		num.DivMod(num, base, mod)
		encoded.WriteByte(base62Chars[mod.Int64()])
	}

	// 4. 反转字符串（因为我们是从低位到高位计算的）
	result := reverseString(encoded.String())

	// 5. 补齐到43位（如果需要）
	if len(result) < 43 {
		result = strings.Repeat("0", 43-len(result)) + result
	}

	return result, nil
}

// Base62ToSHA256 将Base62编码转回SHA256（用于验证）
func Base62ToSHA256(base62Str string) (string, error) {
	if len(base62Str) > 43 {
		return "", fmt.Errorf("invalid Base62 string length: %d", len(base62Str))
	}

	// 1. 转换为big.Int
	num := new(big.Int)
	base := big.NewInt(62)
	power := new(big.Int)

	for i := 0; i < len(base62Str); i++ {
		ch := base62Str[i]
		val := strings.IndexByte(base62Chars, ch)
		if val == -1 {
			return "", fmt.Errorf("invalid character in Base62 string: %c", ch)
		}

		power.Exp(base, big.NewInt(int64(len(base62Str)-i-1)), nil)
		power.Mul(power, big.NewInt(int64(val)))
		num.Add(num, power)
	}

	// 2. 转换为16进制字符串
	hex := num.Text(16)

	// 3. 补齐到64位
	if len(hex) < 64 {
		hex = strings.Repeat("0", 64-len(hex)) + hex
	}

	return hex, nil
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
