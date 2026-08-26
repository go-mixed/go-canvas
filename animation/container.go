package animation

import (
	"github.com/go-mixed/go-canvas/internel/misc"
	"github.com/go-mixed/go-canvas/render"
)

const (
	ModifiedFieldScrollLeft ModifiedField = "scroll_left"
	ModifiedFieldScrollTop  ModifiedField = "scroll_top"
)

type ContainerAttributeAnimation struct {
	*BaseAttributeAnimation
	originalScrollLeft int
	originalScrollTop  int

	targetScrollLeft int
	targetScrollTop  int
}

func NewContainerAttributeAnimation(sprite render.IElement) *ContainerAttributeAnimation {
	container, ok := sprite.(render.IContainer)
	if !ok {
		panic("sprite must be ShapeSprite")
	}
	return BuildAttributeAnimation(container, func(base *BaseAttributeAnimation) *ContainerAttributeAnimation {
		return &ContainerAttributeAnimation{
			BaseAttributeAnimation: base,
			originalScrollLeft:     container.GetScrollLeft(),
			originalScrollTop:      container.GetScrollTop(),
		}
	})
}

func (a *ContainerAttributeAnimation) SetScrollLeft(v int) *ContainerAttributeAnimation {
	a.targetScrollLeft = v
	a.MarkModified(ModifiedFieldScrollLeft)
	return a
}

func (a *ContainerAttributeAnimation) SetScrollTop(v int) *ContainerAttributeAnimation {
	a.targetScrollTop = v
	a.MarkModified(ModifiedFieldScrollTop)
	return a
}

func (a *ContainerAttributeAnimation) Ticking(t float32) {
	// 先调用父级的Ticking
	a.BaseAttributeAnimation.Ticking(t)

	if container, ok := a.sprite.(render.IContainer); ok {
		t = a.easing(t)
		if a.IsModified(ModifiedFieldScrollLeft) {
			container.ScrollLeft(misc.Lerp(a.originalScrollLeft, a.targetScrollLeft, t))
		}

		if a.IsModified(ModifiedFieldScrollTop) {
			container.ScrollTop(misc.Lerp(a.originalScrollTop, a.targetScrollTop, t))
		}

	}

}
