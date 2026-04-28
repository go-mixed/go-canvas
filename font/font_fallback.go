package font

import (
	"fmt"
	"strings"

	"github.com/go-mixed/go-canvas/internel/misc"
	xfont "golang.org/x/image/font"
)

// initFallbackPaths 初始化 3 个 fallback 字体（regular/bold/light）。
// initFallbackPaths initializes 3 fallback fonts (regular/bold/light).
func (fs *FontLibrary) initFallbackPaths() error {
	if fs.fallbackLoaded {
		return nil
	}

	locale := systemLanguage()
	preferred := preferredUnicodeRangesForLocale(locale)
	systemFamilies := normalizeFamilyNames(localeFallbackFamilyTable(locale))

	fs.logf("[font-library]locale \"%s\", font families: %q", locale, systemFamilies)

	var err error
	fs.fallbackRegularInfo, err = fs.findFallbackFont(systemFamilies, preferred, xfont.WeightNormal)
	if err != nil {
		return err
	}
	fs.fallbackBoldInfo, err = fs.findFallbackFont(systemFamilies, preferred, xfont.WeightBold)
	if err != nil {
		return err
	}
	fs.fallbackLightInfo, err = fs.findFallbackFont(systemFamilies, preferred, xfont.WeightLight)
	// Light 失败时回退到 Regular，避免初始化中断。
	if err != nil {
		fs.fallbackLightInfo = fs.fallbackRegularInfo
	}

	fs.fallbackLoaded = true
	return nil
}

// fallbackFontInfo 按字重返回已初始化好的 fallback 字体。
// fallbackFontInfo returns pre-initialized fallback font by weight.
func (fs *FontLibrary) fallbackFontInfo(weight xfont.Weight) *FontInfo {
	if weight == xfont.WeightNormal {
		return fs.fallbackRegularInfo
	}
	if weight > xfont.WeightNormal {
		return fs.fallbackBoldInfo
	}
	return fs.fallbackLightInfo
}

// systemLanguage 读取系统语言并标准化 locale（小写、去掉编码与修饰符、统一分隔符）。
// systemLanguage normalizes locale string to lowercase canonical form.
func systemLanguage() string {
	v := detectSystemLanguage()
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

// findFallbackFont 按顺序选择 fallback：
// 1) locale 对应家族表
// 2) coverage 评分
// 3) systemFallbackFamilies 兜底
func (fs *FontLibrary) findFallbackFont(systemFamilies []string, preferred unicodeRanges, fontWeight xfont.Weight) (*FontInfo, error) {
	fi := fs.findFontByFamilies(systemFamilies, fontWeight)
	if fi != nil {
		fs.logf("[font-library]set fallback %d font: [%s] from system locale", fontWeight, fi.Family)
		return fi, nil
	}

	fi = fs.selectFallbackFontByCoverage(preferred, fontWeight)
	if fi != nil {
		fs.logf("[font-library]set fallback %d font: [%s] from coverage", fontWeight, fi.Family)
		return fi, nil
	}

	systemFamilies = normalizeFamilyNames(systemFallbackFamilies())
	fi = fs.findFontByFamilies(systemFamilies, fontWeight)
	if fi != nil {
		fs.logf("[font-library]set fallback %d font: [%s] from system fallback", fontWeight, fi.Family)
		return fi, nil
	}

	return nil, fmt.Errorf("system font invalid")
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
		return unicodeRanges{{start: 0x0E00, end: 0x0E7F}} // Thai
	case strings.HasPrefix(locale, "ar"), strings.HasPrefix(locale, "fa"), strings.HasPrefix(locale, "ur"):
		return unicodeRanges{
			{start: 0x0600, end: 0x06FF}, // Arabic
			{start: 0x0750, end: 0x077F}, // Arabic Supplement
			{start: 0x08A0, end: 0x08FF}, // Arabic Extended-A
			{start: 0xFB50, end: 0xFDFF}, // Arabic Presentation Forms-A
			{start: 0xFE70, end: 0xFEFF}, // Arabic Presentation Forms-B
		}
	case strings.HasPrefix(locale, "he"):
		return unicodeRanges{{start: 0x0590, end: 0x05FF}} // Hebrew
	case strings.HasPrefix(locale, "hi"), strings.HasPrefix(locale, "mr"):
		return unicodeRanges{{start: 0x0900, end: 0x097F}} // Devanagari
	case strings.HasPrefix(locale, "bn"):
		return unicodeRanges{{start: 0x0980, end: 0x09FF}} // Bengali
	case strings.HasPrefix(locale, "pa"):
		return unicodeRanges{{start: 0x0A00, end: 0x0A7F}} // Gurmukhi
	case strings.HasPrefix(locale, "gu"):
		return unicodeRanges{{start: 0x0A80, end: 0x0AFF}} // Gujarati
	case strings.HasPrefix(locale, "ta"):
		return unicodeRanges{{start: 0x0B80, end: 0x0BFF}} // Tamil
	case strings.HasPrefix(locale, "te"):
		return unicodeRanges{{start: 0x0C00, end: 0x0C7F}} // Telugu
	case strings.HasPrefix(locale, "ml"):
		return unicodeRanges{{start: 0x0D00, end: 0x0D7F}} // Malayalam
	case strings.HasPrefix(locale, "el"):
		return unicodeRanges{{start: 0x0370, end: 0x03FF}} // Greek
	case strings.HasPrefix(locale, "ru"), strings.HasPrefix(locale, "uk"), strings.HasPrefix(locale, "bg"), strings.HasPrefix(locale, "sr"):
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
func (fs *FontLibrary) selectFallbackFontByCoverage(preferred unicodeRanges, fontWeight xfont.Weight) *FontInfo {
	var best *FontInfo = nil
	bestScore := -1

	for _, fonts := range fs.fonts {
		for _, fi := range fonts {
			if fi == nil || fi.FontPath == "" {
				continue
			}
			score := fs.coverageFallbackScore(fi, preferred, fontWeight)
			if fi.Italic {
				score -= 20000
			}
			if score > bestScore {
				best = fi
				bestScore = score
			}
		}
	}
	return best
}

// coverageFallbackScore 计算 fallback 候选评分。
// coverageFallbackScore computes score for fallback candidate selection.
func (fs *FontLibrary) coverageFallbackScore(fi *FontInfo, preferred unicodeRanges, fontWeight xfont.Weight) int {
	preferredCoverage := fi.coverageRanges.IntersectionCount(preferred)
	breadthScore := len(fi.coverageRanges) * 10
	weightScore := int(fontWightDelta - misc.Abs(fi.Weight-fontWeight))
	return preferredCoverage + breadthScore + weightScore
}
