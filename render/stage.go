package render

import (
	"sync"

	"github.com/go-mixed/go-canvas/misc"
	"github.com/go-mixed/go-canvas/ti"
	"github.com/pkg/errors"
)

type Stage struct {
	renderer *Renderer

	container IContainer
	mutex     *sync.RWMutex
}

var _ IParent = (*Stage)(nil)

// NewStage 创建舞台，注意：调用 Stage.Release 时，不会释放 Renderer
func NewStage(renderer *Renderer, width, height uint32) (*Stage, error) {
	s := &Stage{
		renderer: renderer,
		mutex:    &sync.RWMutex{},
	}

	container, err := NewContainer(SelfRelease(renderer), width, height)
	if err != nil {
		return nil, errors.Wrapf(err, "create container of stage failed")
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
}
