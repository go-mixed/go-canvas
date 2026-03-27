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

func (s *TextSprite) RemoveMask(mask IMask) {
	//TODO implement me
	panic("implement me")
}

// NewTextSprite 创建文字精灵
func NewTextSprite(parent IParent, text string, width, height uint32, align ti.Align) (ISprite, error) {
	rt := font.BuildRichTextLines(text)

	img := rt.RenderText(align)
	imgW, imgH := uint32(img.Bounds().Dx()), uint32(img.Bounds().Dy())
	if width == 0 {
		width = imgW
	}
	if height == 0 {
		height = imgH
	}

	texture, err := ti.NewTiImage(parent.Renderer().Runtime(), width, height)
	if err != nil {
		return nil, err
	}

	// 加上裁切代码
	var imgOffset ti.Point[int]
	switch align.HAlign {
	case ti.HAlignCenter:
		imgOffset.X = (int(width) - int(imgW)) / 2
	case ti.HAlignRight:
		imgOffset.X = int(width) - int(imgW)
	default:
	}

	switch align.VAlign {
	case ti.VAlignMiddle:
		imgOffset.Y = (int(height) - int(imgH)) / 2
	case ti.VAlignBottom:
		imgOffset.Y = int(height) - int(imgH)
	default:
	}

	if err = ti.UploadImageToTexture(texture, img, imgOffset); err != nil {
		texture.Release()
		return nil, err
	}

	return BuildSprite(parent, texture, func(s *Sprite) (*TextSprite, error) {
		ts := &TextSprite{
			Sprite:   s,
			richText: rt,
			align:    align,
		}
		return ts, nil
	})
}

// SetText 设置文字内容并重新渲染
func (s *TextSprite) SetText(text string) {
	s.LockForUpdate(func() {
		s.richText = font.BuildRichTextLines(text)
	}, func() bool {
		return s.richText.Equal(text)
	})
}
