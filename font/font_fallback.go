package font

import (
	"fmt"
	"strings"

	"github.com/go-mixed/go-canvas/misc"
)

// initFallbackPaths 初始化 3 个 fallback 字体（regular/bold/light）。
// initFallbackPaths initializes 3 fallback fonts (regular/bold/light).
func (fs *FontLibrary) initFallbackPaths() error {
	if fs.fallbackLoaded {
		return nil
	}
	base, err := fs.mustSystemFallbackBase()
	if err != nil {
		return err
	}
	locale := normalizeLocale(detectSystemLanguage())
	preferred := preferredUnicodeRangesForLocale(locale)
	fs.fallbackRegularInfo = fs.selectFallbackFontByCoverage(FontWeightRegular, preferred, base)
	fs.fallbackBoldInfo = fs.selectFallbackFontByCoverage(FontWeightBold, preferred, base)
	fs.fallbackLightInfo = fs.selectFallbackFontByCoverage(FontWeightLight, preferred, base)
	fs.fallbackLoaded = true
	return nil
}

// fallbackFontInfo 按字重返回已初始化好的 fallback 字体。
// fallbackFontInfo returns pre-initialized fallback font by weight.
func (fs *FontLibrary) fallbackFontInfo(weight FontWeight) *FontInfo {
	if weight == FontWeightRegular {
		return fs.fallbackRegularInfo
	}
	if weight > FontWeightRegular {
		return fs.fallbackBoldInfo
	}
	return fs.fallbackLightInfo
}

// normalizeLocale 标准化 locale（小写、去掉编码与修饰符、统一分隔符）。
// normalizeLocale normalizes locale string to lowercase canonical form.
func normalizeLocale(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		return ""
	}
	if i := strings.IndexByte(v, '.'); i >= 0 {
		v = v[:i]
	}
	if i := strings.IndexByte(v, '@'); i >= 0 {
		v = v[:i]
	}
	return strings.ReplaceAll(v, "-", "_")
}

// preferredUnicodeRangesForLocale 返回 locale 对应的优先脚本区间。
// preferredUnicodeRangesForLocale returns preferred script ranges for locale.
func preferredUnicodeRangesForLocale(locale string) unicodeRanges {
	switch {
	case strings.HasPrefix(locale, "zh"):
		return unicodeRanges{
			{start: 0x3400, end: 0x4DBF}, // CJK Ext-A
			{start: 0x4E00, end: 0x9FFF}, // CJK Unified
			{start: 0xF900, end: 0xFAFF}, // CJK Compatibility Ideographs
			{start: 0x3000, end: 0x303F}, // CJK Symbols and Punctuation
			{start: 0xFF00, end: 0xFFEF}, // Fullwidth/Halfwidth
		}
	case strings.HasPrefix(locale, "ja"):
		return unicodeRanges{
			{start: 0x3040, end: 0x309F}, // Hiragana
			{start: 0x30A0, end: 0x30FF}, // Katakana
			{start: 0x31F0, end: 0x31FF}, // Katakana Phonetic Extensions
			{start: 0x4E00, end: 0x9FFF}, // CJK Unified
			{start: 0x3000, end: 0x303F}, // CJK Symbols and Punctuation
		}
	case strings.HasPrefix(locale, "ko"):
		return unicodeRanges{
			{start: 0x1100, end: 0x11FF}, // Hangul Jamo
			{start: 0x3130, end: 0x318F}, // Hangul Compatibility Jamo
			{start: 0xAC00, end: 0xD7AF}, // Hangul Syllables
		}
	case strings.HasPrefix(locale, "th"):
		return unicodeRanges{
			{start: 0x0E00, end: 0x0E7F}, // Thai
		}
	case strings.HasPrefix(locale, "ar"), strings.HasPrefix(locale, "fa"), strings.HasPrefix(locale, "ur"):
		return unicodeRanges{
			{start: 0x0600, end: 0x06FF}, // Arabic
			{start: 0x0750, end: 0x077F}, // Arabic Supplement
			{start: 0x08A0, end: 0x08FF}, // Arabic Extended-A
			{start: 0xFB50, end: 0xFDFF}, // Arabic Presentation Forms-A
			{start: 0xFE70, end: 0xFEFF}, // Arabic Presentation Forms-B
		}
	case strings.HasPrefix(locale, "ru"):
		return unicodeRanges{
			{start: 0x0400, end: 0x04FF}, // Cyrillic
			{start: 0x0500, end: 0x052F}, // Cyrillic Supplement
		}
	default:
		return unicodeRanges{
			{start: 0x0020, end: 0x007E}, // Basic Latin
			{start: 0x00A0, end: 0x00FF}, // Latin-1 Supplement
			{start: 0x0100, end: 0x024F}, // Latin Extended
		}
	}
}

// selectFallbackFontByCoverage 在可用字体中按覆盖区间与字重评分选 fallback。
// selectFallbackFontByCoverage picks fallback by coverage and weight score.
func (fs *FontLibrary) selectFallbackFontByCoverage(target FontWeight, preferred unicodeRanges, base *FontInfo) *FontInfo {
	best := base
	bestScore := fs.coverageFallbackScore(base, target, preferred)

	for _, fonts := range fs.fonts {
		for _, fi := range fonts {
			if fi == nil || fi.FontPath == "" {
				continue
			}
			score := fs.coverageFallbackScore(fi, target, preferred)
			if fi.Italic {
				score -= 20000
			}
			if score > bestScore {
				best = fi
				bestScore = score
			}
		}
	}
	if best != nil {
		return best
	}
	return base
}

// coverageFallbackScore 计算 fallback 候选评分。
// coverageFallbackScore computes score for fallback candidate selection.
func (fs *FontLibrary) coverageFallbackScore(fi *FontInfo, target FontWeight, preferred unicodeRanges) int {
	preferredCoverage := fi.coverageRanges.IntersectionCount(preferred)
	breadthScore := len(fi.coverageRanges) * 10
	weightScore := int(2000 - misc.Abs(fi.Bold-target))
	if weightScore < 0 {
		weightScore = 0
	}
	return preferredCoverage + breadthScore + weightScore
}

// mustSystemFallbackBase 从系统保底 family 中找到一个可用字体，找不到返回 error。
// mustSystemFallbackBase finds one available font from system fallback families.
func (fs *FontLibrary) mustSystemFallbackBase() (*FontInfo, error) {
	for _, family := range systemFallbackFamilies() {
		if fi := fs.findFamilyFallbackFont(family); fi != nil {
			return fi, nil
		}
	}
	return nil, fmt.Errorf("system fallback fonts not found in scanned list: %v", systemFallbackFamilies())
}

// findFamilyFallbackFont 在已扫描字体中按 family 找到最合适的 regular 近似字体。
// findFamilyFallbackFont finds a regular-like font by exact normalized family match.
func (fs *FontLibrary) findFamilyFallbackFont(family string) *FontInfo {
	want := normalizeFamilyName(family)
	if want == "" {
		return nil
	}

	fonts := fs.fonts[want]
	if len(fonts) == 0 {
		return nil
	}

	var best *FontInfo
	bestScore := -1
	for _, fi := range fonts {
		if fi == nil || fi.FontPath == "" {
			continue
		}
		score := int(1000 - misc.Abs(fi.Bold-FontWeightRegular))
		if fi.Italic {
			score -= 100
		}
		if score > bestScore {
			best = fi
			bestScore = score
		}
	}
	return best
}
