package font

import (
	"fmt"
	"image"
	"sort"
	"strings"
	"time"

	"github.com/go-mixed/go-canvas/internel/misc"
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
	// - italicAlphaBuf: fake italic path (1 byte/pixel).
	italicAlphaBuf *image.Alpha

	maxWidth, maxHeight int // 约束：换行宽度与最大渲染高度
	width, height       int // 缓存内容宽度和高度，避免重复计算
}

// BuildRichTextLines 解析带标签的文字，返回文本片段列表
// 标签格式：<text bold italic color="#RRGGBBAA" font-size="15" font-family="Noto Sans CJK SC">文字</text>
func BuildRichTextLines(fs *FontLibrary, opts *RichTextOptions) *RichText {
	if opts == nil {
		opts = RTOpt()
	}
	if opts.logger == nil {
		opts.SetLogger(fs.logger)
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

// logf 通过可选 logger 输出调试信息；未设置 logger 时静默。
// logf writes debug logs through optional logger; silent when logger is nil.
func (r *RichText) logf(format string, args ...any) {
	if r == nil || r.opts == nil || r.opts.logger == nil {
		return
	}
	r.opts.logger.Printf(format, args...)
}

// logfFn 延迟构造日志内容：仅在 logger 存在时才执行 fn。
// logfFn lazily builds log content only when logger is configured.
func (r *RichText) logfFn(fn func(logger misc.Logger)) {
	if r == nil || r.opts == nil || r.opts.logger == nil || fn == nil {
		return
	}
	fn(r.opts.logger)
}

func (r *RichText) SetText(input string) {
	t := time.Now()
	logStep := func(name string) {
		r.logf("[richtext] %s=%s", name, time.Since(t))
		t = time.Now()
	}

	maxWidth := r.maxWidth

	r.lines.Clear()
	r.original = input
	r.width = -1
	r.height = -1

	segments := r.parseText(input)
	logStep("parse")
	if len(segments) == 0 {
		return
	}

	expanded := make(TextSegments, 0, len(segments))
	for _, seg := range segments {
		parts := splitSegmentByNewline(seg)
		for _, part := range parts {
			if part == nil {
				continue
			}
			if part.BreakLine {
				expanded = append(expanded, part)
				continue
			}
			chunks := r.splitSegmentByFontCoverage(part)
			for _, p := range chunks {
				expanded = append(expanded, p)
			}
		}
	}
	logStep("split")

	for _, seg := range expanded {
		seg.Font.GetOpenTypeFont()
		r.fontLibrary.GetFace(seg.Font, seg.FontSize)
	}

	logStep("initial fonts & faces")

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
	logStep("wrap")

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

	r.measure()
	logStep("measure")

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
		// 只有两侧都已测量时，才可直接累加宽高并保持 measured=true。
		// 当任一侧未测量（如 BiDi 重排后的 CopyWithText 片段）时，必须失效测量，
		// 让后续 measure() 重新计算，避免宽度为 0 导致重叠。
		// Keep measured=true only when both sides are already measured.
		// If either side is unmeasured (e.g. BiDi reordered pieces), invalidate and re-measure later.
		if last.measured && seg.measured {
			last.baseWidth += seg.baseWidth
			last.Width += seg.Width
			if seg.Height > last.Height {
				last.Height = seg.Height
			}
			if (seg.metrics.Ascent + seg.metrics.Descent).Ceil() > (last.metrics.Ascent + last.metrics.Descent).Ceil() {
				last.metrics = seg.metrics
			}
		} else {
			last.baseWidth = 0
			last.Width = 0
			last.Height = 0
			last.metrics = font.Metrics{}
			last.measured = false
		}
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

// splitSegmentByFontCoverage 按字符覆盖能力拆分 segment，并为每段选择可渲染字体。
// splitSegmentByFontCoverage splits segment by rune coverage and assigns drawable fonts.
func (r *RichText) splitSegmentByFontCoverage(seg *TextSegment) TextSegments {
	if seg == nil {
		return nil
	}
	if seg.BreakLine || seg.Text == "" {
		return TextSegments{seg}
	}

	base := seg.Font

	start := 0
	current := base
	var out TextSegments
	var fallbackRunes map[string][]rune = make(map[string][]rune)
	for idx, rn := range seg.Text {
		fi := base
		if !base.coverageRanges.SupportsRune(rn) {
			fi = r.fontLibrary.MatchRuneOrFeedback(base, rn)
			key := fmt.Sprintf("%s -> %s", base.Family, fi.Family)
			fallbackRunes[key] = append(fallbackRunes[key], rn)
		}
		if current == fi {
			continue
		}
		if idx > start {
			out = append(out, makeFontBoundSegment(seg, seg.Text[start:idx], current))
		}
		start = idx
		current = fi
	}

	if start < len(seg.Text) {
		out = append(out, makeFontBoundSegment(seg, seg.Text[start:], current))
	}
	if len(out) == 0 {
		return TextSegments{seg}
	}
	if len(fallbackRunes) > 0 {
		r.logfFn(func(logger misc.Logger) {
			logger.Printf(
				"[richtext.fallback.rune] text=%q chunks=%d details=%s",
				summarizeTextForLog(seg.Text), len(out), summarizeFallbackRunes(fallbackRunes),
			)
		})
	}
	return out
}

func makeFontBoundSegment(base *TextSegment, text string, fi *FontInfo) *TextSegment {
	seg := base.CopyWithText(text)
	seg.Font = fi
	seg.FontFamily = fi.Family
	seg.FakeItalic = seg.Italic && !fi.Italic
	return seg
}

// measure 测量每个文本片段的宽度和高度
func (r *RichText) measure() {
	for _, segments := range r.lines.Range() {
		for _, seg := range segments {
			if seg.measured {
				continue
			}
			face := r.fontLibrary.GetFace(seg.Font, seg.FontSize)
			if face == nil {
				seg.baseWidth = 0
				seg.Width = 0
				seg.Height = 0
				seg.metrics = font.Metrics{}
				seg.measured = true
				continue
			}
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

func summarizeTextForLog(s string) string {
	const maxRunes = 24
	rs := []rune(s)
	if len(rs) <= maxRunes {
		return s
	}
	return string(rs[:maxRunes]) + "..."
}

func summarizeFallbackRunes(m map[string][]rune) string {
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	dedupRunes := func(in []rune) []rune {
		if len(in) == 0 {
			return nil
		}
		seen := make(map[rune]struct{}, len(in))
		out := make([]rune, 0, len(in))
		for _, rn := range in {
			if _, ok := seen[rn]; ok {
				continue
			}
			seen[rn] = struct{}{}
			out = append(out, rn)
		}
		return out
	}

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		rs := m[k]
		uniq := dedupRunes(rs)
		sample := 3
		if len(uniq) < sample {
			sample = len(uniq)
		}
		s := make([]string, 0, sample)
		for i := 0; i < sample; i++ {
			s = append(s, fmt.Sprintf("U+%04X(%q)", uniq[i], string(uniq[i])))
		}
		parts = append(parts, fmt.Sprintf("%s x%d [%s]", k, len(uniq), strings.Join(s, ", ")))
	}
	return strings.Join(parts, "; ")
}
