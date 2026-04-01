package font

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/f64"
	"golang.org/x/image/math/fixed"
)

const syntheticItalicShear = 0.26

// RenderText 渲染文字为图像（支持水平和垂直居中）
func (r *RichText) RenderText() image.Image {
	renderStart := time.Now()
	defer func() {
		r.timing.Render = time.Since(renderStart)
	}()

	var emptyImg = image.NewRGBA(image.Rect(0, 0, 0, 0))
	if r.IsEmpty() {
		return emptyImg
	}

	// 获取 RichText 的最大宽度
	maxWidth := r.Width()
	// 计算总高度
	totalHeight := r.Height()
	if maxWidth <= 0 || totalHeight <= 0 {
		return emptyImg
	}

	// 创建图像并填充透明背景
	img := image.NewRGBA(image.Rect(0, 0, maxWidth, totalHeight))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.Transparent), image.Point{}, draw.Src)

	// 渲染每行
	offsetY := 0
	for _, segments := range r.lines.Range() {
		lineWidth := segments.Width()
		lineHeight := segments.Height()

		// 计算水平起始偏移
		offsetX := 0
		switch r.opts.align.HAlign {
		case ti.HAlignCenter:
			offsetX = (maxWidth - lineWidth) / 2
		case ti.HAlignRight:
			offsetX = maxWidth - lineWidth
		default:
			// 默认左对齐
		}

		// 获取该行最大字号的 Metrics
		maxMetrics := segments.MaxMetrics()
		maxTopPadding := (maxMetrics.Ascent - maxMetrics.Height).Ceil()

		// 渲染该行的每个 segment
		for _, seg := range segments {
			face := r.GetFace(seg.FontFamily, seg.FontSize)
			if face == nil {
				continue
			}

			segMetrics := seg.metrics
			segTopPadding := (segMetrics.Ascent - segMetrics.Height).Ceil()

			// 计算垂直偏移
			offsetYSeg := offsetY
			switch r.opts.align.VAlign {
			case ti.VAlignTop:
				// 顶部对齐：所有字号的顶部对齐
				offsetYSeg = offsetY + maxTopPadding - segTopPadding
			case ti.VAlignMiddle:
				// 垂直居中：所有字号的垂直中心对齐
				// 中心 = baseline - topPadding + height/2
				// 要让 seg 的中心与 max 的中心对齐：
				// offsetY - maxTopPadding + maxHeight/2 = offsetYSeg - segTopPadding + seg.Height/2
				// offsetYSeg = offsetY + maxTopPadding - segTopPadding + (maxHeight - seg.Height)/2
				maxHeight := (maxMetrics.Ascent + maxMetrics.Descent).Ceil()
				offsetYSeg = offsetY + maxTopPadding - segTopPadding + (maxHeight-seg.Height)/2
			case ti.VAlignBottom:
				fallthrough
			default:
				// 底部对齐：所有字号的底部对齐
				offsetYSeg = offsetY + (maxMetrics.Ascent).Ceil() - segMetrics.Ascent.Ceil()
			}

			// 绘制文字
			src := image.NewUniform(seg.Color)
			r.drawSegmentText(img, seg, face, src, offsetX, offsetYSeg)

			if seg.Underline {
				drawUnderline(img, seg, face, offsetX, offsetYSeg)
			}

			// 移动到下一个 segment
			offsetX += seg.Width
		}

		// 移动到下一行
		offsetY += lineHeight
	}

	return img
}

func (r *RichText) drawSegmentText(dst *image.RGBA, seg *TextSegment, face font.Face, src image.Image, offsetX, offsetY int) {
	d := &font.Drawer{
		Dst:  dst,
		Src:  src,
		Face: face,
		Dot:  fixedP(offsetX, offsetY+face.Metrics().Ascent.Ceil()),
	}

	transformer := draw.BiLinear
	base := r.baseTextMatrix(seg, offsetY)
	prevC := rune(-1)
	for _, c := range seg.Text {
		if prevC >= 0 {
			d.Dot.X += d.Face.Kern(prevC, c)
		}
		dr, mask, maskp, advance, ok := d.Face.Glyph(d.Dot, c)
		if !ok {
			continue
		}

		sr := dr.Sub(dr.Min)
		fx, fy := float64(dr.Min.X), float64(dr.Min.Y)
		m := base.Translate(fx, fy)
		s2d := f64.Aff3{m.XX, m.XY, m.X0, m.YX, m.YY, m.Y0}
		transformer.Transform(d.Dst, s2d, d.Src, sr, draw.Over, &draw.Options{
			SrcMask:  mask,
			SrcMaskP: maskp,
		})

		d.Dot.X += advance
		prevC = c
	}
}

func (r *RichText) baseTextMatrix(seg *TextSegment, offsetY int) misc.Matrix {
	if seg.FakeItalic {
		// Shear(-k,0) 会产生 x' = x - k*y，先补偿全局 y 带来的整体左移，
		// 再补一个斜体额外宽，避免行首裁切。
		yComp := int(math.Ceil(float64(offsetY) * syntheticItalicShear))
		comp := float64(yComp + syntheticItalicExtraWidth(seg.Height))
		return r.matrix.Shear(-syntheticItalicShear, 0).Translate(comp, 0)
	}
	return r.matrix
}

func fixedP(x, y int) fixed.Point26_6 {
	return fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
}

func drawUnderline(dst *image.RGBA, seg *TextSegment, face font.Face, offsetX, offsetY int) {
	baseline := offsetY + face.Metrics().Ascent.Ceil()
	thickness := max(1, seg.FontSize/14)
	underlineY := baseline + max(1, seg.metrics.Descent.Ceil()/3)
	underlineWidth := seg.baseWidth
	if underlineWidth <= 0 {
		underlineWidth = seg.Width
	}

	for i := 0; i < thickness; i++ {
		y := underlineY + i
		if y < 0 || y >= dst.Bounds().Dy() {
			continue
		}
		for x := offsetX; x < offsetX+underlineWidth && x < dst.Bounds().Dx(); x++ {
			if x < 0 {
				continue
			}
			dst.Set(x, y, seg.Color)
		}
	}
}
