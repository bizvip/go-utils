package img

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/davidbyttow/govips/v2/vips"

	"tab/pkg/utils/fs"
)

type Info struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Format   string `json:"format"`
	Ext      string `json:"ext"`
	Size     int64  `json:"size"`
	FileName string `json:"filename"`
	FileMD5  string `json:"file_md5"`
}

func GetInfo(path string, withMd5 bool) (*Info, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		// 尝试使用 govips 读取图像信息
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
		fmt.Println(format)
		// 创建一个空的 image.RGBA 以避免返回 nil 图像
		img = &image.RGBA{Rect: image.Rect(0, 0, width, height)}
	}

	ext := filepath.Ext(path)
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileMd5 := ""
	if withMd5 {
		fileMd5, err = fs.GetFileMd5Stream(path)
		if err != nil {
			return nil, err
		}
	}

	info := &Info{
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
