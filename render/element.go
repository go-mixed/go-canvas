package render

import (
	"image/color"
	"math"
	"sync"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
)

type tiElement struct {
	rect           ti.Rectangle[float32]
	scaleX, scaleY float32 // 1.0 for no scaling
	rotation       float32 // 0.0 for no rotation, in radians
	alpha          float32 // 0.0 for no alpha, 1.0 for full alpha

	// 中心点的相对值
	deltaCx, deltaCy float32

	texture *ti.TiImage

	mutex    *sync.RWMutex
	renderer *Renderer

	dirty bool
}

func (e *tiElement) initial(renderer *Renderer) *tiElement {
	e.mutex = &sync.RWMutex{}
	e.renderer = renderer
	return e
}

// LockForUpdate 锁定并更新数据，请勿在 SetDirty 变量中使用RLock, Lock，否则会死锁
func (e *tiElement) LockForUpdate(updateFn func(), triggerDirty func() bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if triggerDirty() {
		e.dirty = true
	}
	updateFn()
}

func (e *tiElement) IsDirty() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	//for _, child := range e.children {
	//	if child.IsDirty() {
	//		return true
	//	}
	//}

	return e.dirty
}

func (e *tiElement) SetDirty(val bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.dirty = val
}

func (e *tiElement) SetX(x float32) {
	e.MoveTo(x, e.Y())
}

func (e *tiElement) X() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.rect.X()
}

func (e *tiElement) SetY(y float32) {
	e.MoveTo(e.X(), y)
}

func (e *tiElement) Y() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.rect.Y()
}

func (e *tiElement) MoveTo(x float32, y float32) {
	e.LockForUpdate(func() {
		e.rect = e.rect.MoveTo(x, y)
	}, func() bool {
		return misc.NumberEqual(e.rect.X(), x, misc.Epsilon) && misc.NumberEqual(e.rect.Y(), y, misc.Epsilon)
	})
}

func (e *tiElement) Width() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.width()
}

func (e *tiElement) width() float32 {
	return e.rect.Dx()
}

func (e *tiElement) Height() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.height()
}

func (e *tiElement) height() float32 {
	return e.rect.Dy()
}

func (e *tiElement) SetScale(scaleX, scaleY float32) {
	e.LockForUpdate(func() {
		e.scaleX = scaleX
		e.scaleY = scaleY
	}, func() bool {
		return misc.NumberEqual(e.scaleX, scaleX, misc.Epsilon) && misc.NumberEqual(e.scaleY, scaleY, misc.Epsilon)
	})
}

func (e *tiElement) ScaleX() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.scaleX
}

func (e *tiElement) ScaleY() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.scaleY
}

func (e *tiElement) SetRotation(rotation float32) {
	e.LockForUpdate(func() {
		e.rotation = rotation
	}, func() bool {
		return misc.NumberEqual(e.rotation, rotation, misc.Epsilon)
	})
}

func (e *tiElement) Rotation() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.rotation
}

func (e *tiElement) SetAlpha(alpha float32) {
	e.LockForUpdate(func() {
		e.alpha = alpha
	}, func() bool {
		return misc.NumberEqual(e.alpha, alpha, misc.Epsilon)
	})
}

func (e *tiElement) Alpha() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.alpha
}

func (e *tiElement) Cx() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.X() + e.deltaCx
}

func (e *tiElement) SetCx(cx float32) {
	e.LockForUpdate(func() {
		e.deltaCx = cx
	}, func() bool {
		return misc.NumberEqual(e.deltaCx, cx, misc.Epsilon)
	})
}

func (e *tiElement) Cy() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.Y() + e.deltaCy
}

func (e *tiElement) SetCy(cy float32) {
	e.LockForUpdate(func() {
		e.deltaCy = cy
	}, func() bool { return misc.NumberEqual(e.deltaCy, cy, misc.Epsilon) })
}

//// SetTexture 设置精灵纹理
//func (s *tiElement) SetTexture(texture *ti.TiImage)  {
//	s.LockForUpdate(func() {
//		if s.texture != nil {
//			s.texture.Release()
//		}
//		s.texture = texture
//	}, func() bool { return s.texture != texture })
//	return s
//}

// FillTexture 填充纯色
func (e *tiElement) FillTexture(color color.Color) {
	e.LockForUpdate(func() {
		e.Renderer().Module().FillTexture(e.texture, color)
	}, func() bool { return false })
}

func (e *tiElement) Texture() *ti.TiImage {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.texture
}

func (e *tiElement) ResizeTo(width, height uint32) error {
	var err error
	e.LockForUpdate(func() {
		if misc.NumberEqual(e.width(), width, misc.Epsilon) && misc.NumberEqual(e.height(), height, misc.Epsilon) {
			return
		}
		e.rect = ti.Rect[float32](0, 0, float32(width), float32(height))
		var newTexture *ti.TiImage
		newTexture, err = ti.NewTiImage(e.Renderer().Runtime(), width, height)
		if err != nil {
			return
		}

		if e.texture == nil {
			e.texture = newTexture
			return
		}

		e.Renderer().Module().Resize(e.texture, newTexture, ti.ResizeOptions{
			FillMode:  ti.FillModeFit,
			ScaleMode: ti.ScaleModeLanczos,
		})
		if e.texture != nil {
			e.texture.Release()
		}
		e.texture = newTexture
	}, func() bool {
		return misc.NumberEqual(e.width(), width, misc.Epsilon) && misc.NumberEqual(e.height(), height, misc.Epsilon)
	})

	return err
}

// ClientRect 获取元素自身旋转+缩放后的边界框
// 注意：不与父级区域求交集，返回的是元素在自身坐标系下的实际边界
// 请查看examples/bounding_box_visualize.png以了解原理
//
// 纹理四个角点（相对于纹理左上角）
// 假设纹理左上角为(0,0)，中心点为(cx, cy)：
// (0, 0)           ●────────● (width, 0)
//
//	│        │
//	│   ●    │  ← 中心点 (cx, cy)
//	│        │
//
// (0, height)      ●────────● (width, height)
//
// corners存储角点相对于中心点的偏移量：
// 左上: (-cx, -cy)，右上: (width-cx, -cy)
// 右下: (width-cx, height-cy)，左下: (-cx, height-cy)
func (e *tiElement) ClientRect() ti.Rectangle[float32] {
	cx, cy := e.Cx(), e.Cy()

	corners := [][2]float32{
		{0 - cx, 0 - cy},
		{e.Width() - cx, 0 - cy},
		{e.Width() - cx, e.Height() - cy},
		{0 - cx, e.Height() - cy},
	}

	cosR := float32(math.Cos(float64(e.rotation)))
	sinR := float32(math.Sin(float64(e.rotation)))

	var minX, maxX, minY, maxY float32
	minX, maxX = math.MaxFloat32, -math.MaxFloat32
	minY, maxY = math.MaxFloat32, -math.MaxFloat32

	for _, corner := range corners {
		dx, dy := corner[0], corner[1]
		scaledDx := dx * e.scaleX
		scaledDy := dy * e.scaleY
		rx := scaledDx*cosR - scaledDy*sinR + cx
		ry := scaledDx*sinR + scaledDy*cosR + cy

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

	return ti.Rect(minX, minY, maxX, maxY)
}

// ClippedRect 获取与父级区域裁剪后的可视区域
//   - parentWidth, parentHeight: 父级显示区域尺寸
//     请查看examples/bounding_box_visualize.png以了解原理
//
// 返回：边界框在屏幕坐标系中的范围（与父级区域求交集后的结果）
func (e *tiElement) ClippedRect(parentWidth, parentHeight float32) ti.Rectangle[float32] {
	parentRect := ti.Rect(0, 0, parentWidth, parentHeight)
	return e.ClientRect().Intersect(parentRect)
}

func (e *tiElement) Renderer() *Renderer {
	return e.renderer
}
