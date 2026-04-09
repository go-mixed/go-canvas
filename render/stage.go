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

	options stageOptions
}

var _ IParent = (*Stage)(nil)

// NewStage 创建舞台，注意：调用 Stage.Release 时，不会释放 Renderer
func NewStage(renderer *Renderer, width, height int, opts ...stageOptFunc) (*Stage, error) {
	options := stageOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	s := &Stage{
		renderer: renderer,
		mutex:    &sync.RWMutex{},
		options:  options,
	}

	container, err := NewContainer(SelfRelease(renderer), ti.Attr().SetWH(width, height))
	if err != nil {
		return nil, errors.Wrapf(err, "create container of stage failed")
	}
	s.container = container

	if options.enabledRAWImage {
		imageTexture, err := ti.NewBgraImage(renderer.Runtime(), uint32(height), uint32(width))
		if err != nil {
			s.container.Release()
			return nil, errors.Wrapf(err, "create image texture failed")
		}
		s.imageTexture = imageTexture
	}

	return s, nil
}

// Render 修改之后，需要调用本函数来渲染，之后才能得到渲染结果
func (s *Stage) Render(frameIndex int) {
	s.container.Render(frameIndex)

	if s.options.enabledRAWImage {
		// 将 ti image 转换为 bgra image
		s.Renderer().Module().TiImageToBgra(s.Texture(), s.imageTexture)
	}

	s.renderer.runtime.Wait()

	// 每帧渲染结束后，释放垃圾纹理
	// 后置释放，这样可以让Resize、Blur等方法异步执行
	s.container.ReleaseGarbageTextures()
}

func (s *Stage) IsDirty() bool {
	return s.container.IsDirty()
}

// HasAnimationAt returns true when stage tree has animation to be evaluated
// at the given absolute frame.
func (s *Stage) HasAnimationAt(frameIndex int) bool {
	return s.container.HasAnimationAt(frameIndex)
}

func (s *Stage) GetBgraImage(buffer []uint32) error {
	if s.imageTexture == nil {
		return errors.New("image texture is nil")
	}
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

func (s *Stage) Width() int {
	return s.container.Width()
}

func (s *Stage) Height() int {
	return s.container.Height()
}

func (s *Stage) Children() *misc.List[ISprite] {
	return s.container.Children()
}

func (s *Stage) Release() {
	if s.container != nil {
		s.container.Release()
	}

	if s.options.enabledRAWImage {
		s.imageTexture.Release()
	}
}
