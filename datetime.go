/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package goutils

import (
	"fmt"
	"time"
)

type TimeUtils struct{}

func NewTimeUtils() *TimeUtils { return &TimeUtils{} }

func (t *TimeUtils) SetTimezone(tz ...string) {
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
func (t *TimeUtils) CompareTimeStrings(t1, t2 string) (int, error) {
	const layout = "2006-01-02 15:04:05"

	time1, err := time.Parse(layout, t1)
	if err != nil {
		return 0, fmt.Errorf("解析时间字符串失败: %v", err)
	}

	time2, err := time.Parse(layout, t2)
	if err != nil {
		return 0, fmt.Errorf("解析时间字符串失败: %v", err)
	}

	if time1.Before(time2) {
		return -1, nil
	} else if time1.After(time2) {
		return 1, nil
	} else {
		return 0, nil
	}
}

// TimeDifference 计算两个时间字符串之间的时间间隔
func (t *TimeUtils) TimeDifference(t1, t2 string) (time.Duration, error) {
	const layout = "2006-01-02 15:04:05"

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
