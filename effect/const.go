package effect

// EffectInOut 转场方向
type EffectInOut uint32

const (
	EffectIn  EffectInOut = 0
	EffectOut EffectInOut = 1
)

type EffectFn func(inOut EffectInOut, opts ...optionFn) IEffect
