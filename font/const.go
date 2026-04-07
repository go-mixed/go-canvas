package font

import "sort"

type FontWeight int16

// Font weight constants (CSS Font Weight 100-900)
const (
	FontWeightThin       FontWeight = 100
	FontWeightExtraLight FontWeight = 200
	FontWeightLight      FontWeight = 300
	FontWeightRegular    FontWeight = 400
	FontWeightMedium     FontWeight = 500
	FontWeightSemiBold   FontWeight = 600
	FontWeightBold       FontWeight = 700
	FontWeightExtraBold  FontWeight = 800
	FontWeightBlack      FontWeight = 900

	FontWeightMax = FontWeightBlack
)

var fontStyles = map[string]FontWeight{
	// Thin/Hairline (100)
	"thin":     FontWeightThin,
	"hairline": FontWeightThin,
	// ExtraLight/UltraLight (200)
	"extralight":  FontWeightExtraLight,
	"ultralight":  FontWeightExtraLight,
	"ultra light": FontWeightExtraLight,
	"ultrathin":   FontWeightExtraLight,
	// Light (300)
	"light": FontWeightLight,
	// Regular/Normal/Medium (400)
	"medium":  FontWeightRegular,
	"regular": FontWeightRegular,
	"normal":  FontWeightRegular,
	"book":    FontWeightRegular,
	"roman":   FontWeightRegular,
	"plain":   FontWeightRegular,
	// SemiBold/DemiBold (600)
	"semibold":  FontWeightSemiBold,
	"demibold":  FontWeightSemiBold,
	"demi bold": FontWeightSemiBold,
	"semi bold": FontWeightSemiBold,
	// Bold/Heavy (700)
	"bold":  FontWeightBold,
	"heavy": FontWeightBold,
	// ExtraBold (800)
	"extrabold":  FontWeightExtraBold,
	"ultra bold": FontWeightExtraBold,
	// Black/Ultra (900)
	"black": FontWeightBlack,
	"ultra": FontWeightBlack,
}

// italicStyles 记录斜体关键词
var italicStyles = map[string]bool{
	"italic":         true,
	"regularoblique": true,
	"oblique":        true,
	"slanted":        true,
	"kursiv":         true,
	"inclined":       true,
	"backslant":      true,
}

type unicodeRange struct {
	start rune
	end   rune
}

type unicodeRanges []unicodeRange

func (ur unicodeRanges) SupportsRune(r rune) bool {
	if r < 0 || r > 0x10FFFF {
		return false
	}
	if len(ur) == 0 {
		// Runtime must not parse coverage on demand.
		// If coverage is missing, treat as unsupported.
		return false
	}

	i := sort.Search(len(ur), func(i int) bool {
		return ur[i].end >= r
	})
	if i >= len(ur) {
		return false
	}
	rng := ur[i]
	return r >= rng.start && r <= rng.end
}

func (ur unicodeRanges) IntersectionCount(preferred unicodeRanges) int {
	if len(ur) == 0 || len(preferred) == 0 {
		return 0
	}
	i, j := 0, 0
	total := 0
	for i < len(ur) && j < len(preferred) {
		a := ur[i]
		b := preferred[j]
		start := a.start
		if b.start > start {
			start = b.start
		}
		end := a.end
		if b.end < end {
			end = b.end
		}
		if end >= start {
			total += int(end - start + 1)
		}
		if a.end < b.end {
			i++
		} else {
			j++
		}
	}
	return total
}

type fontCollection struct {
	baseFont  *FontInfo
	runeFonts []*FontInfo
}

func (c fontCollection) matchRuneFont(rn rune) *FontInfo {
	if c.baseFont != nil && c.baseFont.coverageRanges.SupportsRune(rn) {
		return c.baseFont
	}

	for _, runeFont := range c.runeFonts {
		if runeFont.coverageRanges.SupportsRune(rn) {
			return runeFont
		}
	}
	return nil
}

func (c fontCollection) appendRuneFont(font *FontInfo) {
	c.runeFonts = append(c.runeFonts, font)
}
