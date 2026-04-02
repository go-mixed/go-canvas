package ti

import (
	"image/color"
)

// BGR 0xBBGGRR
type BGR uint32

var _ color.Color = (*BGR)(nil)

// u8ColorTo16 将 8 位颜色转换为 16 位颜色（0-255 -> 0-65535）
func u8ColorTo16(v uint32) uint32 {
	return v * 0xffff / 0xff
}

func u16ColorTo8(v uint32) uint32 {
	return v * 0xff / 0xffff
}

func f32ColorToU8(v float32) uint32 {
	if v > 0.999999 {
		return 255
	} else if v < 0.000001 {
		return 0
	}
	return uint32(v * 255)
}

func (C BGR) RGBA() (r, g, b, a uint32) {
	return u8ColorTo16(uint32(C>>16) & 0xff), u8ColorTo16(uint32(C>>8) & 0xff), u8ColorTo16(uint32(C) & 0xff), 0xffff
}

// BGRA 0xBBGGRRAA
type BGRA uint32

var _ color.Color = (*BGRA)(nil)

func (C BGRA) RGBA() (r, g, b, a uint32) {
	return u8ColorTo16(uint32(C>>8) & 0xff), u8ColorTo16(uint32(C>>16) & 0xff), u8ColorTo16(uint32(C>>24) & 0xff), u8ColorTo16(uint32(C) & 0xff)
}

func ToBGRA(color color.Color) BGRA {
	r, g, b, a := color.RGBA()
	return BGRA(u16ColorTo8(b)<<8 | u16ColorTo8(g)<<16 | u16ColorTo8(r)<<24 | u16ColorTo8(a))
}

// ARGB 0xAARRGGBB
type ARGB uint32

var _ color.Color = (*ARGB)(nil)

func (A ARGB) RGBA() (r, g, b, a uint32) {
	return u8ColorTo16(uint32(A>>16) & 0xff), u8ColorTo16(uint32(A>>8) & 0xff), u8ColorTo16(uint32(A) & 0xff), u8ColorTo16(uint32(A>>24) & 0xff)
}

func ToARGB(color color.Color) ARGB {
	r, g, b, a := color.RGBA()
	return ARGB(u16ColorTo8(a)<<24 | u16ColorTo8(r)<<16 | u16ColorTo8(g)<<8 | u16ColorTo8(b))
}

// RGBA 0xRRGGBBAA
type RGBA uint32

var _ color.Color = (*RGBA)(nil)

func (R RGBA) RGBA() (r, g, b, a uint32) {
	return u8ColorTo16(uint32(R>>24) & 0xff), u8ColorTo16(uint32(R>>16) & 0xff), u8ColorTo16(uint32(R>>8) & 0xff), u8ColorTo16(uint32(R) & 0xff)
}

func ToRGBA(color color.Color) RGBA {
	r, g, b, a := color.RGBA()
	return RGBA(u16ColorTo8(r)<<24 | u16ColorTo8(g)<<16 | u16ColorTo8(b)<<8 | u16ColorTo8(a))
}

// TiColorToColor 将 Taichi 纹理颜色转换为 Go 颜色
func TiColorToColor(r, g, b, a float32) color.Color {
	return RGBA(f32ColorToU8(r)<<24 | f32ColorToU8(g)<<16 | f32ColorToU8(b)<<8 | f32ColorToU8(a))
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

// ColorEqual 比较两个 color.Color 的 RGBA 展开值是否一致。
// ColorEqual compares two color.Color values by expanded RGBA components.
func ColorEqual(a, b color.Color) bool {
	if a == nil || b == nil {
		return a == b
	}
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}
