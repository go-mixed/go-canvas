package render

import (
	"image/color"
	"math"
	"slideshow/misc"
	"slideshow/ti"
)

type Sprite struct {
	rect     ti.Rectangle[float32]
	scale    float32 // 1.0 for no scaling
	rotation float32 // 0.0 for no rotation, in radians
	alpha    float32 // 0.0 for no alpha, 1.0 for full alpha

	texture  *ti.TiImage
	renderer *Renderer
}

type ISpriteOperator interface {
	X() float32
	Y() float32
	Width() float32
	Height() float32
	Scale() float32
	Rotation() float32
	Alpha() float32
	// CenterX 实际的中心点x
	CenterX() float32
	// CenterY 实际的中心点y
	CenterY() float32
	Texture() *ti.TiImage

	SetX(x float32) ISprite
	SetY(y float32) ISprite
	SetScale(scale float32) ISprite
	// 缩放到指定尺寸
	SetScaleTo(width, height float32) ISprite
	SetRotation(rotation float32) ISprite
	SetAlpha(alpha float32) ISprite
}

type ISprite interface {
	ISpriteOperator

	// 替换纹理（之前的纹理会释放）
	SetTexture(texture *ti.TiImage) ISprite
	// 所有像素填充同一个颜色
	FillTexture(rgba color.Color)
	// 通过父级的宽高获取包围盒
	BoundingBox(parentWidth, parentHeight float32) ti.Rectangle[float32]
	// ResizeTo 重置尺寸，会替换纹理
	ResizeTo(width, height uint32) ISprite

	// 获取渲染器
	Renderer() *Renderer

	// Release 释放资源（必须调用，不然GPU显存泄漏）
	Release()
}

var _ ISprite = (*Sprite)(nil)

// NewSprite 创建精灵，需要传入纹理
func NewSprite(renderer *Renderer, texture *ti.TiImage) ISprite {
	shape := texture.Shape()
	return &Sprite{
		rect:     ti.Rect(0, 0, float32(shape[0]), float32(shape[1])),
		alpha:    1.0,
		texture:  texture,
		scale:    1.0,
		renderer: renderer,
	}
}

// NewBlockSprite 创建纯色块精灵，颜色为透明的
func NewBlockSprite(renderer *Renderer, width, height uint32) (ISprite, error) {
	texture, err := ti.NewTiImage(renderer.Runtime(), width, height)
	if err != nil {
		return nil, err
	}

	return NewSprite(renderer, texture), nil
}

func (s *Sprite) SetX(x float32) ISprite {
	s.rect = s.rect.MoveTo(x, s.rect.Y())
	return s
}

func (s *Sprite) X() float32 {
	return s.rect.X()
}

func (s *Sprite) SetY(y float32) ISprite {
	s.rect = s.rect.MoveTo(s.rect.X(), y)
	return s
}

func (s *Sprite) Y() float32 {
	return s.rect.Y()
}

func (s *Sprite) Width() float32 {
	return s.rect.Dx()
}

func (s *Sprite) Height() float32 {
	return s.rect.Dy()
}

func (s *Sprite) SetScale(scale float32) ISprite {
	s.scale = scale
	return s
}

func (s *Sprite) Scale() float32 {
	return s.scale
}

func (s *Sprite) SetRotation(rotation float32) ISprite {
	s.rotation = rotation
	return s
}

func (s *Sprite) Rotation() float32 {
	return s.rotation
}

func (s *Sprite) SetAlpha(alpha float32) ISprite {
	s.alpha = alpha
	return s
}

func (s *Sprite) Alpha() float32 {
	return s.alpha
}

func (s *Sprite) CenterX() float32 {
	return s.rect.Center().X
}

func (s *Sprite) CenterY() float32 {
	return s.rect.Center().Y
}

// SetTexture 设置精灵纹理
func (s *Sprite) SetTexture(texture *ti.TiImage) ISprite {
	if s.texture != nil {
		s.texture.Release()
	}
	s.texture = texture
	return s
}

// FillTexture 填充纯色
func (s *Sprite) FillTexture(color color.Color) {
	s.renderer.Module().FillTexture(s.texture, color)
}

func (s *Sprite) Texture() *ti.TiImage {
	return s.texture
}
func (s *Sprite) Release() {
	if s.texture != nil {
		s.texture.Release()
		s.texture = nil
	}
}

// SetScaleTo 缩放到指定尺寸
func (s *Sprite) SetScaleTo(width, height float32) ISprite {
	s.SetScale(min(width/s.Width(), height/s.Height()))
	return s
}

func (s *Sprite) ResizeTo(width, height uint32) ISprite {
	if misc.NumberEqual(s.Width(), width, misc.Epsilon) && misc.NumberEqual(s.Height(), height, misc.Epsilon) {
		return s
	}
	s.rect = ti.Rect[float32](0, 0, float32(width), float32(height))
	newTexture, err := ti.NewTiImage(s.renderer.Runtime(), width, height)
	if err != nil {
		return s
	}

	if s.texture == nil {
		return s
	}

	s.renderer.Module().Resize(s.texture, newTexture, ti.ResizeOptions{
		FillMode:  ti.FillModeFit,
		ScaleMode: ti.ScaleModeLanczos,
	})
	s.SetTexture(newTexture)
	return s
}

// BoundingBox 计算旋转+缩放后的精灵包围盒
//   - parentWidth, parentHeight: 父级显示区域尺寸
//     请查看examples/bounding_box_visualize.png以了解原理
//
// 返回：包围盒在屏幕坐标系中的范围
func (s *Sprite) BoundingBox(parentWidth, parentHeight float32) ti.Rectangle[float32] {
	cx, cy := s.CenterX(), s.CenterY()
	// 精灵中心在屏幕上的位置
	centerX, centerY := s.X()+cx, s.Y()+cy

	// 纹理四个角点（相对于纹理左上角）
	// 假设纹理左上角为(0,0)，中心点为(cx, cy)：
	// (0, 0)           ●────────● (width, 0)
	//                  │        │
	//                  │   ●    │  ← 中心点 (cx, cy)
	//                  │        │
	// (0, height)      ●────────● (width, height)
	//
	// corners存储角点相对于中心点的偏移量：
	// 左上: (-cx, -cy)，右上: (width-cx, -cy)
	// 右下: (width-cx, height-cy)，左下: (-cx, height-cy)
	corners := [][2]float32{
		{0 - cx, 0 - cy},                  // 左上
		{s.Width() - cx, 0 - cy},          // 右上
		{s.Width() - cx, s.Height() - cy}, // 右下
		{0 - cx, s.Height() - cy},         // 左下
	}

	// 旋转矩阵
	cosR := float32(math.Cos(float64(s.rotation)))
	sinR := float32(math.Sin(float64(s.rotation)))

	// 变换所有角点到屏幕坐标系
	var minX, maxX, minY, maxY float32
	minX, maxX = math.MaxFloat32, -math.MaxFloat32
	minY, maxY = math.MaxFloat32, -math.MaxFloat32

	for _, corner := range corners {
		dx, dy := corner[0], corner[1]

		// 先缩放，再旋转，最后平移
		scaledDx := dx * s.scale
		scaledDy := dy * s.scale
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
	bbox := ti.Rect(max(0, minX), max(0, minY), min(parentWidth-1, maxX), min(parentHeight-1, maxY))

	return bbox
}

// Renderer 获取渲染器
func (s *Sprite) Renderer() *Renderer {
	return s.renderer
}
