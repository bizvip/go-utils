//go:build !libvips

package img

import (
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"

	"github.com/bizvip/go-utils/os/fs"
)

func GetImageInfo(path string, withMd5 bool) (*ImageInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	info := &ImageInfo{
		Width:    cfg.Width,
		Height:   cfg.Height,
		Format:   format,
		Ext:      strings.ToLower(strings.TrimPrefix(filepath.Ext(path), ".")),
		Size:     stat.Size(),
		FileName: filepath.Base(path),
		FrameNum: 1,
	}

	if format == "gif" {
		if _, err := f.Seek(0, 0); err == nil {
			if g, err := gif.DecodeAll(f); err == nil {
				info.FrameNum = len(g.Image)
				info.Animated = len(g.Image) > 1
				info.LoopCount = g.LoopCount
			}
		}
	}

	if withMd5 {
		md5sum, err := fs.GetBigFileMd5(path)
		if err != nil {
			return nil, err
		}
		info.FileMD5 = md5sum
	}

	return info, nil
}
