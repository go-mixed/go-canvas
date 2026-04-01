package font

import (
	"image/color"
	"math"
	"time"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
	"golang.org/x/image/font"
)

// TextSegment 单个文本片段
type TextSegment struct {
	Text       string
	Font       *FontInfo
	FontSize   int
	Color      color.Color
	Bold       FontWeight
	BreakLine  bool
	Italic     bool
	Underline  bool
	FakeItalic bool
	FontFamily string
	baseWidth  int
	Width      int
	Height     int
	metrics    font.Metrics
	measured   bool
}

func (t *TextSegment) CopyWithText(text string) *TextSegment {
	newSegment := *t
	newSegment.Text = text
	// 换行标记只由输入 text 决定，避免继承上一个 segment 的 BreakLine 状态。
	newSegment.BreakLine = text == "\n"

	// 文本变更后，测量相关字段必须失效，等待 wrap/measure 重新填充。
	newSegment.baseWidth = 0
	newSegment.Width = 0
	newSegment.Height = 0
	newSegment.metrics = font.Metrics{}
	newSegment.measured = false

	return &newSegment
}

func (t *TextSegment) MeasureString(face font.Face) (int, int) {
	segWidth := font.MeasureString(face, t.Text).Ceil()
	t.baseWidth = segWidth

	t.Width = segWidth
	// 使用 ascent + |descent| 作为高度，确保能完整渲染
	t.metrics = face.Metrics()
	t.Height = (t.metrics.Ascent + t.metrics.Descent).Ceil()
	if t.FakeItalic {
		t.Width += syntheticItalicExtraWidth(t.Height)
	}
	t.measured = true
	return segWidth, t.Height
}

func syntheticItalicExtraWidth(height int) int {
	if height <= 0 {
		return 0
	}
	// Roughly 14-15 degrees of slant for synthetic italic.
	return int(math.Ceil(float64(height) * 0.26))
}

type TextSegments []*TextSegment

// Height 返回该行最大字号的高度（Metrics.Height）
func (s TextSegments) Height() int {
	var maxHeight int
	for _, seg := range s {
		if maxHeight < seg.Height {
			maxHeight = seg.Height
		}
	}
	return maxHeight
}

// MaxMetrics 返回该行最大字号的 Metrics
func (s TextSegments) MaxMetrics() font.Metrics {
	var maxMetrics font.Metrics
	for _, seg := range s {
		if (maxMetrics.Ascent + maxMetrics.Descent).Ceil() < seg.Height {
			maxMetrics = seg.metrics
		}
	}
	return maxMetrics
}

// Width 总长度
func (s TextSegments) Width() int {
	var sumV int
	for _, seg := range s {
		sumV += seg.Width
	}
	return sumV
}

type RichText struct {
	fontLibrary *FontLibrary
	original    string
	lines       *misc.List[TextSegments]

	opts   *RichTextOptions
	matrix misc.Matrix

	wrapScratch wrapScratch

	width, height int // 缓存宽度和高度，避免重复计算

	timing RichTextTiming
}

type RichTextTiming struct {
	Parse   time.Duration
	Wrap    time.Duration
	Layout  time.Duration
	Measure time.Duration
	Render  time.Duration
	SetText time.Duration
}

// BuildRichTextLines 解析带标签的文字，返回文本片段列表
// 标签格式：<text bold italic color="#RRGGBBAA" font-size="15" font-family="Noto Sans CJK SC">文字</text>
func BuildRichTextLines(fs *FontLibrary, opts *RichTextOptions) *RichText {
	return &RichText{
		fontLibrary: fs,
		lines:       misc.NewList[TextSegments](),
		width:       -1,
		height:      -1,
		opts:        opts,
		matrix:      misc.IdentityMatrix(),
	}
}

func (r *RichText) SetText(input string) {
	setTextStart := time.Now()
	r.timing = RichTextTiming{}

	maxWidth := r.width
	if maxWidth < 0 {
		maxWidth = 0
	}

	r.lines.Clear()
	r.original = input
	r.width = -1
	r.height = -1

	parseStart := time.Now()
	segments := r.parseText(input)
	r.timing.Parse = time.Since(parseStart)
	if len(segments) == 0 {
		r.timing.SetText = time.Since(setTextStart)
		return
	}

	layoutStart := time.Now()
	expanded := make(TextSegments, 0, len(segments))
	for _, seg := range segments {
		r.ensureSegmentFontAndFace(seg)
		expanded = append(expanded, splitSegmentByNewline(seg)...)
	}
	r.timing.Layout = time.Since(layoutStart)

	wrapStart := time.Now()
	var wrapped TextSegments
	switch r.opts.wrapAlgo {
	case WrapAlgorithmFirstFit:
		wrapped = r.wrapFirstFit(expanded, maxWidth, r.opts.breakMode)
	default:
		wrapped = r.wordWrap(expanded, maxWidth, r.opts.breakMode)
	}
	r.timing.Wrap = time.Since(wrapStart)

	var line TextSegments
	for _, seg := range wrapped {
		if seg.BreakLine {
			if len(line) > 0 {
				r.lines.PushBack(line)
				line = nil
			}
			continue
		}
		line = append(line, seg)
	}

	if len(line) > 0 {
		r.lines.PushBack(line)
	}

	measureStart := time.Now()
	r.measure()
	r.timing.Measure = time.Since(measureStart)
	r.timing.SetText = time.Since(setTextStart)

}

// Len 返回文本段落的总行数
func (r *RichText) Len() int {
	return r.lines.Len()
}

// GetSegments 返回所有文本段落
func (r *RichText) GetSegments() TextSegments {
	segments := make(TextSegments, 0, r.lines.Len())
	for el := r.lines.Front(); el != nil; el = el.Next() {
		segments = append(segments, el.Value...)
	}
	return segments
}

func (r *RichText) Equal(text string) bool {
	return r.original == text
}

func (r *RichText) ensureSegmentFontAndFace(seg *TextSegment) {
	if seg.Font == nil {
		seg.Font = r.fontLibrary.MatchOrFeedback(seg.FontFamily, seg.Bold, seg.Italic)
		seg.FakeItalic = seg.Italic && !seg.Font.Italic
	}
	_ = r.fontLibrary.GetFace(seg.Font, seg.FontSize)
}

// measure 测量每个文本片段的宽度和高度
func (r *RichText) measure() {
	for _, segments := range r.lines.Range() {
		for _, seg := range segments {
			if seg.measured {
				continue
			}
			face := r.fontLibrary.GetFace(seg.Font, seg.FontSize)
			seg.MeasureString(face)
		}
	}
}

func (r *RichText) Width() int {
	if r.width > 0 {
		return r.width
	}

	var maxV int
	for _, segments := range r.lines.Range() {
		w := segments.Width()
		if maxV < segments.Width() {
			maxV = w
		}
	}
	r.width = maxV
	return maxV
}

func (r *RichText) Height() int {
	if r.height > 0 {
		return r.height
	}

	var totalHeight int
	for _, segments := range r.lines.Range() {
		totalHeight += segments.Height()
	}
	r.height = totalHeight
	return r.height
}

func (r *RichText) IsEmpty() bool {
	return r.lines.Len() == 0
}

func (r *RichText) Align() ti.Align {
	return r.opts.align
}

func (r *RichText) FontStyle() RichTextFontStyle {
	return r.opts.fontStyle
}

func (r *RichText) Timing() RichTextTiming {
	return r.timing
}

func splitSegmentByNewline(seg *TextSegment) TextSegments {
	var out TextSegments
	start := 0
	for i, rn := range seg.Text {
		if rn != '\n' {
			continue
		}
		out = append(out, seg.CopyWithText(seg.Text[start:i]))
		br := seg.CopyWithText("")
		br.BreakLine = true
		out = append(out, br)
		start = i + 1
	}
	if start <= len(seg.Text) {
		out = append(out, seg.CopyWithText(seg.Text[start:]))
	}
	return out
}
