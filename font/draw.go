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
func (r *RichText) RenderText(align ti.Align) image.Image {
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
		switch align.HAlign {
		case ti.HAlignCenter:
			offsetX = (maxWidth - lineWidth) / 2
		case ti.HAlignRight:
			offsetX = maxWidth - lineWidth
		default:
			// 默认左对齐
		}

		// 渲染该行的每个 segment
		for _, seg := range segments {
			face := r.GetFace(seg.FontFamily, seg.FontSize)

			// 计算垂直偏移（小字号往下偏移以实现垂直居中）
			offsetYSeg := offsetY
			switch align.VAlign {
			case ti.VAlignMiddle:
				offsetYSeg = offsetY + (lineHeight-seg.Height)/2
			case ti.VAlignBottom:
				offsetYSeg = offsetY + lineHeight - seg.Height
			default:
				// 默认顶部对齐
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
