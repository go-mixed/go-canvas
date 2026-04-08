package effect

import (
	"github.com/go-mixed/go-canvas/ti"
)

func transitionFactory(name string, inOut EffectInOut) ti.TargetAttributeFn {
	switch name {
	case "pan_left":
		return Pan(inOut).WithDirection(ti.DirectionLeft).TargetAttributeFn
	case "pan_right":
		return Pan(inOut).WithDirection(ti.DirectionRight).TargetAttributeFn
	case "pan_top":
		return Pan(inOut).WithDirection(ti.DirectionTop).TargetAttributeFn
	case "pan_bottom":
		return Pan(inOut).WithDirection(ti.DirectionBottom).TargetAttributeFn
	case "pan_top_left":
		return Pan(inOut).WithDirection(ti.DirectionTopLeft).TargetAttributeFn
	case "pan_top_right":
		return Pan(inOut).WithDirection(ti.DirectionTopRight).TargetAttributeFn
	case "pan_bottom_left":
		return Pan(inOut).WithDirection(ti.DirectionBottomLeft).TargetAttributeFn
	case "pan_bottom_right":
		return Pan(inOut).WithDirection(ti.DirectionBottomRight).TargetAttributeFn
	case "pan_center":
		return Pan(inOut).WithDirection(ti.DirectionCenter).TargetAttributeFn
	case "rotate":
		return Rotate(inOut).TargetAttributeFn
	case "slide":
		return Slide(inOut).TargetAttributeFn
	case "zoom":
		return Zoom(inOut).TargetAttributeFn
	case "wipe":
		return Wipe(inOut).TargetAttributeFn
	case "heart":
		return Wipe(inOut).WithShapeType(ti.ShapeTypeHeart).TargetAttributeFn
	case "star5":
		return Wipe(inOut).WithShapeType(ti.ShapeTypeStar5).TargetAttributeFn
	case "cross":
		return Wipe(inOut).WithShapeType(ti.ShapeTypeCross).TargetAttributeFn
	case "linear":
		return Wipe(inOut).WithShapeType(ti.ShapeTypeLinear).TargetAttributeFn
	case "circle":
		return Wipe(inOut).WithShapeType(ti.ShapeTypeCircle).TargetAttributeFn
	case "diamond":
		return Wipe(inOut).WithShapeType(ti.ShapeTypeDiamond).TargetAttributeFn
	case "rectangle":
		return Wipe(inOut).WithShapeType(ti.ShapeTypeRectangle).TargetAttributeFn
	case "triangle":
		return Wipe(inOut).WithShapeType(ti.ShapeTypeTriangle).TargetAttributeFn
	default:
		return Fade(inOut).TargetAttributeFn
	}
}

func GetTransitionEffect(name string, inOut EffectInOut) (ti.TargetAttributeFn, error) {
	return transitionFactory(name, inOut), nil
}
