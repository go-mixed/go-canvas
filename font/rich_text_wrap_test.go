package font

import (
	"reflect"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

func TestWordWrap_NoWrapWhenMaxWidthZero(t *testing.T) {
	rt := newWrapTestRichText()
	in := TextSegments{
		{Text: "hello world", Font: &FontInfo{Family: "Test"}, FontFamily: "Test", FontSize: 16},
	}

	out := rt.wordWrap(in, 0, LineBreakWhenNecessary)
	lines := collectLines(out)

	want := []string{"hello world"}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines mismatch: got=%v want=%v", lines, want)
	}
}

func TestWordWrap_WrapAtSpace(t *testing.T) {
	rt := newWrapTestRichText()
	in := TextSegments{
		{Text: "hello world", Font: &FontInfo{Family: "Test"}, FontFamily: "Test", FontSize: 16},
	}

	// basicfont.Face7x13: ASCII 单字宽约 7，42 大约可容纳 "hello "
	out := rt.wordWrap(in, 42, LineBreakWhenNecessary)
	lines := collectLines(out)

	want := []string{"hello ", "world"}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines mismatch: got=%v want=%v", lines, want)
	}
}

func TestWordWrap_ForceBreakWhenNoSemanticBreakpoint(t *testing.T) {
	rt := newWrapTestRichText()
	in := TextSegments{
		{Text: "abcdefghij", Font: &FontInfo{Family: "Test"}, FontFamily: "Test", FontSize: 16},
	}

	// 21 约等于 3 个 ASCII 字符宽度
	out := rt.wordWrap(in, 21, LineBreakWhenNecessary)
	lines := collectLines(out)

	want := []string{"abc", "def", "ghi", "j"}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines mismatch: got=%v want=%v", lines, want)
	}
}

func TestWrapKinsoku_AvoidBreakAfterOpeningPunctuation(t *testing.T) {
	rt := newWrapTestRichText()
	in := TextSegments{
		{Text: "(ABCD", Font: &FontInfo{Family: "Test"}, FontFamily: "Test", FontSize: 16},
	}

	// basicfont 单字符约 7px，14 可容纳 "(A"。
	out := rt.wrapKinsoku(in, 14, LineBreakWhenNecessary)
	lines := collectLines(out)

	want := []string{"(A", "BC", "D"}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines mismatch: got=%v want=%v", lines, want)
	}
}

func TestIsLegalBreak_KinsokuAndStickyPunctuation(t *testing.T) {
	tests := []struct {
		name string
		prev string
		next string
		want bool
	}{
		{name: "no break before cjk comma", prev: "你", next: "，", want: false},
		{name: "no break before sticky ascii punctuation", prev: "A", next: ".", want: false},
		{name: "no break after opening parenthesis", prev: "(", next: "A", want: false},
		{name: "break allowed after whitespace", prev: " ", next: "A", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isLegalBreak(tt.prev, tt.next)
			if got != tt.want {
				t.Fatalf("isLegalBreak(%q,%q)=%v want=%v", tt.prev, tt.next, got, tt.want)
			}
		})
	}
}

func TestSplitGraphemeClusters_KeepZWJEmojiTogether(t *testing.T) {
	in := "A👨‍👩‍👧‍👦B"
	got := splitGraphemeClusters(in)
	want := []string{"A", "👨‍👩‍👧‍👦", "B"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("clusters mismatch: got=%v want=%v", got, want)
	}
}

func newWrapTestRichText() *RichText {
	fs := &FontLibrary{
		faceCache: map[string]font.Face{
			"Test-16": basicfont.Face7x13,
		},
	}
	return &RichText{
		fontLibrary: fs,
	}
}

func collectLines(segments TextSegments) []string {
	lines := make([]string, 0, 4)
	var cur string
	for _, seg := range segments {
		if seg.BreakLine {
			lines = append(lines, cur)
			cur = ""
			continue
		}
		cur += seg.Text
	}
	lines = append(lines, cur)
	return lines
}
