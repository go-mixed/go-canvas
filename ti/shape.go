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

type ShapeOptions struct {
	Direction Direction
	Color     color.Color
}

func WithShapeDirection(direction Direction) func(*ShapeOptions) {
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
