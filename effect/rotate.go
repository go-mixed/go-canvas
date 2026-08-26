package effect

import (
	"github.com/go-mixed/go-canvas/animation"
	"github.com/go-mixed/go-canvas/internel/misc"
	"github.com/go-mixed/go-canvas/render"
	"github.com/go-mixed/go-canvas/ti"
)

type RotateEffect struct {
	inOut      EffectInOut
	angleStart float32
	angleEnd   float32
	scaleStart float32
	scaleEnd   float32
	easing     ti.EasingFunction
}

func Rotate(inOut EffectInOut) *RotateEffect {
	return &RotateEffect{
		inOut:      inOut,
		angleStart: 0,
		angleEnd:   360,
		scaleStart: 0.1,
		scaleEnd:   1,
		easing:     ti.DefaultEasingFunction,
	}
}

func (e *RotateEffect) WithAngleRange(start, end float32) *RotateEffect {
	e.angleStart, e.angleEnd = start, end
	return e
}

func (e *RotateEffect) WithScaleRange(start, end float32) *RotateEffect {
	e.scaleStart, e.scaleEnd = start, end
	return e
}

func (e *RotateEffect) WithEasing(fn ti.EasingFunction) *RotateEffect {
	if fn != nil {
		e.easing = fn
	}
	return e
}

func (e *RotateEffect) WithEasingName(name string) *RotateEffect {
	e.easing = ti.GetEasingFunction(name)
	return e
}

func (e *RotateEffect) AnimateFn(sprite render.IElement) render.TickingFn {
	animation := animation.NewAttributeAnimation(sprite).SetEasing(e.easing)
	if e.inOut == EffectOut {
		animation.SetRotation(misc.Deg2Rad(e.angleStart))
		animation.SetScale(e.scaleStart, e.scaleStart)
	} else {
		animation.SetRotation(misc.Deg2Rad(e.angleEnd))
		animation.SetScale(e.scaleEnd, e.scaleEnd)
	}
	animation.SetAlpha(1.0)
	return animation.Ticking
}
