package render

import (
	"sync"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
	"github.com/go-mixed/go-taichi/taichi"
)

type Stage struct {
	renderer *Renderer

	container IContainer
	mutex     *sync.RWMutex
}

var _ IParent = (*Stage)(nil)

func NewStage(arch taichi.Arch, width, height uint32) (*Stage, error) {
	renderer, err := NewRenderer(arch)
	if err != nil {
		return nil, err
	}

	return NewStageWithRenderer(renderer, width, height)
}

func NewStageWithRenderer(renderer *Renderer, width, height uint32) (*Stage, error) {
	s := &Stage{
		renderer: renderer,
		mutex:    &sync.RWMutex{},
	}

	container, err := NewContainer(SelfRelease(renderer), width, height)
	if err != nil {
		s.Release()
		return nil, err
	}
	s.container = container

	return s, nil
}

// Render 修改之后，需要触发渲染
func (s *Stage) Render() {
	s.container.Render()
}

func (s *Stage) Texture() *ti.TiImage {
	return s.container.Texture()
}

func (s *Stage) Renderer() *Renderer {
	return s.renderer
}

func (s *Stage) AddChild(child ISprite) {
	s.container.AddChild(child)
}

func (s *Stage) RemoveChild(child ISprite) {
	s.container.RemoveChild(child)
}

func (s *Stage) Children() *misc.List[ISprite] {
	return s.container.Children()
}

func (s *Stage) Release() {

	if s.container != nil {
		s.container.Release()
	}

	if s.renderer != nil {
		s.renderer.Release()
	}
}
