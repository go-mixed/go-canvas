package render

import (
	"github.com/go-mixed/go-canvas/font"
	"github.com/go-mixed/go-canvas/ti"
)

var _ ISprite = (*TextSprite)(nil)

type TextSprite struct {
	*Sprite
	richText *font.RichText

	originalWidth, originalHeight uint32
}

// NewTextSprite 创建文字精灵
func NewTextSprite(parent IParent, fontLibrary *font.FontLibrary, width, height uint32, opts ...font.RichTextOptionFn) (*TextSprite, error) {

	rt := font.BuildRichTextLines(fontLibrary, opts...)

	return BuildSprite(parent, nil, func(s *Sprite) (*TextSprite, error) {
		ts := &TextSprite{
			Sprite:         s,
			richText:       rt,
			originalWidth:  width,
			originalHeight: height,
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
		imgW, imgH := uint32(img.Bounds().Dx()), uint32(img.Bounds().Dy())
		var width, height = imgW, imgH
		if s.originalWidth == 0 {
			width = imgW
		}
		if s.originalHeight == 0 {
			height = imgH
		}

		var texture *ti.TiImage
		texture, err = ti.NewTiImage(s.Renderer().Runtime(), width, height)
		if err != nil {
			return
		}

		// 加上裁切代码
		var imgOffset ti.Point[int]
		switch s.richText.Align().HAlign {
		case ti.HAlignCenter:
			imgOffset.X = (int(width) - int(imgW)) / 2
		case ti.HAlignRight:
			imgOffset.X = int(width) - int(imgW)
		default:
		}

		switch s.richText.Align().VAlign {
		case ti.VAlignMiddle:
			imgOffset.Y = (int(height) - int(imgH)) / 2
		case ti.VAlignBottom:
			imgOffset.Y = int(height) - int(imgH)
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
		s.rect.Max = ti.Point[float32]{
			X: float32(width) + s.rect.Min.X,
			Y: float32(height) + s.rect.Min.Y,
		}

	}, func() bool {
		return s.richText.Equal(text)
	})

	return err
}
