package font

import (
	"strings"

	"golang.org/x/image/font"
)

const ellipsisText = "…"

// noWrapEllipsis 对每个显式文本行分别执行 NoWrapEllipsis。
// noWrapEllipsis applies NoWrapEllipsis independently to every explicit input line.
//
// NoWrapEllipsis 禁止自动换行，因此不会把一行拆成多行；只有超宽的显式行
// 才会被截断并追加一个省略号。
// NoWrapEllipsis never creates automatic line breaks; only an overflowing
// explicit line is truncated and gets one trailing ellipsis.
func (r *RichText) noWrapEllipsis(in TextSegments, maxWidth int) TextSegments {
	if maxWidth <= 0 {
		return in
	}

	out := make(TextSegments, 0, len(in))
	line := make(TextSegments, 0, len(in))
	for _, seg := range in {
		if seg == nil {
			continue
		}
		if !seg.BreakLine {
			line = append(line, seg)
			continue
		}

		out = append(out, r.ellipsisLine(line, maxWidth)...)
		out = append(out, seg)
		line = line[:0]
	}
	out = append(out, r.ellipsisLine(line, maxWidth)...)
	return out
}

func (r *RichText) ellipsisLine(line TextSegments, maxWidth int) TextSegments {
	if len(line) == 0 {
		return nil
	}

	// 先测量完整行，只有确实超宽时才创建截断结果。
	// Measure the complete line first and only build a truncated result when it
	// actually exceeds the available width.
	lineWidth := 0
	var overflow *TextSegment
	for _, seg := range line {
		if seg == nil || seg.Text == "" {
			continue
		}
		face := r.fontLibrary.GetFace(seg.Font, seg.FontSize)
		if face == nil {
			continue
		}
		measureSegment(seg, face)
		lineWidth += seg.Width
		if overflow == nil && lineWidth > maxWidth {
			overflow = seg
		}
	}
	if overflow == nil {
		return line
	}

	// 省略号必须先占用宽度。正文可使用的宽度是 maxWidth - width("…")。
	// Reserve the ellipsis width first. The content budget is maxWidth-width("…").
	face := r.fontLibrary.GetFace(overflow.Font, overflow.FontSize)
	if face == nil {
		return nil
	}
	ellipsis := overflow.CopyWithText(ellipsisText)
	measureSegment(ellipsis, face)
	remaining := maxWidth - ellipsis.Width
	if remaining < 0 {
		return nil
	}

	out := make(TextSegments, 0, len(line)+1)
	for _, seg := range line {
		if seg == nil || seg.Text == "" {
			continue
		}
		face := r.fontLibrary.GetFace(seg.Font, seg.FontSize)
		if face == nil {
			continue
		}
		measureSegment(seg, face)
		if seg.Width <= remaining {
			out = append(out, seg)
			remaining -= seg.Width
			continue
		}

		// 之前的 segment 已经完整放入；只需截断第一个放不下的 segment。
		// Earlier segments already fit; only the first overflowing segment is
		// truncated.
		prefix := truncateSegmentPrefix(seg, face, remaining)
		if prefix != nil {
			out = append(out, prefix)
		}
		break
	}
	return append(out, ellipsis)
}

func truncateSegmentPrefix(seg *TextSegment, face font.Face, maxWidth int) *TextSegment {
	if maxWidth <= 0 {
		return nil
	}

	clusters := splitGraphemeClusters(seg.Text)
	if len(clusters) == 0 {
		return nil
	}
	// prefixW[i] 是前 i 个 grapheme cluster 的累计字宽，且单调递增。
	// prefixW[i] is the cumulative width of the first i grapheme clusters.
	prefixW := buildPrefixWidths(face, clusters)
	extraItalic := 0
	if seg.FakeItalic {
		extraItalic = syntheticItalicExtraWidth((face.Metrics().Ascent + face.Metrics().Descent).Ceil())
	}
	// 这里使用二分查找：在 [1, len(clusters)] 中寻找能放下的最大
	// cluster 数量。由于 prefixW 单调递增，不需要逐字符删除并重新测量。
	// Binary search finds the largest fitting cluster count in [1, len(clusters)].
	// Because prefixW is monotonic, we avoid deleting characters and remeasuring.
	count := maxFittingByCluster(prefixW, 0, maxWidth, extraItalic)
	if count == 0 {
		return nil
	}

	prefix := seg.CopyWithText(strings.Join(clusters[:count], ""))
	applySegmentMeasureWithBase(prefix, face, clusterRangeWidth(prefixW, 0, count, 0))
	return prefix
}

func measureSegment(seg *TextSegment, face font.Face) {
	// 普通 segment 复用缓存；截断产生的新 segment 会重新测量。
	// Reuse cached measurements; newly created truncated segments are remeasured.
	if seg.measured && seg.baseWidth > 0 {
		return
	}
	applySegmentMeasureWithBase(seg, face, font.MeasureString(face, seg.Text).Ceil())
}
