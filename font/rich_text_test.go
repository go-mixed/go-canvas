package font

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-mixed/go-canvas/ti"
)

func TestRenderTextWithAlign(t *testing.T) {
	fs, family := mustFontLibraryForRenderTests(t)

	cases := []struct {
		name  string
		align ti.Align
		text  string
	}{
		{
			name:  "left-top single line",
			align: ti.Align{HAlign: ti.HAlignLeft, VAlign: ti.VAlignTop},
			text:  `<text font-size="24" color="#FF0000FF">Hello 世界</text>`,
		},
		{
			name:  "center-middle multiline",
			align: ti.Align{HAlign: ti.HAlignCenter, VAlign: ti.VAlignMiddle},
			text:  `<text font-size="16" color="#FF0000FF">第一行</text>` + "\n" + `<text font-size="20" color="#00FF00FF">Second line</text>`,
		},
		{
			name:  "right-bottom mixed size",
			align: ti.Align{HAlign: ti.HAlignRight, VAlign: ti.VAlignBottom},
			text:  `<text font-size="12" color="#FF0000FF">S</text><text font-size="30" color="#0000FFFF">M</text>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := RTOpt().
				SetAlign(tc.align.HAlign, tc.align.VAlign).
				SetFontFamily(family)
			rt := BuildRichTextLines(fs, opts)
			rt.SetText(tc.text)

			img := rt.RenderText()
			if img == nil {
				t.Fatalf("RenderText returned nil")
			}
			if img.Bounds().Dx() <= 0 || img.Bounds().Dy() <= 0 {
				t.Fatalf("RenderText returned empty image: %v", img.Bounds())
			}
		})
	}
}

func TestRenderTextComplexLayoutCase(t *testing.T) {
	fs, defaultFamily := mustFontLibraryForRenderTests(t)

	type richRun struct {
		Text      string
		SizePt    int
		Color     color.RGBA
		Bold      bool
		Italic    bool
		Underline bool
		Families  []string
	}

	runs := []richRun{
		{Text: "富文本 + 自动降级（字体 fallback）", SizePt: 34, Color: color.RGBA{22, 22, 22, 255}, Bold: true, Families: []string{"Microsoft YaHei", "Noto Sans SC", "DengXian"}},
		{Text: " English bold italic ", SizePt: 30, Color: color.RGBA{20, 95, 190, 255}, Bold: true, Italic: true, Families: []string{"Segoe UI", "Verdana"}},
		{Text: "中文混排 ", SizePt: 40, Color: color.RGBA{190, 45, 45, 255}, Bold: true, Underline: true, Italic: true, Families: []string{"Microsoft YaHei", "Noto Sans SC"}},
		{Text: "日本語テスト ", SizePt: 34, Color: color.RGBA{10, 126, 88, 255}, Families: []string{"Yu Gothic UI", "Meiryo"}},
		{Text: "한국어 테스트 ", SizePt: 34, Color: color.RGBA{120, 68, 170, 255}, Families: []string{"Malgun Gothic", "Gulim"}},
		{Text: "ไทยทดลองการตัดคำ和สระวรรณยุกต์ ", SizePt: 36, Color: color.RGBA{30, 30, 30, 255}, Families: []string{"Leelawadee UI", "Segoe UI"}},
		{Text: " emoji  ", SizePt: 32, Color: color.RGBA{30, 30, 30, 255}, Families: []string{"Segoe UI Emoji", "Segoe UI"}},
		{Text: "长单词降级演示: Pneumonoultramicroscopicsilicovolcanoconiosis_without_space_to_force_break", SizePt: 28, Color: color.RGBA{40, 40, 40, 255}, Families: []string{"Segoe UI", "Verdana"}},
	}

	var b strings.Builder
	for _, run := range runs {
		b.WriteString(buildRunTag(run, pickFamily(fs, run.Families, defaultFamily)))
	}

	opts := RTOpt().
		SetAlign(ti.HAlignLeft, ti.VAlignBottom).
		SetFontFamily(defaultFamily).
		SetWrapAlgorithm(WrapAlgorithmFirstFit).
		SetWidth(980)
	rt := BuildRichTextLines(fs, opts)
	tSet := time.Now()
	rt.SetText(b.String())
	setElapsed := time.Since(tSet)

	tRender := time.Now()
	img := rt.RenderText()
	renderElapsed := time.Since(tRender)
	if img == nil || img.Bounds().Dx() <= 0 || img.Bounds().Dy() <= 0 {
		t.Fatalf("RenderText returned empty image: %v", img.Bounds())
	}
	if rt.Len() < 2 {
		t.Fatalf("expected wrapped multiline text, got lines=%d", rt.Len())
	}

	var foundItalic, foundUnderline bool
	for _, seg := range rt.GetSegments() {
		foundItalic = foundItalic || seg.Italic
		foundUnderline = foundUnderline || seg.Underline
	}
	if !foundItalic {
		t.Fatalf("expected at least one italic segment")
	}
	if !foundUnderline {
		t.Fatalf("expected at least one underline segment")
	}

	tm := rt.Timing()
	t.Logf(
		"timing set=%s render=%s parse=%s layout=%s wrap=%s measure=%s",
		setElapsed, renderElapsed, tm.Parse, tm.Layout, tm.Wrap, tm.Measure,
	)

	if err := savePNG(filepath.Join("test_output", "rich_text_complex_layout.png"), img); err != nil {
		t.Fatalf("save png failed: %v", err)
	}
	t.Logf("saved png: %s", filepath.Join("test_output", "rich_text_complex_layout.png"))
}

func mustFontLibraryForRenderTests(t *testing.T) (*FontLibrary, string) {
	t.Helper()

	fs := NewFontLibrary()
	if len(fs.fonts) == 0 {
		t.Skip("no system fonts found")
	}
	for family := range fs.fonts {
		if family != "" {
			return fs, family
		}
	}
	t.Skip("no usable font family found")
	return nil, ""
}

func pickFamily(fs *FontLibrary, candidates []string, fallback string) string {
	for _, f := range candidates {
		for family := range fs.fonts {
			if strings.EqualFold(family, f) {
				return family
			}
		}
	}
	if len(candidates) > 0 {
		return candidates[0]
	}
	return fallback
}

func buildRunTag(run struct {
	Text      string
	SizePt    int
	Color     color.RGBA
	Bold      bool
	Italic    bool
	Underline bool
	Families  []string
}, family string) string {
	attrs := []string{
		fmt.Sprintf(`font-size="%d"`, run.SizePt),
		fmt.Sprintf(`color="#%02X%02X%02X%02X"`, run.Color.R, run.Color.G, run.Color.B, run.Color.A),
		fmt.Sprintf(`font-family="%s"`, family),
	}
	if run.Bold {
		attrs = append(attrs, "bold")
	}
	if run.Italic {
		attrs = append(attrs, "italic")
	}
	if run.Underline {
		attrs = append(attrs, "underline")
	}
	return "<text " + strings.Join(attrs, " ") + ">" + run.Text + "</text>"
}

func savePNG(path string, img image.Image) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
