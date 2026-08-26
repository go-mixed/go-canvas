package animation

import (
	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/internel/misc"
	"github.com/go-mixed/go-canvas/render"
	"github.com/go-mixed/go-canvas/ti"
)

type BaseAttributeAnimation struct {
	instance IAttributeAnimation

	sprite render.IElement

	baseAttribute   ctypes.Attribute
	targetAttribute ctypes.Attribute
	easing          ti.EasingFunction

	modifiedFields map[ModifiedField]struct{}
	order          []ModifiedField
}

func BuildAttributeAnimation[T IAttributeAnimation](
	sprite render.IElement,
	instanceCreator func(base *BaseAttributeAnimation) T,
) T {
	baseAttribute := *sprite.Attribute()
	a := &BaseAttributeAnimation{
		sprite: sprite,

		modifiedFields:  make(map[ModifiedField]struct{}),
		baseAttribute:   baseAttribute,
		targetAttribute: baseAttribute,
		easing: func(progress float32) float32 {
			return progress
		},
	}

	instance := instanceCreator(a)
	a.instance = instance

	return instance
}

type AttributeAnimation struct {
	*BaseAttributeAnimation
}

func NewAttributeAnimation(sprite render.IElement) *AttributeAnimation {
	return BuildAttributeAnimation(sprite, func(base *BaseAttributeAnimation) *AttributeAnimation {
		return &AttributeAnimation{
			BaseAttributeAnimation: base,
		}
	})
}

func (a *BaseAttributeAnimation) SetEasing(easing ti.EasingFunction) IAttributeAnimation {
	a.easing = easing
	return a.instance
}

func (a *BaseAttributeAnimation) BaseAttribute() *ctypes.Attribute {
	return &a.baseAttribute
}

func (a *BaseAttributeAnimation) TargetAttribute() *ctypes.Attribute {
	return &a.targetAttribute
}

// Ticking 按 modifiedFields 将 from->to 插值并应用到目标精灵。
func (a *BaseAttributeAnimation) Ticking(t float32) {
	if a.sprite == nil {
		return
	}

	t = a.easing(t)

	hasWidth := a.IsModified(ModifiedFieldWidth)
	hasHeight := a.IsModified(ModifiedFieldHeight)
	if hasWidth || hasHeight {
		w := a.baseAttribute.Width()
		h := a.baseAttribute.Height()
		if hasWidth {
			w = misc.Lerp(a.baseAttribute.Width(), a.targetAttribute.Width(), t)
			if w <= 0 {
				w = 1
			}
		}
		if hasHeight {
			h = misc.Lerp(a.baseAttribute.Height(), a.targetAttribute.Height(), t)
			if h <= 0 {
				h = 1
			}
		}
		a.sprite.Resize(w, h)
	}

	if a.IsModified(ModifiedFieldX) || a.IsModified(ModifiedFieldY) {
		x := a.baseAttribute.X()
		y := a.baseAttribute.Y()
		if a.IsModified(ModifiedFieldX) {
			x = misc.Lerp(a.baseAttribute.X(), a.targetAttribute.X(), t)
		}
		if a.IsModified(ModifiedFieldY) {
			y = misc.Lerp(a.baseAttribute.Y(), a.targetAttribute.Y(), t)
		}
		a.sprite.MoveTo(x, y)
	}

	if a.IsModified(ModifiedFieldCx) {
		a.sprite.SetCx(misc.Lerp(a.baseAttribute.Cx(), a.targetAttribute.Cx(), t))
	}
	if a.IsModified(ModifiedFieldCy) {
		a.sprite.SetCy(misc.Lerp(a.baseAttribute.Cy(), a.targetAttribute.Cy(), t))
	}

	if a.IsModified(ModifiedFieldScaleX) || a.IsModified(ModifiedFieldScaleY) {
		scaleX := a.baseAttribute.ScaleX()
		scaleY := a.baseAttribute.ScaleY()
		if a.IsModified(ModifiedFieldScaleX) {
			scaleX = misc.Lerp(a.baseAttribute.ScaleX(), a.targetAttribute.ScaleX(), t)
		}
		if a.IsModified(ModifiedFieldScaleY) {
			scaleY = misc.Lerp(a.baseAttribute.ScaleY(), a.targetAttribute.ScaleY(), t)
		}
		a.sprite.SetScale(scaleX, scaleY)
	}

	if a.IsModified(ModifiedFieldRotation) {
		a.sprite.SetRotation(misc.Lerp(a.baseAttribute.Rotation(), a.targetAttribute.Rotation(), t))
	}
	if a.IsModified(ModifiedFieldAlpha) {
		a.sprite.SetAlpha(misc.Lerp(a.baseAttribute.Alpha(), a.targetAttribute.Alpha(), t))
	}

}
