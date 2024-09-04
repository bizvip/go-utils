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

const layout = "2006-01-02 15:04:05"

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

// GetNanoTimestampStr 当前时间的纳秒时间戳字符串
func GetNanoTimestampStr() string {
	now := time.Now()
	nano := fmt.Sprintf("%06d", now.Nanosecond()/1000)
	return now.Format("20060102150405") + "-" + nano
}

// GetMicroTimestampStr 当前时间的微秒时间戳字符串
func GetMicroTimestampStr() string {
	now := time.Now()
	micro := fmt.Sprintf("%06d", now.Nanosecond()/1000)   // 纳秒转微秒
	return now.Format("20060102150405") + "-" + micro[:3] // 取前3位作为微秒
}

// GetMilliTimestampStr 当前时间的毫秒时间戳字符串
func GetMilliTimestampStr() string {
	now := time.Now()
	milli := fmt.Sprintf("%03d", now.Nanosecond()/1e6) // 纳秒转毫秒
	return now.Format("20060102150405") + "-" + milli
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

// SetTimezone 设置时区，默认上海
func SetTimezone(tz ...string) {
	defaultTimezone := "Asia/Shanghai"

	var timezone string
	if len(tz) > 0 {
		timezone = tz[0]
	} else {
		timezone = defaultTimezone
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		panic("时区设置失败：" + err.Error())
	}
	time.Local = location
}

// CompareTimeStrings 比较两个时间字符串，返回 -1, 0, 1 分别表示第一个时间小于、等于、大于第二个时间
func CompareTimeStrings(t1, t2, layout string) (int, error) {
	time1, err := time.Parse(layout, t1)
	if err != nil {
		return 0, fmt.Errorf("解析时间字符串失败: %v", err)
	}

	time2, err := time.Parse(layout, t2)
	if err != nil {
		return 0, fmt.Errorf("解析时间字符串失败: %v", err)
	}

	// 使用 Before 和 After 方法直接比较
	if time1.Before(time2) {
		return -1, nil
	}
	if time1.After(time2) {
		return 1, nil
	}
	return 0, nil
}

// TimeDifference 计算两个时间字符串之间的时间间隔
func TimeDifference(t1, t2 string) (time.Duration, error) {
	time1, err := time.Parse(layout, t1)
	if err != nil {
		return 0, fmt.Errorf("解析时间字符串失败: %v", err)
	}

	time2, err := time.Parse(layout, t2)
	if err != nil {
		return 0, fmt.Errorf("解析时间字符串失败: %v", err)
	}

	duration := time2.Sub(time1)
	return duration, nil
}
