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

	// 固定服务器时间为 UTC+0
	serverTime := time.Now().UTC()

	// 客户端时间假设为北京时间
	clientTime := time.Unix(0, millis*int64(time.Millisecond)).In(time.FixedZone("UTC+8", 8*3600))

	// 计算时间差（以小时为单位）
	timeDifference := clientTime.Sub(serverTime)
	hoursDifference := int(timeDifference.Hours())

	// 构建时区字符串
	var timezoneName string
	if hoursDifference >= 0 {
		timezoneName = fmt.Sprintf("UTC+%d", hoursDifference)
	} else {
		timezoneName = fmt.Sprintf("UTC%d", hoursDifference)
	}

	return timezoneName, nil
}
