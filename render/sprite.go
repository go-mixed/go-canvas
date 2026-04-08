package render

import (
	"math"

	"github.com/go-mixed/go-canvas/misc"
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
func BuildSprite[T ISprite](parent IParent, attribute *ti.Attribute, texture *ti.TiImage, instanceCreator func(s *Sprite) (T, error)) (T, error) {
	if attribute.Cx() == math.MaxInt64 {
		attribute.SetCx(attribute.Width() / 2)
	}
	if attribute.Cy() == math.MaxInt64 {
		attribute.SetCy(attribute.Height() / 2)
	}

	element := &tiElement{
		attribute: attribute,
		texture:   texture,
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
	}, func() bool {
		return true
	})
}

func (s *Sprite) RemoveMask(mask IMask) {
	s.LockForUpdate(func() {
		s.masks.RemoveAll(func(child IMask) bool {
			return child == mask
		})

		// 递归删除
		mask.Release()
	}, func() bool {
		return true
	})
}

func (s *Sprite) Masks() *misc.List[IMask] {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.masks
}

// NewBlockSprite 创建纯色块精灵，颜色为透明的
func NewBlockSprite(parent IParent, attribute *ti.Attribute) (*Sprite, error) {
	texture, err := ti.NewTiImage(parent.Renderer().Runtime(), uint32(attribute.Width()), uint32(attribute.Height()))
	if err != nil {
		return nil, err
	}

	return BuildSprite(parent, attribute, texture, func(s *Sprite) (*Sprite, error) {
		return s, nil
	})
}

func (s *Sprite) Render(frameIndex int) {
	defer func() {
		s.SetDirty(false)
	}()
}

func (s *Sprite) Animate(targetFn ti.TargetAttributeFn, startAtFrame, durationFrame int) ISprite {
	s.animator.enqueue(targetFn, startAtFrame, durationFrame)

	return s
}

func (s *Sprite) ClearAnimations() ISprite {
	s.animator.clear()

	return s
}

func (s *Sprite) StopAnimation(reset bool) ISprite {
	s.animator.stop(reset)
	return s
}

func (s *Sprite) TickAnimation(frameIndex int) bool {
	return s.animator.tick(frameIndex)
}

// RemoveFromParent 从父级移除精灵
func (s *Sprite) RemoveFromParent() {
	s.parent.RemoveChild(s.instance)
}

func (s *Sprite) Release() {
	s.mutex.RLock()
	if s.texture != nil {
		s.texture.Release()
	}
	s.mutex.RUnlock()

	for _, mask := range s.masks.Range() {
		mask.Release()
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.texture = nil
	s.masks.Clear()
	s.animator.clear()
}
