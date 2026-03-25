package render

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ti"
)

// ISpriteOperator 精灵操作接口，主要是Set/Get，没有复杂操作
type ISpriteOperator interface {
	X() float32
	Y() float32
	Width() float32
	Height() float32
	Scale() float32
	Rotation() float32
	Alpha() float32
	// CenterX 实际的中心点x
	CenterX() float32
	// CenterY 实际的中心点y
	CenterY() float32
	Texture() *ti.TiImage

	SetX(x float32) ISprite
	SetY(y float32) ISprite
	MoveTo(x float32, y float32) ISprite

	SetScale(scale float32) ISprite
	// 缩放到指定尺寸
	SetScaleTo(width, height float32) ISprite
	SetRotation(rotation float32) ISprite
	SetAlpha(alpha float32) ISprite
}

// ISprite 精灵接口，包含操作接口，以及复杂的操作
type ISprite interface {
	ISpriteOperator

	SetMask(mask IMask) ISprite
	Mask() IMask
	// 替换纹理（之前的纹理会释放）
	SetTexture(texture *ti.TiImage) ISprite
	// 所有像素填充同一个颜色
	FillTexture(rgba color.Color)
	// 通过父级的宽高获取包围盒
	BoundingBox(parentWidth, parentHeight float32) ti.Rectangle[float32]
	// ResizeTo 重置尺寸，会替换纹理
	ResizeTo(width, height uint32) ISprite

	// 获取渲染器
	Renderer() *Renderer

	// Release 释放资源（必须调用，不然GPU显存泄漏）
	Release()
}

type IMask interface {
	FillWithTexture(texture *ti.TiImage)
	ApplyFeather(featherRadius uint32, featherMode ti.FeatherMode)
	Release()
	Texture() *ti.TiMask
}
