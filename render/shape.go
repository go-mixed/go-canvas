package render

import (
	"slideshow/ti"

	"github.com/go-mixed/go-taichi/taichi"
)

type ShapeSprite struct {
	ISprite

	renderer *Renderer

	// 归一化的SDF坐标网格
	dx, dy *taichi.NdArray

	// 中心坐标（归一化 0.0-1.0）
	centerX, centerY float32
}

func NewShapeSprite(renderer *Renderer, width, height uint32, cx, cy uint32) (ISprite, error) {
	sprite, err := NewBlockSprite(renderer, width, height)
	if err != nil {
		return nil, err
	}

	dx, err := ti.NewTiGrid(renderer.Runtime(), width, height)
	if err != nil {
		return nil, err
	}

	dy, err := ti.NewTiGrid(renderer.runtime, width, height)
	if err != nil {
		return nil, err
	}

	renderer.Module().ComputeNormalizedCoords(dx, dy, float32(cx), float32(cy))

	return &ShapeSprite{
		ISprite:  sprite,
		renderer: renderer,
		dx:       dx,
		dy:       dy,
		centerX:  float32(cx),
		centerY:  float32(cy),
	}, nil
}

func (s *ShapeSprite) CenterX() float32 {
	return s.centerX
}

func (s *ShapeSprite) CenterY() float32 {
	return s.centerY
}

// DrawShape 绘制形状
// shapeType: 形状类型 (linear, circle, diamond, rectangle, triangle, star5, heart, cross)
// size: 大小参数 0.0-2.0，1.0 表示填充整个屏幕
// fns: 可选参数，如 ti.WithShapeDirection, ti.WithShapeColor
func (s *ShapeSprite) DrawShape(shapeType ti.ShapeType, tVal float32, fns ...func(option *ti.ShapeOptions)) ISprite {
	s.FillTexture(ti.ColorTransparent)
	options := &ti.ShapeOptions{
		Direction: ti.ShapeDirectionCenter,
		Color:     ti.ColorBlack,
	}

	for _, fn := range fns {
		fn(options)
	}

	data := s.Texture()
	s.renderer.Module().ComputeShape(data, s.dx, s.dy, shapeType, tVal, options.Direction, options.Color)

	return s
}
