package img

import (
	"bytes"
	"encoding/base64"
	"image/jpeg"
	"os"
	"path/filepath"
)

func ToBase64(imgPath string) (string, error) {
	// 打开图片文件
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()

	// 解码图片文件
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		return "", err
	}

	// 将图片编码为 JPEG 格式（也可以选择其他格式，如 png.Encode）
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return "", err
	}

	// 将编码后的图片字节序列转换为 Base64 字符串
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Str, nil
}

// WriteBase64ToFile 将Base64编码的图片字符串写入文件
func WriteBase64ToFile(base64Str, fullPath string, overwrite bool) error {
	// 获取文件所在的目录路径
	dir := filepath.Dir(fullPath)

	// 确保目录存在，如果不存在则创建
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err // 创建目录失败
	}

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); err == nil {
		if !overwrite {
			return os.ErrExist // 文件已存在，且不允许覆盖
		}
	} else if !os.IsNotExist(err) {
		return err // os.Stat调用出错
	}

	// 解码Base64字符串
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return err // Base64解码失败
	}

	// 写入文件
	return os.WriteFile(fullPath, data, 0644)
}
