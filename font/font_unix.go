//go:build unix && !darwin

package font

func fallbackFontInfo(fontFamily string, weight FontWeight, italic bool) *FontInfo {
	var path string
	if weight == FontWeightRegular {
		path, _ = FindFont("LiberationSans-Regular.ttf")
		if path == "" {
			path, _ = FindFont("DejaVuSans.ttf")
		}
	} else if weight > FontWeightRegular {
		path, _ = FindFont("LiberationSans-Bold.ttf")
		if path == "" {
			path, _ = FindFont("DejaVuSans-Bold.ttf")
		}
	} else {
		path, _ = FindFont("LiberationSans-Light.ttf")
		if path == "" {
			path, _ = FindFont("DejaVuSans.ttf")
		}
	}
	return &FontInfo{
		Family:   fontFamily,
		Bold:     weight,
		Italic:   false,
		FontPath: path,
	}
}
