package render

import (
	"github.com/go-mixed/go-canvas/font"
	"github.com/go-mixed/go-canvas/ti"
)

var _ ISprite = (*TextSprite)(nil)

type TextSprite struct {
	*Sprite
	richText *font.RichText
	align    ti.Align
}

// NewTextSprite 创建文字精灵
func NewTextSprite(renderer *Renderer, text string, w, h uint32, align ti.Align) (ISprite, error) {
	rt := font.BuildRichTextLines(text)
	ts := &TextSprite{
		richText: rt,
		align:    align,
	}

	img, imgW, imgH := rt.RenderText(w, h, align)

	texture, err := ti.NewTiImage(renderer.runtime, uint32(imgW), uint32(imgH))
	if err != nil {
		return nil, err
	}

	if err = ti.UploadImageToTexture(texture, img); err != nil {
		texture.Release()
		return nil, err
	}

	sprite := NewSprite(renderer, texture)
	ts.Sprite = sprite

	return ts, nil
}

// SetText 设置文字内容并重新渲染
func (s *TextSprite) SetText(text string) {
	s.lockForUpdate(func() {
		s.richText = font.BuildRichTextLines(text)
	}, func() bool {
		return s.richText.Equal(text)
	})
}
