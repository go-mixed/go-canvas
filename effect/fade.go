package effect

import (
	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/internel/misc"
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
func (e *FadeEffect) TargetAttributeFn(base ctypes.Attribute) (*ctypes.Attribute, *ti.TargetAttribute) {
	target := ti.TargetAttr().SetEasing(e.easing)
	if e.inOut == EffectOut {
		if misc.NumberEqual(base.Alpha(), 0, misc.Epsilon) {
			base.SetAlpha(1)
		}
		target.SetAlpha(0.3)
	} else {
		if misc.NumberEqual(base.Alpha(), 1, misc.Epsilon) {
			base.SetAlpha(0.3)
		}
		target.SetAlpha(1.0)
	}
	return &base, target
}
