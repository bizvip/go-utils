package ua

import (
	"github.com/mileusna/useragent"
)

// UserAgent represents parsed user agent information
type UserAgent struct {
	Name      string
	Version   string
	OS        string
	OSVersion string
	Device    string
	Mobile    bool
	Tablet    bool
	Desktop   bool
	Bot       bool
}

// Parse parses a user agent string and returns structured information
func Parse(userAgentString string) UserAgent {
	ua := useragent.Parse(userAgentString)

	return UserAgent{
		Name:      ua.Name,
		Version:   ua.Version,
		OS:        ua.OS,
		OSVersion: ua.OSVersion,
		Device:    ua.Device,
		Mobile:    ua.Mobile,
		Tablet:    ua.Tablet,
		Desktop:   ua.Desktop,
		Bot:       ua.Bot,
	}
}
