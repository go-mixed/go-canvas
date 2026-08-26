package animation

import (
	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/internel/misc"
	"github.com/go-mixed/go-canvas/render"
)

const ModifiedFieldShape ModifiedField = "shape"

type ShapeAttributeAnimation struct {
	*BaseAttributeAnimation

	shapeOpts *ctypes.ShapeMaskOptions
}

func NewShapeAttributeAnimation(sprite render.IElement) *ShapeAttributeAnimation {
	if _, ok := sprite.(render.IShape); !ok {
		panic("sprite must be ShapeSprite")
	}
	return BuildAttributeAnimation(sprite, func(base *BaseAttributeAnimation) *ShapeAttributeAnimation {
		return &ShapeAttributeAnimation{
			BaseAttributeAnimation: base,
			shapeOpts:              nil,
		}
	})
}

func (a *ShapeAttributeAnimation) SetShapeOptions(opts *ctypes.ShapeMaskOptions) *ShapeAttributeAnimation {
	a.shapeOpts = opts
	if opts != nil {
		a.MarkModified(ModifiedFieldShape)
	}
	return a
}

func (a *ShapeAttributeAnimation) Ticking(t float32) {
	// 先调用父级的Ticking
	a.BaseAttributeAnimation.Ticking(t)

	// 检查是否有Shape的更改
	if a.IsModified(ModifiedFieldShape) && a.shapeOpts != nil {
		t = a.easing(t)

		if shapeSprite, ok := a.sprite.(render.IShape); ok {
			opts := a.shapeOpts
			tVal := misc.Lerp(opts.StartT, opts.EndT, t)
			shapeSprite.DrawShape(opts.ShapeType, tVal, opts.ShapeOptions)
		}
	}
}
