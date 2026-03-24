package ti

import (
	"image/color"
)

// BGR 0xBBGGRR
type BGR uint32

var _ color.Color = (*BGR)(nil)

// u8Colorto16 将 8 位颜色转换为 16 位颜色（0-255 -> 0-65535）
func u8Colorto16(v uint32) uint32 {
	return v * 0xffff / 0xff
}

func u16ColorTo8(v uint32) uint32 {
	return v * 0xff / 0xffff
}

func (C BGR) RGBA() (r, g, b, a uint32) {
	return u8Colorto16(uint32(C>>16) & 0xff), u8Colorto16(uint32(C>>8) & 0xff), u8Colorto16(uint32(C) & 0xff), 0xffff
}

// BGRA 0xBBGGRRAA
type BGRA uint32

var _ color.Color = (*BGRA)(nil)

func (C BGRA) RGBA() (r, g, b, a uint32) {
	return u8Colorto16(uint32(C>>8) & 0xff), u8Colorto16(uint32(C>>16) & 0xff), u8Colorto16(uint32(C>>24) & 0xff), u8Colorto16(uint32(C) & 0xff)
}

// ARGB 0xAARRGGBB
type ARGB uint32

var _ color.Color = (*ARGB)(nil)

func (A ARGB) RGBA() (r, g, b, a uint32) {
	return u8Colorto16(uint32(A>>16) & 0xff), u8Colorto16(uint32(A>>8) & 0xff), u8Colorto16(uint32(A) & 0xff), u8Colorto16(uint32(A>>24) & 0xff)
}

// RGBA 0xRRGGBBAA
type RGBA uint32

var _ color.Color = (*RGBA)(nil)

func (R RGBA) RGBA() (r, g, b, a uint32) {
	return u8Colorto16(uint32(R>>24) & 0xff), u8Colorto16(uint32(R>>16) & 0xff), u8Colorto16(uint32(R>>8) & 0xff), u8Colorto16(uint32(R) & 0xff)
}

// TiColor2Color 将 Taichi 纹理颜色转换为 Go 颜色
func TiColor2Color(r, g, b, a float32) color.Color {
	normalized := func(v float32) uint32 {
		if v > 0.999999 {
			return 255
		} else if v < 0.000001 {
			return 0
		}
		return uint32(v * 255)
	}
	return RGBA(normalized(r)<<24 | normalized(g)<<16 | normalized(b)<<8 | normalized(a))
}

// ExpandFColor 将 Go 颜色转换为 0-1 float 颜色
func ExpandFColor(color color.Color) (r, g, b, a float32) {
	ir, ig, ib, ia := color.RGBA()
	return float32(ir) / 65535., float32(ig) / 65535., float32(ib) / 65535., float32(ia) / 65535.
}

// ExpandUColor 将 Go 颜色转换为 0xff 颜色
func ExpandUColor(color color.Color) (r, g, b, a uint32) {
	ir, ig, ib, ia := color.RGBA()
	return u16ColorTo8(ir), u16ColorTo8(ig), u16ColorTo8(ib), u16ColorTo8(ia)
}

func Color2TiColor(color color.Color) TiColor {
	r, g, b, a := ExpandFColor(color)
	return TiColor{r, g, b, a}
}

const ColorWhite = RGBA(0xffffffff)
const ColorBlack = RGBA(0x000000ff)
const ColorRed = RGBA(0xff0000ff)
const ColorGreen = RGBA(0x00ff00ff)
const ColorBlue = RGBA(0x0000ffff)
const ColorYellow = RGBA(0xffff00ff)
const ColorMagenta = RGBA(0xff00ffff)
const ColorCyan = RGBA(0x00ffffff)
const ColorGray = RGBA(0x808080ff)
const ColorTransparent = RGBA(0x00000000)
