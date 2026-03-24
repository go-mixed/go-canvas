package ti

import "image/color"

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

// ShapeDirection 形状的扩展方向
type ShapeDirection int

const (
	ShapeDirectionTop      ShapeDirection = 0
	ShapeDirectionBottom   ShapeDirection = 1
	ShapeDirectionLeft     ShapeDirection = 2
	ShapeDirectionRight    ShapeDirection = 3
	ShapeDirectionTopLeft  ShapeDirection = 4
	ShapeDirectionTopRight ShapeDirection = 5
	ShapeDirectionBotLeft  ShapeDirection = 6
	ShapeDirectionBotRight ShapeDirection = 7
	ShapeDirectionCenter   ShapeDirection = 8
)

// shapeConfig 形状配置，用于 compute_directional kernel
type shapeConfig struct {
	useRadial       float32
	manhattanWeight float32
	chebyshevWeight float32
}

// 形状配置映射
var shapeConfigs = map[ShapeType]shapeConfig{
	ShapeTypeLinear:    {useRadial: 0.0, manhattanWeight: 0.0, chebyshevWeight: 0.0},
	ShapeTypeCircle:    {useRadial: 1.0, manhattanWeight: 0.0, chebyshevWeight: 0.0},
	ShapeTypeDiamond:   {useRadial: 1.0, manhattanWeight: 1.0, chebyshevWeight: 0.0},
	ShapeTypeRectangle: {useRadial: 1.0, manhattanWeight: 0.0, chebyshevWeight: 1.0},
}

const sqrt2Inv = 0.7071067811865476 // 1 / sqrt(2)

// directionVectors 方向向量映射
var directionVectors = map[ShapeDirection][2]float32{
	ShapeDirectionTop:      {0.0, -1.0},
	ShapeDirectionBottom:   {0.0, 1.0},
	ShapeDirectionLeft:     {-1.0, 0.0},
	ShapeDirectionRight:    {1.0, 0.0},
	ShapeDirectionTopLeft:  {-sqrt2Inv, -sqrt2Inv},
	ShapeDirectionTopRight: {sqrt2Inv, -sqrt2Inv},
	ShapeDirectionBotLeft:  {-sqrt2Inv, sqrt2Inv},
	ShapeDirectionBotRight: {sqrt2Inv, sqrt2Inv},
	ShapeDirectionCenter:   {0.0, 0.0},
}

type ShapeOptions struct {
	Direction ShapeDirection
	Color     color.Color
}

func WithShapeDirection(direction ShapeDirection) func(*ShapeOptions) {
	return func(opts *ShapeOptions) {
		opts.Direction = direction
	}
}

func WithShapeColor(color color.Color) func(*ShapeOptions) {
	return func(opts *ShapeOptions) {
		opts.Color = color
	}
}

// FeatherMode 羽化模式
type FeatherMode int

const (
	FeatherModeLinear     FeatherMode = 0
	FeatherModeConic      FeatherMode = 1
	FeatherModeSmoothstep FeatherMode = 2
	FeatherModeSigmoid    FeatherMode = 3
)
