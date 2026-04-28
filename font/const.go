package font

import (
	"sort"
	"strings"

	xfont "golang.org/x/image/font"
)

type namedWeight struct {
	name   string
	weight xfont.Weight
}

var xfontStyles = []namedWeight{
	// Put longer / more specific tokens first to avoid partial shadowing:
	// e.g. "ultralight" must win before "light".
	{name: "ultra light", weight: xfont.WeightExtraLight},
	{name: "ultralight", weight: xfont.WeightExtraLight},
	{name: "extralight", weight: xfont.WeightExtraLight},
	{name: "ultrathin", weight: xfont.WeightExtraLight},
	{name: "semi bold", weight: xfont.WeightSemiBold},
	{name: "demi bold", weight: xfont.WeightSemiBold},
	{name: "semibold", weight: xfont.WeightSemiBold},
	{name: "demibold", weight: xfont.WeightSemiBold},
	{name: "ultra bold", weight: xfont.WeightExtraBold},
	{name: "extrabold", weight: xfont.WeightExtraBold},
	{name: "hairline", weight: xfont.WeightThin},
	{name: "medium", weight: xfont.WeightNormal},
	{name: "regular", weight: xfont.WeightNormal},
	{name: "normal", weight: xfont.WeightNormal},
	{name: "roman", weight: xfont.WeightNormal},
	{name: "plain", weight: xfont.WeightNormal},
	{name: "light", weight: xfont.WeightLight},
	{name: "black", weight: xfont.WeightBlack},
	{name: "heavy", weight: xfont.WeightBold},
	{name: "bold", weight: xfont.WeightBold},
	{name: "thin", weight: xfont.WeightThin},
	{name: "book", weight: xfont.WeightNormal},
	{name: "ultra", weight: xfont.WeightBlack},
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

func normalizeFamilyNames(ss []string) []string {
	var result []string = make([]string, len(ss))
	for i, s := range ss {
		result[i] = normalizeFamilyName(s)
	}
	return result
}

func normalizeFamilyName(s string) string {
	fields := strings.Fields(strings.ToLower(strings.TrimSpace(s)))
	if len(fields) == 0 {
		return ""
	}
	return strings.Join(fields, " ")
}
