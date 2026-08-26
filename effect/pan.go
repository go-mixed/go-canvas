package effect

import (
	"math"

	"github.com/go-mixed/go-canvas/animation"
	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/render"
	"github.com/go-mixed/go-canvas/ti"
)

var panDirections = map[ctypes.Direction][2]float32{
	ctypes.DirectionTop:         {0, -1},
	ctypes.DirectionBottom:      {0, 1},
	ctypes.DirectionLeft:        {-1, 0},
	ctypes.DirectionRight:       {1, 0},
	ctypes.DirectionTopLeft:     {-1, -1},
	ctypes.DirectionTopRight:    {1, -1},
	ctypes.DirectionBottomLeft:  {-1, 1},
	ctypes.DirectionBottomRight: {1, 1},
	ctypes.DirectionCenter:      {0, 0},
}

type PanEffect struct {
	inOut        EffectInOut
	direction    ctypes.Direction
	panIntensity float32
	zoomStart    float32
	zoomEnd      float32
	easing       ti.EasingFunction
}

func Pan(inOut EffectInOut) *PanEffect {
	return &PanEffect{
		inOut:        inOut,
		direction:    ctypes.DirectionCenter,
		panIntensity: 0.1,
		zoomStart:    1.0,
		zoomEnd:      1.2,
		easing:       ti.DefaultEasingFunction,
	}
}

func (e *PanEffect) WithDirection(direction ctypes.Direction) *PanEffect {
	e.direction = direction
	return e
}
func (e *PanEffect) WithPanIntensity(intensity float32) *PanEffect {
	e.panIntensity = intensity
	return e
}
func (e *PanEffect) WithZoomRange(start, end float32) *PanEffect {
	e.zoomStart, e.zoomEnd = start, end
	return e
}
func (e *PanEffect) WithEasing(fn ti.EasingFunction) *PanEffect {
	if fn != nil {
		e.easing = fn
	}
	return e
}
func (e *PanEffect) WithEasingName(name string) *PanEffect {
	e.easing = ti.GetEasingFunction(name)
	return e
}

func (e *PanEffect) AnimateFn(sprite render.IElement) render.TickingFn {
	animation := animation.NewAttributeAnimation(sprite).SetEasing(e.easing)
	baseAttribute := animation.BaseAttribute()
	vec := panDirections[e.direction]
	dx, dy := vec[0], vec[1]
	maxPanX := float32(baseAttribute.Width()) * e.panIntensity
	maxPanY := float32(baseAttribute.Height()) * e.panIntensity

	if e.inOut == EffectOut {
		animation.SetScale(e.zoomStart, e.zoomStart)
		animation.SetX(baseAttribute.X())
		animation.SetY(baseAttribute.Y())
	} else {
		animation.SetScale(e.zoomEnd, e.zoomEnd)
		animation.SetX(baseAttribute.X() + int(math.Round(float64(dx*maxPanX))))
		animation.SetY(baseAttribute.Y() + int(math.Round(float64(dy*maxPanY))))
	}
	animation.SetAlpha(1.0)
	return animation.Ticking
}
