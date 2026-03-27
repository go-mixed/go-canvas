package render

import (
	"github.com/go-mixed/go-canvas/ti"
	"github.com/pkg/errors"
)

type ImageSprite struct {
	*Sprite
}

var _ ISprite = (*ImageSprite)(nil)

// NewImageSprite 从图片文件创建图片精灵
func NewImageSprite(parent IParent, filePath string) (*ImageSprite, error) {
	texture, err := ti.LoadImageToTiImage(parent.Renderer().Runtime(), filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot load image to taichi")
	}

	return BuildSprite(parent, texture, func(s *Sprite) (*ImageSprite, error) {
		return &ImageSprite{
			Sprite: s,
		}, nil
	})
}
