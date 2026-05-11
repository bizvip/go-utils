//go:build libvips

package img

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/davidbyttow/govips/v2/vips"

	"github.com/bizvip/go-utils/os/fs"
)

var vipsOnce sync.Once

// StartupVips initializes libvips. Safe to call multiple times; only first call takes effect.
// Library users may call this explicitly at process startup; otherwise GetImageInfo
// will call it lazily on first use.
func StartupVips() {
	vipsOnce.Do(func() {
		vips.Startup(nil)
	})
}

// ShutdownVips releases libvips resources. Call once at process exit if you want
// a clean teardown. Repeated calls are safe but only the first has effect when
// paired with StartupVips.
func ShutdownVips() {
	vips.Shutdown()
}

func GetImageInfo(path string, withMd5 bool) (*ImageInfo, error) {
	StartupVips()

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	ref, err := vips.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}
	defer ref.Close()

	width := ref.Width()
	height := ref.Height()
	pages := ref.Pages()
	if pages <= 0 {
		pages = 1
	}
	pageHeight := ref.PageHeight()
	if pages > 1 && pageHeight > 0 {
		height = pageHeight
	}
	format := vips.ImageTypes[ref.Format()]

	info := &ImageInfo{
		Width:    width,
		Height:   height,
		Format:   format,
		Ext:      strings.ToLower(strings.TrimPrefix(filepath.Ext(path), ".")),
		Size:     stat.Size(),
		FileName: filepath.Base(path),
		FrameNum: pages,
		Animated: pages > 1,
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
