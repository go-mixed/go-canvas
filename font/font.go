package font

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/go-mixed/go-canvas/internel/misc"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type FontInfo struct {
	Bold              FontWeight
	Italic            bool
	Family, SubFamily string
	FontPath          string
	FaceIndex         int

	OpenTypeFont   *opentype.Font
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

	fontOptions *FontOptions

	mutex *sync.RWMutex
}

func NewFontLibrary(fontOpt *FontOptions, paths ...string) (*FontLibrary, error) {
	fs := &FontLibrary{
		fonts:       make(map[string][]*FontInfo),
		matchCache:  make(map[string]fontCollection),
		faceCache:   make(map[string]xfont.Face),
		mutex:       &sync.RWMutex{},
		fontOptions: fontOpt,
	}
	fs.fonts = fs.loadFonts(paths...)
	if len(fs.fonts) == 0 {
		return nil, fmt.Errorf("no fonts found in %v", paths)
	}
	return fs, fs.initFallbackPaths()
}

type fontScore struct {
	f           *FontInfo
	familyScore float32
	weightScore float32
	italicScore float32
}

func (fs *FontLibrary) logf(format string, args ...any) {
	if fs.fontOptions.logger != nil {
		fs.fontOptions.logger.Printf(format, args...)
	}
}

// MatchOrFeedback 从字体列表中匹配最合适的字体
// weight: 粗细数值 (100-900)，italic: 是否斜体
func (fs *FontLibrary) MatchOrFeedback(fontFamily string, weight FontWeight, italic bool) *FontInfo {
	if fontFamily == "" {
		return fs.fallbackFontInfo(weight)
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
		fi = fs.fallbackFontInfo(weight)
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

// rankFonts 为候选字体做排序匹配：
// - rn == 0：用于常规匹配（不做字符覆盖校验），按 family/weight/italic 排序。
// - rn != 0：用于缺字补字（必须覆盖该 rune），再按 weight/italic/family 排序。
// rankFonts ranks candidates for two paths:
// - rn == 0: base matching without glyph-coverage filtering.
// - rn != 0: rune-aware fallback matching that requires coverage for the rune.
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

	fontStableKey := func(fi *FontInfo) string {
		if fi == nil {
			return ""
		}
		return fi.FontPath + "|" + strconv.Itoa(fi.FaceIndex) + "|" + fi.Family + "|" + fi.SubFamily
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
			familyCmp := cmp.Compare(a.familyScore, b.familyScore)
			if familyCmp != 0 {
				return -familyCmp
			}
			return cmp.Compare(fontStableKey(a.f), fontStableKey(b.f))
		}

		family := cmp.Compare(a.familyScore, b.familyScore)
		if family != 0 {
			return -family
		}
		weightCmp := cmp.Compare(a.weightScore, b.weightScore)
		if weightCmp != 0 {
			return -weightCmp
		}
		italicCmp := cmp.Compare(a.italicScore, b.italicScore)
		if italicCmp != 0 {
			return -italicCmp
		}
		return cmp.Compare(fontStableKey(a.f), fontStableKey(b.f))
	})
	return matches
}

func fontFamilySimilarity(family1, family2 string) float32 {
	n1 := family1
	n2 := family2
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
	fs.mutex.RLock()
	face, ok := fs.faceCache[k]
	fs.mutex.RUnlock()

	if ok {
		return face
	}
	face = fs.CreateFace(fi, fontSize)

	fs.mutex.Lock()
	fs.faceCache[k] = face
	fs.mutex.Unlock()
	return face
}

func (fs *FontLibrary) fontFaceKey(fi *FontInfo, size int) string {
	if fi != nil && fi.FontPath != "" {
		return fi.FontPath + "-" + strconv.Itoa(fi.FaceIndex) + "-" + strconv.Itoa(size)
	}
	if fi != nil {
		return fi.Family + "-" + strconv.Itoa(fi.FaceIndex) + "-" + strconv.Itoa(size)
	}
	return strconv.Itoa(size)
}

func (fs *FontLibrary) CreateFace(fi *FontInfo, fontSize int) xfont.Face {
	of, err := fi.GetOpenTypeFont()
	if err != nil || of == nil {
		return nil
	}
	face, err := opentype.NewFace(of, &opentype.FaceOptions{
		Size: float64(fontSize),
		// px = size * dpi / 72
		DPI:     fs.fontOptions.dpi,
		Hinting: xfont.HintingFull,
	})
	if err != nil {
		return nil
	}
	return face
}

func (fs *FontLibrary) findFontByFamilies(families []string, fontWeight FontWeight) *FontInfo {
	for _, family := range families {
		if fi := fs.findFontByWeight(family, fontWeight); fi != nil {
			return fi
		}
	}
	return nil
}

// findFontByWeight 是 family 内精准匹配入口：
// 仅在指定 family 中按目标字重挑选最接近字体，不参与 rune 覆盖判断与全局排序。
// findFontByWeight is the precise family-local matcher:
// it picks the closest weight inside one family only, without rune coverage checks.
func (fs *FontLibrary) findFontByWeight(family string, fontWeight FontWeight) *FontInfo {
	want := family
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
		score := int(1000 - misc.Abs(fi.Bold-fontWeight))
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
