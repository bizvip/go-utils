/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package goutils

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/bizvip/go-utils/logs"
)

type FFProbeUtils struct{}

func NewFFProbeUtils() *FFProbeUtils {
	return &FFProbeUtils{}
}

func (f *FFProbeUtils) ffProbe(cmdStr, filePath string) string {
	quotedPath := fmt.Sprintf("'%s'", filePath)
	cmd := "ffprobe " + fmt.Sprintf(cmdStr, quotedPath)

	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logs.Logger().Error("error : ", err)
		return ""
	}

	return strings.TrimSpace(string(out))
}

func (f *FFProbeUtils) IsH264(filePath string) bool {
	cmdStr := "-v error -select_streams v:0 -show_entries stream=codec_name -of default=noprint_wrappers=1:nokey=1 %s"
	output := f.ffProbe(cmdStr, filePath)
	return strings.TrimSpace(output) == "h264"
}

func (f *FFProbeUtils) GetVideoDetails(filePath string) string {
	// todo :
	return f.ffProbe("", filePath)
}

func (f *FFProbeUtils) GetFormatName(filePath string) string {
	str := "-v error -show_entries format=format_name -of default=noprint_wrappers=1:nokey=1 %s"
	names := f.ffProbe(str, filePath)
	return strings.TrimSpace(names)
}

func (f *FFProbeUtils) GetCodecNames(filePath string) []string {
	result := f.ffProbe(
		"-v error -show_entries stream=codec_name -of default=noprint_wrappers=1 %s",
		filePath,
	)
	codecNames := strings.Split(result, "\n")
	var names []string
	for _, name := range codecNames {
		name = strings.Replace(name, "codec_name=", "", -1)
		name = strings.TrimSpace(name)
		if name != "" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

func (f *FFProbeUtils) IsValidExt(filePath, ext string) bool {
	formatNames := f.GetFormatName(filePath)
	return strings.Contains(formatNames, ext)
}

func (f *FFProbeUtils) GetBitRates(filePath string) int {
	bitRateStr := f.ffProbe(
		"-v error -select_streams v:0 -show_entries stream=bit_rate -of default=noprint_wrappers=1:nokey=1 %s",
		filePath,
	)
	bitRate, err := strconv.Atoi(bitRateStr)
	if err != nil {
		logs.Logger().Println("Error converting bit rate:", err)
		return 0
	}
	return bitRate
}

func (f *FFProbeUtils) GetResolution(filePath string) string {
	return f.ffProbe(
		"-v error -select_streams v:0 -show_entries stream=width,height -of csv=p=0 %s",
		filePath,
	)
}

func (f *FFProbeUtils) GetDuration(filePath string) int64 {
	durationStr := f.ffProbe(
		"-v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 %s",
		filePath,
	)

	duration, err := NewStrUtils().StrToInt64(durationStr)
	if err != nil {
		logs.Logger().Println(durationStr, " >> 转换播放时长数字类型失败:", err)
		return 0
	}
	return duration
}

// func main() {
// 	ffprobe := NewGoogleTranslationUtils()
// 	fmt.Println(ffprobe.GetFormatName("your_file_path_here"))
// }
