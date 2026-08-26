package effect

import (
	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/render"
)

func transitionFactory(name string, inOut EffectInOut) render.AnimateFn {
	switch name {
	case "pan_left":
		return Pan(inOut).WithDirection(ctypes.DirectionLeft).AnimateFn
	case "pan_right":
		return Pan(inOut).WithDirection(ctypes.DirectionRight).AnimateFn
	case "pan_top":
		return Pan(inOut).WithDirection(ctypes.DirectionTop).AnimateFn
	case "pan_bottom":
		return Pan(inOut).WithDirection(ctypes.DirectionBottom).AnimateFn
	case "pan_top_left", "pan_left_top":
		return Pan(inOut).WithDirection(ctypes.DirectionTopLeft).AnimateFn
	case "pan_top_right", "pan_right_top":
		return Pan(inOut).WithDirection(ctypes.DirectionTopRight).AnimateFn
	case "pan_bottom_left", "pan_left_bottom":
		return Pan(inOut).WithDirection(ctypes.DirectionBottomLeft).AnimateFn
	case "pan_bottom_right", "pan_right_bottom":
		return Pan(inOut).WithDirection(ctypes.DirectionBottomRight).AnimateFn
	case "pan_center":
		return Pan(inOut).WithDirection(ctypes.DirectionCenter).AnimateFn
	case "rotate":
		return Rotate(inOut).AnimateFn
	case "top":
		return Slide(inOut).WithDirection(ctypes.DirectionTop).AnimateFn
	case "bottom":
		return Slide(inOut).WithDirection(ctypes.DirectionBottom).AnimateFn
	case "left":
		return Slide(inOut).WithDirection(ctypes.DirectionLeft).AnimateFn
	case "right":
		return Slide(inOut).WithDirection(ctypes.DirectionRight).AnimateFn
	case "zoom":
		return Zoom(inOut).AnimateFn
	case "heart":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeHeart).AnimateFn
	case "star5":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeStar5).AnimateFn
	case "cross":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeCross).AnimateFn
	case "linear":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeLinear).AnimateFn
	case "circle":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeCircle).AnimateFn
	case "diamond":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeDiamond).AnimateFn
	case "rectangle":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeRectangle).AnimateFn
	case "triangle":
		return Wipe(inOut).WithShapeType(ctypes.ShapeTypeTriangle).AnimateFn
	default:
		return Fade(inOut).AnimateFn
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
func GetTransitionEffect(name string, inOut EffectInOut) render.AnimateFn {
	return transitionFactory(name, inOut)
}
