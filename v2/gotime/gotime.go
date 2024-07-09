/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package gotime

import (
	"fmt"
	"time"
)

// GetTimezoneOffsetByTimestamp 计算客户端时区偏移量和时区名称
func GetTimezoneOffsetByTimestamp(clientTimestampMillis int64) (string, error) {
	serverTime := time.Now()
	clientTime := time.Unix(0, clientTimestampMillis*int64(time.Millisecond))

	// 计算时间差（以小时为单位）
	timeDifference := clientTime.Sub(serverTime)
	hoursDifference := int(timeDifference.Hours())

	// 计算客户端时区偏移量（以小时为单位）
	clientTimezoneOffset := hoursDifference % 24

	// 构建时区字符串
	var timezoneName string
	if clientTimezoneOffset >= 0 {
		timezoneName = fmt.Sprintf("UTC+%d", clientTimezoneOffset)
	} else {
		timezoneName = fmt.Sprintf("UTC%d", clientTimezoneOffset)
	}

	return timezoneName, nil
}
