//go:build unix && !darwin

package font

func (fs *FontLibrary) initFallbackPaths() {
	if fs.fallbackLoaded {
		return
	}
	regular, _ := FindFont("LiberationSans-Regular.ttf")
	if regular == "" {
		regular, _ = FindFont("DejaVuSans.ttf")
	}
	bold, _ := FindFont("LiberationSans-Bold.ttf")
	if bold == "" {
		bold, _ = FindFont("DejaVuSans-Bold.ttf")
	}
	light, _ := FindFont("LiberationSans-Light.ttf")
	if light == "" {
		light, _ = FindFont("DejaVuSans.ttf")
	}

	fs.fallbackRegularInfo = &FontInfo{
		Family:   "fallback",
		Bold:     FontWeightRegular,
		Italic:   false,
		FontPath: regular,
	}
	fs.fallbackBoldInfo = &FontInfo{
		Family:   "fallback",
		Bold:     FontWeightBold,
		Italic:   false,
		FontPath: bold,
	}
	fs.fallbackLightInfo = &FontInfo{
		Family:   "fallback",
		Bold:     FontWeightLight,
		Italic:   false,
		FontPath: light,
	}
	fs.fallbackLoaded = true
}

func (fs *FontLibrary) fallbackFontInfo(fontFamily string, weight FontWeight, italic bool) *FontInfo {
	fs.initFallbackPaths()
	fs.registerFallbackFamilyAlias(fontFamily)

	if fontFamily != "" {
		fs.fallbackRegularInfo.Family = fontFamily
		fs.fallbackBoldInfo.Family = fontFamily
		fs.fallbackLightInfo.Family = fontFamily
	}

	if weight == FontWeightRegular {
		return fs.fallbackRegularInfo
	}
	if weight > FontWeightRegular {
		return fs.fallbackBoldInfo
	}
	return fs.fallbackLightInfo
}
