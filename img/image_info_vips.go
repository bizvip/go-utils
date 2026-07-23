//go:build libvips

package img

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/davidbyttow/govips/v2/vips"

	"github.com/bizvip/go-utils/os/fs"
)

var vipsOnce sync.Once

// StartupVips initializes libvips. Safe to call multiple times; only first call takes effect.
// Library users may call this explicitly at process startup; otherwise GetImageInfo
// will call it lazily on first use.
func StartupVips() {
	vipsOnce.Do(func() {
		// govips 默认把 info 级日志刷到 stderr（每次 thumbnail 一行），生产噪音大——收敛到 warning+。
		// 库使用方如需自定义处理器/级别，可在首次调用前自行覆盖 vips.LoggingSettings。
		vips.LoggingSettings(nil, vips.LogLevelWarning)
		vips.Startup(nil)
	})
}

// ShutdownVips releases libvips resources. Call once at process exit if you want
// a clean teardown. Repeated calls are safe but only the first has effect when
// paired with StartupVips.
func ShutdownVips() {
	vips.Shutdown()
}

func GetImageInfo(path string, withMd5 bool) (*ImageInfo, error) {
	StartupVips()

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	ref, err := vips.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}
	defer ref.Close()

	width := ref.Width()
	height := ref.Height()
	pages := ref.Pages()
	if pages <= 0 {
		pages = 1
	}
	pageHeight := ref.PageHeight()
	if pages > 1 && pageHeight > 0 {
		height = pageHeight
	}
	format := vips.ImageTypes[ref.Format()]

	info := &ImageInfo{
		Width:    width,
		Height:   height,
		Format:   format,
		Ext:      strings.ToLower(strings.TrimPrefix(filepath.Ext(path), ".")),
		Size:     stat.Size(),
		FileName: filepath.Base(path),
		FrameNum: pages,
		Animated: pages > 1,
	}
	info.ColorSpace, info.HasAlpha, info.BitDepth = classifyVipsImage(ref)

	if withMd5 {
		md5sum, err := fs.GetBigFileMd5(path)
		if err != nil {
			return nil, err
		}
		info.FileMD5 = md5sum
	}

	return info, nil
}

// classifyVipsImage 把 vips interpretation/band format 归类为描述性
// ColorSpace/HasAlpha/BitDepth（见 ImageInfo 字段注释）。libvips 只见解码后的
// interpretation，无法区分 NRGBA/YCbCr 等源 color model，标签集与纯 Go 后端不同。
func classifyVipsImage(ref *vips.ImageRef) (colorSpace string, hasAlpha bool, bitDepth int) {
	colorSpace = "RGB"
	switch ref.Interpretation() {
	case vips.InterpretationBW:
		colorSpace = "Gray"
	case vips.InterpretationGrey16:
		colorSpace = "Gray16"
	case vips.InterpretationCMYK:
		colorSpace = "CMYK"
	case vips.InterpretationRGB16:
		colorSpace = "RGB16"
	case vips.InterpretationLAB, vips.InterpretationLABQ, vips.InterpretationLABS:
		colorSpace = "LAB"
	}
	bitDepth = 8
	switch ref.BandFormat() {
	case vips.BandFormatUshort, vips.BandFormatShort:
		bitDepth = 16
	case vips.BandFormatUint, vips.BandFormatInt, vips.BandFormatFloat:
		bitDepth = 32
	}
	return colorSpace, ref.HasAlpha(), bitDepth
}
