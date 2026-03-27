package font

import (
	"image/color"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// parseOptions 解析标签属性
type parseOptions struct {
	bold       bool
	italic     bool
	color      color.Color
	fontSize   int
	fontFamily string
}

func parseText(input string) TextSegments {
	var segments TextSegments
	var text strings.Builder
	currentOpts := parseOptions{
		color:      color.Black,
		fontSize:   16,
		fontFamily: "Default",
	}
	optsStack := []parseOptions{currentOpts}

	i := 0
	for i < len(input) {
		// 检查是否到达标签
		if input[i] == '<' {
			inputToEnd := input[i:]
			// 检查是否是标签：<text>、<text >（带空格）或 </text>
			isCloseTag := strings.HasPrefix(inputToEnd, "</text>")
			isOpenTag := strings.HasPrefix(inputToEnd, "<text ") || strings.HasPrefix(inputToEnd, "<text>")

			if !isOpenTag && !isCloseTag {
				// 不是标签，把 < 当作普通字符处理
				text.WriteByte(input[i])
				i++
				continue
			}

			// 先保存当前累积的文字（保留换行和制表符）
			if text.Len() > 0 {
				seg := createSegment(text.String(), currentOpts)
				segments = append(segments, seg)
				text.Reset()
			}

			// 检查是开标签还是闭标签
			if isCloseTag {
				// 闭标签，恢复之前的状态
				if len(optsStack) > 1 {
					optsStack = optsStack[:len(optsStack)-1]
					currentOpts = optsStack[len(optsStack)-1]
				}
				i += len("</text>")
				continue
			}

			// 开标签，解析属性
			endIdx := strings.Index(inputToEnd, ">")
			if endIdx == -1 {
				i++
				continue
			}

			tagContent := strings.TrimPrefix(input[i+1:i+endIdx], "text")
			tagContent = strings.TrimLeft(tagContent, " \t")
			newOpts := parseOptions{
				color:      currentOpts.color,
				fontSize:   currentOpts.fontSize,
				fontFamily: currentOpts.fontFamily,
				bold:       currentOpts.bold,
				italic:     currentOpts.italic,
			}
			parseAttributes(tagContent, &newOpts)
			optsStack = append(optsStack, newOpts)
			currentOpts = newOpts
			i += endIdx + 1
			continue

		}

		// 普通字符，累积到 text
		text.WriteByte(input[i])
		i++
	}

	// 最后剩余的文字
	if text.Len() > 0 {
		seg := createSegment(text.String(), currentOpts)
		segments = append(segments, seg)
	}

	return segments
}

// attrRegex 匹配属性，如 color="#FF0000"、font-size="24"、font-family="Arial"、bold、italic（无值属性）
var attrRegex = regexp.MustCompile(`([\w-]+)(?:=(?:"([^"]*)"|'([^']*)'|([^"\s>]+)))?`)

// parseAttributes 解析标签属性
func parseAttributes(tag string, opts *parseOptions) {
	// 只在属性存在时设置为 true，避免覆盖继承的值
	if strings.Contains(tag, "bold") {
		opts.bold = true
	}
	if strings.Contains(tag, "italic") {
		opts.italic = true
	}

	// 解析其他属性
	attrs := attrRegex.FindAllStringSubmatch(tag, -1)
	for _, match := range attrs {
		key := match[1]
		var value string
		if len(match[2]) > 0 {
			value = match[2]
		} else if len(match[3]) > 0 {
			value = match[3]
		} else if len(match[4]) > 0 {
			value = match[4]
		}

		switch key {
		case "color":
			if c, err := parseColor(value); err == nil {
				opts.color = c
			}
		case "font-size":
			if size, err := strconv.ParseInt(value, 10, 64); err == nil {
				opts.fontSize = int(size)
			}
		case "font-family":
			opts.fontFamily = value
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
func createSegment(text string, opts parseOptions) *TextSegment {
	font := TryFindFont(opts.fontFamily, opts.bold, opts.italic)
	return &TextSegment{
		Text:       text,
		Font:       font,
		FontSize:   opts.fontSize,
		Color:      opts.color,
		Bold:       opts.bold,
		Italic:     opts.italic,
		FontFamily: opts.fontFamily,
	}
}

func findLastBreak(runes []rune, offset int) int {
	for i := offset; i >= 0; i-- {
		if isBreakable(runes[i]) {
			return i
		}
	}
	return 0
}

func findNextBreak(runes []rune, offset int) int {
	for i := offset; i < len(runes); i++ {
		if isBreakable(runes[i]) {
			return i
		}
	}
	return len(runes)
}

// isBreakable 判断字符是否可断行
func isBreakable(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	}
	if r >= 0x4E00 && r <= 0x9FFF {
		return true
	}
	if r >= 0x3000 && r <= 0x303F {
		return true
	}
	if r >= 0xFF00 && r <= 0xFFEF {
		return true
	}
	return false
}
