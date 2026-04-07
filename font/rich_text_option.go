package font

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ti"
)

type LineBreakPolicy uint8

const (
	// LineBreakNormal 正常换行：仅在超宽时断行，优先语义/合法断点（类似 CSS word-break: normal）。
	// LineBreakNormal wraps only when necessary, preferring semantic/legal breakpoints.
	LineBreakNormal LineBreakPolicy = iota
	// LineBreakNoWrap 不自动换行（类似 CSS white-space: nowrap）。
	// LineBreakNoWrap disables auto wrapping.
	LineBreakNoWrap
	// LineBreakAnywhere 可在任意 cluster 边界断行（类似 CSS overflow-wrap:anywhere）。
	// LineBreakAnywhere allows break at any cluster boundary.
	LineBreakAnywhere
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
	width     int
	height    int
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
		breakMode: LineBreakNormal,
		wrapAlgo:  WrapAlgorithmSmart,
		width:     0,
		height:    0,
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

// SetWidth 设置可用行宽。0 表示无限宽，不自动换行。
// SetWidth sets line wrap width. 0 means unlimited width (no auto wrap).
func (r *RichTextOptions) SetWidth(width int) *RichTextOptions {
	if width < 0 {
		width = 0
	}
	r.width = width
	return r
}

// SetHeight 设置渲染高度上限。0 表示不限制高度。
// SetHeight sets render height limit. 0 means unlimited height.
func (r *RichTextOptions) SetHeight(height int) *RichTextOptions {
	if height < 0 {
		height = 0
	}
	r.height = height
	return r
}

// SetSize 同时设置宽高限制。
// SetSize sets both width and height limits.
func (r *RichTextOptions) SetSize(width, height int) *RichTextOptions {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	r.width = width
	r.height = height
	return r
}
