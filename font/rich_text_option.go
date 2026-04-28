package font

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/internel/misc"
	xfont "golang.org/x/image/font"
)

type RichTextOptions struct {
	align        ctypes.Align
	fontStyle    RichTextFontStyle
	wordWrapMode ctypes.WordWrapMode
	wordWrapAlgo ctypes.WordWrapAlgorithm
	bidi         ctypes.BidiDirection
	logger       misc.Logger

	width  int
	height int
	// 行高，单位px
	lineHeight int
	// 行高比例，比如当前行的字体为20，如果lineHeightRatio为1.2时，则行高为20*1.2=24
	lineHeightRatio float32
}

func RTOpt() *RichTextOptions {
	return &RichTextOptions{
		align: ctypes.Align{
			HAlign: ctypes.HAlignCenter,
			VAlign: ctypes.VAlignMiddle,
		},
		fontStyle: RichTextFontStyle{
			FontFamily: "",
			FontSize:   16,
			Color:      color.Black,
			Weight:     xfont.WeightNormal,
			Underline:  false,
		},
		wordWrapMode:    ctypes.BreakNormal,
		wordWrapAlgo:    ctypes.WrapAlgorithmSmart,
		bidi:            ctypes.BidiAuto,
		width:           misc.NaNInt,
		height:          misc.NaNInt,
		lineHeight:      0,
		lineHeightRatio: 1,
	}
}

// SetVerticalAlign 设置垂直对齐方式。
func (r *RichTextOptions) SetVerticalAlign(vAlign ctypes.VerticalAlign) *RichTextOptions {
	r.align.VAlign = vAlign
	return r
}

// SetHorizontalAlign 设置水平对齐方式。
func (r *RichTextOptions) SetHorizontalAlign(hAlign ctypes.HorizontalAlign) *RichTextOptions {
	r.align.HAlign = hAlign
	return r
}

// SetAlign 同时设置水平、垂直对齐方式。
func (r *RichTextOptions) SetAlign(hAlign ctypes.HorizontalAlign, vAlign ctypes.VerticalAlign) *RichTextOptions {
	r.align.HAlign = hAlign
	r.align.VAlign = vAlign
	return r
}

func (r *RichTextOptions) SetWeight(weight xfont.Weight) *RichTextOptions {
	r.fontStyle.Weight = weight
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
	r.fontStyle.FontFamily = normalizeFamilyName(fontFamily)
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

func (r *RichTextOptions) SetLineHeight(lineHeight int) *RichTextOptions {
	r.lineHeight = lineHeight
	r.lineHeightRatio = 0
	return r
}

func (r *RichTextOptions) SetLineHeightScale(lineHeightScale float32) *RichTextOptions {

	r.lineHeightRatio = max(1, lineHeightScale)
	r.lineHeight = 0
	return r
}

func (r *RichTextOptions) SetWordWrap(mode ctypes.WordWrapMode) *RichTextOptions {
	r.wordWrapMode = mode
	return r
}

func (r *RichTextOptions) SetWrapAlgorithm(algo ctypes.WordWrapAlgorithm) *RichTextOptions {
	r.wordWrapAlgo = algo
	return r
}

// SetBidi 设置 BiDi 段落基础方向。
// SetBidi sets BiDi paragraph base direction.
func (r *RichTextOptions) SetBidi(bidi ctypes.BidiDirection) *RichTextOptions {
	r.bidi = bidi
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

// SetWidthIfNotDefined 仅当当前宽度为默认值时才设置。
func (r *RichTextOptions) SetWidthIfNotDefined(width int) *RichTextOptions {
	if misc.IsNaNInt(r.width) {
		r.SetWidth(width)
	}
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

// SetHeightIfNotDefined 仅当当前高度为默认值时才设置。
// SetHeightIfNotDefined sets height if current height is default.
func (r *RichTextOptions) SetHeightIfNotDefined(height int) *RichTextOptions {
	if misc.IsNaNInt(r.height) {
		r.SetHeight(height)
	}
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

// SetLogger 设置 RichText 日志器；传 nil 表示关闭日志输出。
// SetLogger sets RichText logger; nil disables logging.
func (r *RichTextOptions) SetLogger(logger misc.Logger) *RichTextOptions {
	r.logger = logger
	return r
}
