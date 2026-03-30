package render

import (
	"image/color"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
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

	// Resize 设置尺寸
	Resize(width, height int) error
}

type IElementOperation interface {

	// Blur 模糊纹理（马赛克/高斯/普通）
	Blur(mode ti.BlurMode, radius int32) error
	// Fill 所有像素填充同一个颜色
	Fill(rgba color.Color)
}

// IElement 精灵操作接口，主要是Set/Get，没有复杂操作
type IElement interface {
	IAttribute

	Attribute() *ti.Attribute

	Texture() *ti.TiImage

	// ClientRect 获取元素自身旋转+缩放后的边界
	ClientRect() ti.Rectangle[int]

	IsDirty() bool
	SetDirty(val bool)
	//LockForUpdate(updateFn func(), triggerDirty func() bool)

	// Release 释放资源（必须调用，不然GPU显存泄漏）
	Release()

	Renderer() *Renderer

	IElementOperation
}

// ISprite 精灵接口，包含操作接口，以及复杂的操作
type ISprite interface {
	IElement

	IMaskParent

	RemoveFromParent()
	Render()
}

// IContainer 容器接口
type IContainer interface {
	ISprite

	IParent

	ScrollTop(y int)
	ScrollLeft(x int)
}

type IMask interface {
	FillWithTexture(texture *ti.TiImage)
	ApplyFeather(featherRadius uint32, featherMode ti.FeatherMode)
	Release()
	Texture() *ti.TiMask
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
