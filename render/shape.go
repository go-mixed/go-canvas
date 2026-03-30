package render

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ti"

	"github.com/go-mixed/go-taichi/taichi"
)

type ShapeSprite struct {
	*Sprite

	// 归一化的SDF坐标网格
	dx, dy *taichi.NdArray
}

func NewShapeSprite(parent IParent, attribute *ti.Attribute) (*ShapeSprite, error) {

	width, height := uint32(attribute.Width()), uint32(attribute.Height())
	texture, err := ti.NewTiImage(parent.Renderer().Runtime(), width, height)
	if err != nil {
		return nil, err
	}

	return BuildSprite[*ShapeSprite](parent, attribute, texture, func(sprite *Sprite) (*ShapeSprite, error) {
		dx, err := ti.NewTiGrid(parent.Renderer().Runtime(), width, height)
		if err != nil {
			return nil, err
		}

		dy, err := ti.NewTiGrid(parent.Renderer().Runtime(), width, height)
		if err != nil {
			dx.Release()
			return nil, err
		}

		parent.Renderer().Module().ComputeNormalizedCoords(dx, dy, float32(attribute.Cx()), float32(attribute.Cy()))

		return &ShapeSprite{
			Sprite: sprite,
			dx:     dx,
			dy:     dy,
		}, nil
	})
}

// DrawShape 绘制形状
// shapeType: 形状类型 (linear, circle, diamond, rectangle, triangle, star5, heart, cross)
// size: 大小参数 0.0-2.0，1.0 表示填充整个屏幕
// fns: 可选参数，如 ti.WithShapeDirection, ti.WithShapeColor
func (s *ShapeSprite) DrawShape(shapeType ti.ShapeType, tVal float32, fns ...func(option *ti.ShapeOptions)) ISprite {
	s.Fill(color.Transparent)
	options := &ti.ShapeOptions{
		Direction: ti.DirectionCenter,
		Color:     color.Black,
	}

	for _, fn := range fns {
		fn(options)
	}

	data := s.Texture()
	s.Renderer().Module().ComputeShape(data, s.dx, s.dy, shapeType, tVal, options.Direction, options.Color)

	return s
}
