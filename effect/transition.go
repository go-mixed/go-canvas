package effect

import (
	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/ti"
)

func transitionFactory(name string, inOut EffectInOut) ti.TargetAttributeFn {
	switch name {
	case "pan_left":
		return Pan(inOut).WithDirection(ctypes.DirectionLeft).TargetAttributeFn
	case "pan_right":
		return Pan(inOut).WithDirection(ctypes.DirectionRight).TargetAttributeFn
	case "pan_top":
		return Pan(inOut).WithDirection(ctypes.DirectionTop).TargetAttributeFn
	case "pan_bottom":
		return Pan(inOut).WithDirection(ctypes.DirectionBottom).TargetAttributeFn
	case "pan_top_left", "pan_left_top":
		return Pan(inOut).WithDirection(ctypes.DirectionTopLeft).TargetAttributeFn
	case "pan_top_right", "pan_right_top":
		return Pan(inOut).WithDirection(ctypes.DirectionTopRight).TargetAttributeFn
	case "pan_bottom_left", "pan_left_bottom":
		return Pan(inOut).WithDirection(ctypes.DirectionBottomLeft).TargetAttributeFn
	case "pan_bottom_right", "pan_right_bottom":
		return Pan(inOut).WithDirection(ctypes.DirectionBottomRight).TargetAttributeFn
	case "pan_center":
		return Pan(inOut).WithDirection(ctypes.DirectionCenter).TargetAttributeFn
	case "rotate":
		return Rotate(inOut).TargetAttributeFn
	case "top":
		return Slide(inOut).WithDirection(ctypes.DirectionTop).TargetAttributeFn
	case "bottom":
		return Slide(inOut).WithDirection(ctypes.DirectionBottom).TargetAttributeFn
	case "left":
		return Slide(inOut).WithDirection(ctypes.DirectionLeft).TargetAttributeFn
	case "right":
		return Slide(inOut).WithDirection(ctypes.DirectionRight).TargetAttributeFn
	case "zoom":
		return Zoom(inOut).TargetAttributeFn
	case "heart":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeHeart).TargetAttributeFn
	case "star5":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeStar5).TargetAttributeFn
	case "cross":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeCross).TargetAttributeFn
	case "linear":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeLinear).TargetAttributeFn
	case "circle":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeCircle).TargetAttributeFn
	case "diamond":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeDiamond).TargetAttributeFn
	case "rectangle":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeRectangle).TargetAttributeFn
	case "triangle":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeTriangle).TargetAttributeFn
	default:
		return Fade(inOut).TargetAttributeFn
	}
}

// IsWipeEffect 是否是 wipe 效果
func IsWipeEffect(name string) bool {
	switch name {
	case "heart", "star5", "cross", "linear", "circle", "diamond", "rectangle", "triangle":
		return true
	default:
		return false
	}
}

// GetTransitionEffect 获取过渡效果
func GetTransitionEffect(name string, inOut EffectInOut) ti.TargetAttributeFn {
	return transitionFactory(name, inOut)
}
