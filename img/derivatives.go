//go:build !libvips

package img

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	_ "golang.org/x/image/bmp"
	xdraw "golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

// MakeDerivatives 纯 Go 后端：image + x/image/draw + go-webp，无需系统级 libvips。
// 源图整幅解码进内存（NRGBA 每像素 4B），故沿用 maxPreviewSourcePixels 像素护栏；
// EXIF Orientation 仅对 JPEG 生效。需要全格式自动旋转与更低内存请用 `-tags libvips` 构建。
func MakeDerivatives(ctx context.Context, opts DerivativesOptions) error {
	if err := validateDerivativesOptions(&opts); err != nil {
		return err
	}
	if err := guardSourcePixels(opts.Source); err != nil {
		return err
	}
	base, err := decodeOrientedImage(ctx, opts.Source)
	if err != nil {
		return err
	}
	var mark image.Image
	if opts.Watermark != nil {
		if mark, err = decodeWatermarkImage(opts.Watermark.Source); err != nil {
			return err
		}
	}
	for _, out := range opts.Outputs {
		if err := ctx.Err(); err != nil {
			return err
		}
		im := base
		if out.Width > 0 && out.Width < base.Bounds().Dx() {
			im = scaleToWidth(base, out.Width)
		}
		// 缩放后再加效果，保证各尺寸下的锐化/水印观感独立正确；applyDerivativeEffects 不修改入参。
		im = applyDerivativeEffects(ctx, im, opts, mark)
		if err := encodeDerivativeFile(out, im); err != nil {
			return err
		}
	}
	return nil
}

// guardSourcePixels 全解码前的像素护栏：用 header 尺寸挡住「小体积高像素」压缩炸弹。
func guardSourcePixels(source string) error {
	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("img: open source image: %w", err)
	}
	defer func() { _ = f.Close() }()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return fmt.Errorf("%w: decode config: %v", ErrImageDecode, err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width*cfg.Height > maxPreviewSourcePixels {
		return fmt.Errorf("%w: %dx%d", ErrImageTooLarge, cfg.Width, cfg.Height)
	}
	return nil
}

func decodeOrientedImage(ctx context.Context, source string) (*image.NRGBA, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	f, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("img: open source image: %w", err)
	}
	decoded, format, err := image.Decode(bufio.NewReader(f))
	_ = f.Close()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrImageDecode, err)
	}
	if strings.EqualFold(format, "jpeg") {
		decoded = applyOrientation(decoded, jpegOrientation(source))
	}
	canvas := image.NewNRGBA(image.Rect(0, 0, decoded.Bounds().Dx(), decoded.Bounds().Dy()))
	draw.Draw(canvas, canvas.Bounds(), decoded, decoded.Bounds().Min, draw.Src)
	return canvas, nil
}

func decodeWatermarkImage(source string) (image.Image, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("img: open watermark: %w", err)
	}
	mark, _, err := image.Decode(bufio.NewReader(f))
	_ = f.Close()
	if err != nil {
		return nil, fmt.Errorf("%w: watermark: %v", ErrImageDecode, err)
	}
	return mark, nil
}

func scaleToWidth(src *image.NRGBA, width int) *image.NRGBA {
	height := max(1, int(math.Round(float64(src.Bounds().Dy())*float64(width)/float64(src.Bounds().Dx()))))
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), xdraw.Src, nil)
	return dst
}

// applyDerivativeEffects 不修改 src：锐化输出新缓冲，水印在副本上叠加。
func applyDerivativeEffects(ctx context.Context, src *image.NRGBA, opts DerivativesOptions, mark image.Image) *image.NRGBA {
	canvas := src
	if opts.Sharpen != nil {
		canvas = sharpenImage(ctx, src, *opts.Sharpen)
	}
	if mark != nil {
		if canvas == src {
			canvas = cloneNRGBA(src)
		}
		overlayWatermark(canvas, mark, *opts.Watermark)
	}
	return canvas
}

func cloneNRGBA(src *image.NRGBA) *image.NRGBA {
	dst := image.NewNRGBA(src.Bounds())
	copy(dst.Pix, src.Pix)
	return dst
}

func encodeDerivativeFile(out DerivativeOutput, im image.Image) error {
	var buf bytes.Buffer
	switch out.Format {
	case "webp":
		opts, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(out.Quality))
		if err != nil {
			return fmt.Errorf("img: webp options %s: %w", filepath.Base(out.Path), err)
		}
		if err := webp.Encode(&buf, im, opts); err != nil {
			return fmt.Errorf("img: encode derivative %s: %w", filepath.Base(out.Path), err)
		}
	case "jpeg":
		if err := jpeg.Encode(&buf, im, &jpeg.Options{Quality: out.Quality}); err != nil {
			return fmt.Errorf("img: encode derivative %s: %w", filepath.Base(out.Path), err)
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedSink, out.Format)
	}
	return writeFileAtomic(out.Path, buf.Bytes())
}

// sharpenFlatCutoff 划分平坦区/边缘区的拉普拉斯响应阈值（|4c-l-r-u-d|，量程 0..1020）。
const sharpenFlatCutoff = 32.0

// sharpenImage 近似 libvips sharpen 语义：m1 是平坦区斜率（弱化噪点放大），
// m2 是边缘区斜率。m1 未配置（<=0）时退回旧行为，全图使用 m2 斜率。
func sharpenImage(ctx context.Context, src *image.NRGBA, cfg DerivativeSharpen) *image.NRGBA {
	dst := image.NewNRGBA(src.Bounds())
	amountJagged := cfg.Sigma * cfg.M2 * 0.18
	if amountJagged <= 0 {
		amountJagged = 0.18
	}
	amountFlat := cfg.Sigma * cfg.M1 * 0.18
	if cfg.M1 <= 0 {
		amountFlat = amountJagged
	}
	for y := 0; y < src.Bounds().Dy(); y++ {
		if y%64 == 0 && ctx.Err() != nil {
			return src
		}
		for x := 0; x < src.Bounds().Dx(); x++ {
			center := src.NRGBAAt(x, y)
			left := src.NRGBAAt(max(0, x-1), y)
			right := src.NRGBAAt(min(src.Bounds().Dx()-1, x+1), y)
			up := src.NRGBAAt(x, max(0, y-1))
			down := src.NRGBAAt(x, min(src.Bounds().Dy()-1, y+1))
			dst.SetNRGBA(x, y, color.NRGBA{
				R: sharpenChannel(center.R, left.R, right.R, up.R, down.R, amountFlat, amountJagged),
				G: sharpenChannel(center.G, left.G, right.G, up.G, down.G, amountFlat, amountJagged),
				B: sharpenChannel(center.B, left.B, right.B, up.B, down.B, amountFlat, amountJagged),
				A: center.A,
			})
		}
	}
	return dst
}

func sharpenChannel(c, l, r, u, d uint8, amountFlat, amountJagged float64) uint8 {
	delta := 4*float64(c) - float64(l) - float64(r) - float64(u) - float64(d)
	amount := amountJagged
	if math.Abs(delta) < sharpenFlatCutoff {
		amount = amountFlat
	}
	value := float64(c) + amount*delta
	return uint8(max(0, min(255, int(value+0.5))))
}

func overlayWatermark(dst *image.NRGBA, mark image.Image, cfg DerivativeWatermark) {
	maxWidth := max(1, dst.Bounds().Dx()/3)
	if mark.Bounds().Dx() > maxWidth {
		height := mark.Bounds().Dy() * maxWidth / mark.Bounds().Dx()
		resized := image.NewNRGBA(image.Rect(0, 0, maxWidth, max(1, height)))
		xdraw.CatmullRom.Scale(resized, resized.Bounds(), mark, mark.Bounds(), xdraw.Over, nil)
		mark = resized
	}
	x, y := watermarkOffset(dst.Bounds().Dx(), dst.Bounds().Dy(), mark.Bounds().Dx(), mark.Bounds().Dy(), cfg)
	opacity := uint8(cfg.Opacity * 255 / 100)
	mask := image.NewUniform(color.Alpha{A: opacity})
	draw.DrawMask(dst, image.Rect(x, y, x+mark.Bounds().Dx(), y+mark.Bounds().Dy()),
		mark, mark.Bounds().Min, mask, image.Point{}, draw.Over)
}

func applyOrientation(src image.Image, orientation int) image.Image {
	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	if orientation == 1 || orientation == 0 {
		return src
	}
	var dst *image.NRGBA
	if orientation >= 5 && orientation <= 8 {
		dst = image.NewNRGBA(image.Rect(0, 0, h, w))
	} else {
		dst = image.NewNRGBA(image.Rect(0, 0, w, h))
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var dx, dy int
			switch orientation {
			case 2:
				dx, dy = w-1-x, y
			case 3:
				dx, dy = w-1-x, h-1-y
			case 4:
				dx, dy = x, h-1-y
			case 5:
				dx, dy = y, x
			case 6:
				dx, dy = h-1-y, x
			case 7:
				dx, dy = h-1-y, w-1-x
			case 8:
				dx, dy = y, w-1-x
			}
			dst.Set(dx, dy, src.At(src.Bounds().Min.X+x, src.Bounds().Min.Y+y))
		}
	}
	return dst
}

func jpegOrientation(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 1
	}
	defer func() { _ = f.Close() }()
	var marker [2]byte
	if _, err := io.ReadFull(f, marker[:]); err != nil || marker != [2]byte{0xff, 0xd8} {
		return 1
	}
	for {
		if _, err := io.ReadFull(f, marker[:]); err != nil || marker[0] != 0xff {
			return 1
		}
		if marker[1] == 0xda || marker[1] == 0xd9 {
			return 1
		}
		var sizeBuf [2]byte
		if _, err := io.ReadFull(f, sizeBuf[:]); err != nil {
			return 1
		}
		size := int(binary.BigEndian.Uint16(sizeBuf[:])) - 2
		if size <= 0 || size > 1<<20 {
			return 1
		}
		segment := make([]byte, size)
		if _, err := io.ReadFull(f, segment); err != nil {
			return 1
		}
		if marker[1] == 0xe1 && len(segment) > 14 && string(segment[:6]) == "Exif\x00\x00" {
			return parseTIFFOrientation(segment[6:])
		}
	}
}

func parseTIFFOrientation(data []byte) int {
	if len(data) < 8 {
		return 1
	}
	var order binary.ByteOrder
	switch string(data[:2]) {
	case "II":
		order = binary.LittleEndian
	case "MM":
		order = binary.BigEndian
	default:
		return 1
	}
	offset := int(order.Uint32(data[4:8]))
	if offset < 0 || offset+2 > len(data) {
		return 1
	}
	count := int(order.Uint16(data[offset : offset+2]))
	for i := 0; i < count; i++ {
		entry := offset + 2 + i*12
		if entry+12 > len(data) {
			break
		}
		if order.Uint16(data[entry:entry+2]) == 0x0112 {
			value := int(order.Uint16(data[entry+8 : entry+10]))
			if value >= 1 && value <= 8 {
				return value
			}
		}
	}
	return 1
}
