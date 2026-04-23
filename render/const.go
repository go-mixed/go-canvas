package render

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/internel/misc"
	"github.com/go-mixed/go-canvas/ti"
	"github.com/go-mixed/go-taichi/taichi"
)

type IAttribute interface {
	X() int
	Y() int
	Cx() int
	Cy() int
	Width() int
	Height() int
	ScaleX() float32
	ScaleY() float32
	Rotation() float32
	Alpha() float32

	SetX(x int)
	SetY(y int)
	MoveTo(x int, y int)
	SetCx(cx int)
	SetCy(cy int)

	SetScale(scaleX, scaleY float32)
	SetRotation(rotation float32)
	SetAlpha(alpha float32)
	SetBorderRadius(topLeftRadius, topRightRadius, bottomRightRadius, bottomLeftRadius int)
	SetBorderWidth(top, right, bottom, left int)
	SetAllBorderWidths(width int)
	SetBorderStyle(top, right, bottom, left ctypes.BorderStyle)
	SetAllBorderStyles(style ctypes.BorderStyle)
	SetBorderColor(top, right, bottom, left color.Color)
	SetAllBorderColors(c color.Color)

	// Resize 设置尺寸
	Resize(width, height int) error
}

type IElementOperation interface {
	// Blur 模糊纹理（马赛克/高斯/普通）
	Blur(mode ctypes.BlurMode, radius int32) error
	// Fill 所有像素填充同一个颜色
	Fill(rgba color.Color)
}

type IDirty interface {
	IsDirty() bool
	SetDirty(val ctypes.DirtyMode)
}

type ITexture interface {
	Texture() *taichi.NdArray
}

type IRelease interface {
	// Release 释放资源（必须调用，不然GPU显存泄漏）
	Release()
}

// IElement 精灵操作接口，主要是Set/Get，没有复杂操作
type IElement interface {
	IAttribute
	IElementOperation

	Attribute() *ctypes.Attribute

	// ClientRect 获取元素自身旋转+缩放后的边界
	ClientRect() ctypes.Rectangle[int]

	// LockForUpdate(updateFn func(), triggerDirty func() ctypes.DirtyMode)
}

type IRender interface {
	Render(frameIndex int)
}

type IAnimation interface {
	// HasAnimationAt 在给定绝对帧号下判断是否有动画需要执行。
	HasAnimationAt(frameIndex int) bool
	// Animate 追加一段动画任务。
	// targetFn 在动画实际开始帧被调用，基于当时属性生成目标属性。
	Animate(targetFn ti.TargetAttributeFn, startFrameIndex, durationFrames int) ISprite
	// ClearAnimations 清空当前精灵的动画队列。
	ClearAnimations() ISprite
	// StopAnimation 停止后续动画；reset=true 时回到当前段起点状态。
	StopAnimation(reset bool) ISprite
	// TickAnimation 以绝对帧号推进动画，返回是否仍有动画待执行。
	TickAnimation(frameIndex int) bool
}

type IShape interface {
	DrawShape(shapeType ctypes.ShapeType, tVal float32, options *ctypes.ShapeOptions)
}

type IGarbage interface {
	AddGarbageTexture(texture *taichi.NdArray)
	ReleaseGarbageTextures()
}

// ISprite 精灵接口，包含操作接口，以及复杂的操作
type ISprite interface {
	IElement
	IAnimation
	IRender
	IDirty
	ITexture
	IRelease
	IGarbage

	IMaskParent

	RemoveFromParent()
}

// IContainer 容器接口
type IContainer interface {
	IElement
	IAnimation
	IRender
	IDirty
	ITexture
	IRelease
	IGarbage

	IParent
	ScrollTop(y int)
	ScrollLeft(x int)
}

type IMask interface {
	FillWithTexture(texture *ctypes.TiImage)
	SetFeather(featherRadius uint32, featherMode ctypes.FeatherMode)
	RemoveFromParent()

	IDirty
	IRender
	ITexture
	IRelease
}

type IParent interface {
	AddChild(sprite ISprite)
	RemoveChild(sprite ISprite)
	Children() *misc.List[ISprite]

	Renderer() *Renderer
}

type IMaskParent interface {
	AddMask(mask IMask)
	RemoveMask(mask IMask)
	Masks() *misc.List[IMask]

	Renderer() *Renderer
}

type selfRelease struct {
	renderer *Renderer

	children *misc.List[ISprite]
	masks    *misc.List[IMask]
}

var _ IParent = (*selfRelease)(nil)
var _ IMaskParent = (*selfRelease)(nil)

func SelfRelease(renderer *Renderer) *selfRelease {
	return &selfRelease{
		renderer: renderer,
		children: misc.NewList[ISprite](),
		masks:    misc.NewList[IMask](),
	}
}

func (s *selfRelease) Renderer() *Renderer {
	return s.renderer
}

func (s *selfRelease) AddChild(sprite ISprite) {}

func (s *selfRelease) RemoveChild(sprite ISprite) {}

func (s *selfRelease) Children() *misc.List[ISprite] {
	return s.children
}

func (s *selfRelease) AddMask(mask IMask) {}

func (s *selfRelease) RemoveMask(mask IMask) {}
func (s *selfRelease) Masks() *misc.List[IMask] {
	return s.masks
}
