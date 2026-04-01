package font

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ti"
)

type LineBreakPolicy uint8

const (
	LineBreakWhenNecessary LineBreakPolicy = iota
	LineBreakNever
	LineBreakAlways
)

type WrapAlgorithm uint8

const (
	WrapAlgorithmSmart WrapAlgorithm = iota
	WrapAlgorithmFirstFit
)

type RichTextOptions struct {
	align     ti.Align
	fontStyle RichTextFontStyle
	breakMode LineBreakPolicy
	wrapAlgo  WrapAlgorithm
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
			Underline:  false,
		},
		breakMode: LineBreakWhenNecessary,
		wrapAlgo:  WrapAlgorithmSmart,
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

func (r *RichTextOptions) SetUnderline(underline bool) *RichTextOptions {
	r.fontStyle.Underline = underline
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

func (r *RichTextOptions) SetLineBreakPolicy(mode LineBreakPolicy) *RichTextOptions {
	r.breakMode = mode
	return r
}

func (r *RichTextOptions) SetWrapAlgorithm(algo WrapAlgorithm) *RichTextOptions {
	r.wrapAlgo = algo
	return r
}
