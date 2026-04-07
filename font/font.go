package font

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/golang/freetype/truetype"
	xfont "golang.org/x/image/font"
)

type FontInfo struct {
	Bold              FontWeight
	Italic            bool
	Family, SubFamily string
	FontPath          string

	TruetypeFont   *truetype.Font
	coverageRanges unicodeRanges
}

type FontLibrary struct {
	fonts      map[string][]*FontInfo
	matchCache map[string]fontCollection
	faceCache  map[string]xfont.Face

	fallbackLoaded      bool
	fallbackRegularInfo *FontInfo
	fallbackBoldInfo    *FontInfo
	fallbackLightInfo   *FontInfo

	mutex *sync.RWMutex
}

func NewFontLibrary(paths ...string) *FontLibrary {
	fs := &FontLibrary{
		matchCache: make(map[string]fontCollection),
		faceCache:  make(map[string]xfont.Face),
		mutex:      &sync.RWMutex{},
	}
	fs.fonts = fs.loadFonts(paths...)
	return fs
}

type fontScore struct {
	f           *FontInfo
	familyScore float32
	weightScore float32
	italicScore float32
}

// MatchOrFeedback 从字体列表中匹配最合适的字体
// weight: 粗细数值 (100-900)，italic: 是否斜体
func (fs *FontLibrary) MatchOrFeedback(fontFamily string, weight FontWeight, italic bool) *FontInfo {
	if fontFamily == "" {
		return fs.fallbackFontInfo(fontFamily, weight, italic)
	}

	cacheKey := fontFamily + "|" + strconv.Itoa(int(weight)) + "|" + strconv.FormatBool(italic)
	fs.mutex.RLock()
	fc, ok := fs.matchCache[cacheKey]
	fs.mutex.RUnlock()

	if ok && fc.baseFont != nil {
		return fc.baseFont
	}

	candidates := fs.rankFonts(fontFamily, weight, italic, 0)

	var fi *FontInfo
	if len(candidates) == 0 {
		fi = fs.fallbackFontInfo(fontFamily, weight, italic)
	} else {
		fi = candidates[0].f
	}

	fs.mutex.Lock()
	fc.baseFont = fi
	fs.matchCache[cacheKey] = fc
	fs.mutex.Unlock()
	return fi
}

// MatchRuneOrFeedback 为缺字 rune 选择降级字体。
// MatchRuneOrFeedback picks fallback font for a missing-glyph rune.
func (fs *FontLibrary) MatchRuneOrFeedback(base *FontInfo, rn rune) *FontInfo {
	cacheKey := base.Family + "|" + strconv.Itoa(int(base.Bold)) + "|" + strconv.FormatBool(base.Italic)
	fs.mutex.RLock()
	fc, ok := fs.matchCache[cacheKey]
	fs.mutex.RUnlock()

	if ok {
		if fi := fc.matchRuneFont(rn); fi != nil {
			return fi
		}
	}

	candidates := fs.rankFonts(base.Family, base.Bold, base.Italic, rn)

	var fi *FontInfo
	if len(candidates) == 0 {
		fi = base
	} else {
		fi = candidates[0].f
	}

	fs.mutex.Lock()
	fc.appendRuneFont(fi)
	fs.matchCache[cacheKey] = fc
	fs.mutex.Unlock()
	return fi

}

// rankFonts 为主字体匹配排序：家族优先，其次粗细和斜体。
// rankFonts ranks base font candidates: family first, then weight/italic.
func (fs *FontLibrary) rankFonts(fontFamily string, weight FontWeight, italic bool, rn rune) []fontScore {
	var matches []fontScore

	var detectedRune = rn > 0
	for family, fonts := range fs.fonts {
		familySimilarity := fontFamilySimilarity(family, fontFamily)
		// rune 无值时，对字体的要求严格
		if !detectedRune && familySimilarity <= 0.3 {
			continue
		}
		for _, font := range fonts {
			// rune 有值时，需要必须匹配范围
			if detectedRune && !font.coverageRanges.SupportsRune(rn) {
				continue
			}
			weightSimilarity := 1. - float32(misc.Abs(font.Bold-weight))/1000.
			italicSimilarity := float32(0)
			if italic == font.Italic {
				italicSimilarity = 1
			}
			matches = append(matches, fontScore{
				f:           font,
				familyScore: familySimilarity,
				weightScore: weightSimilarity,
				italicScore: italicSimilarity,
			})
		}
	}

	// 当有rune时，排序为 weight desc, italic desc, family desc
	// 当无rune时，排序为 family desc, weight desc, italic desc
	slices.SortFunc(matches, func(a, b fontScore) int {
		if detectedRune {
			weightCmp := cmp.Compare(a.weightScore, b.weightScore)
			if weightCmp != 0 {
				return -weightCmp
			}
			italicCmp := cmp.Compare(a.italicScore, b.italicScore)
			if italicCmp != 0 {
				return -italicCmp
			}
			return -cmp.Compare(a.familyScore, b.familyScore)
		}

		family := cmp.Compare(a.familyScore, b.familyScore)
		if family != 0 {
			return -family
		}
		weightCmp := cmp.Compare(a.weightScore, b.weightScore)
		if weightCmp != 0 {
			return -weightCmp
		}
		return -cmp.Compare(a.italicScore, b.italicScore)
	})
	return matches
}

func (fs *FontLibrary) registerFallbackFamilyAlias(fontFamily string) {
	key := fontFamily
	if key == "" {
		key = "fallback"
	}
	if _, ok := fs.fonts[key]; ok {
		return
	}
	fs.fonts[key] = []*FontInfo{
		fs.fallbackRegularInfo,
		fs.fallbackBoldInfo,
		fs.fallbackLightInfo,
	}
}

func fontFamilySimilarity(family1, family2 string) float32 {
	n1 := normalizeFamilyName(family1)
	n2 := normalizeFamilyName(family2)
	if n1 == "" || n2 == "" {
		return 0
	}
	if n1 == n2 {
		return 1
	}

	// Progressive token fallback:
	// e.g. "Noto Sans Condensed ExtraLight Italic" -> ... -> "Noto Sans".
	a1 := familyAncestors(n1)
	a2 := familyAncestors(n2)
	best := float32(0)
	for i, x := range a1 {
		for j, y := range a2 {
			if x != y {
				continue
			}
			trimPenalty := 0.08 * float32(i+j)
			score := float32(1.0) - trimPenalty
			if score < 0.1 {
				score = 0.1
			}
			if score > best {
				best = score
			}
		}
	}
	if best > 0 {
		return best
	}

	// Weak fallback for partial contain.
	if strings.Contains(n1, n2) || strings.Contains(n2, n1) {
		return 0.3
	}
	return 0
}

func normalizeFamilyName(s string) string {
	fields := strings.Fields(strings.ToLower(strings.TrimSpace(s)))
	if len(fields) == 0 {
		return ""
	}
	return strings.Join(fields, " ")
}

func familyAncestors(normalized string) []string {
	parts := strings.Split(normalized, " ")
	out := make([]string, 0, len(parts))
	for n := len(parts); n >= 1; n-- {
		out = append(out, strings.Join(parts[:n], " "))
	}
	return out
}

func (fs *FontLibrary) GetFace(fi *FontInfo, fontSize int) xfont.Face {
	if fi == nil {
		return nil
	}
	k := fs.fontFaceKey(fi, fontSize)
	if face, ok := fs.faceCache[k]; ok {
		return face
	}
	face := fs.CreateFace(fi, fontSize)
	fs.faceCache[k] = face
	return face
}

func (fs *FontLibrary) fontFaceKey(fi *FontInfo, size int) string {
	if fi != nil && fi.FontPath != "" {
		return fi.FontPath + "-" + strconv.Itoa(size)
	}
	if fi != nil {
		return fi.Family + "-" + strconv.Itoa(size)
	}
	return strconv.Itoa(size)
}

func (fs *FontLibrary) CreateFace(fi *FontInfo, fontSize int) xfont.Face {
	tf, err := fi.GetTrueTypeFont()
	if err != nil || tf == nil {
		return nil
	}
	return truetype.NewFace(tf, &truetype.Options{
		Size:    float64(fontSize),
		DPI:     120,
		Hinting: xfont.HintingNone,
	})
}
