//go:build !libvips

package img

import (
	"image"
	"image/color"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/bizvip/go-utils/os/fs"
)

func GetImageInfo(path string, withMd5 bool) (*ImageInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	info := &ImageInfo{
		Width:    cfg.Width,
		Height:   cfg.Height,
		Format:   format,
		Ext:      strings.ToLower(strings.TrimPrefix(filepath.Ext(path), ".")),
		Size:     stat.Size(),
		FileName: filepath.Base(path),
		FrameNum: 1,
	}
	info.ColorSpace, info.HasAlpha, info.BitDepth = classifyColorModel(cfg.ColorModel)

	if format == "gif" {
		if _, err := f.Seek(0, 0); err == nil {
			if g, err := gif.DecodeAll(f); err == nil {
				info.FrameNum = len(g.Image)
				info.Animated = len(g.Image) > 1
				info.LoopCount = g.LoopCount
			}
		}
	}

	if withMd5 {
		md5sum, err := fs.GetBigFileMd5(path)
		if err != nil {
			return nil, err
		}
		info.FileMD5 = md5sum
	}

	return info, nil
}

// classifyColorModel 把 image.DecodeConfig 的 color model 归类为描述性
// ColorSpace/HasAlpha/BitDepth（见 ImageInfo 字段注释）。标准库的具名 model
// 均为 *modelFunc 指针，可安全比较；Palette（slice）走类型断言分支。
func classifyColorModel(m color.Model) (colorSpace string, hasAlpha bool, bitDepth int) {
	if m == nil {
		return "RGB", false, 8
	}
	if _, ok := m.(color.Palette); ok {
		return "Paletted", false, 8
	}
	switch m {
	case color.YCbCrModel:
		return "YCbCr", false, 8
	case color.NYCbCrAModel:
		return "YCbCrA", true, 8
	case color.RGBAModel:
		return "RGBA", true, 8
	case color.RGBA64Model:
		return "RGBA64", true, 16
	case color.NRGBAModel:
		return "NRGBA", true, 8
	case color.NRGBA64Model:
		return "NRGBA64", true, 16
	case color.GrayModel:
		return "Gray", false, 8
	case color.Gray16Model:
		return "Gray16", false, 16
	case color.AlphaModel:
		return "Alpha", true, 8
	case color.Alpha16Model:
		return "Alpha16", true, 16
	case color.CMYKModel:
		return "CMYK", false, 8
	}
	return "RGB", false, 8
}
