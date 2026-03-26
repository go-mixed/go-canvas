package render

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ti"
)

// IElement 精灵操作接口，主要是Set/Get，没有复杂操作
type IElement interface {
	X() float32
	Y() float32
	Width() float32
	Height() float32
	ScaleX() float32
	ScaleY() float32
	Rotation() float32
	Alpha() float32
	Cx() float32
	Cy() float32

	//// 替换纹理（之前的纹理会释放）
	//SetTexture(texture *ti.TiImage) ISprite
	// 所有像素填充同一个颜色
	FillTexture(rgba color.Color)
	Texture() *ti.TiImage

	SetX(x float32)
	SetY(y float32)
	MoveTo(x float32, y float32)
	SetCx(cx float32)
	SetCy(cy float32)

	SetScale(scaleX, scaleY float32)
	SetRotation(rotation float32)
	SetAlpha(alpha float32)

	// ResizeTo 重置尺寸，会替换纹理
	ResizeTo(width, height uint32) error
	// ClientRect 获取元素自身旋转+缩放后的边界（不与父级求交集）
	ClientRect() ti.Rectangle[float32]
	// ClippedRect 获取与父级区域裁剪后的可视区域（请在 Container.Render 中使用）
	ClippedRect(parentWidth, parentHeight float32) ti.Rectangle[float32]

	// Release 释放资源（必须调用，不然GPU显存泄漏）
	Release()

	IsDirty() bool
	SetDirty(val bool)
	//lockForUpdate(updateFn func(), triggerDirty func() bool)
}

type IRender interface {
	Render()
}

// ISprite 精灵接口，包含操作接口，以及复杂的操作
type ISprite interface {
	IElement

	SetMask(mask IMask)
	Mask() IMask

	IRender
}

// IContainer 容器接口
type IContainer interface {
	ISprite

	Add(sprite ISprite)
	Remove(sprite ISprite)
	Children() []ISprite

	ClientRect() ti.Rectangle[float32]
}

type IMask interface {
	FillWithTexture(texture *ti.TiImage)
	ApplyFeather(featherRadius uint32, featherMode ti.FeatherMode)
	Release()
	Texture() *ti.TiMask
}
