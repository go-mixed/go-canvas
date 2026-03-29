//go:build darwin

package font

func GetSystemFontDirectories() (paths []string) {
	return []string{
		expandUser("~/Library/Fonts/"),
		"/Library/Fonts/",
		"/System/Library/Fonts/",
	}
}
