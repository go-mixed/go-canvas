package font

import (
	"image/color"
	"testing"
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
		expected []TextSegment
	}{
		{
			name:  "simple text",
			input: "hello",
			expected: []TextSegment{
				{Text: "hello", FontSize: 16, Color: color.Black, FontFamily: "Default"},
			},
		},
		{
			name:  "text with newline preserved",
			input: "hel\nlo",
			expected: []TextSegment{
				{Text: "hel\nlo", FontSize: 16, Color: color.Black, FontFamily: "Default"},
			},
		},
		{
			name:  "text with tab preserved",
			input: "hel\tlo",
			expected: []TextSegment{
				{Text: "hel\tlo", FontSize: 16, Color: color.Black, FontFamily: "Default"},
			},
		},
		{
			name:  "text with angle brackets preserved",
			input: "a < b > c",
			expected: []TextSegment{
				{Text: "a < b > c", FontSize: 16, Color: color.Black, FontFamily: "Default"},
			},
		},
		{
			name:  "text with angle brackets inside tag content",
			input: "<text>a < b</text>",
			expected: []TextSegment{
				{Text: "a < b", FontSize: 16, Color: color.Black, FontFamily: "Default"},
			},
		},
		{
			name:  "text with text tag name as content",
			input: "before <text> after",
			expected: []TextSegment{
				{Text: "before ", FontSize: 16, Color: color.Black, FontFamily: "Default"},
				{Text: " after", FontSize: 16, Color: color.Black, FontFamily: "Default"},
			},
		},
		{
			name:  "text with space preserved",
			input: "hel lo",
			expected: []TextSegment{
				{Text: "hel lo", FontSize: 16, Color: color.Black, FontFamily: "Default"},
			},
		},
		{
			name:  "single tag",
			input: "<text color=\"#FF0000\">hello</text>",
			expected: []TextSegment{
				{Text: "hello", FontSize: 16, Color: color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, FontFamily: "Default"},
			},
		},
		{
			name:  "bold tag",
			input: "<text bold>hello</text>",
			expected: []TextSegment{
				{Text: "hello", FontSize: 16, Color: color.Black, Bold: true, FontFamily: "Default"},
			},
		},
		{
			name:  "italic tag",
			input: "<text italic>hello</text>",
			expected: []TextSegment{
				{Text: "hello", FontSize: 16, Color: color.Black, Italic: true, FontFamily: "Default"},
			},
		},
		{
			name:  "bold and italic tag",
			input: "<text bold italic>hello</text>",
			expected: []TextSegment{
				{Text: "hello", FontSize: 16, Color: color.Black, Bold: true, Italic: true, FontFamily: "Default"},
			},
		},
		{
			name:  "font-size tag",
			input: "<text font-size=\"24\">hello</text>",
			expected: []TextSegment{
				{Text: "hello", FontSize: 24, Color: color.Black, FontFamily: "Default"},
			},
		},
		{
			name:  "font-family tag",
			input: "<text font-family=\"Arial\">hello</text>",
			expected: []TextSegment{
				{Text: "hello", FontSize: 16, Color: color.Black, FontFamily: "Arial"},
			},
		},
		{
			name:  "multiple tags",
			input: "<text color=\"#FF0000\">red</text><text bold>bold</text>",
			expected: []TextSegment{
				{Text: "red", FontSize: 16, Color: color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, FontFamily: "Default"},
				{Text: "bold", FontSize: 16, Color: color.Black, Bold: true, FontFamily: "Default"},
			},
		},
		{
			name:  "nested tags",
			input: "<text color=\"#FF0000\"><text bold>red bold</text> red</text>",
			expected: []TextSegment{
				{Text: "red bold", FontSize: 16, Color: color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, Bold: true, FontFamily: "Default"},
				{Text: " red", FontSize: 16, Color: color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, FontFamily: "Default"},
			},
		},
		{
			name:  "deeply nested tags (properly nested)",
			input: "<text color=\"#0000FF\"><text bold><text italic>bold italic</text></text></text>",
			expected: []TextSegment{
				{Text: "bold italic", FontSize: 16, Color: color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}, Bold: true, Italic: true, FontFamily: "Default"},
			},
		},
		{
			name:  "multiple nested levels with different colors",
			input: "<text color=\"#FF0000\"><text bold>red bold</text><text italic color=\"#00FF00\">red italic</text></text>",
			expected: []TextSegment{
				{Text: "red bold", FontSize: 16, Color: color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, Bold: true, FontFamily: "Default"},
				{Text: "red italic", FontSize: 16, Color: color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}, Italic: true, FontFamily: "Default"},
			},
		},
		{
			name:  "nested override parent property",
			input: "<text bold color=\"#FF0000\"><text italic>override italic</text> still bold</text>",
			expected: []TextSegment{
				{Text: "override italic", FontSize: 16, Color: color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, Bold: true, Italic: true, FontFamily: "Default"},
				{Text: " still bold", FontSize: 16, Color: color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, Bold: true, FontFamily: "Default"},
			},
		},
		{
			name:  "text between nested tags",
			input: "<text>before <text bold>bold</text> after</text>",
			expected: []TextSegment{
				{Text: "before ", FontSize: 16, Color: color.Black, FontFamily: "Default"},
				{Text: "bold", FontSize: 16, Color: color.Black, Bold: true, FontFamily: "Default"},
				{Text: " after", FontSize: 16, Color: color.Black, FontFamily: "Default"},
			},
		},
		{
			name:  "escaped characters preserved",
			input: "&lt;hello&gt; &quot;world&quot;",
			expected: []TextSegment{
				{Text: "&lt;hello&gt; &quot;world&quot;", FontSize: 16, Color: color.Black, FontFamily: "Default"},
			},
		},
		{
			name:  "escaped in tag",
			input: "<text color=\"#&quot;FF0000\">hello</text>",
			expected: []TextSegment{
				{Text: "hello", FontSize: 16, Color: color.Black, FontFamily: "Default"}, // 无法解析 #"FF0000
			},
		},
		{
			name:  "text before and after tag",
			input: "before<text color=\"#00FF00\">inside</text>after",
			expected: []TextSegment{
				{Text: "before", FontSize: 16, Color: color.Black, FontFamily: "Default"},
				{Text: "inside", FontSize: 16, Color: color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}, FontFamily: "Default"},
				{Text: "after", FontSize: 16, Color: color.Black, FontFamily: "Default"},
			},
		},
		{
			name:  "full attributes",
			input: "<text bold italic color=\"#12345678\" font-size=\"20\" font-family=\"Test\">hello</text>",
			expected: []TextSegment{
				{Text: "hello", FontSize: 20, Color: color.RGBA{R: 0x12, G: 0x34, B: 0x56, A: 0x78}, Bold: true, Italic: true, FontFamily: "Test"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments := parseText(tt.input)

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
				if seg.Bold != exp.Bold {
					t.Errorf("segment[%d].Bold = %v, want %v", i, seg.Bold, exp.Bold)
				}
				if seg.Italic != exp.Italic {
					t.Errorf("segment[%d].Italic = %v, want %v", i, seg.Italic, exp.Italic)
				}
				if seg.FontFamily != exp.FontFamily {
					t.Errorf("segment[%d].FontFamily = %v, want %v", i, seg.FontFamily, exp.FontFamily)
				}
			}
		})
	}
}

func TestParseTextEmpty(t *testing.T) {
	segments := parseText("")
	if len(segments) != 0 {
		t.Errorf("parseText(\"\") returned %d segments, want 0", len(segments))
	}
}

func TestParseTextOnlyTags(t *testing.T) {
	segments := parseText("<text bold></text>")
	if len(segments) != 0 {
		t.Errorf("parseText(\"<text bold></text>\") returned %d segments, want 0", len(segments))
	}
}
