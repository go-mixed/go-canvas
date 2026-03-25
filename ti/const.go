package ti

import (
	"strings"

	"github.com/go-mixed/go-taichi/taichi"
)

type CvImage = taichi.NdArray // ndarray[h, w, [b, g, r]]
type TiImage = taichi.NdArray // ndarray[w, h, [r, g, b, a]]
type TiMask = taichi.NdArray  // ndarray[w, h, a]
type TiGrid = taichi.NdArray  // ndarray[w, h, f32]
type TiColor = []float32      // [4]

// Direction 形状的扩展方向
type Direction int

const (
	DirectionTop      Direction = 0
	DirectionBottom   Direction = 1
	DirectionLeft     Direction = 2
	DirectionRight    Direction = 3
	DirectionTopLeft  Direction = 4
	DirectionTopRight Direction = 5
	DirectionBotLeft  Direction = 6
	DirectionBotRight Direction = 7
	DirectionCenter   Direction = 8
)

const Sqrt2Inv = 0.7071067811865476 // 1 / √2

// DirectionVectors 方向向量映射（已归一化对角线方向）
var DirectionVectors = map[Direction][2]float32{
	DirectionTop:      {0.0, -1.0},
	DirectionBottom:   {0.0, 1.0},
	DirectionLeft:     {-1.0, 0.0},
	DirectionRight:    {1.0, 0.0},
	DirectionTopLeft:  {-Sqrt2Inv, -Sqrt2Inv},
	DirectionTopRight: {Sqrt2Inv, -Sqrt2Inv},
	DirectionBotLeft:  {-Sqrt2Inv, Sqrt2Inv},
	DirectionBotRight: {Sqrt2Inv, Sqrt2Inv},
	DirectionCenter:   {0.0, 0.0},
}

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
		return DirectionBotLeft
	case "bottom_right", "bottomright", "br":
		return DirectionBotRight
	case "center", "c":
		return DirectionCenter
	default:
		return DirectionCenter
	}
}
