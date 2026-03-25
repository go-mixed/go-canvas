package ti

import (
	"image"

	"github.com/go-mixed/go-canvas/misc"

	"github.com/go-mixed/go-taichi/taichi"
	"github.com/pkg/errors"

	_ "image/jpeg" // 注册 JPEG 解码器
	_ "image/png"  // 注册 PNG 解码器
)

func NewTiImage(runtime *taichi.Runtime, width, height uint32) (*TiImage, error) {
	return taichi.NewNdArray2DWithElemShape(runtime, width, height, taichi.Shape(4), taichi.DataTypeF32)
}

func NewTiGrid(runtime *taichi.Runtime, width, height uint32) (*TiGrid, error) {
	return taichi.NewNdArray2D(runtime, width, height, taichi.DataTypeF32)
}

func NewTiMask(runtime *taichi.Runtime, width, height uint32) (*TiMask, error) {
	return NewTiGrid(runtime, width, height)
}

// LoadImageToTiImage 将图片加载到 Taichi.NdArray(w, h, (r, g, b, a))
func LoadImageToTiImage(rt *taichi.Runtime, filePath string) (*TiImage, error) {
	img, err := misc.LoadImage(filePath)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	texture, err := NewTiImage(
		rt,
		uint32(width),
		uint32(height),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot create taichi texture")
	}

	err = texture.MapFloat32(func(data []float32) error {
		//misc.ParallelForeach(height, 1, func(yStart, yEnd int) {
		yStart := 0
		yEnd := height
		for y := yStart; y < yEnd; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := ExpandFColor(img.At(x, y))
				idx := (x*height + y) * 4
				data[idx] = r
				data[idx+1] = g
				data[idx+2] = b
				data[idx+3] = a
			}
		}
		//})
		return nil
	})

	if err != nil {
		texture.Release()
		return nil, errors.Wrapf(err, "Cannot upload image to taichi texture")
	}
	return texture, nil
}

// SaveTiImageToFile 将 Taichi.NdArray(w, h, (r, g, b, a)) 保存到图片文件
func SaveTiImageToFile(texture *TiImage, filePath string) error {
	shape := texture.Shape()
	width, height := int(shape[0]), int(shape[1])
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	err := texture.MapFloat32(func(data []float32) error {
		//misc.ParallelForeach(height, 16, func(yStart, yEnd int) {
		yStart := 0
		yEnd := height

		for y := yStart; y < yEnd; y++ {
			for x := 0; x < width; x++ {
				idx := (x*height + y) * 4
				c := TiColor2Color(data[idx], data[idx+1], data[idx+2], data[idx+3])
				img.Set(x, y, c)
			}
		}
		//})
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "Cannot download taichi texture to image")
	}
	return misc.SaveImage(img, filePath)
}
