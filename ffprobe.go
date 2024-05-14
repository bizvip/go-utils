package goutils

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/bizvip/go-utils/logs"
)

type FFProbeService struct{}

func NewFFProbeService() *FFProbeService {
	return &FFProbeService{}
}

func (f *FFProbeService) ffProbe(cmdStr, filePath string) string {
	quotedPath := fmt.Sprintf("'%s'", filePath)
	cmd := "ffprobe " + fmt.Sprintf(cmdStr, quotedPath)

	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logs.Logger().Error("error : ", err)
		return ""
	}

	return strings.TrimSpace(string(out))
}

func (f *FFProbeService) IsH264(filePath string) bool {
	cmdStr := "-v error -select_streams v:0 -show_entries stream=codec_name -of default=noprint_wrappers=1:nokey=1 %s"
	output := f.ffProbe(cmdStr, filePath)
	return strings.TrimSpace(output) == "h264"
}

func (f *FFProbeService) GetVideoDetails(filePath string) string {
	// todo :
	return f.ffProbe("", filePath)
}

func (f *FFProbeService) GetFormatName(filePath string) string {
	str := "-v error -show_entries format=format_name -of default=noprint_wrappers=1:nokey=1 %s"
	names := f.ffProbe(str, filePath)
	return strings.TrimSpace(names)
}

func (f *FFProbeService) GetCodecNames(filePath string) []string {
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

func (f *FFProbeService) IsValidExt(filePath, ext string) bool {
	formatNames := f.GetFormatName(filePath)
	return strings.Contains(formatNames, ext)
}

func (f *FFProbeService) GetBitRates(filePath string) int {
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

func (f *FFProbeService) GetResolution(filePath string) string {
	return f.ffProbe(
		"-v error -select_streams v:0 -show_entries stream=width,height -of csv=p=0 %s",
		filePath,
	)
}

func (f *FFProbeService) GetDuration(filePath string) int64 {
	durationStr := f.ffProbe(
		"-v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 %s",
		filePath,
	)

	duration, err := Str().StrToInt64(durationStr)
	if err != nil {
		logs.Logger().Println(durationStr, " >> 转换播放时长数字类型失败:", err)
		return 0
	}
	return duration
}

// func main() {
// 	ffprobe := NewFFProbeService()
// 	fmt.Println(ffprobe.GetFormatName("your_file_path_here"))
// }
