package font

import (
	"image"
	"image/color"
	"time"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const syntheticItalicShear = 0.26

// RenderText 渲染文字为图像（支持水平和垂直居中）
func (r *RichText) RenderText() image.Image {
	renderStart := time.Now()
	defer func() {
		r.logf("[richtext.render] elapsed=%s", time.Since(renderStart))
	}()

	var emptyImg = image.NewRGBA(image.Rect(0, 0, 0, 0))
	if r.IsEmpty() {
		return emptyImg
	}

	// 获取 RichText 的最大宽度
	maxWidth := r.Width()
	// 计算总高度
	totalHeight := r.Height()
	if r.maxHeight > 0 && totalHeight > r.maxHeight {
		totalHeight = r.maxHeight
	}
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
		if offsetY >= totalHeight {
			break
		}

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
			face := r.fontLibrary.GetFace(seg.Font, seg.FontSize)
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
				drawUnderline(img, seg, offsetX, offsetYSeg)
			}

			// 移动到下一个 segment
			offsetX += seg.Width
		}

		// 移动到下一行
		offsetY += lineHeight
	}

	return img
}

// drawSegmentText 绘制单个文本段，普通文本直绘，假斜体走专用路径。
// drawSegmentText renders one segment; normal text is direct draw, fake italic uses dedicated paths.
func (r *RichText) drawSegmentText(dst *image.RGBA, seg *TextSegment, face font.Face, src image.Image, offsetX, offsetY int) {
	// Fast path for normal text.
	if !seg.FakeItalic {
		dot := fixedP(offsetX, offsetY+seg.metrics.Ascent.Ceil())
		d := &font.Drawer{Dst: dst, Src: src, Face: face, Dot: dot}
		d.DrawString(seg.Text)
		return
	}

	// Fake italic path: use alpha-mask transform.
	r.drawSegmentFakeItalicMask(dst, seg, face, offsetX, offsetY)
}

// drawSegmentFakeItalicMask 用 Alpha mask 渲染假斜体，减少像素带宽与内存开销。
// drawSegmentFakeItalicMask renders fake italic via an alpha mask to reduce memory bandwidth.
func (r *RichText) drawSegmentFakeItalicMask(dst *image.RGBA, seg *TextSegment, face font.Face, offsetX, offsetY int) {
	segHeight := max(1, seg.Height)
	extra := syntheticItalicExtraWidth(segHeight)
	baseWidth := seg.baseWidth
	if baseWidth <= 0 {
		baseWidth = max(1, seg.Width-extra)
	}
	bufW := max(1, baseWidth+extra)
	mask := r.ensureItalicAlphaBuffer(bufW, segHeight)
	// Reuse a shared alpha buffer and clear only the active region.
	clearAlpha(mask, image.Rect(0, 0, bufW, segHeight))
	maskDrawer := &font.Drawer{
		Dst:  mask,
		Src:  image.Opaque,
		Face: face,
		Dot:  fixedP(extra, seg.metrics.Ascent.Ceil()),
	}
	maskDrawer.DrawString(seg.Text)

	drawShearedMaskNearest(dst, image.NewUniform(seg.Color), mask, offsetX, offsetY, -syntheticItalicShear)
}

// ensureItalicAlphaBuffer 确保 Alpha 斜体缓冲尺寸足够，不足则扩容复用。
// ensureItalicAlphaBuffer ensures reusable alpha italic buffer capacity.
func (r *RichText) ensureItalicAlphaBuffer(width, height int) *image.Alpha {
	if r.italicAlphaBuf == nil || r.italicAlphaBuf.Bounds().Dx() < width || r.italicAlphaBuf.Bounds().Dy() < height {
		r.italicAlphaBuf = image.NewAlpha(image.Rect(0, 0, width, height))
	}
	return r.italicAlphaBuf
}

// clearAlpha 清空 Alpha 缓冲中的指定矩形区域。
// clearAlpha clears a target rectangle in an alpha buffer.
func clearAlpha(img *image.Alpha, rect image.Rectangle) {
	rect = rect.Intersect(img.Bounds())
	if rect.Empty() {
		return
	}
	stride := img.Stride
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		row := img.Pix[y*stride+rect.Min.X : y*stride+rect.Max.X]
		for i := range row {
			row[i] = 0
		}
	}
}

// drawShearedMaskNearest 以“逐行整数位移”的方式绘制斜切 mask。
// drawShearedMaskNearest draws a sheared alpha mask by per-row integer shifting.
func drawShearedMaskNearest(dst *image.RGBA, src image.Image, mask *image.Alpha, offsetX, offsetY int, shearX float64) {
	b := mask.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 {
		return
	}
	for y := 0; y < h; y++ {
		shiftX := misc.FloorToInt(shearX * float64(y))
		dr := image.Rect(offsetX+shiftX, offsetY+y, offsetX+shiftX+w, offsetY+y+1)
		if dr.Max.Y <= dst.Bounds().Min.Y || dr.Min.Y >= dst.Bounds().Max.Y {
			continue
		}
		draw.DrawMask(dst, dr, src, image.Point{}, mask, image.Point{b.Min.X, b.Min.Y + y}, draw.Over)
	}
}

// fixedP 将整型像素坐标转换为 fixed.Point26_6。
// fixedP converts integer pixel coordinates to fixed.Point26_6.
func fixedP(x, y int) fixed.Point26_6 {
	return fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
}

// drawUnderline 绘制文本段下划线，使用 segment 预计算的 metrics 对齐基线。
// drawUnderline draws underline using precomputed segment metrics for baseline alignment.
func drawUnderline(dst *image.RGBA, seg *TextSegment, offsetX, offsetY int) {
	baseline := offsetY + seg.metrics.Ascent.Ceil()
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
