package ti

import (
	"image"

	"github.com/go-mixed/go-canvas/internel/misc"
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

func NewBgraImage(runtime *taichi.Runtime, width, height uint32) (*TiImage, error) {
	return taichi.NewNdArray2D(runtime, width, height, taichi.DataTypeU32)
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

	err = UploadImageToTexture(texture, img, Point[int]{})

	if err != nil {
		texture.Release()
		return nil, errors.Wrapf(err, "Cannot upload image to taichi texture")
	}

	return texture, nil
}

// UploadImageToTexture 将 image.Image 上传到 TiImage
func UploadImageToTexture(texture *TiImage, img image.Image, imgOffset Point[int]) error {
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	shape := texture.Shape()
	shapeWidth, shapeHeight := int(shape[0]), int(shape[1])

	err := texture.MapFloat32(func(data []float32) error {
		for y := 0; y < shapeHeight; y++ {
			for x := 0; x < shapeWidth; x++ {

				imgX, imgY := x-imgOffset.X, y-imgOffset.Y
				if imgX < 0 || imgX >= imgWidth || imgY < 0 || imgY >= imgHeight {
					continue
				}

				r, g, b, a := ExpandF32Color(img.At(imgX, imgY))
				idx, _ := texture.GetOffset(x, y)
				data[idx] = r
				data[idx+1] = g
				data[idx+2] = b
				data[idx+3] = a
			}
		}
		return nil
	})
	return errors.Wrapf(err, "Cannot upload image to taichi texture")
}

// SaveTiImageToFile 将 Taichi.NdArray(w, h, (r, g, b, a)) 保存到图片文件
func SaveTiImageToFile(texture *TiImage, filePath string) error {
	shape := texture.Shape()
	width, height := int(shape[0]), int(shape[1])
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	err := DownloadTextureToImage(texture, img)

	if err != nil {
		return errors.Wrapf(err, "Cannot download taichi texture to image")
	}
	return misc.SaveImage(img, filePath)
}

func DownloadTextureToImage(texture *TiImage, img ImageWriter) error {
	shape := texture.Shape()
	width, height := int(shape[0]), int(shape[1])

	err := texture.MapFloat32(func(data []float32) error {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx, _ := texture.GetOffset(x, y)
				c := TiColorToColor(data[idx], data[idx+1], data[idx+2], data[idx+3])
				img.Set(x, y, c)
			}
		}
		//})
		return nil
	})

	return errors.Wrapf(err, "Cannot download taichi texture to image")
}

func CalcResizeWH(originalWidth, originalHeight int, targetWidth, targetHeight int, opts ResizeOptions) (newWidth int, newHeight int) {
	if originalWidth == 0 || originalHeight == 0 {
		return targetWidth, targetHeight
	}

	srcW, srcH := float32(originalWidth), float32(originalHeight)
	dstW, dstH := float32(targetWidth), float32(targetHeight)

	scaleX := dstW / srcW
	scaleY := dstH / srcH

	var scale float32
	switch opts.FillMode {
	case FillModeStretch:
		scale = 1 // 不缩放，用目标尺寸
	case FillModeFit:
		scale = min(scaleX, scaleY)
	case FillModeFill:
		scale = max(scaleX, scaleY)
	}

	newWidth = int(srcW * scale)
	newHeight = int(srcH * scale)

	return newWidth, newHeight
}
