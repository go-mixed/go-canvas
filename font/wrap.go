package font

import (
	"slices"
	"unicode"
	"unicode/utf8"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// isLineBreakNever 判断策略是否为“禁止自动换行”。
// isLineBreakNever reports whether policy disables auto wrapping.
func isLineBreakNever(p WordWrapMode) bool {
	return p == NoWrap
}

// isLineBreakAlways 判断策略是否为“总是按宽度断行”。
// isLineBreakAlways reports whether policy forces width-based wrapping.
func isLineBreakAlways(p WordWrapMode) bool {
	return p == BreakAll
}

// chooseBreakIndex 选择当前行可放下的最佳断点（先估算，再在候选断点上二分）。
// chooseBreakIndex picks the best fitting break index (estimate first, then binary-search candidates).
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

// buildCandidates 返回断点候选列表（语义断点或全部 cluster 边界）。
// buildCandidates returns break candidates (semantic breakpoints or all cluster boundaries).
func buildCandidates(legal []int, start, n int, useSemanticBreak bool) []int {
	return buildCandidatesInto(nil, legal, start, n, useSemanticBreak)
}

// buildCandidatesInto 将候选断点写入可复用切片，减少临时分配。
// buildCandidatesInto writes candidate breakpoints into a reusable slice to reduce allocations.
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

// findLastCandidateLE 返回小于等于给定值的最后一个候选断点。
// findLastCandidateLE returns the last candidate breakpoint that is <= target value.
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

// findLegalBreak 计算文本的合法断点索引。
// findLegalBreak computes legal break indexes for given grapheme clusters.
func findLegalBreaks(clusters []string) []int {
	return findLegalBreaksInto(nil, clusters)
}

// findLegalBreaksInto 将合法断点写入可复用切片，减少分配。
// findLegalBreaksInto writes legal breakpoints into a reusable slice.
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

// isLegalBreak 判断两个 cluster 间是否允许断行，包含 kinsoku 与脚本切换规则。
// isLegalBreak reports whether a line break is allowed between two clusters, including kinsoku/script rules.
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

// isLegalBreakLegacy 保留历史断点逻辑（不含 kinsoku/sticky 约束），仅用于考古对比，不参与运行时分支。
// isLegalBreakLegacy keeps historical break logic (without kinsoku/sticky constraints) for archaeology only.
func isLegalBreakLegacy(prevCluster, nextCluster string) bool {
	if prevCluster == "" || nextCluster == "" {
		return true
	}
	prevLast, _ := utf8LastRune(prevCluster)
	nextFirst, _ := utf8.DecodeRuneInString(nextCluster)

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

// buildPrefixWidths 构建 cluster 前缀宽数组。
// buildPrefixWidths builds prefix-width array for grapheme clusters.
func buildPrefixWidths(face FontFaceAdvance, clusters []string) []fixed.Int26_6 {
	return buildPrefixWidthsInto(nil, face, clusters)
}

// buildPrefixWidthsInto 构建可复用的前缀宽数组（含 kerning）。
// buildPrefixWidthsInto builds reusable prefix widths (including kerning).
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

// FontFaceAdvance 提供 wrap 阶段所需的最小字形 advance/kerning 接口。
// FontFaceAdvance is the minimal glyph advance/kerning interface used by wrapping.
type FontFaceAdvance interface {
	GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool)
	Kern(r0, r1 rune) fixed.Int26_6
}

// applySegmentMeasureWithBase 将测量结果写回 segment，并统一处理伪斜体额外宽度。
// applySegmentMeasureWithBase stores measured values into segment and applies synthetic-italic extra width.
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

// wrapScratch 复用 wrap 过程中的临时切片，降低 GC 压力。
// wrapScratch reuses temporary slices in wrapping to reduce GC pressure.
type wrapScratch struct {
	clusters   []string
	prefixW    []fixed.Int26_6
	legal      []int
	candidates []int
}

// resetForSegment 清空当前 segment 的临时缓冲。
// resetForSegment clears temporary buffers for the current segment.
func (w *wrapScratch) resetForSegment() {
	w.clusters = w.clusters[:0]
	w.prefixW = w.prefixW[:0]
	w.legal = w.legal[:0]
	w.candidates = w.candidates[:0]
}

// clusterRangeWidth 通过前缀宽 O(1) 计算 cluster 子区间宽度。
// clusterRangeWidth computes sub-range width in O(1) using prefix widths.
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

// splitGraphemeClusters 将字符串拆为 grapheme cluster 列表。
// splitGraphemeClusters splits input text into grapheme clusters.
func splitGraphemeClusters(s string) []string {
	return splitGraphemeClustersInto(nil, s)
}

// splitGraphemeClustersInto 将 cluster 写入可复用切片，处理 ZWJ/变体选择符/组合附标等。
// splitGraphemeClustersInto writes clusters into reusable slice, handling ZWJ/VS/combining marks.
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

// isCombiningMark 判断字符是否为组合附标。
// isCombiningMark reports whether rune is a combining mark.
func isCombiningMark(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r) || unicode.Is(unicode.Me, r)
}

// isVariationSelector 判断字符是否为变体选择符。
// isVariationSelector reports whether rune is a variation selector.
func isVariationSelector(r rune) bool {
	return (r >= 0xFE00 && r <= 0xFE0F) || (r >= 0xE0100 && r <= 0xE01EF)
}

// isEmojiModifier 判断字符是否为 emoji 肤色修饰符。
// isEmojiModifier reports whether rune is an emoji skin-tone modifier.
func isEmojiModifier(r rune) bool {
	return r >= 0x1F3FB && r <= 0x1F3FF
}

// isRegionalIndicator 判断字符是否为区域旗帜指示符。
// isRegionalIndicator reports whether rune is a regional indicator symbol.
func isRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

// utf8LastRune 返回字符串最后一个有效 rune。
// utf8LastRune returns the last valid rune of a UTF-8 string.
func utf8LastRune(s string) (rune, int) {
	if s == "" {
		return utf8.RuneError, 0
	}
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && size == 1 {
		return utf8.RuneError, 0
	}
	return r, size
}

// isCJKRune 判断字符是否属于 CJK 常见区间。
// isCJKRune reports whether rune belongs to common CJK ranges.
func isCJKRune(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) ||
		(r >= 0x3400 && r <= 0x4DBF) ||
		(r >= 0x3040 && r <= 0x30FF) ||
		(r >= 0xAC00 && r <= 0xD7AF)
}

// isPunctuationRune 判断字符是否为常见标点（含全角标点块）。
// isPunctuationRune reports whether rune is punctuation (including full-width blocks).
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

// scriptClass 返回字符所属的简化脚本类别，用于脚本切换断点规则。
// scriptClass returns simplified script class used by script-switch break rules.
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

// isKinsokuStartRune 判断字符是否属于行首禁则集合。
// isKinsokuStartRune reports whether rune is prohibited at line start.
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

// isKinsokuEndRune 判断字符是否属于行末禁则集合。
// isKinsokuEndRune reports whether rune is prohibited at line end.
func isKinsokuEndRune(r rune) bool {
	switch r {
	case '"', '(', '[', '{', '“', '‘', '«', '‹',
		'（', '〔', '〈', '《', '「', '『', '【', '〖', '〘', '〚':
		return true
	default:
		return false
	}
}

// isLeftStickyPunctuationRune 判断字符是否应黏在前一个词尾，不应出现在新行行首。
// isLeftStickyPunctuationRune reports whether rune should stick to previous token and not start a new line.
func isLeftStickyPunctuationRune(r rune) bool {
	switch r {
	case '.', ',', '!', '?', ':', ';', ')', ']', '}', '%', '"', '”', '’', '»', '›', '…',
		'،', '؛', '؟', '।', '॥':
		return true
	default:
		return false
	}
}

// newBreakMarker 创建换行标记 segment。
// newBreakMarker creates a segment used as explicit line-break marker.
func newBreakMarker(seg *TextSegment) *TextSegment {
	m := seg.CopyWithText("")
	m.BreakLine = true
	return m
}

// clampInt 将整数限制在闭区间 [lo, hi]。
// clampInt clamps integer to inclusive range [lo, hi].
func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// boolToInt 将 bool 转换为 0/1。
// boolToInt converts bool to 0/1 integer.
func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
