package render

import (
	"math"
	"slideshow/types"

	"github.com/go-mixed/go-taichi/taichi"
)

type Sprite struct {
	types.Rect[float32]
	Scale    float32 // 1.0 for no scaling
	Rotation float32 // 0.0 for no rotation, in radians

	texture *taichi.NdArray
}

type ISprite interface {
	CenterX() float32
	CenterY() float32
	Release()
}

func (s *Sprite) CenterX() float32 {
	panic("not implemented")
}

func (s *Sprite) CenterY() float32 {
	panic("not implemented")
}

func (s Screen) FillColor(color uint32) {

}

func (s *Sprite) Release() {
	if s.texture != nil {
		s.texture.Release()
		s.texture = nil
	}
}

// BoundingBox 计算旋转+缩放后的精灵包围盒
//   - parentWidth, parentHeight: 父级显示区域尺寸
//
// 返回：包围盒在屏幕坐标系中的范围
func (s *Sprite) BoundingBox(
	parentWidth, parentHeight float32,
) types.Rect[float32] {
	cx, cy := s.CenterX(), s.CenterY()
	// 精灵中心在屏幕上的位置
	centerX, centerY := s.X+cx, s.Y+cy

	// 纹理四个角点（相对于纹理中心）
	corners := [][2]float32{
		{0 - cx, 0 - cy},              // 左上
		{s.Width - cx, 0 - cy},        // 右上
		{s.Width - cx, s.Height - cy}, // 右下
		{0 - cx, s.Height - cy},       // 左下
	}

	// 旋转矩阵
	cosR := float32(math.Cos(float64(s.Rotation)))
	sinR := float32(math.Sin(float64(s.Rotation)))

	// 变换所有角点到屏幕坐标系
	var minX, maxX, minY, maxY float32
	minX, maxX = math.MaxFloat32, -math.MaxFloat32
	minY, maxY = math.MaxFloat32, -math.MaxFloat32

	for _, corner := range corners {
		dx, dy := corner[0], corner[1]

		// 先缩放，再旋转，最后平移
		scaledDx := dx * s.Scale
		scaledDy := dy * s.Scale
		rx := scaledDx*cosR - scaledDy*sinR + centerX
		ry := scaledDx*sinR + scaledDy*cosR + centerY

		// 更新边界
		if rx < minX {
			minX = rx
		}
		if rx > maxX {
			maxX = rx
		}
		if ry < minY {
			minY = ry
		}
		if ry > maxY {
			maxY = ry
		}
	}

	// 计算包围盒并与父级区域求交集
	bbox := types.Rect[float32]{
		Position: types.Position[float32]{
			X: max(0, minX),
			Y: max(0, minY),
		},
		Size: types.Size[float32]{
			Width:  min(parentWidth-1, maxX),
			Height: min(parentHeight-1, maxY),
		},
	}

	return bbox
}

type ImageSprite struct {
	Sprite
}

var _ ISprite = (*ImageSprite)(nil)

func (s *ImageSprite) CenterX() float32 {
	return s.Width / 2.
}

func (s *ImageSprite) CenterY() float32 {
	return s.Height / 2.
}
