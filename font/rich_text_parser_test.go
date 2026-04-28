package font

import (
	"image/color"
	"testing"

	xfont "golang.org/x/image/font"
)

func TestParseColor(t *testing.T) {
	tests := []struct {
		input    string
		expected color.Color
	}{
		{"#FF0000FF", &color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}},
		{"#00FF00", &color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}}, // 默认 alpha
		{"#0000FFFF", &color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}},
		{"#12345678", &color.RGBA{R: 0x12, G: 0x34, B: 0x56, A: 0x78}},
	}

	for _, tt := range tests {
		result, err := parseColor(tt.input)
		if err != nil {
			t.Errorf("parseColor(%q) returned error: %v", tt.input, err)
			continue
		}
		rgba := result.(color.RGBA)
		expected := tt.expected.(*color.RGBA)
		if rgba != *expected {
			t.Errorf("parseColor(%q) = %v, want %v", tt.input, rgba, *expected)
		}
	}
}

func TestParseText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []parserExpectedSegment
	}{
		{
			name:  "simple text",
			input: "hello",
			expected: []parserExpectedSegment{
				{Text: "hello", FontSize: 16, Color: color.RGBA{A: 0xFF}, FontFamily: "Default"},
			},
		},
		{
			name:  "text with literals preserved",
			input: "a < b > c\tok",
			expected: []parserExpectedSegment{
				{Text: "a < b > c\tok", FontSize: 16, Color: color.RGBA{A: 0xFF}, FontFamily: "Default"},
			},
		},
		{
			name:  "single color tag",
			input: "<text color=\"#FF0000\">hello</text>",
			expected: []parserExpectedSegment{
				{Text: "hello", FontSize: 16, Color: color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, FontFamily: "Default"},
			},
		},
		{
			name:  "bold tag",
			input: "<text bold>hello</text>",
			expected: []parserExpectedSegment{
				{Text: "hello", FontSize: 16, Color: color.RGBA{A: 0xFF}, Bold: true, FontFamily: "Default"},
			},
		},
		{
			name:  "italic tag",
			input: "<text italic>hello</text>",
			expected: []parserExpectedSegment{
				{Text: "hello", FontSize: 16, Color: color.RGBA{A: 0xFF}, Italic: true, FontFamily: "Default"},
			},
		},
		{
			name:  "font-size tag",
			input: "<text font-size=\"24\">hello</text>",
			expected: []parserExpectedSegment{
				{Text: "hello", FontSize: 24, Color: color.RGBA{A: 0xFF}, FontFamily: "Default"},
			},
		},
		{
			name:  "font-family tag",
			input: "<text font-family=\"Arial\">hello</text>",
			expected: []parserExpectedSegment{
				{Text: "hello", FontSize: 16, Color: color.RGBA{A: 0xFF}, FontFamily: "Arial"},
			},
		},
		{
			name:  "nested tags inherit and override",
			input: "<text bold color=\"#FF0000\"><text italic>inner</text> outer</text>",
			expected: []parserExpectedSegment{
				{Text: "inner", FontSize: 16, Color: color.RGBA{R: 0xFF, A: 0xFF}, Bold: true, Italic: true, FontFamily: "Default"},
				{Text: " outer", FontSize: 16, Color: color.RGBA{R: 0xFF, A: 0xFF}, Bold: true, FontFamily: "Default"},
			},
		},
		{
			name:  "underline tag",
			input: "<text underline>hello</text>",
			expected: []parserExpectedSegment{
				{Text: "hello", FontSize: 16, Color: color.RGBA{A: 0xFF}, Underline: true, FontFamily: "Default"},
			},
		},
		{
			name:  "full attributes",
			input: "<text bold italic color=\"#12345678\" font-size=\"20\" font-family=\"Test\">hello</text>",
			expected: []parserExpectedSegment{
				{Text: "hello", FontSize: 20, Color: color.RGBA{R: 0x12, G: 0x34, B: 0x56, A: 0x78}, Bold: true, Italic: true, FontFamily: "Test"},
			},
		},
	}

	rt := newParserTestRichText()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments := rt.parseText(tt.input)

			if len(segments) != len(tt.expected) {
				t.Errorf("parseText(%q) returned %d segments, want %d", tt.input, len(segments), len(tt.expected))
				t.Errorf("Got: %v", segments)
				return
			}

			for i, seg := range segments {
				exp := tt.expected[i]
				if seg.Text != exp.Text {
					t.Errorf("segment[%d].Text = %q, want %q", i, seg.Text, exp.Text)
				}
				if seg.FontSize != exp.FontSize {
					t.Errorf("segment[%d].FontSize = %v, want %v", i, seg.FontSize, exp.FontSize)
				}
				gotBold := seg.Bold >= xfont.WeightBold
				if gotBold != exp.Bold {
					t.Errorf("segment[%d].Bold = %v, want %v", i, gotBold, exp.Bold)
				}
				if seg.Italic != exp.Italic {
					t.Errorf("segment[%d].Italic = %v, want %v", i, seg.Italic, exp.Italic)
				}
				if seg.Underline != exp.Underline {
					t.Errorf("segment[%d].Underline = %v, want %v", i, seg.Underline, exp.Underline)
				}
				if seg.FontFamily != exp.FontFamily {
					t.Errorf("segment[%d].FontFamily = %v, want %v", i, seg.FontFamily, exp.FontFamily)
				}
				gotColor := colorToRGBA(seg.Color)
				if gotColor != exp.Color {
					t.Errorf("segment[%d].Color = %v, want %v", i, gotColor, exp.Color)
				}
			}
		})
	}
}

func TestParseTextEmpty(t *testing.T) {
	rt := newParserTestRichText()
	segments := rt.parseText("")
	if len(segments) != 0 {
		t.Errorf("parseText(\"\") returned %d segments, want 0", len(segments))
	}
}

func TestParseTextOnlyTags(t *testing.T) {
	rt := newParserTestRichText()
	segments := rt.parseText("<text bold></text>")
	if len(segments) != 0 {
		t.Errorf("parseText(\"<text bold></text>\") returned %d segments, want 0", len(segments))
	}
}

type parserExpectedSegment struct {
	Text       string
	FontSize   int
	Color      color.RGBA
	Bold       bool
	Italic     bool
	Underline  bool
	FontFamily string
}

func newParserTestRichText() *RichText {
	fs := &FontLibrary{
		fonts: map[string][]*FontInfo{
			"Default": {{Family: "Default", Weight: xfont.WeightNormal, Italic: false}},
			"Arial":   {{Family: "Arial", Weight: xfont.WeightNormal, Italic: false}},
			"Test":    {{Family: "Test", Weight: xfont.WeightNormal, Italic: false}},
		},
	}
	opts := RTOpt().SetFontFamily("Default")
	return BuildRichTextLines(fs, opts)
}

func colorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}
