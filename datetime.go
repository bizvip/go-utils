package utils

import "time"

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
