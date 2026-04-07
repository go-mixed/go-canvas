package font

import (
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func (r *RichText) wordWrap(in TextSegments, maxWidth int, breakPolicy LineBreakPolicy) TextSegments {
	t0 := time.Now()
	defer func() {
		fmt.Printf("[richtext.wordWrap] maxWidth=%d breakPolicy=%d in=%d elapsed=%s\n", maxWidth, breakPolicy, len(in), time.Since(t0))
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

		if breakPolicy == LineBreakNever {
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
				semantic := breakPolicy != LineBreakAlways
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
			semantic := breakPolicy != LineBreakAlways
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

func chooseBreakIndex(
	clusters []string,
	legal []int,
	prefixW []fixed.Int26_6,
	start int,
	remainingWidth int,
	extraItalic int,
	useSemanticBreak bool,
	ws *wrapScratch,
) int {
	n := len(clusters)
	if start >= n {
		return n
	}

	candidates := buildCandidatesInto(ws.candidates, legal, start, n, useSemanticBreak)
	ws.candidates = candidates
	if len(candidates) == 0 {
		return start
	}

	totalW := clusterRangeWidth(prefixW, start, n, extraItalic)
	if totalW <= remainingWidth {
		return n
	}

	guess := start + int(float64(remainingWidth)/float64(max(1, totalW))*float64(n-start))
	guess = clampInt(guess, start+1, n-1)

	left := findLastCandidateLE(candidates, guess)
	best := start
	if left > start && clusterRangeWidth(prefixW, start, left, extraItalic) <= remainingWidth {
		best = left
	}

	// 在合法断点上二分找“最大可容纳断点”。
	i, _ := slices.BinarySearch(candidates, max(start+1, best+1))
	if i >= len(candidates) {
		return best
	}
	l, r := i, len(candidates)-1
	ans := -1
	for l <= r {
		m := (l + r) >> 1
		k := candidates[m]
		if clusterRangeWidth(prefixW, start, k, extraItalic) <= remainingWidth {
			ans = m
			l = m + 1
		} else {
			r = m - 1
		}
	}
	if ans >= 0 {
		return candidates[ans]
	}
	return best
}

func buildCandidates(legal []int, start, n int, useSemanticBreak bool) []int {
	return buildCandidatesInto(nil, legal, start, n, useSemanticBreak)
}

func buildCandidatesInto(dst []int, legal []int, start, n int, useSemanticBreak bool) []int {
	if !useSemanticBreak {
		candidates := dst[:0]
		for i := start + 1; i <= n; i++ {
			candidates = append(candidates, i)
		}
		return candidates
	}

	candidates := dst[:0]
	for _, k := range legal {
		if k > start && k <= n {
			candidates = append(candidates, k)
		}
	}
	return candidates
}

func findLastCandidateLE(candidates []int, v int) int {
	l, r := 0, len(candidates)-1
	ans := -1
	for l <= r {
		m := (l + r) >> 1
		if candidates[m] <= v {
			ans = m
			l = m + 1
		} else {
			r = m - 1
		}
	}
	if ans >= 0 {
		return candidates[ans]
	}
	return 0
}

func findLegalBreaks(clusters []string) []int {
	return findLegalBreaksInto(nil, clusters)
}

func findLegalBreaksInto(dst []int, clusters []string) []int {
	if len(clusters) <= 1 {
		return dst[:0]
	}

	breaks := dst[:0]
	for i := 1; i < len(clusters); i++ {
		prev := clusters[i-1]
		next := clusters[i]
		if isLegalBreak(prev, next) {
			breaks = append(breaks, i)
		}
	}
	return breaks
}

func isLegalBreak(prevCluster, nextCluster string) bool {
	if prevCluster == "" || nextCluster == "" {
		return true
	}

	prevLast, _ := utf8LastRune(prevCluster)
	nextFirst, _ := utf8.DecodeRuneInString(nextCluster)

	// 行末禁则：开括号/开引号不要放在行尾。
	// Line-end prohibition: opening punctuations should not be placed at line end.
	if isKinsokuEndRune(prevLast) {
		return false
	}
	// 行首禁则：CJK 句读点、闭括号以及左粘连标点不放在下一行行首。
	// Line-start prohibition: CJK closers and sticky punctuations should not start next line.
	if isKinsokuStartRune(nextFirst) || isLeftStickyPunctuationRune(nextFirst) {
		return false
	}

	if unicode.IsSpace(prevLast) {
		return true
	}
	if isPunctuationRune(prevLast) {
		return true
	}
	if isCJKRune(prevLast) && isCJKRune(nextFirst) {
		return true
	}

	prevScript := scriptClass(prevLast)
	nextScript := scriptClass(nextFirst)
	return prevScript != scriptUnknown && nextScript != scriptUnknown && prevScript != nextScript
}

func buildPrefixWidths(face FontFaceAdvance, clusters []string) []fixed.Int26_6 {
	return buildPrefixWidthsInto(nil, face, clusters)
}

func buildPrefixWidthsInto(dst []fixed.Int26_6, face FontFaceAdvance, clusters []string) []fixed.Int26_6 {
	need := len(clusters) + 1
	prefix := dst[:0]
	if cap(prefix) < need {
		prefix = make([]fixed.Int26_6, need)
	} else {
		prefix = prefix[:need]
	}
	var sum fixed.Int26_6
	prev := rune(-1)

	for i, cl := range clusters {
		for _, c := range cl {
			if prev >= 0 {
				sum += face.Kern(prev, c)
			}
			adv, ok := face.GlyphAdvance(c)
			if ok {
				sum += adv
			}
			prev = c
		}
		prefix[i+1] = sum
	}
	return prefix
}

type FontFaceAdvance interface {
	GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool)
	Kern(r0, r1 rune) fixed.Int26_6
}

func applySegmentMeasureWithBase(seg *TextSegment, face font.Face, baseWidth int) {
	seg.baseWidth = baseWidth
	seg.Width = baseWidth
	seg.metrics = face.Metrics()
	seg.Height = (seg.metrics.Ascent + seg.metrics.Descent).Ceil()
	if seg.FakeItalic {
		seg.Width += syntheticItalicExtraWidth(seg.Height)
	}
	seg.measured = true
}

type wrapScratch struct {
	clusters   []string
	prefixW    []fixed.Int26_6
	legal      []int
	candidates []int
}

func (w *wrapScratch) resetForSegment() {
	w.clusters = w.clusters[:0]
	w.prefixW = w.prefixW[:0]
	w.legal = w.legal[:0]
	w.candidates = w.candidates[:0]
}

func clusterRangeWidth(prefixW []fixed.Int26_6, i, j int, extraItalic int) int {
	if i < 0 {
		i = 0
	}
	if j < i {
		j = i
	}
	if j >= len(prefixW) {
		j = len(prefixW) - 1
	}
	w := (prefixW[j] - prefixW[i]).Ceil()
	if extraItalic > 0 && j > i {
		w += extraItalic
	}
	return w
}

func splitGraphemeClusters(s string) []string {
	return splitGraphemeClustersInto(nil, s)
}

func splitGraphemeClustersInto(dst []string, s string) []string {
	if s == "" {
		return dst[:0]
	}

	out := dst[:0]
	var buf []rune
	riCount := 0
	prev := rune(0)
	prevWasZWJ := false

	flush := func() {
		if len(buf) == 0 {
			return
		}
		out = append(out, string(buf))
		buf = buf[:0]
		riCount = 0
	}

	for _, r := range s {
		if len(buf) == 0 {
			buf = append(buf, r)
			riCount = boolToInt(isRegionalIndicator(r))
			prev = r
			prevWasZWJ = r == '\u200d'
			continue
		}

		attach := false
		switch {
		case isCombiningMark(r):
			attach = true
		case isVariationSelector(r):
			attach = true
		case isEmojiModifier(r):
			attach = true
		case prevWasZWJ || r == '\u200d':
			attach = true
		case isRegionalIndicator(prev) && isRegionalIndicator(r) && riCount%2 == 1:
			attach = true
		}

		if !attach {
			flush()
		}
		buf = append(buf, r)
		if isRegionalIndicator(r) {
			riCount++
		} else {
			riCount = 0
		}
		prev = r
		prevWasZWJ = r == '\u200d'
	}

	flush()
	return out
}

func isCombiningMark(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r) || unicode.Is(unicode.Me, r)
}

func isVariationSelector(r rune) bool {
	return (r >= 0xFE00 && r <= 0xFE0F) || (r >= 0xE0100 && r <= 0xE01EF)
}

func isEmojiModifier(r rune) bool {
	return r >= 0x1F3FB && r <= 0x1F3FF
}

func isRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

func utf8LastRune(s string) (rune, int) {
	r, size := rune(utf8.RuneError), 0
	for i := range len(s) {
		rr, sz := utf8.DecodeRuneInString(s[i:])
		if rr == utf8.RuneError && sz == 1 {
			continue
		}
		r = rr
		size = sz
		i += sz - 1
	}
	if size == 0 {
		return utf8.RuneError, 0
	}
	return r, size
}

func isCJKRune(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) ||
		(r >= 0x3400 && r <= 0x4DBF) ||
		(r >= 0x3040 && r <= 0x30FF) ||
		(r >= 0xAC00 && r <= 0xD7AF)
}

func isPunctuationRune(r rune) bool {
	return unicode.IsPunct(r) || (r >= 0x3000 && r <= 0x303F) || (r >= 0xFF00 && r <= 0xFFEF)
}

type scriptKind int

const (
	scriptUnknown scriptKind = iota
	scriptCJK
	scriptLatinDigit
	scriptOtherLetter
)

func scriptClass(r rune) scriptKind {
	switch {
	case isCJKRune(r):
		return scriptCJK
	case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || unicode.IsDigit(r):
		return scriptLatinDigit
	case unicode.IsLetter(r):
		return scriptOtherLetter
	default:
		return scriptUnknown
	}
}

func isKinsokuStartRune(r rune) bool {
	switch r {
	case '，', '．', '！', '：', '；', '？', '、', '。', '・',
		'）', '〕', '〉', '》', '」', '』', '】', '〗', '〙', '〛',
		'ー', '々', '〻', 'ゝ', 'ゞ', 'ヽ', 'ヾ':
		return true
	default:
		return false
	}
}

func isKinsokuEndRune(r rune) bool {
	switch r {
	case '"', '(', '[', '{', '“', '‘', '«', '‹',
		'（', '〔', '〈', '《', '「', '『', '【', '〖', '〘', '〚':
		return true
	default:
		return false
	}
}

func isLeftStickyPunctuationRune(r rune) bool {
	switch r {
	case '.', ',', '!', '?', ':', ';', ')', ']', '}', '%', '"', '”', '’', '»', '›', '…',
		'،', '؛', '؟', '।', '॥':
		return true
	default:
		return false
	}
}

func newBreakMarker(seg *TextSegment) *TextSegment {
	m := seg.CopyWithText("")
	m.BreakLine = true
	return m
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
