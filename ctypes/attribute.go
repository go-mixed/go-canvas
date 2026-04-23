package ctypes

import (
	"image/color"

	"github.com/go-mixed/go-canvas/internel/misc"
)

type Attribute struct {
	rect Rectangle[int]
	// 中心点的相对值
	cx, cy         int
	scaleX, scaleY float32 // 1.0 for no scaling
	rotation       float32 // 0.0 for no rotation, in radians
	alpha          float32 // 0.0 for no alpha, 1.0 for full alpha
	resizeOptions  ResizeOptions

	border  Border
	blur    Blur
	padding Padding
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
		cx:       misc.NaNInt,
		cy:       misc.NaNInt,
		resizeOptions: ResizeOptions{
			FillMode:  FillModeFill,
			ScaleMode: ScaleModeNearest,
		},
		border: Border{
			LeftStyle:         BorderStyleSolid,
			RightStyle:        BorderStyleSolid,
			TopStyle:          BorderStyleSolid,
			BottomStyle:       BorderStyleSolid,
			LeftColor:         color.Black,
			RightColor:        color.Black,
			TopColor:          color.Black,
			BottomColor:       color.Black,
			TopLeftRadius:     0,
			BottomLeftRadius:  0,
			TopRightRadius:    0,
			BottomRightRadius: 0,
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

func (a *Attribute) ClientWidth() int {
	return a.Width() + a.padding.Left + a.padding.Right + a.border.LeftWidth + a.border.RightWidth
}

func (a *Attribute) ClientHeight() int {
	return a.Height() + a.padding.Top + a.padding.Bottom + a.border.TopWidth + a.border.BottomWidth
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

// SetCxIfNotDefined 仅当未cx设置时，设置中心点x的相对值（0,0)为左上角
func (a *Attribute) SetCxIfNotDefined(cx int) *Attribute {
	if misc.IsNaNInt(a.cx) {
		a.cx = cx
	}
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

// SetCyIfNotDefined 仅当未cy设置时，设置中心点y的相对值（0,0)为左上角
func (a *Attribute) SetCyIfNotDefined(cy int) *Attribute {
	if misc.IsNaNInt(a.cy) {
		a.cy = cy
	}
	return a
}

// Cy 中心点y的相对值
func (a *Attribute) Cy() int {
	return a.cy
}

func (a *Attribute) SetCxy(x int, y int) *Attribute {
	a.SetCx(x)
	a.SetCy(y)
	return a
}

func (a *Attribute) SetResizeOptions(fillMode FillMode, scaleMode ScaleMode) *Attribute {
	a.resizeOptions.FillMode = fillMode
	a.resizeOptions.ScaleMode = scaleMode
	return a
}

func (a *Attribute) ResizeOptions() ResizeOptions {
	return a.resizeOptions
}

func (a *Attribute) SetPaddings(top, right, bottom, left int) *Attribute {
	a.padding = Padding{top, right, bottom, left}
	return a
}

func (a *Attribute) SetLeftPadding(left int) *Attribute {
	a.padding.Left = left
	return a
}

func (a *Attribute) SetRightPadding(right int) *Attribute {
	a.padding.Right = right
	return a
}

func (a *Attribute) SetTopPadding(top int) *Attribute {
	a.padding.Top = top
	return a
}

func (a *Attribute) SetBottomPadding(bottom int) *Attribute {
	a.padding.Bottom = bottom
	return a
}

func (a *Attribute) Padding() Padding {
	return a.padding
}

func (a *Attribute) SetBorder(border Border) *Attribute {
	a.border = border
	return a
}

func (a *Attribute) Border() Border {
	return a.border
}

// SetBorderRadius 设置边框圆角半径（像素），顺序：top-left, top-right, bottom-right, bottom-left。
func (a *Attribute) SetBorderRadius(topLeftRadius, topRightRadius, bottomRightRadius, bottomLeftRadius int) *Attribute {
	a.border.TopLeftRadius = topLeftRadius
	a.border.TopRightRadius = topRightRadius
	a.border.BottomRightRadius = bottomRightRadius
	a.border.BottomLeftRadius = bottomLeftRadius
	return a
}

func (a *Attribute) SetAllBorderRadius(radius int) *Attribute {
	return a.SetBorderRadius(radius, radius, radius, radius)
}

// SetBorderWidth 设置四边边框宽度（像素），顺序：top, right, bottom, left。
func (a *Attribute) SetBorderWidth(top, right, bottom, left int) *Attribute {
	a.border.TopWidth = top
	a.border.RightWidth = right
	a.border.BottomWidth = bottom
	a.border.LeftWidth = left
	return a
}

// SetAllBorderWidths 设置四边统一边框宽度（像素）。
func (a *Attribute) SetAllBorderWidths(width int) *Attribute {
	return a.SetBorderWidth(width, width, width, width)
}

// SetBorderStyle 设置四边边框样式，顺序：top, right, bottom, left。
func (a *Attribute) SetBorderStyle(top, right, bottom, left BorderStyle) *Attribute {
	a.border.TopStyle = top
	a.border.RightStyle = right
	a.border.BottomStyle = bottom
	a.border.LeftStyle = left
	return a
}

// SetAllBorderStyles 设置四边统一边框样式。
func (a *Attribute) SetAllBorderStyles(style BorderStyle) *Attribute {
	return a.SetBorderStyle(style, style, style, style)
}

// SetBorderColor 设置四边边框颜色，顺序：top, right, bottom, left。
func (a *Attribute) SetBorderColor(top, right, bottom, left color.Color) *Attribute {
	a.border.TopColor = top
	a.border.RightColor = right
	a.border.BottomColor = bottom
	a.border.LeftColor = left
	return a
}

// SetAllBorderColors 设置四边统一边框颜色。
func (a *Attribute) SetAllBorderColors(c color.Color) *Attribute {
	return a.SetBorderColor(c, c, c, c)
}

// Blur 获取模糊参数。
func (a *Attribute) Blur() Blur {
	return a.blur
}

// SetBlur 设置模糊参数。
func (a *Attribute) SetBlur(mode BlurMode, radius int) *Attribute {
	a.blur.Mode = mode
	a.blur.Radius = radius
	return a
}

func (a *Attribute) Copy() *Attribute {
	if a == nil {
		return nil
	}
	cp := *a
	return &cp
}

func (a *Attribute) Dirty() DirtyMode {
	var dirtyMode DirtyMode
	if !misc.NumberEqual(a.Alpha(), 1, misc.Epsilon) ||
		!misc.NumberEqual(a.ScaleX(), 1, misc.Epsilon) || !misc.NumberEqual(a.ScaleY(), 1, misc.Epsilon) ||
		!misc.NumberEqual(a.rotation, 0, misc.Epsilon) {
		dirtyMode |= DirtyModeComposite
	}

	if !a.border.IsEmpty() {
		dirtyMode |= DirtyModeLayout
	}

	if !a.padding.IsEmpty() {
		dirtyMode |= DirtyModeLayout
	}

	if a.blur.Radius > 0 {
		dirtyMode |= DirtyModePaint
	}

	return dirtyMode
}
