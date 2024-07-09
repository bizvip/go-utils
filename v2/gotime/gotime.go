/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package gotime

import (
	"time"
)

// CalculateClientTimezoneOffset 计算客户端时区偏移量（以小时为单位）
func CalculateClientTimezoneOffset(clientTimestampMillis int64) (int, error) {
	serverTime := time.Now()
	clientTime := time.Unix(0, clientTimestampMillis*int64(time.Millisecond))
	// 计算时间差（以小时为单位）
	timeDifference := serverTime.Sub(clientTime)
	hoursDifference := int(timeDifference.Hours())
	// 计算客户端时区偏移量
	clientTimezoneOffset := hoursDifference % 24
	return clientTimezoneOffset, nil
}
