package effect

import (
	"slideshow/misc"
	"slideshow/ti"
)

type panOptions struct {
	// 平移强度
	PanIntensity float32
}

// rotateOptions 旋转效果配置
type rotateOptions struct {
	// 旋转起始角度
	AngleStart float32
	// 旋转结束角度
	AngleEnd float32
	// 旋转起始缩放
	ScaleStart float32
	// 旋转结束缩放
	ScaleEnd float32
}

// wipeOptions 擦除效果配置
type wipeOptions struct {
	ShapeType ti.ShapeType
}

type zoomOptions struct {
	ZoomStart float32
	ZoomEnd   float32
}

type effectOptions struct {
	panOptions    panOptions
	rotateOptions rotateOptions
	wipeOptions   wipeOptions
	zoomOptions   zoomOptions

	easingFn  EasingFunction
	direction ti.Direction
}

func buildOptionFnFromMap(opt map[string]string) []optionFn {
	if len(opt) == 0 {
		return nil
	}

	var fns []optionFn

	// direction
	if v, ok := opt["direction"]; ok {
		fns = append(fns, WithDirectionStr(v))
	}

	// easing
	if v, ok := opt["easing"]; ok {
		fns = append(fns, WithEasing(v))
	}

	// panIntensity
	if v, ok := misc.MapGetFloat(opt, "pan_intensity"); ok {
		fns = append(fns, WithPanIntensity(float32(v)))
	}

	// angle_start, angle_end
	if angle, ok := misc.MapMultiGetFloat(opt, "angle_start", "angle_end"); ok {
		fns = append(fns, WithRotateAngle(float32(angle[0]), float32(angle[1])))
	}

	// scale_start, scale_end
	if scale, ok := misc.MapMultiGetFloat(opt, "scale_start", "scale_end"); ok {
		fns = append(fns, WithRotateScale(float32(scale[0]), float32(scale[1])))
	}

	// zoom_start, zoom_end
	if zoom, ok := misc.MapMultiGetFloat(opt, "zoom_start", "zoom_end"); ok {
		fns = append(fns, WithZoomRange(float32(zoom[0]), float32(zoom[1])))
	}

	// shape_type (for wipe effect)
	if v, ok := opt["shape_type"]; ok {
		fns = append(fns, WithShapeTypeStr(v))
	}

	return fns
}

func toOptions(initialOptions effectOptions, opts ...optionFn) effectOptions {
	for _, opt := range opts {
		opt(&initialOptions)
	}
	return initialOptions
}

type optionFn func(options *effectOptions)

func WithDirectionStr(direction string) func(options *effectOptions) {
	return func(options *effectOptions) {
		options.direction = ti.DirectionFromString(direction)
	}
}

func WithDirection(direction ti.Direction) func(options *effectOptions) {
	return func(options *effectOptions) {
		options.direction = direction
	}
}

func WithZoomRange(zoomRangeMin, zoomRangeMax float32) func(options *effectOptions) {
	return func(options *effectOptions) {
		options.zoomOptions.ZoomStart = zoomRangeMin
		options.zoomOptions.ZoomEnd = zoomRangeMax
	}
}

func WithPanIntensity(panIntensity float32) func(options *effectOptions) {
	return func(options *effectOptions) {
		options.panOptions.PanIntensity = panIntensity
	}
}

func WithEasing(easingName string) func(options *effectOptions) {
	return func(options *effectOptions) {
		options.easingFn = GetEasingFunction(easingName)
	}
}

// WithRotateAngle 设置旋转角度范围
func WithRotateAngle(angleStart, angleEnd float32) func(*effectOptions) {
	return func(opt *effectOptions) {
		opt.rotateOptions.AngleStart = angleStart
		opt.rotateOptions.AngleEnd = angleEnd
	}
}

// WithRotateScale 设置旋转时的缩放范围
func WithRotateScale(scaleStart, scaleEnd float32) func(*effectOptions) {
	return func(opt *effectOptions) {
		opt.rotateOptions.ScaleStart = scaleStart
		opt.rotateOptions.ScaleEnd = scaleEnd
	}
}

func WithShapeTypeStr(shapeType string) func(*effectOptions) {
	return func(opt *effectOptions) {
		opt.wipeOptions.ShapeType = ti.ShapeTypeFromString(shapeType)
	}
}
func WithShapeType(shapeType ti.ShapeType) func(*effectOptions) {
	return func(opt *effectOptions) {
		opt.wipeOptions.ShapeType = shapeType
	}
}
