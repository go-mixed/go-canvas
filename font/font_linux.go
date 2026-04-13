//go:build unix && !darwin

package font

import (
	"os"
	"strings"
)

func detectSystemLanguage() string {
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG", "LANGUAGE"} {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return ""
}

func systemFallbackFamilies() []string {
	return []string{"Noto Sans", "DejaVu Sans", "Liberation Sans"}
}

func localeFallbackFamilyTable(locale string) []string {
	var table = map[string][]string{
		"zh_cn":   {"Noto Sans CJK SC", "Source Han Sans SC", "WenQuanYi Micro Hei", "Noto Sans SC"},
		"zh_sg":   {"Noto Sans CJK SC", "Source Han Sans SC", "WenQuanYi Micro Hei", "Noto Sans SC"},
		"zh_tw":   {"Noto Sans CJK TC", "Source Han Sans TC", "Noto Sans TC"},
		"zh_hk":   {"Noto Sans CJK HK", "Source Han Sans HC", "Noto Sans HK"},
		"zh_hant": {"Noto Sans CJK TC", "Noto Sans CJK HK", "Source Han Sans TC"},
		"zh_hans": {"Noto Sans CJK SC", "Source Han Sans SC", "WenQuanYi Micro Hei"},
		"ja_jp":   {"Noto Sans CJK JP", "Source Han Sans", "IPAGothic", "VL Gothic"},
		"ko_kr":   {"Noto Sans CJK KR", "Source Han Sans K", "NanumGothic"},
		"ar":      {"Noto Sans Arabic", "Amiri", "DejaVu Sans"},
		"he":      {"Noto Sans Hebrew", "DejaVu Sans", "Liberation Sans"},
		"th":      {"Noto Sans Thai", "TLWG Typist", "DejaVu Sans"},
		"hi":      {"Noto Sans Devanagari", "Lohit Devanagari", "Kalimati"},
		"bn":      {"Noto Sans Bengali", "Lohit Bengali", "Mukti Narrow"},
		"ta":      {"Noto Sans Tamil", "Lohit Tamil", "Noto Serif Tamil"},
		"te":      {"Noto Sans Telugu", "Lohit Telugu", "Pothana2000"},
		"ml":      {"Noto Sans Malayalam", "Rachana", "Meera"},
		"gu":      {"Noto Sans Gujarati", "Lohit Gujarati", "Rekha"},
		"pa":      {"Noto Sans Gurmukhi", "Saab", "AnmolLipi"},
		"mr":      {"Noto Sans Devanagari", "Lohit Devanagari", "Kalimati"},
		"ru":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"uk":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"bg":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"sr":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"el":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"vi":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"tr":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"pl":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"cs":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"hu":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"ro":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"id":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"ms":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
		"en":      {"Noto Sans", "DejaVu Sans", "Liberation Sans"},
	}

	v, ok := table[locale]
	if !ok {
		if i := strings.IndexByte(locale, '_'); i > 0 {
			return table[locale[:i]]
		}
	}
	return v
}
