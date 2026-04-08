package effect

import "github.com/go-mixed/go-canvas/ti"

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
func (e *FadeEffect) TargetAttributeFn(base ti.Attribute) (*ti.Attribute, *ti.TargetAttribute) {
	target := ti.TargetAttr().SetEasing(e.easing)
	if e.inOut == EffectOut {
		target.SetAlpha(0.0)
	} else {
		target.SetAlpha(1.0)
	}
	return &base, target
}
