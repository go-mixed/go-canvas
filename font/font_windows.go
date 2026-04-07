//go:build windows

package font

import (
	"os"
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
	return []string{"Arial"}
}
