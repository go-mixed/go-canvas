package font

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func defaultSuffixes() []string {
	return []string{".ttf", ".ttc", ".otf"}
}

var fontPathCache map[string]string = make(map[string]string)

// Find tries to locate the specified font file in the current directory as
// well as in platform specific user and system font directories; if there is
// no exact match, Find tries substring matching - files with the standard font suffixes (.ttf, .ttc, .otf) are considered.
func FindFont(fileName string) (filePath string, err error) {
	var ok bool
	if filePath, ok = fontPathCache[fileName]; ok {
		return filePath, nil
	}

	if filePath, err = FindWithSuffixes(fileName, defaultSuffixes()); err == nil {
		fontPathCache[fileName] = filePath
		return filePath, nil
	}
	return "", errors.Errorf("cannot find font '%s' in user or system directories", fileName)
}

// FindWithSuffixes tries to locate the specified font file in the current directory as
// well as in platform specific user and system font directories; if there is
// no exact match, Find tries substring matching - only font files with the give suffixes are considered.
func FindWithSuffixes(fileName string, suffixes []string) (filePath string, err error) {
	// check if fileName already points to a readable file
	if _, err := os.Stat(fileName); err == nil {
		return fileName, nil
	}

	// search in user and system directories
	return find(filepath.Base(fileName), suffixes)
}

// ListFont returns a list of all font files (determined by standard suffixes: .ttf, .ttc, .otf) found on the system.
func ListFont(dirs []string) (filePaths []string) {
	return ListFontWithSuffixes(dirs, defaultSuffixes())
}

// ListFontWithSuffixes returns a list of all font files (determined by given file suffixes) found on the system.
func ListFontWithSuffixes(dirs []string, suffixes []string) (filePaths []string) {
	var pathList []string

	walkF := func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if !info.IsDir() && isFontFile(path, suffixes) {
				pathList = append(pathList, path)
			}
		}
		return nil
	}

	for _, dir := range dirs {
		filepath.Walk(dir, walkF)
	}

	return pathList
}

func isFontFile(fileName string, suffixes []string) bool {
	lower := strings.ToLower(fileName)
	for _, suffix := range suffixes {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	return false
}

func stripExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func expandUser(path string) (expandedPath string) {
	if strings.HasPrefix(path, "~") {
		if u, err := user.Current(); err == nil {
			return strings.Replace(path, "~", u.HomeDir, -1)
		}
	}
	return path
}

func find(needle string, suffixes []string) (filePath string, err error) {
	lowerNeedle := strings.ToLower(needle)
	lowerNeedleBase := stripExtension(lowerNeedle)

	match := ""
	partial := ""
	partialScore := -1

	walkF := func(path string, info os.FileInfo, err error) error {
		// we have already found a match -> nothing to do
		if match != "" {
			return nil
		}
		if err != nil {
			return nil
		}

		lowerPath := strings.ToLower(info.Name())

		if !info.IsDir() && isFontFile(lowerPath, suffixes) {
			lowerBase := stripExtension(lowerPath)
			if lowerPath == lowerNeedle {
				// exact match
				match = path
			} else if strings.Contains(lowerBase, lowerNeedleBase) {
				// partial match
				score := len(lowerBase) - len(lowerNeedle)
				if partialScore < 0 || score < partialScore {
					partialScore = score
					partial = path
				}
			}
		}
		return nil
	}

	for _, dir := range GetSystemFontDirectories() {
		filepath.Walk(dir, walkF)
		if match != "" {
			return match, nil
		}
	}

	if partial != "" {
		return partial, nil
	}

	return "", fmt.Errorf("cannot find font '%s' in user or system directories", needle)
}
