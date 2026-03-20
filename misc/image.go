package misc

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// LoadImage 从文件加载图片
func LoadImage(filePath string) (image.Image, error) {
	// 打开图片文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot open the image")
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot decode the image")
	}

	return img, nil
}

func SaveImage(img image.Image, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "Cannot create image file")
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".png":
		err = png.Encode(file, img)
	case ".jpg", ".jpeg":
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
	default:
		return errors.New("Unsupported image format")
	}

	if err != nil {
		return errors.Wrapf(err, "Cannot encode image")
	}

	return nil
}
