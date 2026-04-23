package ctypes

import "image/color"

// FillMode 缩放填充模式
type FillMode int32

const (
	FillModeStretch FillMode = 0 // 拉伸（不保持宽高比）
	FillModeFit     FillMode = 1 // 等比适应（可能有黑边）
	FillModeFill    FillMode = 2 // 等比填充（可能裁剪）
)

// ScaleMode 缩放算法模式
type ScaleMode int32

const (
	ScaleModeNearest ScaleMode = 0 // 最近邻
	ScaleModeLinear  ScaleMode = 1 // 双线性
	ScaleModeCubic   ScaleMode = 2 // 双三次
	ScaleModeLanczos ScaleMode = 3 // Lanczos4（质量最高）
)

type ResizeOptions struct {
	FillMode  FillMode  // fit, fill
	ScaleMode ScaleMode // nearest, linear, cubic
}

type Blur struct {
	Mode   BlurMode
	Radius int
}

func (b Blur) IsEmpty() bool {
	return b.Radius <= 0
}

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

type Border struct {
	// 边框宽度
	LeftWidth, RightWidth, TopWidth, BottomWidth int
	// 边框样式
	LeftStyle, RightStyle, TopStyle, BottomStyle BorderStyle
	// 边框颜色
	LeftColor, RightColor, TopColor, BottomColor color.Color
	// 圆角像素
	TopLeftRadius, BottomLeftRadius, TopRightRadius, BottomRightRadius int
}

func (b Border) IsEmpty() bool {
	return b.LeftWidth == 0 && b.RightWidth == 0 && b.TopWidth == 0 && b.BottomWidth == 0 &&
		b.TopLeftRadius == 0 && b.BottomLeftRadius == 0 && b.TopRightRadius == 0 && b.BottomRightRadius == 0
}
