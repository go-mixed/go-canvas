package render

import (
	"github.com/go-mixed/go-canvas/ti"
	"github.com/pkg/errors"
)

type ImageSprite struct {
	*Sprite
}

var _ ISprite = (*ImageSprite)(nil)

// NewImageSprite 从图片文件创建图片精灵，请勿设置rect的W、H（可以传入0，设置了也会被覆盖），表示使用图片的宽、高
func NewImageSprite(parent IParent, attribute *ti.Attribute, filePath string) (*ImageSprite, error) {
	texture, err := ti.LoadImageToTiImage(parent.Renderer().Runtime(), filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot load image to taichi")
	}
	shape := texture.Shape()
	w, h := shape[0], shape[1]
	attribute.SetWH(int(w), int(h))

	return BuildSprite(parent, attribute, texture, func(s *Sprite) (*ImageSprite, error) {
		return &ImageSprite{
			Sprite: s,
		}, nil
	})
}
