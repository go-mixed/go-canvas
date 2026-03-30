package font

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ti"
)

type RichTextOptions struct {
	align     ti.Align
	fontStyle RichTextFontStyle
}

func RTOpt() *RichTextOptions {
	return &RichTextOptions{
		align: ti.Align{
			HAlign: ti.HAlignCenter,
			VAlign: ti.VAlignMiddle,
		},
		fontStyle: RichTextFontStyle{
			FontFamily: "sans-serif",
			FontSize:   16,
			Color:      color.Black,
			Bold:       false,
		},
	}
}

func (r *RichTextOptions) SetVerticalAlign(vAlign ti.VerticalAlign) *RichTextOptions {
	r.align.VAlign = vAlign
	return r
}

func (r *RichTextOptions) SetAlign(hAlign ti.HorizontalAlign, vAlign ti.VerticalAlign) *RichTextOptions {
	r.align.HAlign = hAlign
	r.align.VAlign = vAlign
	return r
}

func (r *RichTextOptions) SetBold(bold bool) *RichTextOptions {
	r.fontStyle.Bold = bold
	return r
}

func (r *RichTextOptions) SetItalic(italic bool) *RichTextOptions {
	r.fontStyle.Italic = italic
	return r
}

func (r *RichTextOptions) SetFontSize(fontSize int) *RichTextOptions {
	r.fontStyle.FontSize = fontSize
	return r
}

func (r *RichTextOptions) SetFontFamily(fontFamily string) *RichTextOptions {
	r.fontStyle.FontFamily = fontFamily
	return r
}

func (r *RichTextOptions) SetFontColor(color color.Color) *RichTextOptions {
	r.fontStyle.Color = color
	return r
}

func (r *RichTextOptions) SetFontStyle(font RichTextFontStyle) *RichTextOptions {
	r.fontStyle = font
	return r
}
