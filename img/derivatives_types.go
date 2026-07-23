package img

import (
	"errors"
	"strings"
)

// ErrNoDerivativeOutputs MakeDerivatives 的 Outputs 为空。
var ErrNoDerivativeOutputs = errors.New("img: derivative outputs must not be empty")

// DerivativeOutput 描述一个派生物输出。
type DerivativeOutput struct {
	// Path 输出文件路径（父目录不存在时自动创建，原子写入）。
	Path string
	// Width >0 且小于源宽时等比缩放到该宽度；0 或 ≥源宽 = 保持原尺寸。不放大。
	Width int
	// Format 输出格式："webp" | "jpeg"。
	Format string
	// Quality 1..100，0/越界回退 80。
	Quality int
}

// DerivativeSharpen 锐化参数，语义对齐 libvips sharpen（sigma + 平坦区/边缘区双斜率）。
//
// 后端差异：纯 Go 后端用 4 邻域拉普拉斯近似实现，m1/m2 均生效；libvips 后端走原生
// vips sharpen，但 govips 未暴露 m1（固定为 vips 默认 0 = 平坦区不锐化，抗噪更强），
// x1 阈值取 vips 默认 2.0。两后端观感存在轻微差异，属已接受的后端差异。
type DerivativeSharpen struct {
	Sigma float64
	M1    float64 // 平坦区斜率（<=0 时纯 Go 后端退回全图使用 M2）
	M2    float64 // 边缘区斜率
}

// DerivativeWatermark 水印参数。水印图最大宽度为目标图宽度的 1/3（超出等比缩小）。
type DerivativeWatermark struct {
	// Source 水印图路径。
	Source string
	// Position 九宫格位置：topleft|top|topright|left|center|right|bottomleft|bottom|bottomright；
	// 无法识别的值回退 bottomright。
	Position string
	// Opacity 不透明度 0..100（越界钳制）。
	Opacity  int
	PaddingX int
	PaddingY int
}

// DerivativesOptions MakeDerivatives 的输入。
//
// 语义：源图单次解码 → 每个输出按 Width 等比缩放 → 锐化 → 水印 → 编码落盘。
// 缩放先于锐化/水印，保证各尺寸下的锐化与水印观感独立正确。
// 动图（多帧 GIF/WebP）不受支持：两后端都只处理首帧，调用方应先用 GetImageInfo
// 的 Animated 自行拒绝。EXIF Orientation：libvips 后端全格式自动旋转；纯 Go 后端仅 JPEG。
type DerivativesOptions struct {
	Source    string
	Sharpen   *DerivativeSharpen   // nil = 不锐化
	Watermark *DerivativeWatermark // nil = 不加水印
	Outputs   []DerivativeOutput
}

// validateDerivativesOptions 归一化并校验参数：钳制 Quality/Opacity，拒绝空源/空输出/未知格式。
func validateDerivativesOptions(opts *DerivativesOptions) error {
	if opts.Source == "" {
		return ErrEmptySource
	}
	if len(opts.Outputs) == 0 {
		return ErrNoDerivativeOutputs
	}
	for i := range opts.Outputs {
		out := &opts.Outputs[i]
		if out.Path == "" {
			return ErrEmptySource
		}
		switch out.Format {
		case "webp", "jpeg":
		default:
			return ErrUnsupportedSink
		}
		if out.Quality <= 0 || out.Quality > 100 {
			out.Quality = 80
		}
	}
	if wm := opts.Watermark; wm != nil {
		if wm.Source == "" {
			return ErrEmptySource
		}
		if wm.Opacity < 0 {
			wm.Opacity = 0
		}
		if wm.Opacity > 100 {
			wm.Opacity = 100
		}
	}
	return nil
}

// watermarkOffset 按九宫格位置与 padding 计算水印左上角坐标（两后端共用的纯计算）。
func watermarkOffset(dstW, dstH, markW, markH int, wm DerivativeWatermark) (int, int) {
	left, centerX, right := wm.PaddingX, (dstW-markW)/2, dstW-markW-wm.PaddingX
	top, centerY, bottom := wm.PaddingY, (dstH-markH)/2, dstH-markH-wm.PaddingY
	switch strings.ToLower(wm.Position) {
	case "topleft":
		return left, top
	case "top":
		return centerX, top
	case "topright":
		return right, top
	case "left":
		return left, centerY
	case "center":
		return centerX, centerY
	case "right":
		return right, centerY
	case "bottomleft":
		return left, bottom
	case "bottom":
		return centerX, bottom
	default:
		return right, bottom
	}
}
