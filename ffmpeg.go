/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package goutils

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/bizvip/go-utils/logs"
)

type FFMpegUtils struct{}

func NewFFMpegUtils() *FFMpegUtils {
	return &FFMpegUtils{}
}

func (f *FFMpegUtils) ffMpeg(shellCmd string) error {
	cmdStr := "ffmpeg " + shellCmd
	cmd := exec.Command("bash", "-c", cmdStr)

	// 获取标准错误输出
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logs.Logger().Error("获取标准错误失败 : ", err)
		return err
	}

	// 开始执行命令
	if err := cmd.Start(); err != nil {
		logs.Logger().Error("启动FFMpeg命令失败 : ", err)
		return err
	}

	errorBytes, _ := io.ReadAll(stderr)
	ffOutputStr := string(errorBytes)
	// if len(ffOutputStr) > 0 {
	// Logger().Infoln("FFMpeg输出 : ", ffOutputStr)
	// }

	if err := cmd.Wait(); err != nil {
		logs.Logger().Error("执行FFMpeg命令失败 : ", err, ffOutputStr)
		return err
	}

	return nil
}

func (f *FFMpegUtils) H264ToHls(filePath string, tsDir string, idxFilePath string, tsSeconds uint8) error {
	str := fmt.Sprintf(
		"-y -i '%s' -codec copy -map 0 -f segment -segment_list '%s' -segment_time %d ",
		filePath,
		idxFilePath,
		tsSeconds,
	)
	str += tsDir
	return f.ffMpeg(str)
}

func (f *FFMpegUtils) ToHlsNonH264(filePath string, tsDir string, idxFilePath string, tsSeconds uint8) error {
	cmd := fmt.Sprintf(
		"-y -i '%s' -c:v libx264 -c:a aac -map 0 -f segment -segment_list '%s' -segment_time %d ",
		filePath,
		idxFilePath,
		tsSeconds,
	)
	cmd += tsDir
	return f.ffMpeg(cmd)
}
