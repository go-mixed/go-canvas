package ti

import "github.com/go-mixed/go-canvas/ctypes"

type TargetAttributeFn func(base ctypes.Attribute) (from *ctypes.Attribute, to *TargetAttribute)

const (
	ModifiedFieldX        = "x"
	ModifiedFieldY        = "y"
	ModifiedFieldWidth    = "width"
	ModifiedFieldHeight   = "height"
	ModifiedFieldCx       = "cx"
	ModifiedFieldCy       = "cy"
	ModifiedFieldScaleX   = "scale_x"
	ModifiedFieldScaleY   = "scale_y"
	ModifiedFieldRotation = "rotation"
	ModifiedFieldAlpha    = "alpha"
	ModifiedFieldShape    = "shape"
)

type TargetAttribute struct {
	*ctypes.Attribute

	modifiedFields map[string]struct{}
	order          []string
	easing         EasingFunction
	ShapeOpts      *ctypes.ShapeMaskOptions
}

func TargetAttr() *TargetAttribute {
	return &TargetAttribute{
		Attribute:      ctypes.Attr(),
		modifiedFields: map[string]struct{}{},
		easing: func(progress float32) float32 {
			return progress
		},
	}
}

func (a *TargetAttribute) mark(fields ...string) {
	for _, field := range fields {
		if _, ok := a.modifiedFields[field]; ok {
			continue
		}
		a.modifiedFields[field] = struct{}{}
		a.order = append(a.order, field)
	}
}

func (a *TargetAttribute) ModifiedFields() []string {
	out := make([]string, 0, len(a.order))
	out = append(out, a.order...)
	return out
}

func (a *TargetAttribute) IsModified(field string) bool {
	_, ok := a.modifiedFields[field]
	return ok
}

func (a *TargetAttribute) SetEasing(fn EasingFunction) *TargetAttribute {
	if fn != nil {
		a.easing = fn
	}
	return a
}

func (a *TargetAttribute) Easing(t float32) float32 {
	return a.easing(t)
}

func (a *TargetAttribute) SetRect(rect ctypes.Rectangle[int]) *TargetAttribute {
	a.Attribute.SetRect(rect)
	a.mark(ModifiedFieldX, ModifiedFieldY, ModifiedFieldWidth, ModifiedFieldHeight)
	return a
}

func (a *TargetAttribute) SetXYWH(x, y, width, height int) *TargetAttribute {
	a.Attribute.SetXYWH(x, y, width, height)
	a.mark(ModifiedFieldX, ModifiedFieldY, ModifiedFieldWidth, ModifiedFieldHeight)
	return a
}

func (a *TargetAttribute) SetXY(x, y int) *TargetAttribute {
	a.Attribute.SetXY(x, y)
	a.mark(ModifiedFieldX, ModifiedFieldY)
	return a
}

func (a *TargetAttribute) SetWH(width, height int) *TargetAttribute {
	a.Attribute.SetWH(width, height)
	a.mark(ModifiedFieldWidth, ModifiedFieldHeight)
	return a
}

func (a *TargetAttribute) SetX(x int) *TargetAttribute {
	a.Attribute.SetX(x)
	a.mark(ModifiedFieldX)
	return a
}

func (a *TargetAttribute) SetY(y int) *TargetAttribute {
	a.Attribute.SetY(y)
	a.mark(ModifiedFieldY)
	return a
}

func (a *TargetAttribute) SetWidth(width int) *TargetAttribute {
	a.Attribute.SetWidth(width)
	a.mark(ModifiedFieldWidth)
	return a
}

func (a *TargetAttribute) SetHeight(height int) *TargetAttribute {
	a.Attribute.SetHeight(height)
	a.mark(ModifiedFieldHeight)
	return a
}

func (a *TargetAttribute) MoveTo(x, y int) *TargetAttribute {
	a.Attribute.MoveTo(x, y)
	a.mark(ModifiedFieldX, ModifiedFieldY)
	return a
}

func (a *TargetAttribute) SetScale(x, y float32) *TargetAttribute {
	a.Attribute.SetScale(x, y)
	a.mark(ModifiedFieldScaleX, ModifiedFieldScaleY)
	return a
}

func (a *TargetAttribute) SetRotation(rotation float32) *TargetAttribute {
	a.Attribute.SetRotation(rotation)
	a.mark(ModifiedFieldRotation)
	return a
}

func (a *TargetAttribute) SetAlpha(alpha float32) *TargetAttribute {
	a.Attribute.SetAlpha(alpha)
	a.mark(ModifiedFieldAlpha)
	return a
}

func (a *TargetAttribute) SetCx(cx int) *TargetAttribute {
	a.Attribute.SetCx(cx)
	a.mark(ModifiedFieldCx)
	return a
}

func (a *TargetAttribute) SetCy(cy int) *TargetAttribute {
	a.Attribute.SetCy(cy)
	a.mark(ModifiedFieldCy)
	return a
}

func (a *TargetAttribute) SetCxy(x int, y int) *TargetAttribute {
	a.Attribute.SetCxy(x, y)
	a.mark(ModifiedFieldCx, ModifiedFieldCy)
	return a
}

func (a *TargetAttribute) SetShapeOptions(opts *ctypes.ShapeMaskOptions) *TargetAttribute {
	a.ShapeOpts = opts
	if opts != nil {
		a.mark(ModifiedFieldShape)
	}
	return a
}
