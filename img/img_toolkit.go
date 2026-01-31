package img

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
)

// ImageInfo represents information about an image
type ImageInfo struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Format   string `json:"format"`
	Ext      string `json:"ext"`
	Size     int64  `json:"size"`
	FileName string `json:"filename"`
	FileMD5  string `json:"file_md5"`
}

// ResizeOptions defines parameters for the Resize function
type ResizeOptions struct {
	Src      string
	SavePath string
	Width    int
	Height   int
	Quality  int
}

// ImageToBase64 converts an image file to a base64 string
func ImageToBase64(imgPath string) (string, error) {
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

	// 将图片编码为 JPEG 格式
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return "", err
	}

	// 将编码后的图片字节序列转换为 Base64 字符串
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Str, nil
}

// Base64ToFile writes a base64 encoded image string to a file
func Base64ToFile(base64Str, fullPath string, overwrite bool) error {
	// 获取文件所在的目录路径
	dir := filepath.Dir(fullPath)

	// 确保目录存在，如果不存在则创建
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); err == nil {
		if !overwrite {
			return os.ErrExist
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	// 解码Base64字符串
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(fullPath, data, 0644)
}

// ResizeImage resizes an image according to the provided options
func ResizeImage(options ResizeOptions) error {
	if options.Src == "" {
		return fmt.Errorf("source image path must be provided")
	}
	if options.SavePath == "" {
		options.SavePath = options.Src
	}
	if options.Width == 0 || options.Height == 0 {
		return fmt.Errorf("both width and height must be non-zero")
	}
	if options.Quality == 0 {
		options.Quality = 75
	}

	// 打开源图像文件
	file, err := os.Open(options.Src)
	if err != nil {
		return err
	}
	defer file.Close()

	// 解码图像，自动检测格式
	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	// 创建一个新的目标图像
	dstImage := image.NewRGBA(image.Rect(0, 0, options.Width, options.Height))

	draw.CatmullRom.Scale(dstImage, dstImage.Bounds(), img, img.Bounds(), draw.Over, nil)

	outputFile, err := os.Create(options.SavePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	switch format {
	case "jpeg":
		err = jpeg.Encode(outputFile, dstImage, &jpeg.Options{Quality: options.Quality})
	case "png":
		err = png.Encode(outputFile, dstImage)
	case "gif":
		err = gif.Encode(outputFile, dstImage, nil)
	case "bmp":
		err = bmp.Encode(outputFile, dstImage)
	case "webp":
		webpOptions, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(options.Quality))
		if err != nil {
			return fmt.Errorf("error creating WebP encoder options: %v", err)
		}
		err = webp.Encode(outputFile, dstImage, webpOptions)
	default:
		err = jpeg.Encode(outputFile, dstImage, &jpeg.Options{Quality: options.Quality})
	}

	if err != nil {
		return err
	}

	return nil
}
