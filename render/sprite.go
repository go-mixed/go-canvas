package render

import (
	"github.com/go-mixed/go-canvas/ti"
)

type Sprite struct {
	*Element
	mask IMask
}

var _ ISprite = (*Sprite)(nil)

// NewSprite 创建非容器的精灵，需要传入纹理
func NewSprite(renderer *Renderer, texture *ti.TiImage) *Sprite {
	shape := texture.Shape()
	w, h := shape[0], shape[1]

	element := &Element{
		rect:     ti.Rect(0, 0, float32(w), float32(h)),
		alpha:    1.0,
		texture:  texture,
		scaleX:   1.0,
		scaleY:   1.0,
		deltaCx:  float32(w / 2),
		deltaCy:  float32(h / 2),
		renderer: renderer,
	}
	return &Sprite{
		Element: element.initial(renderer),
		mask:    nil,
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

func (s *Sprite) SetMask(mask IMask) {
	s.lockForUpdate(func() {
		s.mask = mask
	}, func() bool {
		return s.mask != mask
	})
}

func (s *Sprite) Mask() IMask {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.mask
}

func (s *Sprite) Render() {
	defer func() {
		s.SetDirty(false)
	}()

}
