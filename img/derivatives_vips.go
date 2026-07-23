//go:build libvips

package img

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/davidbyttow/govips/v2/vips"
)

// MakeDerivatives libvips 后端：原生 vips resize/sharpen/composite 管线。
// 相比纯 Go 后端：EXIF Orientation 全格式自动旋转（autorot）、SIMD 编码更快、
// 内存峰值更低，无需像素护栏。语义与纯 Go 后端(derivatives.go)一致：
// 每个输出按 Width 等比缩放（不放大）→ 锐化 → 水印 → 编码落盘，输出剥离元数据。
// 需以 `-tags libvips` 构建，构建/运行主机须装 libvips。
func MakeDerivatives(ctx context.Context, opts DerivativesOptions) error {
	if err := validateDerivativesOptions(&opts); err != nil {
		return err
	}
	StartupVips()

	base, err := vips.NewImageFromFile(opts.Source)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrImageDecode, err)
	}
	defer base.Close()
	if err := base.AutoRotate(); err != nil {
		return fmt.Errorf("img: autorotate source: %w", err)
	}

	var mark *vips.ImageRef
	if opts.Watermark != nil {
		mark, err = vips.NewImageFromFile(opts.Watermark.Source)
		if err != nil {
			return fmt.Errorf("%w: watermark: %v", ErrImageDecode, err)
		}
		defer mark.Close()
		// 统一成 sRGB+alpha，方便按不透明度调制 alpha 通道后 composite。
		if err := mark.ToColorSpace(vips.InterpretationSRGB); err != nil {
			return fmt.Errorf("img: watermark colorspace: %w", err)
		}
		if !mark.HasAlpha() {
			if err := mark.AddAlpha(); err != nil {
				return fmt.Errorf("img: watermark alpha: %w", err)
			}
		}
	}

	for _, out := range opts.Outputs {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := makeVipsDerivative(base, mark, opts, out); err != nil {
			return err
		}
	}
	return nil
}

func makeVipsDerivative(base, mark *vips.ImageRef, opts DerivativesOptions, out DerivativeOutput) error {
	im, err := base.Copy()
	if err != nil {
		return fmt.Errorf("img: copy source: %w", err)
	}
	defer im.Close()

	if out.Width > 0 && out.Width < im.Width() {
		if err := im.Resize(float64(out.Width)/float64(im.Width()), vips.KernelLanczos3); err != nil {
			return fmt.Errorf("img: resize derivative %s: %w", filepath.Base(out.Path), err)
		}
	}
	if s := opts.Sharpen; s != nil {
		// govips 未暴露 vips_sharpen 的 m1（平坦区斜率），vips 按默认 m1=0 处理 =
		// 平坦区不锐化（抗噪更强）；x1 阈值取 vips 默认 2.0。见 DerivativeSharpen 注释。
		if err := im.Sharpen(s.Sigma, 2.0, s.M2); err != nil {
			return fmt.Errorf("img: sharpen derivative %s: %w", filepath.Base(out.Path), err)
		}
	}
	if mark != nil {
		if err := compositeVipsWatermark(im, mark, *opts.Watermark); err != nil {
			return fmt.Errorf("img: watermark derivative %s: %w", filepath.Base(out.Path), err)
		}
	}

	var buf []byte
	switch out.Format {
	case "webp":
		params := vips.NewWebpExportParams()
		params.Quality = out.Quality
		params.StripMetadata = true
		if buf, _, err = im.ExportWebp(params); err != nil {
			return fmt.Errorf("img: encode derivative %s: %w", filepath.Base(out.Path), err)
		}
	case "jpeg":
		// jpeg 无 alpha：黑底拍平，对齐纯 Go 后端 jpeg.Encode 丢弃 alpha 的行为。
		if im.HasAlpha() {
			if err := im.Flatten(&vips.Color{R: 0, G: 0, B: 0}); err != nil {
				return fmt.Errorf("img: flatten derivative %s: %w", filepath.Base(out.Path), err)
			}
		}
		params := vips.NewJpegExportParams()
		params.Quality = out.Quality
		params.StripMetadata = true
		if buf, _, err = im.ExportJpeg(params); err != nil {
			return fmt.Errorf("img: encode derivative %s: %w", filepath.Base(out.Path), err)
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedSink, out.Format)
	}
	return writeFileAtomic(out.Path, buf)
}

func compositeVipsWatermark(dst, mark *vips.ImageRef, cfg DerivativeWatermark) error {
	mc, err := mark.Copy()
	if err != nil {
		return err
	}
	defer mc.Close()

	maxWidth := max(1, dst.Width()/3)
	if mc.Width() > maxWidth {
		if err := mc.Resize(float64(maxWidth)/float64(mc.Width()), vips.KernelLanczos3); err != nil {
			return err
		}
	}
	if cfg.Opacity < 100 {
		// 只调制 alpha 通道（mark 已归一为 sRGB+alpha = 4 band）。
		factor := float64(cfg.Opacity) / 100
		if err := mc.Linear([]float64{1, 1, 1, factor}, []float64{0, 0, 0, 0}); err != nil {
			return err
		}
	}
	x, y := watermarkOffset(dst.Width(), dst.Height(), mc.Width(), mc.Height(), cfg)
	return dst.Composite(mc, vips.BlendModeOver, x, y)
}
