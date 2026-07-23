package img

// 本测试文件不带 build tag：`go test ./img` 走纯 Go 后端，`go test -tags libvips ./img`
// 走 libvips 后端，同一份断言即两后端行为对齐测试（尺寸/格式/水印/EXIF 定向语义级一致；
// 像素级不要求 byte 相同）。

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	_ "golang.org/x/image/webp"
)

func writeGradientPNG(t *testing.T, path string, w, h int) {
	t.Helper()
	im := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			im.SetNRGBA(x, y, color.NRGBA{R: uint8(x % 255), G: uint8(y % 255), B: 96, A: 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, im); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func decodeFile(t *testing.T, path string) (image.Image, string) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	im, format, err := image.Decode(f)
	if err != nil {
		t.Fatalf("decode %s: %v", filepath.Base(path), err)
	}
	return im, format
}

func TestMakeDerivatives(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source.png")
	writeGradientPNG(t, source, 640, 360)

	markPath := filepath.Join(root, "mark.png")
	mark := image.NewNRGBA(image.Rect(0, 0, 300, 80))
	for y := 0; y < 80; y++ {
		for x := 0; x < 300; x++ {
			mark.SetNRGBA(x, y, color.NRGBA{R: 255, G: 0, B: 0, A: 200})
		}
	}
	fm, err := os.Create(markPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(fm, mark); err != nil {
		_ = fm.Close()
		t.Fatal(err)
	}
	if err := fm.Close(); err != nil {
		t.Fatal(err)
	}

	fullWebp := filepath.Join(root, "out", "image_full.webp")
	fullJpeg := filepath.Join(root, "out", "image_full.jpg")
	thumbWebp := filepath.Join(root, "out", "image_320.webp")
	err = MakeDerivatives(context.Background(), DerivativesOptions{
		Source:  source,
		Sharpen: &DerivativeSharpen{Sigma: 0.6, M1: 1.0, M2: 2.0},
		Watermark: &DerivativeWatermark{
			Source: markPath, Position: "bottomright", Opacity: 60, PaddingX: 10, PaddingY: 10,
		},
		Outputs: []DerivativeOutput{
			{Path: fullWebp, Width: 0, Format: "webp", Quality: 80},
			{Path: fullJpeg, Width: 0, Format: "jpeg", Quality: 85},
			{Path: thumbWebp, Width: 320, Format: "webp", Quality: 60},
		},
	})
	if err != nil {
		t.Fatalf("MakeDerivatives: %v", err)
	}

	for _, tc := range []struct {
		path         string
		format       string
		wantW, wantH int
	}{
		{fullWebp, "webp", 640, 360},
		{fullJpeg, "jpeg", 640, 360},
		{thumbWebp, "webp", 320, 180},
	} {
		stat, err := os.Stat(tc.path)
		if err != nil || stat.Size() == 0 {
			t.Fatalf("missing derivative %s: stat=%v err=%v", filepath.Base(tc.path), stat, err)
		}
		im, format := decodeFile(t, tc.path)
		if format != tc.format {
			t.Fatalf("%s format=%s want %s", filepath.Base(tc.path), format, tc.format)
		}
		if im.Bounds().Dx() != tc.wantW || im.Bounds().Dy() != tc.wantH {
			t.Fatalf("%s dims=%dx%d want %dx%d",
				filepath.Base(tc.path), im.Bounds().Dx(), im.Bounds().Dy(), tc.wantW, tc.wantH)
		}
	}

	// 右下角应带红色水印色偏（水印区域内取样）。
	full, _ := decodeFile(t, fullJpeg)
	r, g, _, _ := full.At(600, 330).RGBA()
	if r <= g {
		t.Fatalf("expected red watermark tint at bottom-right, got r=%d g=%d", r>>8, g>>8)
	}
}

// TestMakeDerivativesOrientation 验证 EXIF Orientation=6（旋转 90°）后输出宽高互换。
// 两后端路径不同（纯 Go 手写 EXIF 解析 / libvips autorot），断言相同。
func TestMakeDerivativesOrientation(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "oriented.jpg")

	im := image.NewNRGBA(image.Rect(0, 0, 640, 360))
	for y := 0; y < 360; y++ {
		for x := 0; x < 640; x++ {
			im.SetNRGBA(x, y, color.NRGBA{R: 200, G: 120, B: 40, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, im, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(source, injectExifOrientation(t, buf.Bytes(), 6), 0o644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(root, "rotated.jpg")
	err := MakeDerivatives(context.Background(), DerivativesOptions{
		Source:  source,
		Outputs: []DerivativeOutput{{Path: out, Format: "jpeg", Quality: 85}},
	})
	if err != nil {
		t.Fatalf("MakeDerivatives: %v", err)
	}
	rotated, _ := decodeFile(t, out)
	if rotated.Bounds().Dx() != 360 || rotated.Bounds().Dy() != 640 {
		t.Fatalf("orientation 6 dims=%dx%d, want 360x640",
			rotated.Bounds().Dx(), rotated.Bounds().Dy())
	}
}

// injectExifOrientation 在 JPEG SOI 后插入只含 Orientation 标签的最小 EXIF APP1 段。
func injectExifOrientation(t *testing.T, jpg []byte, orientation uint16) []byte {
	t.Helper()
	if len(jpg) < 2 || jpg[0] != 0xff || jpg[1] != 0xd8 {
		t.Fatal("not a jpeg")
	}
	tiff := make([]byte, 0, 26)
	tiff = append(tiff, 'I', 'I', 0x2a, 0x00, 0x08, 0x00, 0x00, 0x00) // TIFF header, IFD0 @8
	tiff = append(tiff, 0x01, 0x00)                                   // 1 entry
	entry := make([]byte, 12)
	binary.LittleEndian.PutUint16(entry[0:2], 0x0112) // Orientation
	binary.LittleEndian.PutUint16(entry[2:4], 3)      // SHORT
	binary.LittleEndian.PutUint32(entry[4:8], 1)
	binary.LittleEndian.PutUint16(entry[8:10], orientation)
	tiff = append(tiff, entry...)
	tiff = append(tiff, 0x00, 0x00, 0x00, 0x00) // next IFD = none

	payload := append([]byte("Exif\x00\x00"), tiff...)
	seg := make([]byte, 0, len(payload)+4)
	seg = append(seg, 0xff, 0xe1)
	seg = binary.BigEndian.AppendUint16(seg, uint16(len(payload)+2))
	seg = append(seg, payload...)

	out := make([]byte, 0, len(jpg)+len(seg))
	out = append(out, jpg[:2]...)
	out = append(out, seg...)
	out = append(out, jpg[2:]...)
	return out
}

func TestMakeDerivativesValidation(t *testing.T) {
	ctx := context.Background()
	if err := MakeDerivatives(ctx, DerivativesOptions{}); !errors.Is(err, ErrEmptySource) {
		t.Fatalf("empty source err=%v", err)
	}
	if err := MakeDerivatives(ctx, DerivativesOptions{Source: "x.png"}); !errors.Is(err, ErrNoDerivativeOutputs) {
		t.Fatalf("no outputs err=%v", err)
	}
	err := MakeDerivatives(ctx, DerivativesOptions{
		Source:  "x.png",
		Outputs: []DerivativeOutput{{Path: "y.gif", Format: "gif"}},
	})
	if !errors.Is(err, ErrUnsupportedSink) {
		t.Fatalf("bad format err=%v", err)
	}
}
