/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package dt

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// GetTimezoneOffsetByMillis 计算客户端时区偏移量和时区名称
func GetTimezoneOffsetByMillis(millis int64) (string, error) {
	// 将传入的毫秒时间戳转换为 time.Time 对象
	timestamp := time.Unix(0, millis*int64(time.Millisecond))

	// 检查时间范围（时间戳应在 Unix 纪元之后且不超过当前时间+24小时）
	if timestamp.Before(time.Unix(0, 0)) || timestamp.After(time.Now().Add(24*time.Hour)) {
		return "", fmt.Errorf("timestamp is out of valid range")
	}

	// 获取当前服务器时区
	serverLocation := time.Now().Location()
	serverTime := time.Now().In(serverLocation)

	// 客户端时间
	clientTime := timestamp

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

// AdjustMilliTimestamp 根据 13 位数字毫秒时间戳和加减秒数，返回计算后的时间戳
func AdjustMilliTimestamp(timestamp uint64, seconds int64) uint64 {
	// 将毫秒时间戳转换为 time.Time
	t := time.Unix(0, int64(timestamp)*int64(time.Millisecond))

	// 加减指定的秒数
	t = t.Add(time.Duration(seconds) * time.Second)

	// 返回计算后的毫秒时间戳
	return uint64(t.UnixNano() / int64(time.Millisecond))
}

var (
	ErrInvalidDurationFormat = fmt.Errorf("invalid duration format")
	ErrInvalidNumber         = fmt.Errorf("invalid number")
	ErrInvalidTimeUnit       = fmt.Errorf("invalid time unit")
)

// AdjustMilliTimestampByStr 根据传入的时间单位（如 "1d", "-1m","1y"）对当前时间戳进行加减，并返回结果毫秒时间戳
func AdjustMilliTimestampByStr(timestamp uint64, shift string) (uint64, error) {
	re := regexp.MustCompile(`^(-?\d+)([dmy])$`)
	matches := re.FindStringSubmatch(shift)
	if len(matches) != 3 {
		return 0, ErrInvalidDurationFormat
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, ErrInvalidNumber
	}

	current := time.UnixMilli(int64(timestamp))
	var newTime time.Time

	switch matches[2] {
	case "d":
		newTime = current.AddDate(0, 0, value)
	case "m":
		newTime = current.AddDate(0, value, 0)
	case "y":
		newTime = current.AddDate(value, 0, 0)
	default:
		return 0, ErrInvalidTimeUnit
	}

	return uint64(newTime.UnixMilli()), nil
}

// ConvertStrMillisToTime 将字符串格式的毫秒时间戳转换为 time.Time
func ConvertStrMillisToTime(millis string) (time.Time, error) {
	ms, err := strconv.ParseInt(millis, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid millisecond timestamp: %w", err)
	}
	t := time.UnixMilli(ms)
	return t, nil
}
