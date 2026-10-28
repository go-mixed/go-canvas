package render

import (
	"github.com/go-mixed/go-canvas/ctypes"
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
func NewTextSprite(parent IParent, fontLibrary *font.FontLibrary, attribute *ctypes.Attribute, opts *font.RichTextOptions) (*TextSprite, error) {

	// richText没有单独设置宽高，则使用attribute的宽高
	if attribute.Width() != 0 {
		opts.SetWidthIfNotDefined(attribute.Width())
	}
	if attribute.Height() != 0 {
		opts.SetHeightIfNotDefined(attribute.Height())
	}

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
		if s.richText.Equal(text) {
			return
		}

		s.richText.SetText(text)

		// 当上一次有文本、本次清空文本时，回收上一次生成的纹理和画布，避免 Texture() 返回残留像素。
		// 清空后不再保留纹理，Container.Render 会跳过没有纹理的精灵。
		// When clearing previously non-empty text, recycle the previously generated texture and canvas so Texture() cannot return stale pixels.
		// No texture is retained after clearing; Container.Render skips texture-less sprites.
		if s.richText.IsEmpty() {
			s.addGarbageTexture(s.texture)
			s.addGarbageTexture(s.canvas)
			s.texture = nil
			s.canvas = nil
			s.attribute.SetWH(s.originalWidth, s.originalHeight)
			s.attribute.SetCx(s.originalWidth / 2)
			s.attribute.SetCy(s.originalHeight / 2)
			return
		}

		img := s.richText.RenderText()
		imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()
		var width, height = s.originalWidth, s.originalHeight
		if s.originalWidth == 0 {
			width = imgW
		}
		if s.originalHeight == 0 {
			height = imgH
		}

		var texture *ctypes.TiImage
		texture, err = ti.NewTiImage(s.Renderer().Runtime(), uint32(width), uint32(height))
		if err != nil {
			return
		}

		// 加上裁切代码
		var imgOffset ctypes.Point[int]
		switch s.richText.Align().HAlign {
		case ctypes.HAlignCenter:
			imgOffset.X = (width - imgW) / 2
		case ctypes.HAlignRight:
			imgOffset.X = width - imgW
		default:
		}

		switch s.richText.Align().VAlign {
		case ctypes.VAlignMiddle:
			imgOffset.Y = ((height) - (imgH)) / 2
		case ctypes.VAlignBottom:
			imgOffset.Y = (height) - (imgH)
		default:
		}

		if err = ti.UploadImageToTexture(texture, img, imgOffset); err != nil {
			texture.Release()
			return
		}

		s.addGarbageTexture(s.texture)
		s.texture = texture
		s.attribute.SetWH(width, height)
		s.attribute.SetCx(width / 2)
		s.attribute.SetCy(height / 2)

	}, func() ctypes.DirtyMode {
		// TextSprite 比较特殊，修改文字之后 s.texture 已经修改，只用触发 DirtyModeComposite，
		// 如果还有Border之类的设置，才能触发 DirtyModeLayout、 DirtyModePainting
		if !s.richText.Equal(text) {
			return ctypes.DirtyModeComposite | s.attribute.Dirty()
		}
		return ctypes.DirtyModeNone
	})

	return err
}
