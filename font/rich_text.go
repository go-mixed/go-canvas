package font

import (
	"image"
	"time"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
	"golang.org/x/image/font"
)

type RichText struct {
	fontLibrary *FontLibrary
	original    string
	lines       *misc.List[TextSegments]

	opts *RichTextOptions

	wrapScratch wrapScratch
	// Fake italic buffers are owned by RichText and reused across segments
	// to avoid per-segment allocations during render.
	// - italicAlphaBuf: default non-emoji fake italic path (1 byte/pixel).
	// - italicRGBABuf: emoji fallback path for better glyph compatibility.
	italicRGBABuf  *image.RGBA
	italicAlphaBuf *image.Alpha

	maxWidth, maxHeight int // 约束：换行宽度与最大渲染高度
	width, height       int // 缓存内容宽度和高度，避免重复计算

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
	if opts == nil {
		opts = RTOpt()
	}
	return &RichText{
		fontLibrary: fs,
		lines:       misc.NewList[TextSegments](),
		maxWidth:    opts.width,
		maxHeight:   opts.height,
		width:       -1,
		height:      -1,
		opts:        opts,
	}
}

func (r *RichText) SetText(input string) {
	setTextStart := time.Now()
	r.timing = RichTextTiming{}

	maxWidth := r.maxWidth

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
	if maxWidth <= 0 {
		// 无宽度限制：不自动换行，直接使用展开后的 segments。
		// Unlimited width: skip wrapping and keep expanded segments as-is.
		wrapped = expanded
	} else {
		// 如果只有一行，就快速返回
		if fast, ok := r.fastPathNoWrap(expanded, maxWidth); ok {
			wrapped = fast
		}
		if wrapped == nil {
			switch r.opts.wrapAlgo {
			case WrapAlgorithmFirstFit:
				wrapped = r.wrapFirstFit(expanded, maxWidth, r.opts.breakMode)
			default:
				wrapped = r.wordWrap(expanded, maxWidth, r.opts.breakMode)
			}
		}
	}
	r.timing.Wrap = time.Since(wrapStart)

	var line TextSegments
	for _, seg := range wrapped {
		if seg.BreakLine {
			if len(line) > 0 {
				r.lines.PushBack(coalesceLineSegments(line))
				line = nil
			}
			continue
		}
		line = append(line, seg)
	}

	if len(line) > 0 {
		r.lines.PushBack(coalesceLineSegments(line))
	}

	measureStart := time.Now()
	r.measure()
	r.timing.Measure = time.Since(measureStart)
	r.timing.SetText = time.Since(setTextStart)

}

// coalesceLineSegments 合并同一行内样式完全一致的相邻 segment，减少 DrawString 调用次数。
// coalesceLineSegments merges adjacent segments with identical style in one line to reduce DrawString calls.
func coalesceLineSegments(line TextSegments) TextSegments {
	if len(line) <= 1 {
		return line
	}
	out := make(TextSegments, 0, len(line))
	for _, seg := range line {
		if seg == nil || seg.Text == "" {
			continue
		}
		n := len(out)
		if n == 0 {
			out = append(out, seg)
			continue
		}
		last := out[n-1]
		if !last.CanMergeAdjacent(seg) {
			out = append(out, seg)
			continue
		}

		last.Text += seg.Text
		last.baseWidth += seg.baseWidth
		last.Width += seg.Width
		if seg.Height > last.Height {
			last.Height = seg.Height
		}
		if (seg.metrics.Ascent + seg.metrics.Descent).Ceil() > (last.metrics.Ascent + last.metrics.Descent).Ceil() {
			last.metrics = seg.metrics
		}
		// Width/height 已在合并时更新，避免后续重复测量。
		// Width/height are already updated during merge; keep measured to avoid re-measure.
		last.measured = true
	}
	return out
}

// fastPathNoWrap 快速判断是否可整行直过，避免进入 wrap 分段流程。
// fastPathNoWrap checks whether all segments fit in one line and can skip wrapping.
func (r *RichText) fastPathNoWrap(in TextSegments, maxWidth int) (TextSegments, bool) {
	lineWidth := 0
	for _, seg := range in {
		if seg == nil {
			continue
		}
		if seg.BreakLine {
			return nil, false
		}
		if seg.Text == "" {
			continue
		}

		face := r.fontLibrary.GetFace(seg.Font, seg.FontSize)
		if face == nil {
			return nil, false
		}

		if !seg.measured || seg.baseWidth <= 0 {
			baseWidth := font.MeasureString(face, seg.Text).Ceil()
			applySegmentMeasureWithBase(seg, face, baseWidth)
		}

		lineWidth += seg.Width
		if lineWidth > maxWidth {
			return nil, false
		}
	}
	return in, true
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
	if r.width >= 0 {
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
	if r.height >= 0 {
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
