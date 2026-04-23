package ctypes

const Sqrt2Inv = 0.7071067811865476 // 1 / √2

type DirtyMode uint32

const DirtyModeNone DirtyMode = 0
const (
	// DirtyModeLayout 影响 CSS 盒模型布局（margin/border/padding/尺寸等布局输入变化）。
	DirtyModeLayout DirtyMode = 1 << iota
	// DirtyModePainting 本节点内容像素变化（text/image/crop/resize、borderColor/borderRadius 等），需要重绘/重算绘制。
	DirtyModePainting
	// DirtyModeComposite 仅合成参数变化（affine/rotate/alpha/scroll），在 Container 合成阶段生效。
	DirtyModeComposite
	// DirtyModeChildren 子节点增删改导致父容器需要重新合成。
	DirtyModeChildren
	// DirtyModeMask 遮罩增删改或遮罩参数变化导致需要重做 mask 流程。
	DirtyModeMask
)

// Backward compatible aliases.
const (
	DirtyModeCanvas    = DirtyModePainting
	DirtyModePaint     = DirtyModePainting
	DirtyModeTransform = DirtyModeComposite
)
