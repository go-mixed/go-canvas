package font

import (
	"slices"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// wrapFirstFit: textwrap::wrap_first_fit 风格的贪心分行。
// 对每段文本按可断点做 first-fit，超长无断点时再回退到 grapheme 强制截断。
func (r *RichText) wrapFirstFit(in TextSegments, maxWidth int, breakPolicy LineBreakPolicy) TextSegments {
	t0 := time.Now()
	defer func() {
		r.logf("[richtext.wrap.first_fit] maxWidth=%d breakPolicy=%d in=%d elapsed=%s", maxWidth, breakPolicy, len(in), time.Since(t0))
	}()

	if maxWidth <= 0 {
		return in
	}

	out := make(TextSegments, 0, len(in))
	lineWidth := 0
	lineHasContent := false
	ws := &r.wrapScratch

	for _, seg := range in {
		if seg == nil {
			continue
		}
		if seg.BreakLine {
			out = append(out, seg)
			lineWidth = 0
			lineHasContent = false
			continue
		}
		if seg.Text == "" {
			continue
		}

		face := r.fontLibrary.GetFace(seg.Font, seg.FontSize)
		if face == nil {
			out = append(out, seg)
			lineWidth += seg.Width
			lineHasContent = seg.Text != ""
			continue
		}

		extraItalic := 0
		if seg.FakeItalic {
			extraItalic = syntheticItalicExtraWidth((face.Metrics().Ascent + face.Metrics().Descent).Ceil())
		}

		segWidth := seg.Width
		if !seg.measured || seg.baseWidth <= 0 {
			segWidth = font.MeasureString(face, seg.Text).Ceil()
		} else if seg.baseWidth > 0 {
			segWidth = seg.baseWidth
		}
		applySegmentMeasureWithBase(seg, face, segWidth)
		if extraItalic > 0 {
			segWidth += extraItalic
		}

		remaining := maxWidth - lineWidth
		if segWidth <= remaining {
			out = append(out, seg)
			lineWidth += segWidth
			lineHasContent = true
			continue
		}

		if isLineBreakNever(breakPolicy) {
			if lineHasContent {
				out = append(out, newBreakMarker(seg))
				lineWidth = 0
				lineHasContent = false
			}
			out = append(out, seg)
			lineWidth += segWidth
			lineHasContent = true
			continue
		}

		ws.resetForSegment()
		clusters := splitGraphemeClustersInto(ws.clusters, seg.Text)
		ws.clusters = clusters
		if len(clusters) == 0 {
			continue
		}
		prefixW := buildPrefixWidthsInto(ws.prefixW, face, clusters)
		ws.prefixW = prefixW
		legal := findLegalBreaksInto(ws.legal, clusters)
		ws.legal = legal

		start := 0
		for start < len(clusters) {
			remaining = maxWidth - lineWidth
			if remaining <= 0 && lineHasContent {
				out = append(out, newBreakMarker(seg))
				lineWidth = 0
				lineHasContent = false
				continue
			}

			totalW := clusterRangeWidth(prefixW, start, len(clusters), extraItalic)
			if totalW <= remaining {
				chunkText := strings.Join(clusters[start:], "")
				chunkBase := clusterRangeWidth(prefixW, start, len(clusters), 0)
				chunk := seg.CopyWithText(chunkText)
				applySegmentMeasureWithBase(chunk, face, chunkBase)
				out = append(out, chunk)
				lineWidth += totalW
				lineHasContent = chunk.Text != ""
				break
			}

			k := 0
			if isLineBreakAlways(breakPolicy) {
				k = maxFittingByCluster(prefixW, start, remaining, extraItalic)
			} else {
				k = maxFittingByLegalBreak(legal, prefixW, start, remaining, extraItalic)
			}

			if k <= start {
				if lineHasContent {
					out = append(out, newBreakMarker(seg))
					lineWidth = 0
					lineHasContent = false
					continue
				}
				k = maxFittingByCluster(prefixW, start, remaining, extraItalic)
				if k <= start {
					k = min(start+1, len(clusters))
				}
			}

			chunkW := clusterRangeWidth(prefixW, start, k, extraItalic)
			chunkText := strings.Join(clusters[start:k], "")
			chunkBase := clusterRangeWidth(prefixW, start, k, 0)
			chunk := seg.CopyWithText(chunkText)
			applySegmentMeasureWithBase(chunk, face, chunkBase)
			out = append(out, chunk)
			start = k
			if start < len(clusters) {
				out = append(out, newBreakMarker(seg))
				lineWidth = 0
				lineHasContent = false
			} else {
				lineWidth += chunkW
				lineHasContent = chunk.Text != ""
			}
		}
	}

	return out
}

func maxFittingByLegalBreak(legal []int, prefixW []fixed.Int26_6, start, remainingWidth, extraItalic int) int {
	if len(legal) == 0 {
		return start
	}
	i, _ := slices.BinarySearch(legal, start+1)
	if i >= len(legal) {
		return start
	}
	lo, hi := i, len(legal)-1
	best := start
	for lo <= hi {
		mid := (lo + hi) >> 1
		k := legal[mid]
		if clusterRangeWidth(prefixW, start, k, extraItalic) <= remainingWidth {
			best = k
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return best
}

func maxFittingByCluster(prefixW []fixed.Int26_6, start, remainingWidth, extraItalic int) int {
	lo, hi := start+1, len(prefixW)-1
	best := start
	for lo <= hi {
		mid := (lo + hi) >> 1
		w := clusterRangeWidth(prefixW, start, mid, extraItalic)
		if w <= remainingWidth {
			best = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return best
}
