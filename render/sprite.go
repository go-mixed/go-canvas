package render

import (
	"image/color"

	"github.com/go-mixed/go-canvas/ctypes"
	"github.com/go-mixed/go-canvas/internel/misc"
	"github.com/go-mixed/go-canvas/ti"
)

type Sprite struct {
	*tiElement

	parent IParent
	// 真实的Sprite的实例，
	instance ISprite

	masks    *misc.List[IMask]
	animator *spriteAnimator
}

var _ ISprite = (*Sprite)(nil)
var _ IMaskParent = (*Sprite)(nil)

// BuildSprite 创建非容器的精灵，需要传入纹理
func BuildSprite[T ISprite](parent IParent, attribute *ctypes.Attribute, texture *ctypes.TiImage, instanceCreator func(s *Sprite) (T, error)) (T, error) {
	attribute.SetCxIfNotDefined(attribute.Width() / 2)
	attribute.SetCyIfNotDefined(attribute.Height() / 2)

	element := &tiElement{
		attribute: attribute,
		texture:   texture,
		dirty:     attribute.Dirty(),
	}
	s := &Sprite{
		tiElement: element.initial(parent.Renderer()),
		parent:    parent,
		masks:     misc.NewList[IMask](),
	}

	instance, err := instanceCreator(s)
	if err != nil {
		var nilT T
		return nilT, err
	}

	s.animator = newSpriteAnimator(instance)

	// 添加到父级
	s.parent.AddChild(instance)
	s.instance = instance

	return instance, nil
}

func (s *Sprite) AddMask(mask IMask) {
	s.LockForUpdate(func() {
		if s.masks.Index(func(child IMask) bool {
			return child == mask
		}) >= 0 {
			return
		}
		s.masks.PushBack(mask)
	}, func() ctypes.DirtyMode {
		return ctypes.DirtyModeMask | ctypes.DirtyModeComposite
	})
}

func (s *Sprite) RemoveMask(mask IMask) {
	s.LockForUpdate(func() {
		s.masks.RemoveAll(func(child IMask) bool {
			return child == mask
		})

		// 递归删除
		mask.Release()
	}, func() ctypes.DirtyMode {
		return ctypes.DirtyModeMask | ctypes.DirtyModeComposite
	})
}

func (s *Sprite) Masks() *misc.List[IMask] {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.masks
}

// NewBlockSprite 创建纯色块精灵，颜色为透明的
func NewBlockSprite(parent IParent, attribute *ctypes.Attribute) (*Sprite, error) {
	texture, err := ti.NewTiImage(parent.Renderer().Runtime(), uint32(attribute.Width()), uint32(attribute.Height()))
	if err != nil {
		return nil, err
	}

	return BuildSprite(parent, attribute, texture, func(s *Sprite) (*Sprite, error) {
		return s, nil
	})
}

func (s *Sprite) SetDirty(val ctypes.DirtyMode) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.dirty = val
}

func (s *Sprite) IsDirty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.dirty != ctypes.DirtyModeNone {
		return true
	}

	// 检测mask是否 dirty
	for _, mask := range s.masks.Range() {
		if mask.IsDirty() {
			return true
		}
	}

	return false
}

func (s *Sprite) Render(frameIndex int) error {
	defer func() {
		s.SetDirty(ctypes.DirtyModeNone)
	}()

	for _, mask := range s.masks.Range() {
		if err := mask.Render(frameIndex); err != nil {
			return err
		}
	}

	if err := s.renderCanvas(); err != nil {
		return err
	}

	return nil
}

// renderCanvas 在 Layout/Painting 脏时重建并重绘展示 canvas。
// canvas 尺寸变化时会新建纹理，并将旧纹理加入垃圾回收列表。
func (s *Sprite) renderCanvas() error {
	s.mutex.RLock()
	dirty := s.dirty
	clientW := s.attribute.ClientWidth()
	clientH := s.attribute.ClientHeight()
	src := s.texture
	border := s.attribute.Border()
	padding := s.attribute.Padding()
	blur := s.attribute.Blur()
	s.mutex.RUnlock()

	if src == nil || dirty&(ctypes.DirtyModeLayout|ctypes.DirtyModePainting) == 0 {
		return nil
	}

	if clientW <= 0 || clientH <= 0 {
		return nil
	}

	var originalCanvasW, originalCanvasH int
	if s.canvas != nil {
		shape := s.canvas.Shape()
		originalCanvasW = int(shape[0])
		originalCanvasH = int(shape[1])
	}

	// 新建新的canvas
	if originalCanvasW != clientW || originalCanvasH != clientH || s.canvas == nil {
		newCanvas, err := ti.NewTiImage(s.renderer.Runtime(), uint32(clientW), uint32(clientH))
		if err != nil {
			return err
		}
		s.mutex.Lock()
		s.addGarbageTexture(s.canvas)
		s.canvas = newCanvas
		s.mutex.Unlock()
	} else {
		// 清空
		s.Renderer().Module().FillColor(s.canvas, color.Transparent)
	}

	contentX := max(0, border.LeftWidth+padding.Left)
	contentY := max(0, border.TopWidth+padding.Top)
	contentW := max(0, s.attribute.Width())
	contentH := max(0, s.attribute.Height())

	shape := s.texture.Shape()
	srcW, srcH := int(shape[0]), int(shape[1])

	// 需要Resize
	if s.attribute.Width() != srcW || s.attribute.Height() != srcH {
		s.Renderer().Module().AsyncResize(src, s.canvas,
			s.attribute.ResizeOptions(),
			ctypes.RectWH(0, 0, srcW, srcH),
			ctypes.RectWH(contentX, contentY, contentW, contentH),
		)
	} else {
		s.Renderer().Module().AsyncCopy(src, s.canvas,
			ctypes.RectWH(0, 0, srcW, srcH),
			ctypes.RectWH(contentX, contentY, contentW, contentH),
		)
	}

	// 需要Blur
	if !blur.IsEmpty() {
		newCanvas, err := ti.NewTiImage(s.renderer.Runtime(), uint32(clientW), uint32(clientH))
		if err != nil {
			return err
		}
		s.Renderer().Module().AsyncBlur(s.canvas, newCanvas, blur.Mode, int32(blur.Radius))
		s.mutex.Lock()
		s.addGarbageTexture(s.canvas)
		s.canvas = newCanvas
		s.mutex.Unlock()
	}

	// 渲染盒模型：border 叠加 + 圆角裁剪（content 已由 AsyncResize 写入）
	if !s.attribute.Border().IsEmpty() {
		s.Renderer().Module().AsyncRenderBorder(s.canvas, border)
	}

	return nil
}

// HasAnimationAt returns true when an animation segment should be evaluated
// at the given absolute frame.
func (s *Sprite) HasAnimationAt(frameIndex int) bool {
	if s.animator.hasAnimationAt(frameIndex) {
		return true
	}

	for _, mask := range s.masks.Range() {
		_maskSprite, ok := mask.(IAnimation)
		if ok && _maskSprite.HasAnimationAt(frameIndex) {
			return true
		}
	}
	return false
}

// Animate 向当前精灵追加一段动画。
func (s *Sprite) Animate(targetFn ti.TargetAttributeFn, startFrameIndex, durationFrames int) ISprite {
	s.animator.enqueue(targetFn, startFrameIndex, durationFrames)

	return s
}

// ClearAnimations 清空当前精灵所有未完成动画。
func (s *Sprite) ClearAnimations() ISprite {
	s.animator.clear()

	return s
}

// StopAnimation 停止动画；reset=true 时恢复到当前段的起始状态。
func (s *Sprite) StopAnimation(reset bool) ISprite {
	s.animator.stop(reset)
	return s
}

// TickAnimation 按绝对帧号推进动画。
func (s *Sprite) TickAnimation(frameIndex int) bool {
	v := s.animator.tick(frameIndex)
	for _, mask := range s.masks.Range() {
		_maskSprite, ok := mask.(IAnimation)
		if ok && _maskSprite.TickAnimation(frameIndex) {
			v = true
		}
	}
	return v
}

// RemoveFromParent 从父级移除精灵
func (s *Sprite) RemoveFromParent() {
	s.parent.RemoveChild(s.instance)
}

// SetInstance 设置实例
func (s *Sprite) SetInstance(m ISprite) {
	s.LockForUpdate(func() {
		s.instance = m
		s.animator.setSprite(m)
	}, func() ctypes.DirtyMode {
		if s.instance != m {
			return ctypes.DirtyModeComposite
		}
		return ctypes.DirtyModeNone
	})
}

func (s *Sprite) Release() {
	s.mutex.RLock()
	if s.texture != nil {
		s.texture.Release()
	}
	if s.canvas != nil {
		s.canvas.Release()
	}
	s.mutex.RUnlock()

	for _, mask := range s.masks.Range() {
		mask.Release()
	}

	s.ReleaseGarbageTextures()

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.texture = nil
	s.canvas = nil
	s.masks.Clear()
	s.animator.clear()
}
