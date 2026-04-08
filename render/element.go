package render

import (
	"image/color"
	"math"
	"sync"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
	"github.com/pkg/errors"
)

type tiElement struct {
	attribute *ti.Attribute

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

	return e.dirty
}

func (e *tiElement) SetDirty(val bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.dirty = val
}

func (e *tiElement) SetX(x int) {
	e.MoveTo(x, e.attribute.Y())
}

func (e *tiElement) X() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.attribute.X()
}

func (e *tiElement) SetY(y int) {
	e.MoveTo(e.attribute.X(), y)
}

func (e *tiElement) Y() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.attribute.Y()
}

func (e *tiElement) MoveTo(x int, y int) {
	e.LockForUpdate(func() {
		e.attribute.MoveTo(x, y)
	}, func() bool {
		return !misc.NumberEqual(e.attribute.X(), x, misc.Epsilon) || !misc.NumberEqual(e.attribute.Y(), y, misc.Epsilon)
	})
}

func (e *tiElement) Width() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.attribute.Width()
}

func (e *tiElement) Height() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.attribute.Height()
}

func (e *tiElement) SetScale(scaleX, scaleY float32) {
	e.LockForUpdate(func() {
		e.attribute.SetScale(scaleX, scaleY)
	}, func() bool {
		return !misc.NumberEqual(e.attribute.ScaleX(), scaleX, misc.Epsilon) || !misc.NumberEqual(e.attribute.ScaleY(), scaleY, misc.Epsilon)
	})
}

func (e *tiElement) ScaleX() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.attribute.ScaleX()
}

func (e *tiElement) ScaleY() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.attribute.ScaleY()
}

func (e *tiElement) SetRotation(rotation float32) {
	e.LockForUpdate(func() {
		e.attribute.SetRotation(rotation)
	}, func() bool {
		return !misc.NumberEqual(e.attribute.Rotation(), rotation, misc.Epsilon)
	})
}

func (e *tiElement) Rotation() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.attribute.Rotation()
}

func (e *tiElement) SetAlpha(alpha float32) {
	e.LockForUpdate(func() {
		e.attribute.SetAlpha(alpha)
	}, func() bool {
		return !misc.NumberEqual(e.attribute.Alpha(), alpha, misc.Epsilon)
	})
}

func (e *tiElement) Alpha() float32 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.attribute.Alpha()
}

func (e *tiElement) Cx() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.attribute.Cx()
}

func (e *tiElement) SetCx(cx int) {
	e.LockForUpdate(func() {
		e.attribute.SetCx(cx)
	}, func() bool {
		return e.attribute.Cx() != cx
	})
}

func (e *tiElement) Cy() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.attribute.Cy()
}

func (e *tiElement) SetCy(cy int) {
	e.LockForUpdate(func() {
		e.attribute.SetCy(cy)
	}, func() bool { return e.attribute.Cy() != cy })
}

func (e *tiElement) Attribute() *ti.Attribute {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.attribute
}

// Fill 填充纯色
func (e *tiElement) Fill(color color.Color) {
	e.LockForUpdate(func() {
		e.Renderer().Module().FillColor(e.texture, color)
	}, func() bool { return false })
}

func (e *tiElement) Texture() *ti.TiImage {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.texture
}

func (e *tiElement) Resize(width, height int) error {
	if width <= 0 || height <= 0 {
		return errors.New("width or height must be greater than 0")
	}

	var err error
	e.LockForUpdate(func() {
		if e.attribute.Width() == width && e.attribute.Width() == height {
			return
		}

		nW, nH := ti.CalcResizeWH(e.attribute.Width(), e.attribute.Width(), width, height, e.attribute.ResizeOptions())

		var newTexture *ti.TiImage
		newTexture, err = ti.NewTiImage(e.Renderer().Runtime(), uint32(nW), uint32(nH))
		if err != nil {
			return
		}

		if e.texture == nil {
			e.texture = newTexture
			return
		}

		e.Renderer().Module().Resize(e.texture, newTexture, e.attribute.ResizeOptions())

		e.texture.Release()
		e.texture = newTexture
		e.attribute.SetWH(nW, nH)
		e.attribute.SetCxy(width/2, height/2)
	}, func() bool {
		return e.attribute.Width() != width || e.attribute.Width() != height
	})

	if err != nil {
		return errors.Wrap(err, "resize texture failed")
	}
	return nil
}

// Blur 模糊纹理（马赛克/高斯/普通），原地修改
func (e *tiElement) Blur(mode ti.BlurMode, radius int32) error {
	var err error

	e.LockForUpdate(func() {
		if e.texture == nil {
			return
		}

		var newTexture *ti.TiImage
		shape := e.texture.Shape()
		width, height := shape[0], shape[1]
		newTexture, err = ti.NewTiImage(e.Renderer().Runtime(), width, height)
		if err != nil {
			return
		}

		e.Renderer().Module().Blur(e.texture, newTexture, mode, radius)
		e.texture.Release()
		e.texture = newTexture
	}, func() bool { return true })

	if err != nil {
		return errors.Wrap(err, "blur texture failed")
	}
	return nil
}

// ClientRect 获取元素自身旋转+缩放后的边界框
// 注意：不与父级区域求交集。
// 计算过程分两步：
// 1) 先在元素局部坐标系（左上角为原点）计算旋转后的包围盒；
// 2) 再将该包围盒整体平移到父坐标系（通过 Add(x, y)）。
// 请查看examples/bounding_box_visualize.png以了解原理
//
// 纹理四个角点（相对于纹理左上角）
// 假设纹理左上角为(0,0)，中心点为(cx, cy)：
// (0, 0)           ●────────● (width, 0)
//
//		            │        │
//	             	│   ●    │  ← 中心点 (cx, cy)
//	             	│        │
//
// (0, height)      ●────────● (width, height)
//
// corners存储角点相对于中心点的偏移量：
// 左上: (-cx, -cy)，右上: (width-cx, -cy)
// 右下: (width-cx, height-cy)，左下: (-cx, height-cy)
func (e *tiElement) ClientRect() ti.Rectangle[int] {
	// Local bbox: all computations below assume element origin at (0, 0).
	// World/parent bbox is obtained by translating local bbox with (attribute.X, attribute.Y).
	cx, cy := float32(e.attribute.Cx()), float32(e.attribute.Cy())
	w, h := float32(e.attribute.Width()), float32(e.attribute.Height())

	corners := [][2]float32{
		{0 - cx, 0 - cy},
		{w - cx, 0 - cy},
		{w - cx, h - cy},
		{0 - cx, h - cy},
	}

	cosR := float32(math.Cos(float64(e.attribute.Rotation())))
	sinR := float32(math.Sin(float64(e.attribute.Rotation())))
	// Snap trigonometric results near 0/1 to avoid bbox jitter on right angles.
	const trigEps = 1e-6
	if math.Abs(float64(cosR)) < trigEps {
		cosR = 0
	} else if math.Abs(math.Abs(float64(cosR))-1.0) < trigEps {
		cosR = float32(math.Copysign(1.0, float64(cosR)))
	}
	if math.Abs(float64(sinR)) < trigEps {
		sinR = 0
	} else if math.Abs(math.Abs(float64(sinR))-1.0) < trigEps {
		sinR = float32(math.Copysign(1.0, float64(sinR)))
	}

	var minX, maxX, minY, maxY float32
	minX, maxX = math.MaxFloat32, -math.MaxFloat32
	minY, maxY = math.MaxFloat32, -math.MaxFloat32

	for _, corner := range corners {
		dx, dy := corner[0], corner[1]
		scaledDx := dx * e.attribute.ScaleX()
		scaledDy := dy * e.attribute.ScaleY()
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

	bbox := ti.RectXY(int(minX), int(minY), int(maxX), int(maxY))
	// Translate local bbox into parent coordinates.
	// Use Add(x, y) instead of MoveTo(x, y): MoveTo would overwrite Min and
	// break negative local bounds produced by rotation.
	return bbox.Add(ti.Pt(e.attribute.X(), e.attribute.Y()))
}

func (e *tiElement) Renderer() *Renderer {
	return e.renderer
}
