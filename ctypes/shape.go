package ctypes

import (
	"image/color"
	"slices"
)

type ShapeType string

const (
	// 特殊形状（使用专用 kernel）
	ShapeTypeTriangle ShapeType = "triangle"
	ShapeTypeStar5    ShapeType = "star5"
	ShapeTypeHeart    ShapeType = "heart"
	ShapeTypeCross    ShapeType = "cross"

	// 方向性形状（使用 compute_directional kernel）
	// 线性方向性（8方向）
	ShapeTypeLinear ShapeType = "linear"
	// 圆形（欧几里得距离）
	ShapeTypeCircle ShapeType = "circle"
	// 菱形（曼哈顿距离）
	ShapeTypeDiamond ShapeType = "diamond"
	// 矩形（切比雪夫距离）
	ShapeTypeRectangle ShapeType = "rectangle"
)

// ShapeTypeFromString 从字符串获取形状类型
func ShapeTypeFromString(s string) ShapeType {
	if slices.Contains([]ShapeType{
		ShapeTypeTriangle, ShapeTypeStar5, ShapeTypeHeart, ShapeTypeCross,
		ShapeTypeLinear, ShapeTypeCircle, ShapeTypeDiamond, ShapeTypeRectangle,
	}, ShapeType(s)) {
		return ShapeType(s)
	}

	return ShapeTypeLinear // Default to linear if no match found
}

type ShapeOptions struct {
	Direction Direction
	Color     color.Color
}

func ShapeOpt() *ShapeOptions {
	return &ShapeOptions{
		Direction: DirectionCenter,
		Color:     color.Black,
	}
}

func (o *ShapeOptions) SetDirection(direction Direction) *ShapeOptions {
	o.Direction = direction
	return o
}

func (o *ShapeOptions) SetColor(c color.Color) *ShapeOptions {
	o.Color = c
	return o
}

// FeatherMode 羽化模式
type FeatherMode int

const (
	FeatherModeLinear     FeatherMode = 0
	FeatherModeConic      FeatherMode = 1
	FeatherModeSmoothstep FeatherMode = 2
	FeatherModeSigmoid    FeatherMode = 3
)

type ShapeMaskOptions struct {
	*ShapeOptions

	ShapeType ShapeType

	StartT float32
	EndT   float32

	FeatherRadius uint32
	FeatherMode   FeatherMode
}

func ShapeMaskOpt() *ShapeMaskOptions {
	return &ShapeMaskOptions{
		ShapeOptions: ShapeOpt(),
		ShapeType:    ShapeTypeRectangle,
		StartT:       0.0,
		EndT:         2.0,
	}
}

func (o *ShapeMaskOptions) SetShapeType(shapeType ShapeType) *ShapeMaskOptions {
	o.ShapeType = shapeType
	return o
}

func (o *ShapeMaskOptions) SetDirection(direction Direction) *ShapeMaskOptions {
	o.Direction = direction
	return o
}

func (o *ShapeMaskOptions) SetColor(c color.Color) *ShapeMaskOptions {
	o.Color = c
	return o
}

func (o *ShapeMaskOptions) SetTRange(startT, endT float32) *ShapeMaskOptions {
	o.StartT = startT
	o.EndT = endT
	return o
}

func (o *ShapeMaskOptions) SetFeather(radius uint32, mode FeatherMode) *ShapeMaskOptions {
	o.FeatherRadius = radius
	o.FeatherMode = mode
	return o
}
