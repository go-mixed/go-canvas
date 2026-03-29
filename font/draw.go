package font

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/go-mixed/go-canvas/ti"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// RenderText 渲染文字为图像（支持水平和垂直居中）
func (r *RichText) RenderText() image.Image {
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
			d := &font.Drawer{
				Dst:  img,
				Src:  src,
				Face: face,
				Dot: fixed.Point26_6{
					X: fixed.I(offsetX),
					Y: fixed.I(offsetYSeg + face.Metrics().Ascent.Ceil()),
				},
			}
			d.DrawString(seg.Text)

			// 移动到下一个 segment
			offsetX += seg.Width
		}

		// 移动到下一行
		offsetY += lineHeight
	}

	return img
}
