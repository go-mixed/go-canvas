package effect

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ti"
)

type WipeEffect struct {
	inOut EffectInOut

	shapeMaskOptions *ti.ShapeMaskOptions
	easing           ti.EasingFunction
}

func Wipe(inOut EffectInOut) *WipeEffect {
	return &WipeEffect{inOut: inOut, shapeMaskOptions: ti.ShapeMaskOpt(), easing: ti.DefaultEasingFunction}
}

func (e *WipeEffect) WithShapeType(shapeType ti.ShapeType) *WipeEffect {
	e.shapeMaskOptions.SetShapeType(shapeType)
	return e
}
func (e *WipeEffect) WithDirection(direction ti.Direction) *WipeEffect {
	e.shapeMaskOptions.SetDirection(direction)
	return e
}

func (e *WipeEffect) WithColor(c color.Color) *WipeEffect {
	e.shapeMaskOptions.SetColor(c)
	return e
}

func (e *WipeEffect) WithTRange(startT, endT float32) *WipeEffect {
	e.shapeMaskOptions.SetTRange(startT, endT)
	return e
}

func (e *WipeEffect) WithFeather(radius uint32, mode ti.FeatherMode) *WipeEffect {
	e.shapeMaskOptions.SetFeather(radius, mode)
	return e
}

func (e *WipeEffect) WithShapeMaskOptions(opts *ti.ShapeMaskOptions) *WipeEffect {
	if opts == nil {
		return e
	}
	copied := *opts
	if opts.ShapeOptions != nil {
		shapeOptCopy := *opts.ShapeOptions
		copied.ShapeOptions = &shapeOptCopy
	}
	e.shapeMaskOptions = &copied
	return e
}

func (e *WipeEffect) WithEasing(fn ti.EasingFunction) *WipeEffect {
	if fn != nil {
		e.easing = fn
	}
	return e
}

func (e *WipeEffect) WithEasingName(name string) *WipeEffect {
	e.easing = ti.GetEasingFunction(name)
	return e
}

func (e *WipeEffect) TargetAttributeFn(base ti.Attribute) (*ti.Attribute, *ti.TargetAttribute) {
	target := ti.TargetAttr().SetEasing(e.easing)
	opts := *e.shapeMaskOptions
	if e.shapeMaskOptions.ShapeOptions != nil {
		shapeOptCopy := *e.shapeMaskOptions.ShapeOptions
		opts.ShapeOptions = &shapeOptCopy
	}
	if e.inOut == EffectOut {
		opts.SetTRange(2.0, 0.0)
	} else {
		opts.SetTRange(0.0, 2.0)
	}
	target.SetShapeOptions(&opts)
	return &base, target
}
