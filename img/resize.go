/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package img

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
)

// ResizeOptions 定义了 Resize 方法的参数
type ResizeOptions struct {
	Src      string
	SavePath string
	Width    int
	Height   int
	Quality  int
}

// Resize 根据提供的 ResizeOptions 调整图像大小
func (i *ImageUtils) Resize(ro ResizeOptions) error {
	if ro.Src == "" {
		return fmt.Errorf("要resize的图片源地址必须提供")
	}
	if ro.SavePath == "" {
		ro.SavePath = ro.Src
	}
	if ro.Width == 0 || ro.Height == 0 {
		return fmt.Errorf("必须提供要缩放的新图片宽高大小数字")
	}
	if ro.Quality == 0 {
		ro.Quality = 75
	}

	// 打开源图像文件
	file, err := os.Open(ro.Src)
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
	dstImage := image.NewRGBA(image.Rect(0, 0, ro.Width, ro.Height))

	draw.CatmullRom.Scale(dstImage, dstImage.Bounds(), img, img.Bounds(), draw.Over, nil)

	outputFile, err := os.Create(ro.SavePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	switch format {
	case "jpeg":
		err = jpeg.Encode(outputFile, dstImage, &jpeg.Options{Quality: ro.Quality})
	case "png":
		err = png.Encode(outputFile, dstImage)
	case "gif":
		err = gif.Encode(outputFile, dstImage, nil)
	case "bmp":
		err = bmp.Encode(outputFile, dstImage)
	case "webp":
		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(ro.Quality))
		if err != nil {
			return fmt.Errorf("创建 WebP 编码选项时出错: %v", err)
		}
		err = webp.Encode(outputFile, dstImage, options)
	default:
		err = jpeg.Encode(outputFile, dstImage, &jpeg.Options{Quality: ro.Quality})
	}

	if err != nil {
		return err
	}

	return nil
}
