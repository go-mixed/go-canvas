package misc

import (
	"os"
	"path/filepath"
)

func GetCurrentDir() string {
	p, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(p)
}
