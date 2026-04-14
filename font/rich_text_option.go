package font

import (
	"image/color"

	"github.com/go-mixed/go-canvas/internel/misc"
	"github.com/go-mixed/go-canvas/ti"
)

type WordWrapMode uint8

const (
	// BreakNormal 正常换行：仅在超宽时断行，优先语义/合法断点（类似 CSS word-break: normal）。
	// BreakNormal wraps only when necessary, preferring semantic/legal breakpoints.
	BreakNormal WordWrapMode = iota
	// NoWrap 不自动换行（类似 CSS white-space: nowrap）。
	// NoWrap disables auto wrapping.
	NoWrap
	// BreakAll 可在任意 cluster 边界断行（类似 CSS overflow-wrap:anywhere）。
	// BreakAll allows break at any cluster boundary.
	BreakAll
)

type WordWrapAlgorithm uint8

const (
	WrapAlgorithmSmart WordWrapAlgorithm = iota
	WrapAlgorithmFirstFit
)

// BidiBaseDirection 定义 BiDi 段落基础方向。
// BidiBaseDirection defines paragraph base direction for BiDi reordering.
type BidiBaseDirection uint8

const (
	// BidiBaseAuto 自动根据首个强方向字符决定段落方向。
	// BidiBaseAuto auto-detects paragraph direction from strong characters.
	BidiBaseAuto BidiBaseDirection = iota
	// BidiBaseLTR 强制段落基础方向为左到右。
	// BidiBaseLTR forces paragraph base direction to left-to-right.
	BidiBaseLTR
	// BidiBaseRTL 强制段落基础方向为右到左。
	// BidiBaseRTL forces paragraph base direction to right-to-left.
	BidiBaseRTL
)

type RichTextOptions struct {
	align     ti.Align
	fontStyle RichTextFontStyle
	breakMode LineBreakPolicy
	wrapAlgo  WrapAlgorithm
	width     int
	height    int
	align        ti.Align
	fontStyle    RichTextFontStyle
	wordWrapMode WordWrapMode
	wordWrapAlgo WordWrapAlgorithm
	logger       misc.Logger

	width  int
	height int
}

func RTOpt() *RichTextOptions {
	return &RichTextOptions{
		align: ti.Align{
			HAlign: ti.HAlignCenter,
			VAlign: ti.VAlignMiddle,
		},
		fontStyle: RichTextFontStyle{
			FontFamily: "",
			FontSize:   16,
			Color:      color.Black,
			Bold:       false,
			Underline:  false,
		},
		breakMode: LineBreakNormal,
		wrapAlgo:  WrapAlgorithmSmart,
		width:     0,
		height:    0,
		wordWrapMode: BreakNormal,
		wordWrapAlgo: WrapAlgorithmSmart,
		width:        0,
		height:       0,
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

func (r *RichTextOptions) SetWordWrap(mode WordWrapMode) *RichTextOptions {
	r.wordWrapMode = mode
	return r
}

func (r *RichTextOptions) SetWrapAlgorithm(algo WordWrapAlgorithm) *RichTextOptions {
	r.wordWrapAlgo = algo
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
