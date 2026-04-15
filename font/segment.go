package font

import (
	"image/color"
	"math"

	"github.com/go-mixed/go-canvas/ctypes"
	"golang.org/x/image/font"
)

// TextSegment 单个文本片段
type TextSegment struct {
	Text       string
	Font       *FontInfo
	FontSize   int
	Color      color.Color
	Bold       FontWeight
	BreakLine  bool
	Italic     bool
	Underline  bool
	FakeItalic bool
	FontFamily string
	baseWidth  int
	Width      int
	Height     int
	metrics    font.Metrics
	measured   bool
}

func (t *TextSegment) CopyWithText(text string) *TextSegment {
	newSegment := *t
	newSegment.Text = text
	// 换行标记只由输入 text 决定，避免继承上一个 segment 的 BreakLine 状态。
	newSegment.BreakLine = text == "\n"

	// 文本变更后，测量相关字段必须失效，等待 wrap/measure 重新填充。
	newSegment.baseWidth = 0
	newSegment.Width = 0
	newSegment.Height = 0
	newSegment.metrics = font.Metrics{}
	newSegment.measured = false

	return &newSegment
}

func (t *TextSegment) MeasureString(face font.Face) (int, int) {
	segWidth := font.MeasureString(face, t.Text).Ceil()
	t.baseWidth = segWidth

	t.Width = segWidth
	// 使用 ascent + |descent| 作为高度，确保能完整渲染
	t.metrics = face.Metrics()
	t.Height = (t.metrics.Ascent + t.metrics.Descent).Ceil()
	if t.FakeItalic {
		t.Width += syntheticItalicExtraWidth(t.Height)
	}
	t.measured = true
	return segWidth, t.Height
}

// CanMergeAdjacent 判断当前 segment 是否可与相邻 segment 合并（样式需完全一致）。
// CanMergeAdjacent reports whether current segment can merge with an adjacent one (exact style match).
func (t *TextSegment) CanMergeAdjacent(next *TextSegment) bool {
	if t == nil || next == nil {
		return false
	}
	if t.BreakLine || next.BreakLine {
		return false
	}
	if t.Font != next.Font ||
		t.FontSize != next.FontSize ||
		t.Bold != next.Bold ||
		t.Italic != next.Italic ||
		t.FakeItalic != next.FakeItalic ||
		t.Underline != next.Underline ||
		t.FontFamily != next.FontFamily {
		return false
	}
	return ctypes.ColorEqual(t.Color, next.Color)
}

func syntheticItalicExtraWidth(height int) int {
	if height <= 0 {
		return 0
	}
	// Roughly 14-15 degrees of slant for synthetic italic.
	return int(math.Ceil(float64(height) * 0.26))
}

type TextSegments []*TextSegment

// Height 返回该行最大字号的高度（Metrics.Height）
func (s TextSegments) Height() int {
	var maxHeight int
	for _, seg := range s {
		if maxHeight < seg.Height {
			maxHeight = seg.Height
		}
	}
	return maxHeight
}

// MaxMetrics 返回该行最大字号的 Metrics
func (s TextSegments) MaxMetrics() font.Metrics {
	var maxMetrics font.Metrics
	for _, seg := range s {
		if (maxMetrics.Ascent + maxMetrics.Descent).Ceil() < seg.Height {
			maxMetrics = seg.metrics
		}
	}
	return maxMetrics
}

// Width 总长度
func (s TextSegments) Width() int {
	var sumV int
	for _, seg := range s {
		sumV += seg.Width
	}
	return sumV
}
