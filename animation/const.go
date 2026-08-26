package animation

import (
	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/ti"
)

type IAttributeAnimation interface {
	BaseAttribute() *ctypes.Attribute
	TargetAttribute() *ctypes.Attribute

	ModifiedFields() []ModifiedField
	IsModified(field ModifiedField) bool

	Ticking(t float32)

	SetEasing(easing ti.EasingFunction) IAttributeAnimation
	SetRect(rect ctypes.Rectangle[int]) IAttributeAnimation
	SetXYWH(x, y, width, height int) IAttributeAnimation
	SetXY(x, y int) IAttributeAnimation
	SetWH(width, height int) IAttributeAnimation
	SetX(x int) IAttributeAnimation
	SetY(y int) IAttributeAnimation
	SetWidth(width int) IAttributeAnimation
	SetHeight(height int) IAttributeAnimation
	MoveTo(x, y int) IAttributeAnimation
	SetScale(x, y float32) IAttributeAnimation
	SetRotation(rotation float32) IAttributeAnimation
	SetAlpha(alpha float32) IAttributeAnimation
	SetCx(cx int) IAttributeAnimation
	SetCy(cy int) IAttributeAnimation
	SetCxy(x int, y int) IAttributeAnimation
}

type ModifiedField string

const (
	ModifiedFieldX        ModifiedField = "x"
	ModifiedFieldY        ModifiedField = "y"
	ModifiedFieldWidth    ModifiedField = "width"
	ModifiedFieldHeight   ModifiedField = "height"
	ModifiedFieldCx       ModifiedField = "cx"
	ModifiedFieldCy       ModifiedField = "cy"
	ModifiedFieldScaleX   ModifiedField = "scale_x"
	ModifiedFieldScaleY   ModifiedField = "scale_y"
	ModifiedFieldRotation ModifiedField = "rotation"
	ModifiedFieldAlpha    ModifiedField = "alpha"
)
