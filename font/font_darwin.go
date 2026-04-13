//go:build darwin

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
	return []string{"Helvetica Neue", "Arial Unicode MS", "Arial"}
}

func localeFallbackFamilyTable(locale string) []string {
	var table = map[string][]string{
		"zh_cn":   {"PingFang SC", "Heiti SC", "STHeiti", "Arial Unicode MS"},
		"zh_sg":   {"PingFang SC", "Heiti SC", "STHeiti", "Arial Unicode MS"},
		"zh_tw":   {"PingFang TC", "Heiti TC", "STHeiti", "Arial Unicode MS"},
		"zh_hk":   {"PingFang HK", "Heiti TC", "STHeiti", "Arial Unicode MS"},
		"zh_hant": {"PingFang TC", "PingFang HK", "Heiti TC", "Arial Unicode MS"},
		"zh_hans": {"PingFang SC", "Heiti SC", "STHeiti", "Arial Unicode MS"},
		"ja_jp":   {"Hiragino Sans", "Hiragino Kaku Gothic ProN", "Osaka"},
		"ko_kr":   {"Apple SD Gothic Neo", "NanumGothic"},
		"ar":      {"Geeza Pro", "Arial Unicode MS", "Helvetica Neue"},
		"he":      {"Arial Hebrew", "Arial", "Helvetica Neue"},
		"th":      {"Thonburi", "Helvetica Neue", "Arial Unicode MS"},
		"hi":      {"Kohinoor Devanagari", "Devanagari Sangam MN", "Arial Unicode MS"},
		"bn":      {"Kohinoor Bangla", "Bangla Sangam MN", "Arial Unicode MS"},
		"ta":      {"Kohinoor Tamil", "Tamil Sangam MN", "Arial Unicode MS"},
		"te":      {"Kohinoor Telugu", "Telugu Sangam MN", "Arial Unicode MS"},
		"gu":      {"Kohinoor Gujarati", "Gujarati Sangam MN", "Arial Unicode MS"},
		"pa":      {"Kohinoor Gurmukhi", "Gurmukhi MN", "Arial Unicode MS"},
		"ml":      {"Malayalam Sangam MN", "Arial Unicode MS", "Helvetica Neue"},
		"mr":      {"Kohinoor Devanagari", "Devanagari Sangam MN", "Arial Unicode MS"},
		"ru":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"uk":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"bg":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"sr":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"el":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"vi":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"tr":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"pl":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"cs":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"hu":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"ro":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"id":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"ms":      {"Helvetica Neue", "Arial", "Arial Unicode MS"},
		"en":      {"SF Pro Text", "Helvetica Neue", "Arial"},
	}

	v, ok := table[locale]
	if !ok {
		if i := strings.IndexByte(locale, '_'); i > 0 {
			return table[locale[:i]]
		}
	}
	return v
}
