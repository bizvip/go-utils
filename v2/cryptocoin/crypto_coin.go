/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package cryptocoin

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"strings"
)

// IsValidERC20Address erc20地址字符串合法性校验
func IsValidERC20Address(address string) bool {
	return len(address) == 42 && strings.HasPrefix(address, "0x")
}

// IsValidTRC20Address trc20地址字符串合法性校验
func IsValidTRC20Address(address string) bool {
	return len(address) == 34 && strings.HasPrefix(address, "T")
}

// IsValidWtcAddress wtc地址字符串合法性校验
func IsValidWtcAddress(address string) bool {
	return len(address) == 42 && strings.HasPrefix(address, "wtc")
}

// GenWBAddress 生成wtc币的收款地址
func GenWBAddress() (string, error) {
	addressLength := 39
	bytes := make([]byte, addressLength/2)
	_, err := cryptorand.Read(bytes)
	if err != nil {
		return "", err
	}
	address := hex.EncodeToString(bytes)
	return "wtc" + address, nil
}
