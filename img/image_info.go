//go:build !libvips

package img

import (
	"image"
	"os"
	"path/filepath"

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
		return nil, err
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
