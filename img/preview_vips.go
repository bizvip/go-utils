//go:build libvips

package img

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

// MakeWebpPreview（libvips 后端）：用 vips_thumbnail 的 shrink-on-load 在解码阶段即降采样，
// 峰值内存受限、天然抵御解压炸弹；原生应用 EXIF Orientation（自动旋转并清除方向标记），
// 最终导出为有损 WebP 并剥离元数据（顺带去掉预览里的 EXIF GPS 等）。
//
// 语义与纯 Go 后端(preview.go)一致：等比缩放到最长边 ≤maxSide、不放大、quality(1..100)。
// 需以 `-tags libvips` 构建，构建/运行主机须装 libvips-dev。
func MakeWebpPreview(data []byte, maxSide, quality int) ([]byte, error) {
	if maxSide <= 0 {
		maxSide = defaultPreviewMaxSide
	}
	if quality <= 0 || quality > 100 {
		quality = defaultPreviewQuality
	}
	StartupVips()
	// InterestingNone = 只缩放不裁剪、纳入 maxSide×maxSide 框内保持纵横比；SizeDown = 绝不放大。
	ref, err := vips.NewThumbnailWithSizeFromBuffer(data, maxSide, maxSide, vips.InterestingNone, vips.SizeDown)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrImageDecode, err)
	}
	defer ref.Close()

	params := vips.NewWebpExportParams()
	params.Quality = quality
	params.StripMetadata = true
	out, _, err := ref.ExportWebp(params)
	if err != nil {
		return nil, fmt.Errorf("img: export webp: %w", err)
	}
	return out, nil
}
