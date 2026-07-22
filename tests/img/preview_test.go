package img_test

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"testing"

	xwebp "golang.org/x/image/webp"

	"github.com/bizvip/go-utils/img"
)

func jpegBytes(t *testing.T, w, h int) []byte {
	t.Helper()
	src := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			src.Set(x, y, color.RGBA{R: uint8(x % 255), G: uint8(y % 255), B: 128, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, src, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("encode source jpeg: %v", err)
	}
	return buf.Bytes()
}

// MakeWebpPreview 两后端(纯 Go / -tags libvips)行为契约相同：输出可解码 WebP、
// 最长边缩到 ≤maxSide、保持纵横比、不放大。
func TestMakeWebpPreviewDownscales(t *testing.T) {
	out, err := img.MakeWebpPreview(jpegBytes(t, 1600, 900), 320, 45)
	if err != nil {
		t.Fatalf("MakeWebpPreview: %v", err)
	}
	cfg, err := xwebp.DecodeConfig(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("output is not decodable webp: %v", err)
	}
	if cfg.Width != 320 {
		t.Fatalf("longest edge = %d, want 320", cfg.Width)
	}
	if cfg.Height != 180 { // 320 * 900/1600
		t.Fatalf("height = %d, want 180 (aspect preserved)", cfg.Height)
	}
}

func TestMakeWebpPreviewNoUpscale(t *testing.T) {
	out, err := img.MakeWebpPreview(jpegBytes(t, 100, 80), 320, 45)
	if err != nil {
		t.Fatalf("MakeWebpPreview: %v", err)
	}
	cfg, err := xwebp.DecodeConfig(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("output is not decodable webp: %v", err)
	}
	if cfg.Width != 100 || cfg.Height != 80 {
		t.Fatalf("small image was resized to %dx%d, want 100x80 (no upscale)", cfg.Width, cfg.Height)
	}
}

func TestMakeWebpPreviewRejectsGarbage(t *testing.T) {
	if _, err := img.MakeWebpPreview([]byte("not an image"), 320, 45); !errors.Is(err, img.ErrImageDecode) {
		t.Fatalf("garbage input error = %v, want ErrImageDecode", err)
	}
}
