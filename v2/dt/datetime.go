/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package dt

import (
	"fmt"
	"time"
)

// GetTimezoneOffsetByMillis 计算客户端时区偏移量和时区名称
func GetTimezoneOffsetByMillis(millis int64) (string, error) {
	// 先验证毫秒时间戳是否合法
	timestamp := time.Unix(0, millis*int64(time.Millisecond))
	// 检查时间范围（时间戳应在 Unix 纪元之后且不超过当前时间+24小时）
	if timestamp.Before(time.Unix(0, 0)) || timestamp.After(time.Now().Add(24*time.Hour)) {
		return "", fmt.Errorf("timestamp is out of valid range")
	}

	serverTime := time.Now().UTC() // 使用 UTC 时间进行计算
	clientTime := time.Unix(0, millis*int64(time.Millisecond)).UTC()

	// 计算时间差（以小时为单位）
	timeDifference := clientTime.Sub(serverTime)
	hoursDifference := int(timeDifference.Hours())

	// 计算客户端时区偏移量（以小时为单位）
	clientTimezoneOffset := (hoursDifference + 24) % 24

	// 构建时区字符串
	var timezoneName string
	if clientTimezoneOffset >= 0 {
		timezoneName = fmt.Sprintf("UTC+%d", clientTimezoneOffset)
	} else {
		timezoneName = fmt.Sprintf("UTC%d", clientTimezoneOffset)
	}

	return timezoneName, nil
}
