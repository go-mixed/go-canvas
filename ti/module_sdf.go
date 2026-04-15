package ti

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ctypes"
)

// AsyncComputeNormalizedCoords 计算归一化坐标网格
func (m *AotModule) AsyncComputeNormalizedCoords(dx, dy *ctypes.TiGrid, cx, cy float32) {
	kernel := m.getCache("compute_normalized_coords")
	kernel.Launch().
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(cx).
		ArgFloat32(cy).
		RunAsync()
}

func (m *AotModule) ComputeNormalizedCoords(dx, dy *ctypes.TiGrid, cx, cy float32) {
	m.AsyncComputeNormalizedCoords(dx, dy, cx, cy)
	m.runtime.Wait()
}

// AsyncComputeCircle 计算圆形遮罩
func (m *AotModule) AsyncComputeCircle(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_circle")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(ctypes.Color2TiColor(color)...).
		RunAsync()
}

func (m *AotModule) ComputeCircle(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	m.AsyncComputeCircle(data, dx, dy, tVal, color)
	m.runtime.Wait()
}

// AsyncComputeDiamond 计算菱形遮罩（曼哈顿距离）
func (m *AotModule) AsyncComputeDiamond(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_diamond")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(ctypes.Color2TiColor(color)...).
		RunAsync()
}

func (m *AotModule) ComputeDiamond(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	m.AsyncComputeDiamond(data, dx, dy, tVal, color)
	m.runtime.Wait()
}

// AsyncComputeRect 计算矩形遮罩
// dirVal: 0=TOP, 1=BOTTOM, 2=LEFT, 3=RIGHT, 4=TOP_LEFT, 5=TOP_RIGHT, 6=BOTTOM_LEFT, 7=BOTTOM_RIGHT, 8=CENTER
func (m *AotModule) AsyncComputeRect(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, dirVal int32, color color.Color) {
	kernel := m.getCache("compute_rect")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgInt32(dirVal).
		ArgVectorFloat32(ctypes.Color2TiColor(color)...).
		RunAsync()
}

func (m *AotModule) ComputeRect(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, dirVal int32, color color.Color) {
	m.AsyncComputeRect(data, dx, dy, tVal, dirVal, color)
	m.runtime.Wait()
}

// AsyncComputeDirectional 通用方向性遮罩计算
// useRadial: 0.0=线性投影, 1.0=径向距离
// manhattanWeight: 0.0=欧几里得, 1.0=曼哈顿
// chebyshevWeight: 0.0=不使用, 1.0=切比雪夫
func (m *AotModule) AsyncComputeDirectional(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal, dirX, dirY, useRadial, manhattanWeight, chebyshevWeight float32, color color.Color) {
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
		ArgVectorFloat32(ctypes.Color2TiColor(color)...).
		RunAsync()
}

func (m *AotModule) ComputeDirectional(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal, dirX, dirY, useRadial, manhattanWeight, chebyshevWeight float32, color color.Color) {
	m.AsyncComputeDirectional(data, dx, dy, tVal, dirX, dirY, useRadial, manhattanWeight, chebyshevWeight, color)
	m.runtime.Wait()
}

// AsyncComputeTriangle 计算三角形遮罩
func (m *AotModule) AsyncComputeTriangle(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_triangle")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(ctypes.Color2TiColor(color)...).
		RunAsync()
}

func (m *AotModule) ComputeTriangle(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	m.AsyncComputeTriangle(data, dx, dy, tVal, color)
	m.runtime.Wait()
}

// AsyncComputeStar 计算五角星遮罩
func (m *AotModule) AsyncComputeStar(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_star")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(ctypes.Color2TiColor(color)...).
		RunAsync()
}

func (m *AotModule) ComputeStar(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	m.AsyncComputeStar(data, dx, dy, tVal, color)
	m.runtime.Wait()
}

// AsyncComputeHeart 计算心形遮罩
func (m *AotModule) AsyncComputeHeart(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_heart")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(ctypes.Color2TiColor(color)...).
		RunAsync()
}

func (m *AotModule) ComputeHeart(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	m.AsyncComputeHeart(data, dx, dy, tVal, color)
	m.runtime.Wait()
}

// AsyncComputeCross 计算十字形遮罩
func (m *AotModule) AsyncComputeCross(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	kernel := m.getCache("compute_cross")
	kernel.Launch().
		ArgNdArray(data).
		ArgNdArray(dx).
		ArgNdArray(dy).
		ArgFloat32(tVal).
		ArgVectorFloat32(ctypes.Color2TiColor(color)...).
		RunAsync()
}

func (m *AotModule) ComputeCross(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, tVal float32, color color.Color) {
	m.AsyncComputeCross(data, dx, dy, tVal, color)
	m.runtime.Wait()
}

// sdfDirectionVectors SDF形状扩展方向向量映射（已归一化对角线方向）
var sdfDirectionVectors = map[ctypes.Direction][2]float32{
	ctypes.DirectionTop:         {0.0, -1.0},
	ctypes.DirectionBottom:      {0.0, 1.0},
	ctypes.DirectionLeft:        {-1.0, 0.0},
	ctypes.DirectionRight:       {1.0, 0.0},
	ctypes.DirectionTopLeft:     {-ctypes.Sqrt2Inv, -ctypes.Sqrt2Inv},
	ctypes.DirectionTopRight:    {ctypes.Sqrt2Inv, -ctypes.Sqrt2Inv},
	ctypes.DirectionBottomLeft:  {-ctypes.Sqrt2Inv, ctypes.Sqrt2Inv},
	ctypes.DirectionBottomRight: {ctypes.Sqrt2Inv, ctypes.Sqrt2Inv},
	ctypes.DirectionCenter:      {0.0, 0.0},
}

// 形状配置映射
var shapeConfigs = map[ctypes.ShapeType]shapeConfig{
	ctypes.ShapeTypeLinear:    {useRadial: 0.0, manhattanWeight: 0.0, chebyshevWeight: 0.0},
	ctypes.ShapeTypeCircle:    {useRadial: 1.0, manhattanWeight: 0.0, chebyshevWeight: 0.0},
	ctypes.ShapeTypeDiamond:   {useRadial: 1.0, manhattanWeight: 1.0, chebyshevWeight: 0.0},
	ctypes.ShapeTypeRectangle: {useRadial: 1.0, manhattanWeight: 0.0, chebyshevWeight: 1.0},
}

// shapeConfig 形状配置，用于 compute_directional kernel
type shapeConfig struct {
	useRadial       float32
	manhattanWeight float32
	chebyshevWeight float32
}

// AsyncComputeShape 统一的形状计算方法
// 根据 shapeType 自动选择合适的 kernel 和参数
// dir 参数仅对线性方向性形状有效
func (m *AotModule) AsyncComputeShape(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, shapeType ctypes.ShapeType, tVal float32, dir ctypes.Direction, color color.Color) {
	switch shapeType {
	case ctypes.ShapeTypeTriangle:
		m.AsyncComputeTriangle(data, dx, dy, tVal, color)
	case ctypes.ShapeTypeStar5:
		m.AsyncComputeStar(data, dx, dy, tVal, color)
	case ctypes.ShapeTypeHeart:
		m.AsyncComputeHeart(data, dx, dy, tVal, color)
	case ctypes.ShapeTypeCross:
		m.AsyncComputeCross(data, dx, dy, tVal, color)
	default:
		// 使用 compute_directional kernel
		cfg, ok := shapeConfigs[shapeType]
		if !ok {
			return
		}
		vec, ok := sdfDirectionVectors[dir]
		if !ok {
			vec = sdfDirectionVectors[ctypes.DirectionCenter]
		}
		m.AsyncComputeDirectional(data, dx, dy, tVal, vec[0], vec[1], cfg.useRadial, cfg.manhattanWeight, cfg.chebyshevWeight, color)
	}
}

func (m *AotModule) ComputeShape(data *ctypes.TiImage, dx, dy *ctypes.TiGrid, shapeType ctypes.ShapeType, tVal float32, dir ctypes.Direction, color color.Color) {
	m.AsyncComputeShape(data, dx, dy, shapeType, tVal, dir, color)
	m.runtime.Wait()
}
