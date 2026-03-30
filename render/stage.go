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

	imageTexture *ti.BgraImage
}

var _ IParent = (*Stage)(nil)

// NewStage 创建舞台，注意：调用 Stage.Release 时，不会释放 Renderer
func NewStage(renderer *Renderer, width, height int) (*Stage, error) {
	s := &Stage{
		renderer: renderer,
		mutex:    &sync.RWMutex{},
	}

	container, err := NewContainer(SelfRelease(renderer), ti.Attr().SetWH(width, height))
	if err != nil {
		return nil, errors.Wrapf(err, "create container of stage failed")
	}
	s.container = container

	imageTexture, err := ti.NewBgraImage(renderer.Runtime(), uint32(height), uint32(width))
	if err != nil {
		return nil, errors.Wrapf(err, "create image texture failed")
	}

	s.imageTexture = imageTexture

	return s, nil
}

// Render 修改之后，需要触发渲染
func (s *Stage) Render() {
	s.container.Render()
}

func (s *Stage) ToBgraImage(buffer []uint32) error {
	// 将 ti image 转换为 bgra image
	s.Renderer().Module().TiImageToBgra(s.Texture(), s.imageTexture)

	return s.imageTexture.MapUint32(func(data []uint32) error {
		copy(buffer, data)
		return nil
	})
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

	if s.imageTexture != nil {
		s.imageTexture.Release()
	}
}
