package effect

import "github.com/go-mixed/go-canvas/ti"

type ZoomEffect struct {
	inOut     EffectInOut
	zoomStart float32
	zoomEnd   float32
	easing    ti.EasingFunction
}

func Zoom(inOut EffectInOut) *ZoomEffect {
	return &ZoomEffect{inOut: inOut, zoomStart: 0.5, zoomEnd: 1, easing: ti.DefaultEasingFunction}
}
func (e *ZoomEffect) WithZoomRange(start, end float32) *ZoomEffect {
	e.zoomStart, e.zoomEnd = start, end
	return e
}
func (e *ZoomEffect) WithEasing(fn ti.EasingFunction) *ZoomEffect {
	if fn != nil {
		e.easing = fn
	}
	return e
}
func (e *ZoomEffect) WithEasingName(name string) *ZoomEffect {
	e.easing = ti.GetEasingFunction(name)
	return e
}
func (e *ZoomEffect) TargetAttributeFn(base *ti.Attribute) *ti.TargetAttribute {
	target := ti.TargetAttr().SetEasing(e.easing)
	if e.inOut == EffectOut {
		target.SetScale(e.zoomStart, e.zoomStart)
	} else {
		target.SetScale(e.zoomEnd, e.zoomEnd)
	}
	target.SetAlpha(1.0)
	return target
}
