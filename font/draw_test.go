package font

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"

	"github.com/go-mixed/go-canvas/ti"
)

// TestVerticalAlign 测试垂直对齐（Top, Middle, Bottom）
func TestVerticalAlign(t *testing.T) {
	text := `<text font-size="12" color="#FF0000FF">中文abc123</text>` +
		`<text font-size="24" color="#00FF00FF">中文abc123</text>` +
		`<text font-size="36" color="#0000FFFF">中文abc123</text>`

	rt := BuildRichTextLines(text)

	// 3种垂直对齐（上下排列）
	alignments := []struct {
		align ti.Align
		name  string
	}{
		{ti.Align{HAlign: ti.HAlignLeft, VAlign: ti.VAlignTop}, "Top"},
		{ti.Align{HAlign: ti.HAlignLeft, VAlign: ti.VAlignMiddle}, "Middle"},
		{ti.Align{HAlign: ti.HAlignLeft, VAlign: ti.VAlignBottom}, "Bottom"},
	}

	const cellWidth = 400
	const cellHeight = 100
	const rows = 3
	const padding = 20

	imgWidth := cellWidth + 2*padding
	imgHeight := rows*cellHeight + (rows+1)*padding

	gridImg := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(gridImg, gridImg.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	for i, tc := range alignments {
		row := i

		img := rt.RenderText(tc.align)
		if img.Bounds().Empty() {
			t.Errorf("RenderText returned empty image for %s", tc.name)
			continue
		}

		x := padding
		y := row*cellHeight + padding

		draw.Draw(gridImg, image.Rect(x, y, x+img.Bounds().Dx(), y+img.Bounds().Dy()), img, image.Point{}, draw.Over)
	}

	f, err := os.Create("test_vertical_align.png")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, gridImg); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	t.Logf("Test image saved to test_vertical_align.png")
}

// TestHorizontalAlign 测试水平对齐（Left, Center, Right）
func TestHorizontalAlign(t *testing.T) {
	text := `<text font-size="16" color="#FF0000FF">文字内容测试</text>`

	rt := BuildRichTextLines(text)

	// 3种水平对齐（上下排列）
	alignments := []struct {
		align ti.Align
		name  string
	}{
		{ti.Align{HAlign: ti.HAlignLeft, VAlign: ti.VAlignTop}, "Left"},
		{ti.Align{HAlign: ti.HAlignCenter, VAlign: ti.VAlignTop}, "Center"},
		{ti.Align{HAlign: ti.HAlignRight, VAlign: ti.VAlignTop}, "Right"},
	}

	const cellWidth = 300
	const cellHeight = 80
	const rows = 3
	const padding = 20

	imgWidth := cellWidth + 2*padding
	imgHeight := rows*cellHeight + (rows+1)*padding

	gridImg := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(gridImg, gridImg.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	for i, tc := range alignments {
		row := i

		img := rt.RenderText(tc.align)
		if img.Bounds().Empty() {
			t.Errorf("RenderText returned empty image for %s", tc.name)
			continue
		}

		x := padding
		y := row*cellHeight + padding

		draw.Draw(gridImg, image.Rect(x, y, x+img.Bounds().Dx(), y+img.Bounds().Dy()), img, image.Point{}, draw.Over)
	}

	f, err := os.Create("test_horizontal_align.png")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, gridImg); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	t.Logf("Test image saved to test_horizontal_align.png")
}

// TestRenderTextSingleLine 测试单行多字体混排
func TestRenderTextSingleLine(t *testing.T) {
	text := `<text font-size="12" color="#FF0000FF">小字12号红色</text>` +
		`<text font-size="24" color="#00FF00FF">中字24号绿色</text>` +
		`<text font-size="36" color="#0000FFFF">大字36号蓝色</text>`

	rt := BuildRichTextLines(text)

	align := ti.Align{HAlign: ti.HAlignCenter, VAlign: ti.VAlignMiddle}
	img := rt.RenderText(align)

	f, err := os.Create("test_single_line.png")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	t.Logf("Single line test image saved to test_single_line.png")
}

// TestRenderTextMultiLine 测试多行文本
func TestRenderTextMultiLine(t *testing.T) {
	text := `<text font-size="16" color="#FF0000FF">这是第一行比较长的文字测试</text>` + "\n" +
		`<text font-size="20" color="#00FF00FF">这是第二行比较长的文字测试</text>` + "\n" +
		`<text font-size="24" color="#0000FFFF">这是第三行比较长的文字测试</text>`

	rt := BuildRichTextLines(text)

	align := ti.Align{HAlign: ti.HAlignCenter, VAlign: ti.VAlignMiddle}
	img := rt.RenderText(align)

	f, err := os.Create("test_multi_line.png")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	t.Logf("Multi line test image saved to test_multi_line.png")
}

// TestMultiLineMultiSize 测试多行多字号（用于查看行距）
func TestMultiLineMultiSize(t *testing.T) {
	text := `<text font-size="12" color="#FF0000FF">小12</text>` +
		`<text font-size="24" color="#00FF00FF">中24</text>` +
		`<text font-size="36" color="#0000FFFF">大36</text>` + "\n" +
		`<text font-size="12" color="#FF0000FF">小12</text>` +
		`<text font-size="24" color="#00FF00FF">中24</text>` +
		`<text font-size="36" color="#0000FFFF">大36</text>` + "\n" +
		`<text font-size="12" color="#FF0000FF">小12</text>` +
		`<text font-size="24" color="#00FF00FF">中24</text>` +
		`<text font-size="36" color="#0000FFFF">大36</text>`

	rt := BuildRichTextLines(text)

	align := ti.Align{HAlign: ti.HAlignLeft, VAlign: ti.VAlignBottom}
	img := rt.RenderText(align)

	f, err := os.Create("test_multi_line_multi_size.png")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	t.Logf("Multi line multi size test image saved to test_multi_line_multi_size.png")
}

// TestMultiLineSameSize 测试多行相同字号不同长度
func TestMultiLineSameSize(t *testing.T) {
	text := `<text font-size="16" color="#FF0000FF">第一行较短</text>` + "\n" +
		`<text font-size="16" color="#FF0000FF">第二行中等长度文字</text>` + "\n" +
		`<text font-size="16" color="#FF0000FF">第三行非常非常长的文字内容测试</text>`

	rt := BuildRichTextLines(text)

	align := ti.Align{HAlign: ti.HAlignLeft, VAlign: ti.VAlignBottom}
	img := rt.RenderText(align)

	f, err := os.Create("test_multi_line_same_size.png")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	t.Logf("Multi line same size test image saved to test_multi_line_same_size.png")
}
