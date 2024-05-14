/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package goutils

import "time"

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
