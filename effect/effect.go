package effect

import (
	"slideshow/misc"
	"slideshow/render"
)

type IEffect interface {
	Apply(sprite render.ISprite, progress float32)
}

type Effect struct {
	options   effectOptions
	direction EffectInOut
}

func newEffect(inOut EffectInOut, options effectOptions) *Effect {
	if options.easingFn == nil {
		options.easingFn = DefaultEasingFunction
	}
	return &Effect{
		options:   options,
		direction: inOut,
	}
}

func (e *Effect) getEaseProgress(progress float32) float32 {
	progress = misc.Clamp(progress)
	eased := e.options.easingFn(progress)
	if e.direction == EffectOut {
		eased = 1.0 - eased
	}

	return eased
}
