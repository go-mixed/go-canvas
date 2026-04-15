package ctypes

type Align struct {
	HAlign HorizontalAlign
	VAlign VerticalAlign
}

// HorizontalAlign 水平对齐方式
type HorizontalAlign int

const (
	HAlignLeft HorizontalAlign = iota
	HAlignCenter
	HAlignRight
)

// VerticalAlign 垂直对齐方式
type VerticalAlign int

const (
	VAlignTop VerticalAlign = iota
	VAlignMiddle
	VAlignBottom
)
