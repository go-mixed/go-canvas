package types

type Unit interface {
	~float32 | ~int32
}

type Position[T Unit] struct {
	X, Y T
}

type Size[T Unit] struct {
	Width, Height T
}

type Rect[T Unit] struct {
	Position[T]
	Size[T]
}

type Color[T Unit] struct {
	Red, Green, Blue, Alpha T
}

type IntColor Color[int32]
type FloatColor Color[float32]

func (c IntColor) ToFloatColor() FloatColor {
	return FloatColor{
		Red:   float32(c.Red) / 255.,
		Green: float32(c.Green) / 255.,
		Blue:  float32(c.Blue) / 255.,
		Alpha: float32(c.Alpha) / 255.,
	}
}

func (c FloatColor) ToIntColor() IntColor {
	return IntColor{
		Red:   int32(c.Red * 255.),
		Green: int32(c.Green * 255.),
		Blue:  int32(c.Blue * 255.),
		Alpha: int32(c.Alpha * 255.),
	}
}

// UInt32ToColor 将 uint32 0xAARRGGBB 颜色 转换为 IntColor
func UInt32ToColor(color uint32) IntColor {
	return IntColor{
		Red:   int32((color >> 16) & 0xff),
		Green: int32((color >> 8) & 0xff),
		Blue:  int32(color & 0xff),
		Alpha: int32((color >> 24) & 0xff),
	}
}
