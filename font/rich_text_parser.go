package font

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/go-mixed/go-canvas/internel/misc"
)

// RichTextFontStyle 富文本的样式
type RichTextFontStyle struct {
	Bold       bool
	Italic     bool
	Underline  bool
	Color      color.Color
	FontSize   int
	FontFamily string
}

func (r *RichText) parseText(input string) TextSegments {
	segments := make(TextSegments, 0, 16)
	currentStyle := r.opts.fontStyle
	styleStack := []RichTextFontStyle{currentStyle}

	i := 0
	textStart := 0
	for i < len(input) {
		// 检查是否到达标签
		if input[i] == '<' {
			inputToEnd := input[i:]
			// 检查是否是标签：<text>、<text >（带空格）或 </text>
			isCloseTag := strings.HasPrefix(inputToEnd, "</text>")
			isOpenTag := strings.HasPrefix(inputToEnd, "<text ") || strings.HasPrefix(inputToEnd, "<text>")

			if !isOpenTag && !isCloseTag {
				// 不是标签，把 < 当作普通字符处理（保持在原字符串片段中）
				i++
				continue
			}

			// 先保存当前累积的文字（保留换行和制表符）
			if i > textStart {
				seg := r.createSegment(input[textStart:i], currentStyle)
				segments = append(segments, seg)
			}

			// 检查是开标签还是闭标签
			if isCloseTag {
				// 闭标签，恢复之前的状态
				if len(styleStack) > 1 {
					styleStack = styleStack[:len(styleStack)-1]
					currentStyle = styleStack[len(styleStack)-1]
				}
				i += len("</text>")
				textStart = i
				continue
			}

			// 开标签，解析属性
			endIdx := strings.IndexByte(inputToEnd, '>')
			if endIdx == -1 {
				i++
				continue
			}

			tagContent := strings.TrimPrefix(input[i+1:i+endIdx], "text")
			tagContent = strings.TrimLeft(tagContent, " \t")
			newStyle := RichTextFontStyle{
				Color:      currentStyle.Color,
				FontSize:   currentStyle.FontSize,
				FontFamily: currentStyle.FontFamily,
				Bold:       currentStyle.Bold,
				Italic:     currentStyle.Italic,
				Underline:  currentStyle.Underline,
			}
			r.parseAttributes(tagContent, &newStyle)
			styleStack = append(styleStack, newStyle)
			currentStyle = newStyle
			i += endIdx + 1
			textStart = i
			continue

		}

		i++
	}

	// 最后剩余的文字
	if textStart < len(input) {
		seg := r.createSegment(input[textStart:], currentStyle)
		segments = append(segments, seg)
	}

	return segments
}

// parseAttributes 解析标签属性
func (r *RichText) parseAttributes(tag string, opts *RichTextFontStyle) {
	i := 0
	for i < len(tag) {
		for i < len(tag) && isASCIISpace(tag[i]) {
			i++
		}
		if i >= len(tag) {
			break
		}

		ks := i
		for i < len(tag) && isAttrKeyChar(tag[i]) {
			i++
		}
		if ks == i {
			i++
			continue
		}
		key := tag[ks:i]

		for i < len(tag) && isASCIISpace(tag[i]) {
			i++
		}

		value := ""
		if i < len(tag) && tag[i] == '=' {
			i++
			for i < len(tag) && isASCIISpace(tag[i]) {
				i++
			}
			if i < len(tag) && (tag[i] == '"' || tag[i] == '\'') {
				quote := tag[i]
				i++
				vs := i
				for i < len(tag) && tag[i] != quote {
					i++
				}
				value = tag[vs:i]
				if i < len(tag) && tag[i] == quote {
					i++
				}
			} else {
				vs := i
				for i < len(tag) && !isASCIISpace(tag[i]) {
					i++
				}
				value = tag[vs:i]
			}
		}
		if len(value) >= 2 {
			first := value[0]
			last := value[len(value)-1]
			if (first == '"' || first == '\'') && first == last {
				value = value[1 : len(value)-1]
			}
		}

		switch key {
		case "bold":
			opts.Bold = value == "" || misc.ToBool(value)
		case "italic":
			opts.Italic = value == "" || misc.ToBool(value)
		case "underline":
			opts.Underline = value == "" || misc.ToBool(value)
		case "color":
			if c, err := parseColor(value); err == nil {
				opts.Color = c
			}
		case "font-size":
			if size, err := strconv.ParseInt(value, 10, 64); err == nil {
				opts.FontSize = int(size)
			}
		case "font-family":
			opts.FontFamily = normalizeFamilyName(value)
		}
	}
}

// parseColor 解析 #RRGGBBAA 格式的颜色
func parseColor(s string) (color.Color, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) == 6 {
		s += "FF" // 没有 alpha，默认 FF
	}
	if len(s) != 8 {
		return color.Black, nil
	}

	r, _ := strconv.ParseUint(s[0:2], 16, 8)
	g, _ := strconv.ParseUint(s[2:4], 16, 8)
	b, _ := strconv.ParseUint(s[4:6], 16, 8)
	a, _ := strconv.ParseUint(s[6:8], 16, 8)

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}, nil
}

// createSegment 根据当前选项创建文本片段
func (r *RichText) createSegment(text string, opts RichTextFontStyle) *TextSegment {
	weight := FontWeightRegular
	if opts.Bold {
		weight = FontWeightBold
	}
	fi := r.fontLibrary.MatchOrFeedback(opts.FontFamily, weight, opts.Italic)
	if opts.FontFamily != fi.Family {
		r.logf(
			"[richtext.fallback] req=%q got=%q bold=%d italic=%t text=%q",
			opts.FontFamily, fi.Family, weight, opts.Italic, summarizeTextForLog(text),
		)
	}
	return &TextSegment{
		Text:       text,
		Font:       fi,
		FontSize:   opts.FontSize,
		Color:      opts.Color,
		Bold:       weight,
		Italic:     opts.Italic,
		Underline:  opts.Underline,
		FakeItalic: opts.Italic && !fi.Italic,
		FontFamily: opts.FontFamily,
	}
}

func isASCIISpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func isAttrKeyChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '_'
}
