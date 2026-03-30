package ti

import (
	"image/color"
)

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

type Attribute struct {
	rect Rectangle[int]
	// 中心点的相对值
	cx, cy         int
	scaleX, scaleY float32 // 1.0 for no scaling
	rotation       float32 // 0.0 for no rotation, in radians
	alpha          float32 // 0.0 for no alpha, 1.0 for full alpha
	resizeOptions  ResizeOptions

	Border BorderAttribute
}

type BorderAttribute struct {
	LeftWidth, RightWidth, TopWidth, BottomWidth float32
	LeftStyle, RightStyle, TopStyle, BottomStyle BorderStyle
	LeftColor, RightColor, TopColor, BottomColor color.Color
}

type BorderStyle uint

const (
	BorderStyleNone BorderStyle = iota
	BorderStyleSolid
	BorderStyleDotted
	BorderStyleDashed
	BorderStyleDouble
	BorderStyleGroove
	BorderStyleRidge
	BorderStyleInset
	BorderStyleOutset
)

func Attr() *Attribute {
	return &Attribute{
		rect:     Rectangle[int]{},
		scaleX:   1.0,
		scaleY:   1.0,
		rotation: 0.0,
		alpha:    1.0,
		resizeOptions: ResizeOptions{
			FillMode:  FillModeFit,
			ScaleMode: ScaleModeNearest,
		},
	}
}

func (a *Attribute) Rect() Rectangle[int] {
	return a.rect
}

func (a *Attribute) SetRect(rect Rectangle[int]) *Attribute {
	a.rect = rect
	return a
}

func (a *Attribute) SetXYWH(x, y, width, height int) *Attribute {
	a.rect = RectWH(x, y, width, height)
	return a
}

func (a *Attribute) SetXY(x, y int) *Attribute {
	a.rect = RectWH(x, y, a.Width(), a.Height())
	return a
}

func (a *Attribute) SetWH(width, height int) *Attribute {
	a.rect = RectWH(a.X(), a.Y(), width, height)
	return a
}

func (a *Attribute) SetX(x int) *Attribute {
	a.rect = a.rect.MoveTo(x, a.Y())
	return a
}

func (a *Attribute) X() int {
	return a.rect.Min.X
}

func (a *Attribute) SetY(y int) *Attribute {
	a.rect = a.rect.MoveTo(a.X(), y)
	return a
}

func (a *Attribute) Y() int {
	return a.rect.Min.Y
}

func (a *Attribute) SetWidth(width int) *Attribute {
	a.rect.Max.X = a.rect.Min.X + width
	return a
}

func (a *Attribute) Width() int {
	return a.rect.Max.X - a.rect.Min.X
}

func (a *Attribute) SetHeight(height int) *Attribute {
	a.rect.Max.Y = a.rect.Min.Y + height
	return a
}

func (a *Attribute) Height() int {
	return a.rect.Max.Y - a.rect.Min.Y
}

func (a *Attribute) MoveTo(x, y int) *Attribute {
	a.rect = a.rect.MoveTo(x, y)
	return a
}

func (a *Attribute) SetScale(x, y float32) *Attribute {
	a.scaleX = x
	a.scaleY = y
	return a
}

func (a *Attribute) ScaleX() float32 {
	return a.scaleX
}

func (a *Attribute) ScaleY() float32 {
	return a.scaleY
}

func (a *Attribute) SetRotation(rotation float32) *Attribute {
	a.rotation = rotation
	return a
}

func (a *Attribute) Rotation() float32 {
	return a.rotation
}

// SetAlpha 设置透明度。 [0, 1]， 0为完全透明, 1为完全不透明
func (a *Attribute) SetAlpha(alpha float32) *Attribute {
	a.alpha = alpha
	return a
}

// Alpha 透明度。 [0, 1], 0为完全透明, 1为完全不透明
func (a *Attribute) Alpha() float32 {
	return a.alpha
}

// SetCx 设置中心点x的相对值（0,0)为左上角
func (a *Attribute) SetCx(cx int) *Attribute {
	a.cx = cx
	return a
}

// Cx 中心点x的相对值
func (a *Attribute) Cx() int {
	return a.cx
}

// SetCy 设置中心点y的相对值（0,0)为左上角
func (a *Attribute) SetCy(cy int) *Attribute {
	a.cy = cy
	return a
}

// Cy 中心点y的相对值
func (a *Attribute) Cy() int {
	return a.cy
}

func (a *Attribute) SetResizeOptions(fillMode FillMode, scaleMode ScaleMode) *Attribute {
	a.resizeOptions.FillMode = fillMode
	a.resizeOptions.ScaleMode = scaleMode
	return a
}

func (a *Attribute) ResizeOptions() ResizeOptions {
	return a.resizeOptions
}
