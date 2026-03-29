package ti

import (
	"image/color"
	"strings"

	"github.com/go-mixed/go-taichi/taichi"
)

type TiImage = taichi.NdArray   // ndarray[w, h, [r, g, b, a]]
type BgraImage = taichi.NdArray // ndarray[w, h, u32] u32 == [0xbbggrraa]
type TiMask = taichi.NdArray    // ndarray[w, h, a]
type TiGrid = taichi.NdArray    // ndarray[w, h, f32]
type TiColor = []float32        // [4]

// Direction 形状的扩展方向
type Direction int

const (
	DirectionTop         Direction = 0
	DirectionBottom      Direction = 1
	DirectionLeft        Direction = 2
	DirectionRight       Direction = 3
	DirectionTopLeft     Direction = 4
	DirectionTopRight    Direction = 5
	DirectionBottomLeft  Direction = 6
	DirectionBottomRight Direction = 7
	DirectionCenter      Direction = 8
)

const Sqrt2Inv = 0.7071067811865476 // 1 / √2

// DirectionFromString 将字符串转换为 Direction
//
// 支持的字符串: "top", "bottom", "left", "right", "top_left", "top_right", "bottom_left", "bottom_right", "center"（不区分大小写）
func DirectionFromString(s string) Direction {
	switch strings.ToLower(s) {
	case "top", "t":
		return DirectionTop
	case "bottom", "b":
		return DirectionBottom
	case "left", "l":
		return DirectionLeft
	case "right", "r":
		return DirectionRight
	case "top_left", "topleft", "tf":
		return DirectionTopLeft
	case "top_right", "topright", "tr":
		return DirectionTopRight
	case "bottom_left", "bottomleft", "bl":
		return DirectionBottomLeft
	case "bottom_right", "bottomright", "br":
		return DirectionBottomRight
	case "center", "c":
		return DirectionCenter
	default:
		return DirectionCenter
	}
}

type Align struct {
	HAlign HorizontalAlign
	VAlign VerticalAlign
}

// HorizontalAlign 水平对齐方式
type HorizontalAlign int

const (
	HAlignLeft HorizontalAlign = iota
	HAlignCenter
	HAlignRight
)

// VerticalAlign 垂直对齐方式
type VerticalAlign int

const (
	VAlignTop VerticalAlign = iota
	VAlignMiddle
	VAlignBottom
)

// BlurMode 模糊模式
type BlurMode int32

const (
	BlurModeBox      BlurMode = 0 // 普通模糊
	BlurModeGaussian BlurMode = 1 // 高斯模糊
	BlurModeMosaic   BlurMode = 2 // 马赛克
)

type ImageWriter interface {
	Set(x, y int, c color.Color)
}
