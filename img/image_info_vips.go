//go:build libvips

package img

import (
	"image"
	"os"
	"path/filepath"

	"github.com/davidbyttow/govips/v2/vips"

	"github.com/bizvip/go-utils/os/fs"
)

func GetImageInfo(path string, withMd5 bool) (*ImageInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		vips.Startup(nil)
		defer vips.Shutdown()

		ref, err := vips.NewImageFromFile(path)
		if err != nil {
			return nil, err
		}
		defer ref.Close()

		width := ref.Width()
		height := ref.Height()
		format = vips.ImageTypes[ref.Format()]

		img = &image.RGBA{Rect: image.Rect(0, 0, width, height)}
	}

	ext := filepath.Ext(path)
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileMd5 := ""
	if withMd5 {
		fileMd5, err = fs.GetBigFileMd5(path)
		if err != nil {
			return nil, err
		}
	}

	info := &ImageInfo{
		Width:    img.Bounds().Dx(),
		Height:   img.Bounds().Dy(),
		Format:   format,
		Ext:      ext,
		Size:     fileInfo.Size(),
		FileName: filepath.Base(path),
		FileMD5:  fileMd5,
	}

	return info, nil
}
