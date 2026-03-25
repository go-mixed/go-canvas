package render

import (
	"image/color"
	"math"
	"slices"
	"sync"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
)

type Sprite struct {
	rect           ti.Rectangle[float32]
	scaleX, scaleY float32 // 1.0 for no scaling
	rotation       float32 // 0.0 for no rotation, in radians
	alpha          float32 // 0.0 for no alpha, 1.0 for full alpha

	// 中心点的相对值
	deltaCx, deltaCy float32

	texture *ti.TiImage

	mask IMask

	children    []ISprite
	isContainer bool

	dirty    bool // 是否需要重新绘制
	mutex    sync.RWMutex
	renderer *Renderer
}

var _ ISprite = (*Sprite)(nil)

// NewNonContainerSprite 创建非容器的精灵，需要传入纹理
func NewNonContainerSprite(renderer *Renderer, texture *ti.TiImage) ISprite {
	shape := texture.Shape()
	w, h := shape[0], shape[1]
	return &Sprite{
		rect:    ti.Rect(0, 0, float32(w), float32(h)),
		alpha:   1.0,
		texture: texture,
		scaleX:  1.0,
		scaleY:  1.0,
		deltaCx: float32(w / 2),
		deltaCy: float32(h / 2),
		mask:    nil,
		dirty:   false,

		mutex:       sync.RWMutex{},
		isContainer: false, // 非容器
		renderer:    renderer,
	}
}

// NewBlockSprite 创建纯色块精灵，颜色为透明的
func NewBlockSprite(renderer *Renderer, width, height uint32) (ISprite, error) {
	texture, err := ti.NewTiImage(renderer.Runtime(), width, height)
	if err != nil {
		return nil, err
	}

	return NewNonContainerSprite(renderer, texture), nil
}

// NewContainerSprite 创建容器精灵，只能添加子精灵
func NewContainerSprite(renderer *Renderer, width, height uint32) (ISprite, error) {
	texture, err := ti.NewTiImage(renderer.Runtime(), width, height)
	if err != nil {
		return nil, err
	}
	return &Sprite{
		rect:    ti.Rect(0, 0, float32(width), float32(height)),
		alpha:   1.0,
		texture: texture,
		scaleX:  1.0,
		scaleY:  1.0,
		deltaCx: float32(width / 2),
		deltaCy: float32(height / 2),
		mask:    nil,
		dirty:   false,

		isContainer: true, // 容器
		renderer:    renderer,
	}, nil
}

func (s *Sprite) lock(fn func(), triggerDirty func() bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if triggerDirty() {
		s.dirty = true
	}
	fn()
}

func (s *Sprite) IsDirty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, child := range s.children {
		if child.IsDirty() {
			return true
		}
	}

	return s.dirty
}

func (s *Sprite) SetX(x float32) ISprite {
	return s.MoveTo(x, s.Y())
}

func (s *Sprite) X() float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.rect.X()
}

func (s *Sprite) SetY(y float32) ISprite {
	return s.MoveTo(s.X(), y)
}

func (s *Sprite) Y() float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.rect.Y()
}

func (s *Sprite) MoveTo(x float32, y float32) ISprite {
	s.lock(func() {
		s.rect = s.rect.MoveTo(x, y)
	}, func() bool {
		return misc.NumberEqual(s.rect.X(), x, misc.Epsilon) && misc.NumberEqual(s.rect.Y(), y, misc.Epsilon)
	})
	return s
}

func (s *Sprite) Width() float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.width()
}

func (s *Sprite) width() float32 {
	return s.rect.Dx()
}

func (s *Sprite) Height() float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.height()
}

func (s *Sprite) height() float32 {
	return s.rect.Dy()
}

func (s *Sprite) SetScale(scaleX, scaleY float32) ISprite {
	s.lock(func() {
		s.scaleX = scaleX
		s.scaleY = scaleY
	}, func() bool {
		return misc.NumberEqual(s.scaleX, scaleX, misc.Epsilon) && misc.NumberEqual(s.scaleY, scaleY, misc.Epsilon)
	})
	return s
}

func (s *Sprite) ScaleX() float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.scaleX
}

func (s *Sprite) ScaleY() float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.scaleY
}

func (s *Sprite) SetRotation(rotation float32) ISprite {
	s.lock(func() {
		s.rotation = rotation
	}, func() bool {
		return misc.NumberEqual(s.rotation, rotation, misc.Epsilon)
	})
	return s
}

func (s *Sprite) Rotation() float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.rotation
}

func (s *Sprite) SetAlpha(alpha float32) ISprite {
	s.lock(func() {
		s.alpha = alpha
	}, func() bool {
		return misc.NumberEqual(s.alpha, alpha, misc.Epsilon)
	})
	return s
}

func (s *Sprite) Alpha() float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.alpha
}

func (s *Sprite) Cx() float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.X() + s.deltaCx
}

func (s *Sprite) SetCx(cx float32) ISprite {
	s.lock(func() {
		s.deltaCx = cx
	}, func() bool {
		return misc.NumberEqual(s.deltaCx, cx, misc.Epsilon)
	})
	return s
}

func (s *Sprite) Cy() float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Y() + s.deltaCy
}

func (s *Sprite) SetCy(cy float32) ISprite {
	s.lock(func() {
		s.deltaCy = cy
	}, func() bool { return misc.NumberEqual(s.deltaCy, cy, misc.Epsilon) })
	return s
}

//// SetTexture 设置精灵纹理
//func (s *Sprite) SetTexture(texture *ti.TiImage) ISprite {
//	s.lock(func() {
//		if s.texture != nil {
//			s.texture.Release()
//		}
//		s.texture = texture
//	}, func() bool { return s.texture != texture })
//	return s
//}

// FillTexture 填充纯色
func (s *Sprite) FillTexture(color color.Color) {
	s.lock(func() {
		s.renderer.Module().FillTexture(s.texture, color)
	}, func() bool { return false })
}

func (s *Sprite) Texture() *ti.TiImage {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.texture
}

func (s *Sprite) Release() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.texture != nil {
		s.texture.Release()
		s.texture = nil
	}
}

func (s *Sprite) SetMask(mask IMask) ISprite {
	s.lock(func() {
		s.mask = mask
	}, func() bool {
		return s.mask != mask
	})
	return s
}

func (s *Sprite) Mask() IMask {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.mask
}

func (s *Sprite) ResizeTo(width, height uint32) (ISprite, error) {
	var err error
	s.lock(func() {
		if misc.NumberEqual(s.width(), width, misc.Epsilon) && misc.NumberEqual(s.height(), height, misc.Epsilon) {
			return
		}
		s.rect = ti.Rect[float32](0, 0, float32(width), float32(height))
		var newTexture *ti.TiImage
		newTexture, err = ti.NewTiImage(s.renderer.Runtime(), width, height)
		if err != nil {
			return
		}

		if s.texture == nil {
			s.texture = newTexture
			return
		}

		s.renderer.Module().Resize(s.texture, newTexture, ti.ResizeOptions{
			FillMode:  ti.FillModeFit,
			ScaleMode: ti.ScaleModeLanczos,
		})
		if s.texture != nil {
			s.texture.Release()
		}
		s.texture = newTexture
	}, func() bool {
		return misc.NumberEqual(s.width(), width, misc.Epsilon) && misc.NumberEqual(s.height(), height, misc.Epsilon)
	})

	return s, err
}

// BoundingBox 计算旋转+缩放后的精灵包围盒
//   - parentWidth, parentHeight: 父级显示区域尺寸
//     请查看examples/bounding_box_visualize.png以了解原理
//
// 返回：包围盒在屏幕坐标系中的范围
func (s *Sprite) BoundingBox(parentWidth, parentHeight float32) ti.Rectangle[float32] {
	// 精灵中心在屏幕上的位置
	cx, cy := s.Cx(), s.Cy()

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
		scaledDx := dx * s.scaleX
		scaledDy := dy * s.scaleY
		rx := scaledDx*cosR - scaledDy*sinR + cx
		ry := scaledDx*sinR + scaledDy*cosR + cy

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

func (s *Sprite) Add(sprite ISprite) ISprite {
	if !s.isContainer {
		panic("sprite is not a container")
	}
	s.lock(func() {
		s.children = append(s.children, sprite)
	}, func() bool { return true })
	return s
}

func (s *Sprite) Remove(sprite ISprite) ISprite {
	if !s.isContainer {
		panic("sprite is not a container")
	}
	s.lock(func() {
		s.children = slices.DeleteFunc(s.children, func(child ISprite) bool {
			return child == sprite
		})
	}, func() bool { return true })
	return s
}

func (s *Sprite) Children() []ISprite {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.children
}

func (s *Sprite) Render() {
	defer func() {
		s.mutex.Lock()
		s.dirty = false
		s.mutex.Unlock()
	}()

	// 跳过未改变的精灵
	if !s.isContainer || !s.IsDirty() {
		return
	}

	s.mutex.RLock()
	children := s.children
	s.mutex.RUnlock()

	// 置空na
	s.renderer.Module().FillTexture(s.texture, color.Transparent)

	for _, child := range children {
		// 渲染子级
		child.Render()
		childTexture := child.Texture()

		bbox := child.BoundingBox(s.Width(), s.Height())

		mask := child.Mask()

		options := ti.RenderLayerOptions{
			X:        child.X(),
			Y:        child.Y(),
			Width:    child.Width(),
			Height:   child.Height(),
			Cx:       child.Cx(),
			Cy:       child.Cy(),
			ScaleX:   child.ScaleX(),
			ScaleY:   child.ScaleY(),
			Rotation: child.Rotation(),
			Alpha:    child.Alpha(),
			MinX:     int32(bbox.Min.X),
			MaxX:     int32(bbox.Max.X),
			MinY:     int32(bbox.Min.Y),
			MaxY:     int32(bbox.Max.Y),
		}

		if mask != nil {
			s.renderer.Module().RenderLayerWithMask(
				childTexture,
				mask.Texture(),
				s.texture,
				options,
			)
		} else {
			s.renderer.Module().RenderLayerNoMask(
				childTexture,
				s.texture,
				options,
			)
		}
	}

}
