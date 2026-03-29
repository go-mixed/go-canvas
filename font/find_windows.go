//go:build windows

package font

import (
	"os"
	"path/filepath"
)

func GetSystemFontDirectories() (paths []string) {
	return []string{
		filepath.Join(os.Getenv("windir"), "Fonts"),
		filepath.Join(os.Getenv("localappdata"), "Microsoft", "Windows", "Fonts"),
	}
}
