//go:build windows

package font

func fallbackFontInfo(fontFamily string, weight FontWeight, italic bool) *FontInfo {
	var path string
	if weight == FontWeightRegular {
		path, _ = FindFont("msyh.ttc")
	} else if weight > FontWeightRegular {
		path, _ = FindFont("msyhbd.ttc")
	} else {
		path, _ = FindFont("msyhl.ttc")
	}
	return &FontInfo{
		Family:   fontFamily,
		Bold:     weight,
		Italic:   false,
		FontPath: path,
	}
}
