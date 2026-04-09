//go:build windows

package font

import (
	"os"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

func detectSystemLanguage() string {
	if v := userDefaultLocaleName(); v != "" {
		return v
	}
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG", "LANGUAGE"} {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return ""
}

func userDefaultLocaleName() string {
	const localeNameMaxLength = 85
	k32 := windows.NewLazySystemDLL("kernel32.dll")
	proc := k32.NewProc("GetUserDefaultLocaleName")

	var buf [localeNameMaxLength]uint16
	r1, _, _ := proc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(localeNameMaxLength))
	if r1 == 0 {
		return ""
	}
	n := int(r1)
	if n <= 1 {
		return ""
	}
	return windows.UTF16ToString(buf[:n])
}

func systemFallbackFamilies() []string {
	return []string{"Segoe UI", "Arial", "Tahoma"}
}

func localeFallbackFamilyTable(locale string) []string {
	var table = map[string][]string{
		"zh_cn":   {"Microsoft YaHei UI", "Microsoft YaHei", "DengXian", "SimSun", "NSimSun"},
		"zh_sg":   {"Microsoft YaHei UI", "Microsoft YaHei", "DengXian", "SimSun", "NSimSun"},
		"zh_tw":   {"Microsoft JhengHei UI", "Microsoft JhengHei", "PMingLiU", "MingLiU"},
		"zh_hk":   {"Microsoft JhengHei UI", "Microsoft JhengHei", "MingLiU_HKSCS", "MingLiU"},
		"zh_hant": {"Microsoft JhengHei UI", "Microsoft JhengHei", "PMingLiU", "MingLiU"},
		"zh_hans": {"Microsoft YaHei UI", "Microsoft YaHei", "DengXian", "SimSun", "NSimSun"},
		"ja_jp":   {"Yu Gothic UI", "Meiryo", "MS UI Gothic"},
		"ko_kr":   {"Malgun Gothic", "Gulim", "Dotum"},
		"ar":      {"Segoe UI", "Tahoma", "Arial"},
		"he":      {"Segoe UI", "Arial", "Tahoma"},
		"th":      {"Leelawadee UI", "Tahoma", "Segoe UI"},
		"hi":      {"Nirmala UI", "Mangal", "Segoe UI"},
		"bn":      {"Nirmala UI", "Vrinda", "Segoe UI"},
		"ta":      {"Nirmala UI", "Latha", "Segoe UI"},
		"te":      {"Nirmala UI", "Gautami", "Segoe UI"},
		"ml":      {"Nirmala UI", "Kartika", "Segoe UI"},
		"gu":      {"Nirmala UI", "Shruti", "Segoe UI"},
		"pa":      {"Nirmala UI", "Raavi", "Segoe UI"},
		"mr":      {"Nirmala UI", "Mangal", "Segoe UI"},
		"ru":      {"Segoe UI", "Arial", "Tahoma"},
		"uk":      {"Segoe UI", "Arial", "Tahoma"},
		"bg":      {"Segoe UI", "Arial", "Tahoma"},
		"sr":      {"Segoe UI", "Arial", "Tahoma"},
		"el":      {"Segoe UI", "Arial", "Tahoma"},
		"vi":      {"Segoe UI", "Arial", "Tahoma"},
		"tr":      {"Segoe UI", "Arial", "Tahoma"},
		"pl":      {"Segoe UI", "Arial", "Tahoma"},
		"cs":      {"Segoe UI", "Arial", "Tahoma"},
		"hu":      {"Segoe UI", "Arial", "Tahoma"},
		"ro":      {"Segoe UI", "Arial", "Tahoma"},
		"id":      {"Segoe UI", "Arial", "Tahoma"},
		"ms":      {"Segoe UI", "Arial", "Tahoma"},
		"en":      {"Segoe UI", "Arial", "Tahoma"},
	}

	v, ok := table[locale]
	if !ok {
		if i := strings.IndexByte(locale, '_'); i > 0 {
			return table[locale[:i]]
		}
	}
	return v
}
