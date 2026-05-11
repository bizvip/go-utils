package img_test

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/bizvip/go-utils/img"
)

func writeJPEG(t *testing.T, path string, w, h int) {
	t.Helper()
	src := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			src.Set(x, y, color.RGBA{R: uint8(x % 255), G: uint8(y % 255), B: 128, A: 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := jpeg.Encode(f, src, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatal(err)
	}
}

func writePNG(t *testing.T, path string, w, h int) {
	t.Helper()
	src := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			src.Set(x, y, color.RGBA{R: 255, G: uint8(x), B: uint8(y), A: 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, src); err != nil {
		t.Fatal(err)
	}
}

func writeAnimatedGIF(t *testing.T, path string, w, h, frames int) {
	t.Helper()
	g := &gif.GIF{
		LoopCount: 0,
		Config: image.Config{
			ColorModel: color.Palette(palette.Plan9),
			Width:      w,
			Height:     h,
		},
	}
	for i := 0; i < frames; i++ {
		pal := image.NewPaletted(image.Rect(0, 0, w, h), palette.Plan9)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				c := color.RGBA{R: uint8((x + i*30) % 255), G: uint8((y + i*30) % 255), B: 128, A: 255}
				pal.Set(x, y, c)
			}
		}
		g.Image = append(g.Image, pal)
		g.Delay = append(g.Delay, 10)
		g.Disposal = append(g.Disposal, gif.DisposalNone)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := gif.EncodeAll(f, g); err != nil {
		t.Fatal(err)
	}
}

func decodeImageSize(t *testing.T, path string) (int, int, string) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		t.Fatal(err)
	}
	return cfg.Width, cfg.Height, format
}

func TestResizeJPEGStretch(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.jpg")
	dst := filepath.Join(dir, "out.jpg")
	writeJPEG(t, src, 200, 100)

	if err := img.ResizeImage(img.ResizeOptions{
		Src: src, SavePath: dst, Width: 64, Height: 64, Fit: img.FitStretch,
	}); err != nil {
		t.Fatalf("resize: %v", err)
	}
	w, h, format := decodeImageSize(t, dst)
	if w != 64 || h != 64 {
		t.Fatalf("expected 64x64, got %dx%d", w, h)
	}
	if format != "jpeg" {
		t.Fatalf("expected jpeg format, got %s", format)
	}
}

func TestResizeKeepRatioOnlyWidth(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.png")
	dst := filepath.Join(dir, "out.png")
	writePNG(t, src, 400, 200)

	if err := img.ResizeImage(img.ResizeOptions{
		Src: src, SavePath: dst, Width: 100,
	}); err != nil {
		t.Fatalf("resize: %v", err)
	}
	w, h, format := decodeImageSize(t, dst)
	if w != 100 || h != 50 {
		t.Fatalf("expected 100x50 keep-ratio, got %dx%d", w, h)
	}
	if format != "png" {
		t.Fatalf("expected png format, got %s", format)
	}
}

func TestResizeKeepRatioOnlyHeight(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.png")
	dst := filepath.Join(dir, "out.png")
	writePNG(t, src, 400, 200)

	if err := img.ResizeImage(img.ResizeOptions{
		Src: src, SavePath: dst, Height: 50,
	}); err != nil {
		t.Fatalf("resize: %v", err)
	}
	w, h, _ := decodeImageSize(t, dst)
	if w != 100 || h != 50 {
		t.Fatalf("expected 100x50 keep-ratio, got %dx%d", w, h)
	}
}

func TestResizeFitContain(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.jpg")
	dst := filepath.Join(dir, "out.jpg")
	writeJPEG(t, src, 400, 200)

	if err := img.ResizeImage(img.ResizeOptions{
		Src: src, SavePath: dst, Width: 100, Height: 100, Fit: img.FitContain,
	}); err != nil {
		t.Fatalf("resize: %v", err)
	}
	w, h, _ := decodeImageSize(t, dst)
	if w != 100 || h != 50 {
		t.Fatalf("expected contain 100x50, got %dx%d", w, h)
	}
}

func TestResizeFitCover(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.jpg")
	dst := filepath.Join(dir, "out.jpg")
	writeJPEG(t, src, 400, 200)

	if err := img.ResizeImage(img.ResizeOptions{
		Src: src, SavePath: dst, Width: 100, Height: 100, Fit: img.FitCover,
	}); err != nil {
		t.Fatalf("resize: %v", err)
	}
	w, h, _ := decodeImageSize(t, dst)
	if w != 200 || h != 100 {
		t.Fatalf("expected cover 200x100, got %dx%d", w, h)
	}
}

func TestResizeAnimatedGIFPreservesFramesAndLoopCount(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.gif")
	dst := filepath.Join(dir, "out.gif")
	writeAnimatedGIF(t, src, 80, 60, 5)

	if err := img.ResizeImage(img.ResizeOptions{
		Src: src, SavePath: dst, Width: 40, Height: 30,
	}); err != nil {
		t.Fatalf("resize: %v", err)
	}

	f, err := os.Open(dst)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	g, err := gif.DecodeAll(f)
	if err != nil {
		t.Fatalf("decode out gif: %v", err)
	}
	if len(g.Image) != 5 {
		t.Fatalf("expected 5 frames preserved, got %d", len(g.Image))
	}
	if g.Config.Width != 40 || g.Config.Height != 30 {
		t.Fatalf("expected canvas 40x30, got %dx%d", g.Config.Width, g.Config.Height)
	}
	for i, d := range g.Delay {
		if d != 10 {
			t.Fatalf("frame %d delay = %d, expected 10", i, d)
		}
	}
}

func TestResizeFormatOverride(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.jpg")
	dst := filepath.Join(dir, "out.png")
	writeJPEG(t, src, 60, 60)

	if err := img.ResizeImage(img.ResizeOptions{
		Src: src, SavePath: dst, Width: 30, Height: 30, Format: "png",
	}); err != nil {
		t.Fatalf("resize: %v", err)
	}
	_, _, format := decodeImageSize(t, dst)
	if format != "png" {
		t.Fatalf("expected png after format override, got %s", format)
	}
}

func TestImageToBase64AndBack(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.png")
	writePNG(t, src, 16, 16)

	enc, err := img.ImageToBase64(src, true)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytes.HasPrefix([]byte(enc), []byte("data:image/png;base64,")) {
		t.Fatalf("expected data URI prefix, got prefix %q", enc[:30])
	}

	plain, err := img.ImageToBase64(src, false)
	if err != nil {
		t.Fatalf("encode plain: %v", err)
	}
	if _, err := base64.StdEncoding.DecodeString(plain); err != nil {
		t.Fatalf("plain base64 invalid: %v", err)
	}

	out := filepath.Join(dir, "round.png")
	if err := img.Base64ToFile(enc, out, false); err != nil {
		t.Fatalf("decode: %v", err)
	}
	original, _ := os.ReadFile(src)
	roundTripped, _ := os.ReadFile(out)
	if !bytes.Equal(original, roundTripped) {
		t.Fatalf("round-trip bytes differ (%d vs %d)", len(original), len(roundTripped))
	}
}

func TestBase64ToFileRejectsNonImage(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "bad.bin")
	notImage := base64.StdEncoding.EncodeToString([]byte("just some text, not an image"))
	if err := img.Base64ToFile(notImage, out, false); err == nil {
		t.Fatal("expected error rejecting non-image bytes")
	}
}
