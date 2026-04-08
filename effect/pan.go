package effect

import (
	"math"

	"github.com/go-mixed/go-canvas/ti"
)

var panDirections = map[ti.Direction][2]float32{
	ti.DirectionTop:         {0, -1},
	ti.DirectionBottom:      {0, 1},
	ti.DirectionLeft:        {-1, 0},
	ti.DirectionRight:       {1, 0},
	ti.DirectionTopLeft:     {-1, -1},
	ti.DirectionTopRight:    {1, -1},
	ti.DirectionBottomLeft:  {-1, 1},
	ti.DirectionBottomRight: {1, 1},
	ti.DirectionCenter:      {0, 0},
}

type PanEffect struct {
	inOut        EffectInOut
	direction    ti.Direction
	panIntensity float32
	zoomStart    float32
	zoomEnd      float32
	easing       ti.EasingFunction
}

func Pan(inOut EffectInOut) *PanEffect {
	return &PanEffect{
		inOut:        inOut,
		direction:    ti.DirectionCenter,
		panIntensity: 0.1,
		zoomStart:    1.0,
		zoomEnd:      1.2,
		easing:       ti.DefaultEasingFunction,
	}
}

func (e *PanEffect) WithDirection(direction ti.Direction) *PanEffect {
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

func (e *PanEffect) TargetAttributeFn(base *ti.Attribute) *ti.TargetAttribute {
	if base == nil {
		base = ti.Attr()
	}
	vec := panDirections[e.direction]
	dx, dy := vec[0], vec[1]
	maxPanX := float32(base.Width()) * e.panIntensity
	maxPanY := float32(base.Height()) * e.panIntensity

	target := ti.TargetAttr().SetEasing(e.easing)
	if e.inOut == EffectOut {
		target.SetScale(e.zoomStart, e.zoomStart)
		target.SetX(base.X())
		target.SetY(base.Y())
	} else {
		target.SetScale(e.zoomEnd, e.zoomEnd)
		target.SetX(base.X() + int(math.Round(float64(dx*maxPanX))))
		target.SetY(base.Y() + int(math.Round(float64(dy*maxPanY))))
	}
	target.SetAlpha(1.0)
	return target
}
