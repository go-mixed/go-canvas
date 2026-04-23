package ctypes

import (
	"image/color"

	"github.com/go-mixed/go-taichi/f16"
)

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

func f16ColorToU8(v f16.Float16) uint32 {
	vv := v.Float32()
	if vv > 0.999999 {
		return 255
	} else if vv < 0.000001 {
		return 0
	}
	return uint32(vv * 255)
}

// BGR 0xBBGGRR
type BGR uint32

var _ color.Color = (*BGR)(nil)

func (C BGR) RGBA() (r, g, b, a uint32) {
	return u8ColorTo16(uint32(C>>16) & 0xff), u8ColorTo16(uint32(C>>8) & 0xff), u8ColorTo16(uint32(C) & 0xff), 0xffff
}

// NBGRA 0xBBGGRRAA
type NBGRA uint32

var _ color.Color = (*NBGRA)(nil)

func (C NBGRA) RGBA() (r, g, b, a uint32) {
	// Go color.Color contract requires RGBA() to return alpha-premultiplied channels.
	// If we return straight RGB here, PNG/export paths may show colorful artifacts.
	r8 := uint32(C>>8) & 0xff
	g8 := uint32(C>>16) & 0xff
	b8 := uint32(C>>24) & 0xff
	a8 := uint32(C) & 0xff
	a = u8ColorTo16(a8)
	r = u8ColorTo16(r8) * a / 0xffff
	g = u8ColorTo16(g8) * a / 0xffff
	b = u8ColorTo16(b8) * a / 0xffff
	return
}

func ToNBGRA(color color.Color) NBGRA {
	r, g, b, a := color.RGBA()
	return NBGRA(u16ColorTo8(b)<<8 | u16ColorTo8(g)<<16 | u16ColorTo8(r)<<24 | u16ColorTo8(a))
}

// NARGB 0xAARRGGBB
type NARGB uint32

var _ color.Color = (*NARGB)(nil)

func (A NARGB) RGBA() (r, g, b, a uint32) {
	// Go color.Color contract requires RGBA() to return alpha-premultiplied channels.
	r8 := uint32(A>>16) & 0xff
	g8 := uint32(A>>8) & 0xff
	b8 := uint32(A) & 0xff
	a8 := uint32(A>>24) & 0xff
	a = u8ColorTo16(a8)
	r = u8ColorTo16(r8) * a / 0xffff
	g = u8ColorTo16(g8) * a / 0xffff
	b = u8ColorTo16(b8) * a / 0xffff
	return
}

func ToNARGB(color color.Color) NARGB {
	r, g, b, a := color.RGBA()
	return NARGB(u16ColorTo8(a)<<24 | u16ColorTo8(r)<<16 | u16ColorTo8(g)<<8 | u16ColorTo8(b))
}

// NRGBA 0xRRGGBBAA
type NRGBA uint32

var _ color.Color = (*NRGBA)(nil)

func (R NRGBA) RGBA() (r, g, b, a uint32) {
	// Go color.Color contract requires RGBA() to return alpha-premultiplied channels.
	r8 := uint32(R>>24) & 0xff
	g8 := uint32(R>>16) & 0xff
	b8 := uint32(R>>8) & 0xff
	a8 := uint32(R) & 0xff
	a = u8ColorTo16(a8)
	r = u8ColorTo16(r8) * a / 0xffff
	g = u8ColorTo16(g8) * a / 0xffff
	b = u8ColorTo16(b8) * a / 0xffff
	return
}

func ToRGBA(color color.Color) NRGBA {
	r, g, b, a := color.RGBA()
	return NRGBA(u16ColorTo8(r)<<24 | u16ColorTo8(g)<<16 | u16ColorTo8(b)<<8 | u16ColorTo8(a))
}

// TiColorToColor 将 Taichi 纹理颜色转换为 Go 颜色
func TiColorToColor(r, g, b, a float32) color.Color {
	return NRGBA(f32ColorToU8(r)<<24 | f32ColorToU8(g)<<16 | f32ColorToU8(b)<<8 | f32ColorToU8(a))
}

// ExpandF32Color 将 Go 颜色转换为 0-1 float32 颜色
func ExpandF32Color(c color.Color) (r, g, b, a float32) {
	// Use straight RGBA channels (not premultiplied) so callers like Color2TiColor
	// preserve intuitive color values across APIs.
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	return float32(nrgba.R) / 255., float32(nrgba.G) / 255., float32(nrgba.B) / 255., float32(nrgba.A) / 255.
}

// ExpandF16Color 将 Go 颜色转换为 0-1 float16 颜色
func ExpandF16Color(c color.Color) (r, g, b, a float32) {
	fr, fg, fb, fa := ExpandF32Color(c)
	return fr, fg, fb, fa
}

// ExpandU8Color 将 Go 颜色转换为 0xff 颜色
func ExpandU8Color(c color.Color) (r, g, b, a uint32) {
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	return uint32(nrgba.R), uint32(nrgba.G), uint32(nrgba.B), uint32(nrgba.A)
}

func Color2TiColor(c color.Color) TiColor {
	r, g, b, a := ExpandF32Color(c)
	return TiColor{r, g, b, a}
}

// ColorEqual 比较两个 color.Color 的 NRGBA 展开值是否一致。
// ColorEqual compares two color.Color values by expanded NRGBA components.
func ColorEqual(a, b color.Color) bool {
	if a == nil || b == nil {
		return a == b
	}
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}

// OrTransparentColor 如果 c 为 nil，则返回 color.Transparent
func OrTransparentColor(c color.Color) color.Color {
	if c == nil {
		return color.Transparent
	}
	return c
}
