package ctypes

type WordWrapMode uint8

const (
	// BreakNormal 正常换行：仅在超宽时断行，优先语义/合法断点（类似 CSS word-break: normal）。
	// BreakNormal wraps only when necessary, preferring semantic/legal breakpoints.
	BreakNormal WordWrapMode = iota
	// NoWrap 不自动换行（类似 CSS white-space: nowrap）。
	// NoWrap disables auto wrapping.
	NoWrap
	// BreakAll 可在任意 cluster 边界断行（类似 CSS overflow-wrap:anywhere）。
	// BreakAll allows break at any cluster boundary.
	BreakAll
)

type WordWrapAlgorithm uint8

const (
	WrapAlgorithmSmart WordWrapAlgorithm = iota
	WrapAlgorithmFirstFit
)

// BidiDirection 定义 BiDi 段落基础方向。
// BidiDirection defines paragraph base direction for BiDi reordering.
type BidiDirection uint8

const (
	// BidiAuto 自动根据首个强方向字符决定段落方向。
	// BidiAuto auto-detects paragraph direction from strong characters.
	BidiAuto BidiDirection = iota
	// BidiLTR 强制段落基础方向为左到右。
	// BidiLTR forces paragraph base direction to left-to-right.
	BidiLTR
	// BidiRTL 强制段落基础方向为右到左。
	// BidiRTL forces paragraph base direction to right-to-left.
	BidiRTL
)
