package font

import (
	"os"

	"github.com/golang/freetype/truetype"
)

// FontInfo represents a system font.
type FontInfo struct {
	// Family contains name of the font family.
	Family string

	// Name contains the full name of the font.
	Name string

	// Filename contains the path of the font file.
	Filename string
}

type TrueTypeFont struct {
	*truetype.Font
	fontInfo *FontInfo
}

func (f *TrueTypeFont) Initial() error {
	if f.Font != nil {
		return nil
	}

	path, err := Find(f.fontInfo.Filename)
	if err != nil {
		return err
	}
	f.fontInfo.Filename = path
	data, err := os.ReadFile(f.fontInfo.Filename)
	if err != nil {
		return err
	}

	f.Font, err = truetype.Parse(data)
	return err
}

func TryFindFont(fontName string, bold, italic bool) *TrueTypeFont {
	f := &TrueTypeFont{
		fontInfo: &FontInfo{
			Family:   fontName,
			Name:     fontName,
			Filename: "msyh.ttf",
		},
	}
	f.Initial()
	return f
}
