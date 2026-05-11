package img

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"golang.org/x/image/bmp"
	xdraw "golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

func ResizeImage(opts ResizeOptions) error {
	if opts.Src == "" {
		return ErrEmptySource
	}
	if opts.SavePath == "" {
		opts.SavePath = opts.Src
	}
	if opts.Width < 0 || opts.Height < 0 {
		return ErrInvalidSize
	}
	if opts.Quality <= 0 || opts.Quality > 100 {
		opts.Quality = 80
	}

	format := opts.Format
	if format == "" {
		detected, err := detectFormatByHead(opts.Src)
		if err != nil {
			return err
		}
		format = detected
	}

	if format == "gif" {
		return resizeAnimatedGIF(opts)
	}

	src, err := os.Open(opts.Src)
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	srcImg, _, err := image.Decode(src)
	if err != nil {
		return fmt.Errorf("img: decode %s: %w", opts.Src, err)
	}

	dstW, dstH := computeTargetSize(srcImg.Bounds().Dx(), srcImg.Bounds().Dy(), opts.Width, opts.Height, opts.Fit)
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), srcImg, srcImg.Bounds(), xdraw.Over, nil)

	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: opts.Quality}); err != nil {
			return fmt.Errorf("img: encode jpeg: %w", err)
		}
	case "png":
		enc := &png.Encoder{CompressionLevel: png.DefaultCompression}
		if err := enc.Encode(&buf, dst); err != nil {
			return fmt.Errorf("img: encode png: %w", err)
		}
	case "bmp":
		if err := bmp.Encode(&buf, dst); err != nil {
			return fmt.Errorf("img: encode bmp: %w", err)
		}
	case "webp":
		webpOpts, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(opts.Quality))
		if err != nil {
			return fmt.Errorf("img: webp options: %w", err)
		}
		if err := webp.Encode(&buf, dst, webpOpts); err != nil {
			return fmt.Errorf("img: encode webp: %w", err)
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedSink, format)
	}

	return writeFileAtomic(opts.SavePath, buf.Bytes())
}

func resizeAnimatedGIF(opts ResizeOptions) error {
	src, err := os.Open(opts.Src)
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	g, err := gif.DecodeAll(src)
	if err != nil {
		return fmt.Errorf("img: decode animated gif: %w", err)
	}
	if len(g.Image) == 0 {
		return errors.New("img: gif has no frames")
	}

	canvasW := g.Config.Width
	canvasH := g.Config.Height
	if canvasW == 0 || canvasH == 0 {
		canvasW = g.Image[0].Bounds().Dx()
		canvasH = g.Image[0].Bounds().Dy()
	}
	dstW, dstH := computeTargetSize(canvasW, canvasH, opts.Width, opts.Height, opts.Fit)

	scaleX := float64(dstW) / float64(canvasW)
	scaleY := float64(dstH) / float64(canvasH)

	prevCanvas := image.NewRGBA(image.Rect(0, 0, canvasW, canvasH))
	backupCanvas := image.NewRGBA(image.Rect(0, 0, canvasW, canvasH))

	resized := make([]*image.Paletted, len(g.Image))
	for i, frame := range g.Image {
		if i > 0 {
			switch g.Disposal[i-1] {
			case gif.DisposalBackground:
				draw.Draw(prevCanvas, g.Image[i-1].Bounds(), image.Transparent, image.Point{}, draw.Src)
			case gif.DisposalPrevious:
				draw.Draw(prevCanvas, prevCanvas.Bounds(), backupCanvas, image.Point{}, draw.Src)
			}
		}

		if i < len(g.Disposal) && g.Disposal[i] == gif.DisposalPrevious {
			draw.Draw(backupCanvas, backupCanvas.Bounds(), prevCanvas, image.Point{}, draw.Src)
		}

		composed := image.NewRGBA(prevCanvas.Bounds())
		draw.Draw(composed, composed.Bounds(), prevCanvas, image.Point{}, draw.Src)
		draw.Draw(composed, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)

		draw.Draw(prevCanvas, prevCanvas.Bounds(), composed, image.Point{}, draw.Src)

		dstFrame := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
		xdraw.CatmullRom.Scale(dstFrame, dstFrame.Bounds(), composed, composed.Bounds(), xdraw.Over, nil)

		framePalette := frame.Palette
		if framePalette == nil || len(framePalette) == 0 {
			framePalette = palette.Plan9
		}
		paletted := image.NewPaletted(dstFrame.Bounds(), framePalette)
		draw.FloydSteinberg.Draw(paletted, dstFrame.Bounds(), dstFrame, image.Point{})
		resized[i] = paletted
	}

	scaleRect := func(r image.Rectangle) image.Rectangle {
		x0 := int(float64(r.Min.X)*scaleX + 0.5)
		y0 := int(float64(r.Min.Y)*scaleY + 0.5)
		x1 := int(float64(r.Max.X)*scaleX + 0.5)
		y1 := int(float64(r.Max.Y)*scaleY + 0.5)
		if x1 > dstW {
			x1 = dstW
		}
		if y1 > dstH {
			y1 = dstH
		}
		if x0 >= x1 {
			x0 = x1 - 1
		}
		if y0 >= y1 {
			y0 = y1 - 1
		}
		if x0 < 0 {
			x0 = 0
		}
		if y0 < 0 {
			y0 = 0
		}
		return image.Rect(x0, y0, x1, y1)
	}

	for i := range resized {
		resized[i].Rect = scaleRect(g.Image[i].Bounds())
	}

	out := &gif.GIF{
		Image:           resized,
		Delay:           append([]int(nil), g.Delay...),
		LoopCount:       g.LoopCount,
		Disposal:        append([]byte(nil), g.Disposal...),
		BackgroundIndex: g.BackgroundIndex,
		Config: image.Config{
			ColorModel: color.Palette(palette.Plan9),
			Width:      dstW,
			Height:     dstH,
		},
	}

	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, out); err != nil {
		return fmt.Errorf("img: encode animated gif: %w", err)
	}
	return writeFileAtomic(opts.SavePath, buf.Bytes())
}
