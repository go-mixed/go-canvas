package effect

import "slideshow/ti"

type panOptions struct {
	Direction    ti.Direction
	ZoomRangeMin float32
	ZoomRangeMax float32
	PanIntensity float32
	Easing       EasingFunction
}

func PanDirection(direction string) func(options *panOptions) {
	return func(options *panOptions) {
		options.Direction = ti.DirectionFromString(direction)
	}
}

func PanZoomRange(zoomRangeMin, zoomRangeMax float32) func(options *panOptions) {
	return func(options *panOptions) {
		options.ZoomRangeMin = zoomRangeMin
		options.ZoomRangeMax = zoomRangeMax
	}
}

func PanPanIntensity(panIntensity float32) func(options *panOptions) {
	return func(options *panOptions) {
		options.PanIntensity = panIntensity
	}
}

func PanEasing(easingName string) func(options *panOptions) {
	return func(options *panOptions) {
		options.Easing = GetEasingFunction(easingName)
	}
}
