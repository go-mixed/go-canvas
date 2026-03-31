//go:build darwin

package font

func fallbackFontInfo(fontFamily string, weight FontWeight, italic bool) *FontInfo {
	var path string
	if weight == FontWeightRegular {
		path, _ = FindFont("Helvetica.ttc")
	} else if weight > FontWeightRegular {
		path, _ = FindFont("Helvetica Bold.ttc")
	} else {
		path, _ = FindFont("Helvetica Light.ttc")
	}
	return &FontInfo{
		Family:   fontFamily,
		Bold:     weight,
		Italic:   italic,
		FontPath: path,
	}
}
