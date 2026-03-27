package font

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// TextSegment 单个文本片段
type TextSegment struct {
	Text       string
	Font       *TrueTypeFont
	FontSize   int
	Color      color.Color
	Bold       bool
	Italic     bool
	FontFamily string
	Width      int
	Height     int
	metrics    font.Metrics
}

func (t *TextSegment) CopyWithText(text string) *TextSegment {
	var newSegment TextSegment
	newSegment = *t
	newSegment.Text = text

	return &newSegment
}

// CreateFace 创建字体Face
func (t *TextSegment) CreateFace() font.Face {
	return truetype.NewFace(t.Font.Font, &truetype.Options{
		Size:    float64(t.FontSize),
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func (t *TextSegment) MeasureString(face font.Face) (int, int) {
	segWidth := font.MeasureString(face, t.Text).Ceil()
	defer face.Close()

	t.Width = segWidth
	// 使用 ascent + |descent| 作为高度，确保能完整渲染
	t.metrics = face.Metrics()
	t.Height = (t.metrics.Ascent + t.metrics.Descent).Ceil()
	return segWidth, t.Height
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

type RichText struct {
	original  string
	lines     *misc.List[TextSegments]
	faceCache map[string]font.Face

	width, height int // 缓存宽度和高度，避免重复计算
}

// BuildRichTextLines 解析带标签的文字，返回文本片段列表
// 标签格式：<text bold italic color="#RRGGBBAA" font-size="15" font-family="Noto Sans CJK SC">文字</text>
func BuildRichTextLines(input string) *RichText {

	segments := parseText(input)

	richText := &RichText{
		original:  input,
		lines:     misc.NewList[TextSegments](),
		faceCache: make(map[string]font.Face),
		width:     -1,
		height:    -1,
	}

	if len(segments) == 0 {
		return richText
	}

	var lastSegments TextSegments
	for _, seg := range segments {
		// 创建字体Face
		richText.createFaceOrNot(seg.FontFamily, seg.FontSize, seg)

		parts := strings.Split(seg.Text, "\n")
		// 没有回车
		if len(parts) == 1 {
			lastSegments = append(lastSegments, seg)
		} else { // 有回车
			// 保存第0条记录
			lastSegments = append(lastSegments, seg.CopyWithText(parts[0]))
			richText.lines.PushBack(lastSegments)
			lastSegments = nil

			// 保存从 1~len(parts)-2
			// 因为最后一条会到下一个段落
			for i := 1; i < len(parts)-1; i++ {
				lastSegments = append(lastSegments, seg.CopyWithText(parts[i]))
				richText.lines.PushBack(lastSegments)
			}
		}
	}

	// 收尾
	if len(lastSegments) > 0 {
		richText.lines.PushBack(lastSegments)
		lastSegments = nil
	}

	richText.Initial()

	return richText
}

// Len 返回文本段落的总行数
func (r *RichText) Len() int {
	return r.lines.Len()
}

// GetSegments 返回所有文本段落
func (r *RichText) GetSegments() TextSegments {
	segments := make(TextSegments, 0, r.lines.Len())
	for el := r.lines.Front(); el != nil; el = el.Next() {
		segments = append(segments, el.Value...)
	}
	return segments
}

func (r *RichText) Equal(text string) bool {
	return r.original == text
}

func (r *RichText) createFaceOrNot(family string, size int, seg *TextSegment) {
	k := family + "-" + strconv.Itoa(size)
	if _, ok := r.faceCache[k]; !ok {
		r.faceCache[k] = seg.CreateFace()
	}
}

// GetFace 根据 fontFamily 和 fontSize 获取缓存的 face
func (r *RichText) GetFace(fontFamily string, fontSize int) font.Face {
	k := fontFamily + "-" + strconv.Itoa(fontSize)
	if face, ok := r.faceCache[k]; ok {
		return face
	}
	return nil
}

// Initial 初始化，测量每个文本片段的宽度和高度
func (r *RichText) Initial() {
	for el := r.lines.Front(); el != nil; el = el.Next() {
		for i := range el.Value {
			face := r.GetFace(el.Value[i].FontFamily, int(el.Value[i].FontSize))
			el.Value[i].MeasureString(face)
		}
	}
}

func (r *RichText) Width() int {
	if r.width > 0 {
		return r.width
	}

	var maxV int
	for _, segments := range r.lines.Range() {
		w := segments.Width()
		if maxV < segments.Width() {
			maxV = w
		}
	}
	r.width = maxV
	return maxV
}

func (r *RichText) Height() int {
	if r.height > 0 {
		return r.height
	}

	var totalHeight int
	for _, segments := range r.lines.Range() {
		totalHeight += segments.Height()
	}
	r.height = totalHeight
	return r.height
}

func (r *RichText) IsEmpty() bool {
	return r.lines.Len() == 0
}
