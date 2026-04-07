//go:build !windows

package font

import "os"

func detectSystemLanguage() string {
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG", "LANGUAGE"} {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return ""
}

func systemFallbackFamilies() []string {
	return []string{"DejaVu Sans", "Liberation Sans", "Noto Sans"}
}
