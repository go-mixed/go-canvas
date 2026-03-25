package effect

import (
	"slideshow/render"
	"slideshow/ti"
)

// kenBurnsDirections Ken Burns 效果方向映射
// 与 Python 版本保持一致，对角线方向未归一化
var kenBurnsDirections = map[ti.Direction][2]float32{
	ti.DirectionTop:      {0, -1},
	ti.DirectionBottom:   {0, 1},
	ti.DirectionLeft:     {-1, 0},
	ti.DirectionRight:    {1, 0},
	ti.DirectionTopLeft:  {-1, -1},
	ti.DirectionTopRight: {1, -1},
	ti.DirectionBotLeft:  {-1, 1},
	ti.DirectionBotRight: {1, 1},
	ti.DirectionCenter:   {0, 0},
}

// KenBurnsEffect Ken Burns 效果
// 通过缩放和平移创建经典的 Ken Burns 幻灯片效果
type KenBurnsEffect struct {
	*Effect
	options *panOptions
}

var _ IEffect = (*KenBurnsEffect)(nil)

// NewKenBurnsEffect 创建 Ken Burns 效果
//
// Args:
//
//	durationMs: 持续时间（毫秒）
//	options: 平移选项
func NewKenBurnsEffect(durationMs int32, options *panOptions) *KenBurnsEffect {
	return &KenBurnsEffect{
		Effect:  NewEffect(durationMs, options.Easing),
		options: options,
	}
}

// Apply 应用 Ken Burns 效果
//
// 实现 IEffect 接口
func (e *KenBurnsEffect) Apply(sprite render.ISprite, progress float32) {
	eased := e.getEaseProgress(progress)

	// 计算缩放
	scale := e.options.ZoomRangeMin + (e.options.ZoomRangeMax-e.options.ZoomRangeMin)*eased

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
	maxPanX := width * e.options.PanIntensity
	maxPanY := height * e.options.PanIntensity

	vec := kenBurnsDirections[e.options.Direction]
	dx, dy := vec[0], vec[1]

	panX := dx * maxPanX * progress
	panY := dy * maxPanY * progress
	return panX, panY
}
