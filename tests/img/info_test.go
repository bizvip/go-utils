//go:build !libvips

package img_test

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/bizvip/go-utils/img"
)

func TestGetImageInfoJPEG(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "x.jpg")
	srcImg := image.NewRGBA(image.Rect(0, 0, 320, 240))
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			srcImg.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	f, _ := os.Create(src)
	if err := jpeg.Encode(f, srcImg, nil); err != nil {
		t.Fatal(err)
	}
	f.Close()

	info, err := img.GetImageInfo(src, false)
	if err != nil {
		t.Fatal(err)
	}
	if info.Width != 320 || info.Height != 240 {
		t.Fatalf("expected 320x240, got %dx%d", info.Width, info.Height)
	}
	if info.Format != "jpeg" {
		t.Fatalf("expected jpeg, got %s", info.Format)
	}
	if info.Animated || info.FrameNum != 1 {
		t.Fatalf("static jpeg should have FrameNum=1 Animated=false, got %d %v", info.FrameNum, info.Animated)
	}
}

func TestGetImageInfoAnimatedGIF(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.gif")

	g := &gif.GIF{LoopCount: 3, Config: image.Config{Width: 16, Height: 16}}
	for i := 0; i < 4; i++ {
		p := image.NewPaletted(image.Rect(0, 0, 16, 16), palette.Plan9)
		g.Image = append(g.Image, p)
		g.Delay = append(g.Delay, 5)
		g.Disposal = append(g.Disposal, gif.DisposalNone)
	}
	f, _ := os.Create(src)
	if err := gif.EncodeAll(f, g); err != nil {
		t.Fatal(err)
	}
	f.Close()

	info, err := img.GetImageInfo(src, false)
	if err != nil {
		t.Fatal(err)
	}
	if info.FrameNum != 4 {
		t.Fatalf("expected 4 frames, got %d", info.FrameNum)
	}
	if !info.Animated {
		t.Fatal("expected animated=true")
	}
	if info.LoopCount != 3 {
		t.Fatalf("expected loopCount=3, got %d", info.LoopCount)
	}
}
