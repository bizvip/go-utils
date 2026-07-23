package img

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ImageInfo struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Format    string `json:"format"`
	Ext       string `json:"ext"`
	Size      int64  `json:"size"`
	FileName  string `json:"filename"`
	FileMD5   string `json:"file_md5"`
	Animated  bool   `json:"animated"`
	FrameNum  int    `json:"frame_num"`
	LoopCount int    `json:"loop_count"`
	// ColorSpace/HasAlpha/BitDepth 为描述性元数据，标签集因后端而异：纯 Go 后端可区分
	// NRGBA/YCbCr 等精确 color model；libvips 后端只见解码后的 interpretation（如 jpeg → RGB）。
	// 仅供展示/记录，不要用于业务判断。
	ColorSpace string `json:"color_space,omitempty"`
	HasAlpha   bool   `json:"has_alpha,omitempty"`
	BitDepth   int    `json:"bit_depth,omitempty"`
}

type FitMode int

const (
	FitStretch FitMode = iota
	FitContain
	FitCover
)

type ResizeOptions struct {
	Src      string
	SavePath string
	Width    int
	Height   int
	Quality  int
	Fit      FitMode
	// Format overrides the output format. Empty value reuses the source format.
	// Supported: "jpeg", "png", "gif", "bmp", "webp".
	Format string
}

var (
	ErrEmptySource     = errors.New("img: source path must be provided")
	ErrInvalidSize     = errors.New("img: width and height must both be zero (keep) or one/both positive")
	ErrUnknownFormat   = errors.New("img: unknown image format")
	ErrUnsupportedSink = errors.New("img: unsupported output format")
	// ErrImageDecode 图像解码失败（损坏 / 不支持的格式）。
	ErrImageDecode = errors.New("img: decode image failed")
	// ErrImageTooLarge 源图像素数超过预算（纯 Go 后端全解码前的护栏；libvips 后端靠 shrink-on-load）。
	ErrImageTooLarge = errors.New("img: image exceeds max pixel budget")
)

// 预览/缩略图默认参数（MakeWebpPreview 入参为 0/越界时回退）。
const (
	defaultPreviewMaxSide = 320
	defaultPreviewQuality = 80
)

func ImageToBase64(imgPath string, withDataURI bool) (string, error) {
	data, err := os.ReadFile(imgPath)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	if !withDataURI {
		return encoded, nil
	}
	mime := http.DetectContentType(data)
	if !strings.HasPrefix(mime, "image/") {
		mime = "application/octet-stream"
	}
	return "data:" + mime + ";base64," + encoded, nil
}

func Base64ToFile(base64Str, fullPath string, overwrite bool) error {
	if strings.HasPrefix(base64Str, "data:") {
		if comma := strings.IndexByte(base64Str, ','); comma >= 0 {
			base64Str = base64Str[comma+1:]
		}
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(base64Str))
	if err != nil {
		return err
	}
	if mime := http.DetectContentType(data); !strings.HasPrefix(mime, "image/") {
		return fmt.Errorf("img: decoded base64 is not an image (%s)", mime)
	}
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(fullPath); err == nil {
		if !overwrite {
			return os.ErrExist
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.WriteFile(fullPath, data, 0o644)
}

func detectFormatByBytes(data []byte) string {
	mime := http.DetectContentType(data)
	switch mime {
	case "image/jpeg":
		return "jpeg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	case "image/bmp":
		return "bmp"
	case "image/webp":
		return "webp"
	}
	return ""
}

func detectFormatByHead(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	head := make([]byte, 512)
	n, _ := f.Read(head)
	if n == 0 {
		return "", ErrUnknownFormat
	}
	format := detectFormatByBytes(head[:n])
	if format == "" {
		return "", ErrUnknownFormat
	}
	return format, nil
}

func computeTargetSize(srcW, srcH, dstW, dstH int, fit FitMode) (int, int) {
	if dstW <= 0 && dstH <= 0 {
		return srcW, srcH
	}
	if dstW <= 0 {
		ratio := float64(dstH) / float64(srcH)
		return max1(int(float64(srcW)*ratio + 0.5)), dstH
	}
	if dstH <= 0 {
		ratio := float64(dstW) / float64(srcW)
		return dstW, max1(int(float64(srcH)*ratio + 0.5))
	}
	switch fit {
	case FitStretch:
		return dstW, dstH
	case FitContain:
		ratio := minF(float64(dstW)/float64(srcW), float64(dstH)/float64(srcH))
		return max1(int(float64(srcW)*ratio + 0.5)), max1(int(float64(srcH)*ratio + 0.5))
	case FitCover:
		ratio := maxF(float64(dstW)/float64(srcW), float64(dstH)/float64(srcH))
		return max1(int(float64(srcW)*ratio + 0.5)), max1(int(float64(srcH)*ratio + 0.5))
	}
	return dstW, dstH
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func max1(v int) int {
	if v < 1 {
		return 1
	}
	return v
}

func writeFileAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".img-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}
