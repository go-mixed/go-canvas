package font

type FontWeight int16

// Font weight constants (CSS Font Weight 100-900)
const (
	FontWeightThin       FontWeight = 100
	FontWeightExtraLight FontWeight = 200
	FontWeightLight      FontWeight = 300
	FontWeightRegular    FontWeight = 400
	FontWeightMedium     FontWeight = 500
	FontWeightSemiBold   FontWeight = 600
	FontWeightBold       FontWeight = 700
	FontWeightExtraBold  FontWeight = 800
	FontWeightBlack      FontWeight = 900

	FontWeightMax = FontWeightBlack
)

var fontStyles = map[string]FontWeight{
	// Thin/Hairline (100)
	"thin":     FontWeightThin,
	"hairline": FontWeightThin,
	// ExtraLight/UltraLight (200)
	"extralight":  FontWeightExtraLight,
	"ultralight":  FontWeightExtraLight,
	"ultra light": FontWeightExtraLight,
	"ultrathin":   FontWeightExtraLight,
	// Light (300)
	"light": FontWeightLight,
	// Regular/Normal/Medium (400)
	"medium":  FontWeightRegular,
	"regular": FontWeightRegular,
	"normal":  FontWeightRegular,
	"book":    FontWeightRegular,
	"roman":   FontWeightRegular,
	"plain":   FontWeightRegular,
	// SemiBold/DemiBold (600)
	"semibold":  FontWeightSemiBold,
	"demibold":  FontWeightSemiBold,
	"demi bold": FontWeightSemiBold,
	"semi bold": FontWeightSemiBold,
	// Bold/Heavy (700)
	"bold":  FontWeightBold,
	"heavy": FontWeightBold,
	// ExtraBold (800)
	"extrabold":  FontWeightExtraBold,
	"ultra bold": FontWeightExtraBold,
	// Black/Ultra (900)
	"black": FontWeightBlack,
	"ultra": FontWeightBlack,
}

// italicStyles 记录斜体关键词
var italicStyles = map[string]bool{
	"italic":         true,
	"regularoblique": true,
	"oblique":        true,
	"slanted":        true,
	"kursiv":         true,
	"inclined":       true,
	"backslant":      true,
}
