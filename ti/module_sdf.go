package ti

import "image/color"

// ComputeNormalizedCoords 计算归一化坐标网格
func (m *AotModule) ComputeNormalizedCoords(dx, dy *TiGrid, cx, cy float32) {
	kernel := m.getCache("compute_normalized_coords")
	kernel.Launch().
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(cx).
		ArgFloat32(cy).
		Run()
}

// ComputeCircle 计算圆形遮罩
func (m *AotModule) ComputeCircle(data *TiImage, dx, dy *TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_circle")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(Color2TiColor(color)...).
		Run()
}

// ComputeDiamond 计算菱形遮罩（曼哈顿距离）
func (m *AotModule) ComputeDiamond(data *TiImage, dx, dy *TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_diamond")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(Color2TiColor(color)...).
		Run()
}

// ComputeRect 计算矩形遮罩
// dirVal: 0=TOP, 1=BOTTOM, 2=LEFT, 3=RIGHT, 4=TOP_LEFT, 5=TOP_RIGHT, 6=BOTTOM_LEFT, 7=BOTTOM_RIGHT, 8=CENTER
func (m *AotModule) ComputeRect(data *TiImage, dx, dy *TiGrid, tVal float32, dirVal int32, color color.Color) {
	kernel := m.getCache("compute_rect")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgInt32(dirVal).
		ArgVectorFloat32(Color2TiColor(color)...).
		Run()
}

// ComputeDirectional 通用方向性遮罩计算
// useRadial: 0.0=线性投影, 1.0=径向距离
// manhattanWeight: 0.0=欧几里得, 1.0=曼哈顿
// chebyshevWeight: 0.0=不使用, 1.0=切比雪夫
func (m *AotModule) ComputeDirectional(data *TiImage, dx, dy *TiGrid, tVal, dirX, dirY, useRadial, manhattanWeight, chebyshevWeight float32, color color.Color) {
	kernel := m.getCache("compute_directional")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgFloat32(dirX).
		ArgFloat32(dirY).
		ArgFloat32(useRadial).
		ArgFloat32(manhattanWeight).
		ArgFloat32(chebyshevWeight).
		ArgVectorFloat32(Color2TiColor(color)...).
		Run()
}

// ComputeTriangle 计算三角形遮罩
func (m *AotModule) ComputeTriangle(data *TiImage, dx, dy *TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_triangle")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(Color2TiColor(color)...).
		Run()
}

// ComputeStar 计算五角星遮罩
func (m *AotModule) ComputeStar(data *TiImage, dx, dy *TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_star")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(Color2TiColor(color)...).
		Run()
}

// ComputeHeart 计算心形遮罩
func (m *AotModule) ComputeHeart(data *TiImage, dx, dy *TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_heart")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(Color2TiColor(color)...).
		Run()
}

// ComputeCross 计算十字形遮罩
func (m *AotModule) ComputeCross(data *TiImage, dx, dy *TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_cross")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(Color2TiColor(color)...).
		Run()
}

// ComputeShape 统一的形状计算方法
// 根据 shapeType 自动选择合适的 kernel 和参数
// dir 参数仅对线性方向性形状有效
func (m *AotModule) ComputeShape(data *TiImage, dx, dy *TiGrid, shapeType ShapeType, tVal float32, dir ShapeDirection, color color.Color) {
	switch shapeType {
	case ShapeTypeTriangle:
		m.ComputeTriangle(data, dx, dy, tVal, color)
	case ShapeTypeStar5:
		m.ComputeStar(data, dx, dy, tVal, color)
	case ShapeTypeHeart:
		m.ComputeHeart(data, dx, dy, tVal, color)
	case ShapeTypeCross:
		m.ComputeCross(data, dx, dy, tVal, color)
	default:
		// 使用 compute_directional kernel
		cfg, ok := shapeConfigs[shapeType]
		if !ok {
			return
		}
		vec, ok := directionVectors[dir]
		if !ok {
			vec = directionVectors[ShapeDirectionCenter]
		}
		m.ComputeDirectional(data, dx, dy, tVal, vec[0], vec[1], cfg.useRadial, cfg.manhattanWeight, cfg.chebyshevWeight, color)
	}
}
