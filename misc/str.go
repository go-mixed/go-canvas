package misc

// ContainsEmojiLikeRunes 粗略判断文本是否包含 emoji 风格序列。
// ContainsEmojiLikeRunes heuristically detects emoji-like sequences.
func ContainsEmojiLikeRunes(s string) bool {
	for _, r := range s {
		if r == '\u200d' || IsVariationSelector(r) || IsRegionalIndicator(r) || IsEmojiRune(r) {
			return true
		}
	}
	return false
}

// IsEmojiRune 判断 rune 是否落在常见 emoji Unicode 区段。
// IsEmojiRune reports whether rune is in common emoji Unicode blocks.
func IsEmojiRune(r rune) bool {
	switch {
	case r >= 0x1F300 && r <= 0x1FAFF:
		return true
	case r >= 0x2600 && r <= 0x27BF:
		return true
	case r >= 0x1F1E6 && r <= 0x1F1FF:
		return true
	default:
		return false
	}
}

// IsVariationSelector 判断 rune 是否为变体选择符。
// IsVariationSelector reports whether rune is a variation selector.
func IsVariationSelector(r rune) bool {
	return (r >= 0xFE00 && r <= 0xFE0F) || (r >= 0xE0100 && r <= 0xE01EF)
}

// IsRegionalIndicator 判断 rune 是否为区域指示符（国旗序列）。
// IsRegionalIndicator reports whether rune is a regional indicator.
func IsRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}
