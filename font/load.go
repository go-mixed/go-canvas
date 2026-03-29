package font

import (
	"strings"

	"github.com/golang/freetype/truetype"
)

// LoadFonts 读取给的目录+系统目录的所有字体
// userFontPaths 可以为字体路径、或者包含字体的目录
func LoadFonts(userFontPaths ...string) map[string][]*FontInfo {
	paths := ListFont(append(userFontPaths, GetSystemFontDirectories()...))

	var fontInfos = make(map[string][]*FontInfo)

	for _, path := range paths {
		info, err := ReadFontInfo(path)
		if err != nil {
			continue
		}
		v := fontInfos[info.Family]
		v = append(v, info)
		fontInfos[info.Family] = v
	}

	return fontInfos
}

func ReadFontInfo(path string) (*FontInfo, error) {
	info := &FontInfo{
		FontPath: path,
	}

	f, err := info.GetTrueTypeFont()
	if err != nil {
		return nil, err
	}

	// 从 truetype.Font 中提取 Family 名称
	// NameID: 1=FontFamily, 4=FullName, 6=PostScriptName
	info.Family = f.Name(truetype.NameIDFontFamily)
	if info.Family == "" {
		info.Family = f.Name(4) // FullName
	}

	// 判断粗体和斜体
	if info.SubFamily = f.Name(truetype.NameIDFontSubfamily); info.SubFamily != "" {
		info.Italic = isItalic(info.SubFamily)
	}

	// 从 fontStyles 匹配粗细数值
	info.Bold = matchWeight(info.SubFamily)
	return info, nil
}

func isItalic(subFamily string) bool {
	lower := strings.ToLower(subFamily)
	for style := range italicStyles {
		if strings.Contains(lower, style) {
			return true
		}
	}
	return false
}

func matchWeight(subFamily string) FontWeight {
	if subFamily == "" {
		return FontWeightRegular // 默认 Regular
	}
	lower := strings.ToLower(subFamily)
	var maxWeight FontWeight = 0
	for style, weight := range fontStyles {
		if strings.Contains(lower, style) && weight > maxWeight {
			maxWeight = weight
		}
	}
	if maxWeight == 0 {
		maxWeight = FontWeightRegular
	}
	return maxWeight
}
