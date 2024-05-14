package utils

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func DownImage(url, name, savePath string) (string, error) {
	// 创建一个跳过证书验证的http.Client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// 创建一个http.Request对象
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建http请求失败: %w", err)
	}

	// 设置自定义的Referer和User-Agent头部
	req.Header.Set("Referer", url)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")

	// 使用自定义的client和请求对象发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("下载图片时候创建http请求失败：%w", err)
	}
	defer resp.Body.Close()

	// 从 URL 获取文件扩展名
	ext := filepath.Ext(strings.Split(url, "?")[0])
	if ext == "" {
		ext = ".jpg" // 默认扩展名
	}

	// 如果没有提供名称，则从URL中提取文件名
	if name == "" {
		urlParts := strings.Split(strings.Split(url, "?")[0], "/")
		name = urlParts[len(urlParts)-1] // 获取URL的最后一部分作为文件名
		// 移除扩展名
		name = strings.TrimSuffix(name, ext)
	}

	// 确保保存路径以斜杠结尾
	if !strings.HasSuffix(savePath, "/") {
		savePath += "/"
	}

	// 检查并创建目录
	err = os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		assErr := fmt.Errorf("下载图片时候创建本地存储目录失败：%w", err)
		return "", assErr
	}

	// 创建文件路径
	replacer := strings.NewReplacer(" ", "-", "/", "-", "\\", "|")
	filename := savePath + replacer.Replace(name) + ext

	// 将多个-改成一个
	re := regexp.MustCompile("-+")
	filename = re.ReplaceAllString(filename, "-")

	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		assErr := fmt.Errorf("下载图片时候创建本地文件失败：%w", err)
		return "", assErr
	}
	defer file.Close()

	// 将响应体写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		assErr := fmt.Errorf("下载图片时候写入图片到本地失败：%w", err)
		return "", assErr
	}

	return filename, nil
}
