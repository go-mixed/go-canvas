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
		if v > 0.9999999 {
			return 255
		} else if v < 0.00001 {
			return 0
		}
		return uint32(v * 255)
	}
	return RGBA(normalized(r)<<24 | normalized(g)<<16 | normalized(b)<<8 | normalized(a))
}

// Color2TiColor 将 Go 颜色转换为 Taichi 纹理颜色
func Color2TiColor(color color.Color) (r, g, b, a float32) {
	ir, ig, ib, ia := color.RGBA()
	return float32(ir) / 65535., float32(ig) / 65535., float32(ib) / 65535., float32(ia) / 65535.
}
