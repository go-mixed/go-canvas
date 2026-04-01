package font

import (
	"cmp"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/golang/freetype/truetype"
	xfont "golang.org/x/image/font"
)

type FontInfo struct {
	Bold              FontWeight
	Italic            bool
	Family, SubFamily string
	FontPath          string

	TruetypeFont *truetype.Font
}

type FontLibrary struct {
	fonts      map[string][]*FontInfo
	matchCache map[string]*FontInfo
	faceCache  map[string]xfont.Face

	fallbackLoaded      bool
	fallbackRegularInfo *FontInfo
	fallbackBoldInfo    *FontInfo
	fallbackLightInfo   *FontInfo
}

func NewFontLibrary(paths ...string) *FontLibrary {
	list := LoadFonts(paths...)

	return &FontLibrary{
		fonts:      list,
		matchCache: make(map[string]*FontInfo),
		faceCache:  make(map[string]xfont.Face),
	}
}

// MatchOrFeedback 从字体列表中匹配最合适的字体
// weight: 粗细数值 (100-900)，italic: 是否斜体
func (fs *FontLibrary) MatchOrFeedback(fontFamily string, weight FontWeight, italic bool) *FontInfo {
	cacheKey := fontFamily + "|" + strconv.Itoa(int(weight)) + "|" + strconv.FormatBool(italic)
	if fi, ok := fs.matchCache[cacheKey]; ok {
		return fi
	}

	if fontFamily == "" {
		return fs.fallbackFontInfo(fontFamily, weight, italic)
	}

	type fontScore struct {
		f           *FontInfo
		familyScore float32
		weightScore float32
		italicScore float32
	}
	var matches []fontScore

	var familySimilarity, weightSimilarity, italicSimilarity float32
	// font family 完全匹配 10
	for family, fonts := range fs.fonts {
		if familySimilarity = FontFamilySimilarity(family, fontFamily); familySimilarity > 0.5 {
			for _, font := range fonts {
				weightSimilarity = 1. - float32(misc.Abs(font.Bold-weight))/1000.
				italicSimilarity = 0
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
	}

	if len(matches) == 0 {
		fi := fs.fallbackFontInfo(fontFamily, weight, italic)
		fs.matchCache[cacheKey] = fi
		return fi
	}

	// 多字段排序（family DESC, weight DESC, italic DESC）
	slices.SortFunc(matches, func(a, b fontScore) int {
		family := cmp.Compare(a.familyScore, b.familyScore)
		if family != 0 {
			return -family
		}
		weight := cmp.Compare(a.weightScore, b.weightScore)
		if weight != 0 {
			return -weight
		}
		return -cmp.Compare(a.italicScore, b.italicScore)
	})

	if len(matches) > 0 {
		fi := matches[0].f
		fs.matchCache[cacheKey] = fi
		return fi
	}

	fi := fs.fallbackFontInfo(fontFamily, weight, italic)
	fs.matchCache[cacheKey] = fi
	return fi
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

func (f *FontInfo) GetTrueTypeFont() (*truetype.Font, error) {
	if f.TruetypeFont != nil {
		return f.TruetypeFont, nil
	}

	data, err := os.ReadFile(f.FontPath)
	if err != nil {
		return nil, err
	}

	tf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}
	f.TruetypeFont = tf
	return tf, nil
}

func FontFamilySimilarity(family1, family2 string) float32 {
	if strings.EqualFold(family1, family2) {
		return 1
	}

	return 0
}

func (fs *FontLibrary) GetFace(fi *FontInfo, fontSize int) xfont.Face {
	if fi == nil {
		return nil
	}
	k := fontFaceKey(fi, fontSize)
	if face, ok := fs.faceCache[k]; ok {
		return face
	}
	face := fs.CreateFace(fi, fontSize)
	fs.faceCache[k] = face
	return face
}

func fontFaceKey(fi *FontInfo, size int) string {
	if fi != nil && fi.FontPath != "" {
		return fi.FontPath + "-" + strconv.Itoa(size)
	}
	if fi != nil {
		return fi.Family + "-" + strconv.Itoa(size)
	}
	return strconv.Itoa(size)
}

func (fs *FontLibrary) CreateFace(fi *FontInfo, fontSize int) xfont.Face {
	tf, _ := fi.GetTrueTypeFont()
	return truetype.NewFace(tf, &truetype.Options{
		Size:    float64(fontSize),
		DPI:     120,
		Hinting: xfont.HintingFull,
	})
}
