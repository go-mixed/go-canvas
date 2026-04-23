package render

import (
	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/ti"
	"github.com/pkg/errors"
)

type ImageSprite struct {
	*Sprite
}

var _ ISprite = (*ImageSprite)(nil)

// NewImageSprite 从图片文件创建图片精灵，请勿设置rect的W、H（可以传入0，设置了也会被覆盖），表示使用图片的宽、高
func NewImageSprite(parent IParent, attribute *ctypes.Attribute, filePath string) (*ImageSprite, error) {
	texture, err := ti.LoadImageToTiImage(parent.Renderer().Runtime(), filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot load image to taichi")
	}
	shape := texture.Shape()
	w, h := shape[0], shape[1]
	// 自适应宽高
	if attribute.Width() == 0 {
		attribute.SetWidth(int(w))
	}
	if attribute.Height() == 0 {
		attribute.SetHeight(int(h))
	}

	// Resize Image if image.w/h != attribute.w/h
	if attribute.Width() != int(w) || attribute.Height() != int(h) {
		nW, nH := ti.CalcResizeWH(int(w), int(h), attribute.Width(), attribute.Height(), attribute.ResizeOptions())
		newTexture, err := ti.NewTiImage(parent.Renderer().Runtime(), uint32(nW), uint32(nH))
		if err != nil {
			texture.Release()
			return nil, errors.Wrapf(err, "Cannot create new texture while image resizing")
		}

		parent.Renderer().Module().Resize(texture, newTexture, attribute.ResizeOptions(), ctypes.Rectangle[int]{}, ctypes.Rectangle[int]{})
		texture.Release()
		texture = newTexture
		attribute.SetWH(nW, nH)
	}

	return BuildSprite(parent, attribute, texture, func(s *Sprite) (*ImageSprite, error) {
		return &ImageSprite{
			Sprite: s,
		}, nil
	})
}
