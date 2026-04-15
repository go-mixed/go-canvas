package effect

import (
	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/ti"
)

type SlideEffect struct {
	inOut     EffectInOut
	direction ctypes.Direction
	easing    ti.EasingFunction
}

func Slide(inOut EffectInOut) *SlideEffect {
	return &SlideEffect{inOut: inOut, direction: ctypes.DirectionRight, easing: ti.DefaultEasingFunction}
}

func (e *SlideEffect) WithDirection(direction ctypes.Direction) *SlideEffect {
	e.direction = direction
	return e
}

func (e *SlideEffect) WithEasing(fn ti.EasingFunction) *SlideEffect {
	if fn != nil {
		e.easing = fn
	}
	return e
}

func (e *SlideEffect) WithEasingName(name string) *SlideEffect {
	e.easing = ti.GetEasingFunction(name)
	return e
}

func (e *SlideEffect) TargetAttributeFn(base ctypes.Attribute) (*ctypes.Attribute, *ti.TargetAttribute) {
	w := base.Width()
	h := base.Height()
	offsetX, offsetY := 0, 0
	switch e.direction {
	case ctypes.DirectionTop:
		offsetY = -h
	case ctypes.DirectionBottom:
		offsetY = h
	case ctypes.DirectionLeft:
		offsetX = -w
	case ctypes.DirectionRight:
		offsetX = w
	}

	target := ti.TargetAttr().SetEasing(e.easing)
	if e.inOut == EffectOut {
		target.MoveTo(base.X()+offsetX, base.Y()+offsetY)
	} else {
		// EffectIn: treat current position as the final on-screen position.
		// Move animation start(base) off-screen first, then animate back to current.
		finalX, finalY := base.X(), base.Y()
		base.MoveTo(finalX+offsetX, finalY+offsetY)
		target.MoveTo(finalX, finalY)
	}
	return &base, target
}
