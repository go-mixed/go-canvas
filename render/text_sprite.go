package render

import (
	"github.com/go-mixed/go-canvas/font"
	"github.com/go-mixed/go-canvas/ti"
)

var _ ISprite = (*TextSprite)(nil)

type TextSprite struct {
	*Sprite
	richText *font.RichText

	originalWidth, originalHeight int
}

// NewTextSprite 创建文字精灵
func NewTextSprite(parent IParent, fontLibrary *font.FontLibrary, attribute *ti.Attribute, opts *font.RichTextOptions) (*TextSprite, error) {

	rt := font.BuildRichTextLines(fontLibrary, opts)

	return BuildSprite(parent, attribute, nil, func(s *Sprite) (*TextSprite, error) {
		ts := &TextSprite{
			Sprite:         s,
			richText:       rt,
			originalWidth:  attribute.Width(),
			originalHeight: attribute.Height(),
		}
		return ts, nil
	})
}

// SetText 设置文字内容并重新渲染
func (s *TextSprite) SetText(text string) error {
	var err error
	s.LockForUpdate(func() {
		s.richText.SetText(text)

		img := s.richText.RenderText()
		imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()
		var width, height = s.originalWidth, s.originalHeight
		if s.originalWidth == 0 {
			width = imgW
		}
		if s.originalHeight == 0 {
			height = imgH
		}

		var texture *ti.TiImage
		texture, err = ti.NewTiImage(s.Renderer().Runtime(), uint32(width), uint32(height))
		if err != nil {
			return
		}

		// 加上裁切代码
		var imgOffset ti.Point[int]
		switch s.richText.Align().HAlign {
		case ti.HAlignCenter:
			imgOffset.X = (width - imgW) / 2
		case ti.HAlignRight:
			imgOffset.X = width - imgW
		default:
		}

		switch s.richText.Align().VAlign {
		case ti.VAlignMiddle:
			imgOffset.Y = ((height) - (imgH)) / 2
		case ti.VAlignBottom:
			imgOffset.Y = (height) - (imgH)
		default:
		}

		if err = ti.UploadImageToTexture(texture, img, imgOffset); err != nil {
			texture.Release()
			return
		}

		if s.texture != nil {
			s.texture.Release()
		}
		s.texture = texture
		s.attribute.SetWH(width, height)
		s.attribute.SetCx(width / 2)
		s.attribute.SetCy(height / 2)

	}, func() bool {
		return !s.richText.Equal(text)
	})

	return err
}
