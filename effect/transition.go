package effect

import (
	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/render"
	"github.com/go-mixed/go-canvas/ti"
)

// kenBurnsDirections Ken Burns 效果方向映射
// 与 Python 版本保持一致，对角线方向未归一化
var kenBurnsDirections = map[ti.Direction][2]float32{
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

// KenBurnsEffect Ken Burns 效果
// 通过缩放和平移创建经典的 Ken Burns 幻灯片效果
type KenBurnsEffect struct {
	*Effect
}

var _ IEffect = (*KenBurnsEffect)(nil)

// NewKenBurnsEffect 创建 Ken Burns 效果
//
// Args:
//
//	durationMs: 持续时间（毫秒）
//	options: 平移选项
func NewKenBurnsEffect(effectInOut EffectInOut, opts ...optionFn) IEffect {
	options := effectOptions{
		zoomOptions: zoomOptions{
			ZoomStart: 1.0,
			ZoomEnd:   1.2,
		},
		panOptions: panOptions{
			PanIntensity: 0.1,
		},
	}
	return &KenBurnsEffect{
		Effect: newEffect(effectInOut, toOptions(options, opts...)),
	}
}

// Apply 应用 Ken Burns 效果
//
// 实现 IEffect 接口
func (e *KenBurnsEffect) Apply(sprite render.ISpriteOperator, progress float32) {
	eased := e.getEaseProgress(progress)

	// 计算缩放
	scale := e.options.zoomOptions.ZoomStart + (e.options.zoomOptions.ZoomEnd-e.options.zoomOptions.ZoomStart)*eased

	// 计算平移量
	panX, panY := e.calculatePan(sprite.Width(), sprite.Height(), eased)

	// 应用到 sprite
	sprite.SetScale(scale)
	sprite.SetX(panX)
	sprite.SetY(panY)
	sprite.SetAlpha(1.0)
}

// calculatePan 根据方向计算平移量
func (e *KenBurnsEffect) calculatePan(width, height, progress float32) (float32, float32) {
	maxPanX := width * e.options.panOptions.PanIntensity
	maxPanY := height * e.options.panOptions.PanIntensity

	vec := kenBurnsDirections[e.options.direction]
	dx, dy := vec[0], vec[1]

	panX := dx * maxPanX * progress
	panY := dy * maxPanY * progress
	return panX, panY
}

// FadeEffect 淡入淡出特效
type FadeEffect struct {
	*Effect
}

var _ IEffect = (*FadeEffect)(nil)

// NewFadeEffect 创建淡入淡出特效
func NewFadeEffect(inOut EffectInOut, opts ...optionFn) IEffect {
	var options = effectOptions{}
	return &FadeEffect{
		Effect: newEffect(inOut, toOptions(options, opts...)),
	}
}

func (e *FadeEffect) Apply(sprite render.ISpriteOperator, progress float32) {
	alpha := e.getEaseProgress(progress)
	sprite.SetAlpha(alpha)
}

// RotateEffect 旋转特效
type RotateEffect struct {
	*Effect
}

var _ IEffect = (*RotateEffect)(nil)

// NewRotateEffect 创建旋转特效
func NewRotateEffect(inOut EffectInOut, opts ...optionFn) IEffect {
	var options = effectOptions{
		direction: ti.DirectionRight,
		rotateOptions: rotateOptions{
			AngleStart: 0,
			AngleEnd:   360,
			ScaleStart: 0.5,
			ScaleEnd:   1,
		},
	}
	return &RotateEffect{
		Effect: newEffect(inOut, toOptions(options, opts...)),
	}
}

func (e *RotateEffect) Apply(sprite render.ISpriteOperator, progress float32) {
	eased := e.getEaseProgress(progress)

	// 计算旋转角度
	angleDeg := e.options.rotateOptions.AngleStart + (e.options.rotateOptions.AngleEnd-e.options.rotateOptions.AngleStart)*eased
	angleRad := misc.Deg2Rad(angleDeg) // 转换为弧度

	// 计算缩放
	scale := e.options.rotateOptions.ScaleStart + (e.options.rotateOptions.ScaleEnd-e.options.rotateOptions.ScaleStart)*eased

	sprite.SetRotation(angleRad)
	sprite.SetScale(scale)

	// 当缩放过小时隐藏，避免透视问题
	if scale < 0.01 {
		sprite.SetAlpha(0)
	} else {
		sprite.SetAlpha(1.0)
	}
}

// SlideEffect 移动特效
type SlideEffect struct {
	*Effect
}

var _ IEffect = (*SlideEffect)(nil)

// NewSlideEffect 创建移动特效
func NewSlideEffect(inOut EffectInOut, opts ...optionFn) IEffect {
	var options = effectOptions{
		direction: ti.DirectionRight,
	}
	return &SlideEffect{
		Effect: newEffect(inOut, toOptions(options, opts...)),
	}
}

func (e *SlideEffect) Apply(sprite render.ISpriteOperator, progress float32) {
	eased := e.getEaseProgress(progress)

	h := sprite.Height()
	w := sprite.Width()

	if e.direction == EffectOut {
		eased = 1.0 - eased
	}

	// 计算偏移量
	var offsetX, offsetY float32
	switch e.options.direction {
	case ti.DirectionTop:
		// 从上往下
		offsetX = 0
		offsetY = -h * (1 - eased)
	case ti.DirectionBottom:
		// 从下往上
		offsetX = 0
		offsetY = h * (1 - eased)
	case ti.DirectionLeft:
		// 从左往右
		offsetX = -w * (1 - eased)
		offsetY = 0
	case ti.DirectionRight:
		// 从右往左
		offsetX = w * (1 - eased)
		offsetY = 0
	default:
		offsetX = 0
		offsetY = 0
	}

	sprite.MoveTo(offsetX, offsetY)
	sprite.SetAlpha(1.0)
}

// ZoomEffect 缩放特效
type ZoomEffect struct {
	*Effect
}

var _ IEffect = (*ZoomEffect)(nil)

// NewZoomEffect 创建缩放特效
func NewZoomEffect(inOut EffectInOut, opts ...optionFn) IEffect {
	var options = effectOptions{
		zoomOptions: zoomOptions{
			ZoomStart: 0.5,
			ZoomEnd:   1,
		},
	}
	return &ZoomEffect{
		Effect: newEffect(inOut, toOptions(options, opts...)),
	}
}
func (e *ZoomEffect) Apply(sprite render.ISpriteOperator, progress float32) {
	eased := e.getEaseProgress(progress)

	if e.direction == EffectOut {
		eased = 1.0 - eased
	}

	// 计算缩放比例
	scale := e.options.zoomOptions.ZoomStart + (e.options.zoomOptions.ZoomEnd-e.options.zoomOptions.ZoomStart)*eased

	sprite.SetScale(scale)

	// 当缩放过小时隐藏
	if scale < 0.01 {
		sprite.SetAlpha(0)
	} else {
		sprite.SetAlpha(1.0)
	}
}

// WipeEffect Mask擦除效果
type WipeEffect struct {
	*Effect
	t float32
}

var _ IEffect = (*WipeEffect)(nil)

func NewWipeEffect(inOut EffectInOut, opts ...optionFn) IEffect {

	var options = effectOptions{
		direction: ti.DirectionRight,
		wipeOptions: wipeOptions{
			ShapeType: ti.ShapeTypeRectangle,
		},
	}
	return &WipeEffect{
		Effect: newEffect(inOut, toOptions(options, opts...)),
	}
}

func (e *WipeEffect) Apply(sprite render.ISpriteOperator, progress float32) {
	shapeSprite, ok := sprite.(*render.ShapeSprite)
	if !ok {
		return
	}

	t := e.getEaseProgress(progress)

	// 计算 tVal: 0.0-2.0 范围，1.0 表示填充整个屏幕
	// t=0 时完全隐藏，t=1 时完全显示
	tVal := t * 2.0

	// 绘制形状（ShapeSprite 同时也是 Mask）
	shapeSprite.DrawShape(e.options.wipeOptions.ShapeType, tVal)
}

func buildPanFn(direction ti.Direction) EffectFn {
	return func(inOut EffectInOut, opts ...optionFn) IEffect {
		opts = append(opts, WithDirection(direction))
		return NewKenBurnsEffect(inOut, opts...)
	}
}

func buildWipeFn(shapeType ti.ShapeType) EffectFn {
	return func(inOut EffectInOut, opts ...optionFn) IEffect {
		opts = append(opts, WithShapeType(shapeType))
		return NewWipeEffect(inOut, opts...)
	}
}

var transitionEffects = map[string]EffectFn{
	"pan_left":         buildPanFn(ti.DirectionLeft),
	"pan_right":        buildPanFn(ti.DirectionRight),
	"pan_top":          buildPanFn(ti.DirectionTop),
	"pan_bottom":       buildPanFn(ti.DirectionBottom),
	"pan_top_left":     buildPanFn(ti.DirectionTopLeft),
	"pan_top_right":    buildPanFn(ti.DirectionTopRight),
	"pan_bottom_left":  buildPanFn(ti.DirectionBottomLeft),
	"pan_bottom_right": buildPanFn(ti.DirectionBottomRight),
	"pan_center":       buildPanFn(ti.DirectionCenter),

	"rotate": NewRotateEffect,
	"slide":  NewSlideEffect,
	"zoom":   NewZoomEffect,
	"wipe":   NewWipeEffect,
	"fade":   NewFadeEffect,

	"heart":     buildWipeFn(ti.ShapeTypeHeart),
	"star5":     buildWipeFn(ti.ShapeTypeStar5),
	"cross":     buildWipeFn(ti.ShapeTypeCross),
	"linear":    buildWipeFn(ti.ShapeTypeLinear),
	"circle":    buildWipeFn(ti.ShapeTypeCircle),
	"diamond":   buildWipeFn(ti.ShapeTypeDiamond),
	"rectangle": buildWipeFn(ti.ShapeTypeRectangle),
	"triangle":  buildWipeFn(ti.ShapeTypeTriangle),
}

func GetTransitionEffect(name string) (EffectFn, error) {
	fn, ok := transitionEffects[name]
	if !ok {
		fn = transitionEffects["fade"]
	}
	return fn, nil
}
