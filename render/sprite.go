package render

import (
	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
)

type Sprite struct {
	*tiElement

	parent IParent
	// 真实的Sprite的实例，
	instance ISprite

	masks *misc.List[IMask]
}

var _ ISprite = (*Sprite)(nil)
var _ IMaskParent = (*Sprite)(nil)

// BuildSprite 创建非容器的精灵，需要传入纹理
func BuildSprite[T ISprite](parent IParent, texture *ti.TiImage, instanceCreator func(s *Sprite) (T, error)) (T, error) {
	var w, h uint32

	if texture != nil {
		shape := texture.Shape()
		w, h = shape[0], shape[1]
	}

	element := &tiElement{
		rect:    ti.Rect(0, 0, float32(w), float32(h)),
		alpha:   1.0,
		texture: texture,
		scaleX:  1.0,
		scaleY:  1.0,
		deltaCx: float32(w / 2),
		deltaCy: float32(h / 2),
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
func NewBlockSprite(parent IParent, width, height uint32) (*Sprite, error) {
	texture, err := ti.NewTiImage(parent.Renderer().Runtime(), width, height)
	if err != nil {
		return nil, err
	}

	return BuildSprite(parent, texture, func(s *Sprite) (*Sprite, error) {
		return s, nil
	})
}

func (s *Sprite) Render() {
	defer func() {
		s.SetDirty(false)
	}()
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

}
