/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package cryptocoin

import (
	cryptorand "crypto/rand"
	"encoding/hex"
)

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
