package font

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ti"
)

type richTextOptions struct {
	align     ti.Align
	fontStyle RichTextFontStyle
}

type RichTextOptionFn func(options *richTextOptions)

func WithVerticalAlign(vAlign ti.VerticalAlign) RichTextOptionFn {
	return func(opts *richTextOptions) {
		opts.align.VAlign = vAlign
	}
}

func WithAlign(hAlign ti.HorizontalAlign, vAlign ti.VerticalAlign) RichTextOptionFn {
	return func(opts *richTextOptions) {
		opts.align.HAlign = hAlign
		opts.align.VAlign = vAlign
	}
}

func WithBold(bold bool) RichTextOptionFn {
	return func(opts *richTextOptions) {
		opts.fontStyle.Bold = bold
	}
}

func WithItalic(italic bool) RichTextOptionFn {
	return func(opts *richTextOptions) {
		opts.fontStyle.Italic = italic
	}
}

func WithFontSize(fontSize int) RichTextOptionFn {
	return func(opts *richTextOptions) {
		opts.fontStyle.FontSize = fontSize
	}
}

func WithFontFamily(fontFamily string) RichTextOptionFn {
	return func(opts *richTextOptions) {
		opts.fontStyle.FontFamily = fontFamily
	}
}

func WithFontColor(color color.Color) RichTextOptionFn {
	return func(opts *richTextOptions) {
		opts.fontStyle.Color = color
	}
}

func WithFontStyle(font RichTextFontStyle) RichTextOptionFn {
	return func(opts *richTextOptions) {
		opts.fontStyle = font
	}
}
