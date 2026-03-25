package render

import (
	"github.com/go-mixed/go-canvas/ti"
	"github.com/pkg/errors"
)

type ImageSprite struct {
	ISprite
}

var _ ISprite = (*ImageSprite)(nil)

// NewImageSprite 从图片文件创建图片精灵
func NewImageSprite(renderer *Renderer, filePath string) (ISprite, error) {
	texture, err := ti.LoadImageToTiImage(renderer.Runtime(), filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot load image to taichi")
	}

	s := NewNonContainerSprite(renderer, texture)

	return &ImageSprite{
		ISprite: s,
	}, nil
}
