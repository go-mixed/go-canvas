package animation

import (
	"github.com/go-mixed/go-canvas/ctypes"
)

// MarkModified 标记修改后的字段类型
func (a *BaseAttributeAnimation) MarkModified(fields ...ModifiedField) {
	for _, field := range fields {
		if _, ok := a.modifiedFields[field]; ok {
			continue
		}
		a.modifiedFields[field] = struct{}{}
		a.order = append(a.order, field)
	}
}

func (a *BaseAttributeAnimation) ModifiedFields() []ModifiedField {
	out := make([]ModifiedField, 0, len(a.order))
	out = append(out, a.order...)
	return out
}

func (a *BaseAttributeAnimation) IsModified(field ModifiedField) bool {
	_, ok := a.modifiedFields[field]
	return ok
}

func (a *BaseAttributeAnimation) SetRect(rect ctypes.Rectangle[int]) IAttributeAnimation {
	a.targetAttribute.SetRect(rect)
	a.MarkModified(ModifiedFieldX, ModifiedFieldY, ModifiedFieldWidth, ModifiedFieldHeight)
	return a.instance
}

func (a *BaseAttributeAnimation) SetXYWH(x, y, width, height int) IAttributeAnimation {
	a.targetAttribute.SetXYWH(x, y, width, height)
	a.MarkModified(ModifiedFieldX, ModifiedFieldY, ModifiedFieldWidth, ModifiedFieldHeight)
	return a.instance
}

func (a *BaseAttributeAnimation) SetXY(x, y int) IAttributeAnimation {
	a.targetAttribute.SetXY(x, y)
	a.MarkModified(ModifiedFieldX, ModifiedFieldY)
	return a.instance
}

func (a *BaseAttributeAnimation) SetWH(width, height int) IAttributeAnimation {
	a.targetAttribute.SetWH(width, height)
	a.MarkModified(ModifiedFieldWidth, ModifiedFieldHeight)
	return a.instance
}

func (a *BaseAttributeAnimation) SetX(x int) IAttributeAnimation {
	a.targetAttribute.SetX(x)
	a.MarkModified(ModifiedFieldX)
	return a.instance
}

func (a *BaseAttributeAnimation) SetY(y int) IAttributeAnimation {
	a.targetAttribute.SetY(y)
	a.MarkModified(ModifiedFieldY)
	return a.instance
}

func (a *BaseAttributeAnimation) SetWidth(width int) IAttributeAnimation {
	a.targetAttribute.SetWidth(width)
	a.MarkModified(ModifiedFieldWidth)
	return a.instance
}

func (a *BaseAttributeAnimation) SetHeight(height int) IAttributeAnimation {
	a.targetAttribute.SetHeight(height)
	a.MarkModified(ModifiedFieldHeight)
	return a.instance
}

func (a *BaseAttributeAnimation) MoveTo(x, y int) IAttributeAnimation {
	a.targetAttribute.MoveTo(x, y)
	a.MarkModified(ModifiedFieldX, ModifiedFieldY)
	return a.instance
}

func (a *BaseAttributeAnimation) SetScale(x, y float32) IAttributeAnimation {
	a.targetAttribute.SetScale(x, y)
	a.MarkModified(ModifiedFieldScaleX, ModifiedFieldScaleY)
	return a.instance
}

func (a *BaseAttributeAnimation) SetRotation(rotation float32) IAttributeAnimation {
	a.targetAttribute.SetRotation(rotation)
	a.MarkModified(ModifiedFieldRotation)
	return a.instance
}

func (a *BaseAttributeAnimation) SetAlpha(alpha float32) IAttributeAnimation {
	a.targetAttribute.SetAlpha(alpha)
	a.MarkModified(ModifiedFieldAlpha)
	return a.instance
}

func (a *BaseAttributeAnimation) SetCx(cx int) IAttributeAnimation {
	a.targetAttribute.SetCx(cx)
	a.MarkModified(ModifiedFieldCx)
	return a.instance
}

func (a *BaseAttributeAnimation) SetCy(cy int) IAttributeAnimation {
	a.targetAttribute.SetCy(cy)
	a.MarkModified(ModifiedFieldCy)
	return a.instance
}

func (a *BaseAttributeAnimation) SetCxy(x int, y int) IAttributeAnimation {
	a.targetAttribute.SetCxy(x, y)
	a.MarkModified(ModifiedFieldCx, ModifiedFieldCy)
	return a.instance
}

//func (a *AttributeAnimation) SetShapeOptions(opts *ctypes.ShapeMaskOptions) IAnimation{
//	a.ShapeOpts = opts
//	if opts != nil {
//		a.mark(ModifiedFieldShape)
//	}
//	return a.instance
//}
