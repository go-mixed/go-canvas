package font

import (
	"strings"
	"time"

	"github.com/go-mixed/go-canvas/ctypes"
	"golang.org/x/image/font"
)

func (r *RichText) wordWrap(in TextSegments, maxWidth int, breakPolicy ctypes.WordWrapMode) TextSegments {
	t0 := time.Now()
	defer func() {
		r.logf("[richtext.wrap.smart] maxWidth=%d breakPolicy=%d in=%d elapsed=%s", maxWidth, breakPolicy, len(in), time.Since(t0))
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
			if seg.Text != "" {
				lineHasContent = true
			}
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
			remaining := maxWidth - lineWidth
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

			if lineHasContent {
				semantic := !isLineBreakAlways(breakPolicy)
				k := chooseBreakIndex(clusters, legal, prefixW, start, remaining, extraItalic, semantic, ws)
				if k <= start {
					out = append(out, newBreakMarker(seg))
					lineWidth = 0
					lineHasContent = false
					continue
				}
				chunkText := strings.Join(clusters[start:k], "")
				chunkBase := clusterRangeWidth(prefixW, start, k, 0)
				chunk := seg.CopyWithText(chunkText)
				applySegmentMeasureWithBase(chunk, face, chunkBase)
				out = append(out, chunk)
				out = append(out, newBreakMarker(seg))
				start = k
				lineWidth = 0
				lineHasContent = false
				continue
			}

			// 当前行为空：先语义断点，失败后强制按 cluster 截断。
			semantic := !isLineBreakAlways(breakPolicy)
			k := chooseBreakIndex(clusters, legal, prefixW, start, remaining, extraItalic, semantic, ws)
			if k <= start {
				k = chooseBreakIndex(clusters, nil, prefixW, start, remaining, extraItalic, false, ws)
			}
			if k <= start {
				k = min(start+1, len(clusters))
			}

			chunkText := strings.Join(clusters[start:k], "")
			chunkBase := clusterRangeWidth(prefixW, start, k, 0)
			chunk := seg.CopyWithText(chunkText)
			applySegmentMeasureWithBase(chunk, face, chunkBase)
			out = append(out, chunk)
			start = k
			lineWidth = 0
			lineHasContent = false
			if start < len(clusters) {
				out = append(out, newBreakMarker(seg))
			}
		}
	}

	return out
}
