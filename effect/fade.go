package effect

import (
	"github.com/go-mixed/go-canvas/animation"
	"github.com/go-mixed/go-canvas/internel/misc"
	"github.com/go-mixed/go-canvas/render"
	"github.com/go-mixed/go-canvas/ti"
)

type FadeEffect struct {
	inOut  EffectInOut
	easing ti.EasingFunction
}

func Fade(inOut EffectInOut) *FadeEffect {
	return &FadeEffect{inOut: inOut, easing: ti.DefaultEasingFunction}
}
func (e *FadeEffect) WithEasing(fn ti.EasingFunction) *FadeEffect {
	if fn != nil {
		e.easing = fn
	}
	return e
}
func (e *FadeEffect) WithEasingName(name string) *FadeEffect {
	e.easing = ti.GetEasingFunction(name)
	return e
}
func (e *FadeEffect) AnimateFn(sprite render.IElement) render.TickingFn {
	animation := animation.NewAttributeAnimation(sprite).SetEasing(e.easing)
	baseAttribute := animation.BaseAttribute()
	if e.inOut == EffectOut {
		if misc.NumberEqual(baseAttribute.Alpha(), 0, misc.Epsilon) {
			baseAttribute.SetAlpha(1)
		}
		animation.SetAlpha(0.3)
	} else {
		if misc.NumberEqual(baseAttribute.Alpha(), 1, misc.Epsilon) {
			baseAttribute.SetAlpha(0.3)
		}
		animation.SetAlpha(1.0)
	}
	return animation.Ticking
}
