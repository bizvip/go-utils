//go:build !libvips

package img

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	_ "golang.org/x/image/bmp"
	xdraw "golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

// maxPreviewSourcePixels 纯 Go 后端全解码前的像素护栏：image.Decode 会把整图解到内存
// （NRGBA 每像素 4B，~60MP ≈ 240MB 峰值），故用 header 尺寸先挡住「小体积高像素」压缩炸弹。
// libvips 后端靠 shrink-on-load 天然限制内存，不需要此护栏。
const maxPreviewSourcePixels = 60_000_000

// MakeWebpPreview 解码内存中的图像、等比缩放到最长边 ≤maxSide（不放大），编码为 quality(1..100)
// 的有损 WebP 返回。用于生成低清预览 / 缩略图。
//
// 本文件是纯 Go 后端（image + x/image/draw + go-webp），无需系统级 libvips。
// 注意：纯 Go 后端**不应用 EXIF Orientation**——手机竖拍照片(Orientation 6/8)预览可能侧转；
// 需要 EXIF 自动旋转与更低内存请用 `-tags libvips` 构建（见 preview_vips.go）。
func MakeWebpPreview(data []byte, maxSide, quality int) ([]byte, error) {
	if maxSide <= 0 {
		maxSide = defaultPreviewMaxSide
	}
	if quality <= 0 || quality > 100 {
		quality = defaultPreviewQuality
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, fmt.Errorf("%w: decode config", ErrImageDecode)
	}
	if int64(cfg.Width)*int64(cfg.Height) > maxPreviewSourcePixels {
		return nil, fmt.Errorf("%w: %dx%d", ErrImageTooLarge, cfg.Width, cfg.Height)
	}
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrImageDecode, err)
	}
	w, h := fitWithin(cfg.Width, cfg.Height, maxSide)
	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), xdraw.Src, nil)

	opts, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(quality))
	if err != nil {
		return nil, fmt.Errorf("img: webp options: %w", err)
	}
	var out bytes.Buffer
	if err := webp.Encode(&out, dst, opts); err != nil {
		return nil, fmt.Errorf("img: encode webp: %w", err)
	}
	return out.Bytes(), nil
}

// fitWithin 返回最长边 ≤maxSide、保持纵横比的最大尺寸；小于框则原样返回（不放大）。
func fitWithin(w, h, maxSide int) (int, int) {
	if w <= maxSide && h <= maxSide {
		return w, h
	}
	if w >= h {
		return maxSide, max1(h * maxSide / w)
	}
	return max1(w * maxSide / h), maxSide
}
