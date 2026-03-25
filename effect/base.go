package effect

import (
	"slideshow/misc"
	"slideshow/render"
)

type IEffect interface {
	Apply(sprite render.ISprite, progress float32)
}

type TransitionType int32

const (
	TransitionTypeNone TransitionType = iota
	TransitionTypeFade
	TransitionTypeZoom
	TransitionTypeRotate
	TransitionTypeSlide
)

type Effect struct {
	durationMs int32
	easingFn   EasingFunction
}

func NewEffect(durationMs int32, easingFn EasingFunction) *Effect {
	if easingFn == nil {
		easingFn = DefaultEasingFunction
	}
	return &Effect{
		durationMs: durationMs,
		easingFn:   easingFn,
	}
}

func (e *Effect) getEaseProgress(progress float32) float32 {
	progress = misc.Clamp(progress)
	return e.easingFn(progress)
}

type TransitionEffect struct {
	*Effect
	transition TransitionType
}

var _ IEffect = (*TransitionEffect)(nil)

func NewTransitionEffect(durationMs int32, easingFn EasingFunction, transition TransitionType) *TransitionEffect {
	return &TransitionEffect{
		Effect:     NewEffect(durationMs, easingFn),
		transition: transition,
	}
}

func (e *TransitionEffect) Apply(sprite render.ISprite, progress float32) {
	panic("not implemented")
}
