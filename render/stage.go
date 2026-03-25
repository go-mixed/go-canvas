package render

import (
	"github.com/go-mixed/go-canvas/ti"
	"iter"
	"slices"
	"sync"
)

// Stage 舞台
type Stage struct {
	children []ISprite

	renderer *Renderer
	screen   ISprite
	mutex    sync.Mutex
}

func NewStage(r *Renderer, width, height uint32) (*Stage, error) {
	screen, err := NewBlockSprite(r, width, height)
	if err != nil {
		return nil, err
	}
	return &Stage{
		renderer: r,
		screen:   screen,
	}, nil
}

func (s *Stage) Add(sprite ISprite) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.children = append(s.children, sprite)
}

func (s *Stage) Remove(sprite ISprite) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.children = slices.DeleteFunc(s.children, func(child ISprite) bool {
		return child == sprite
	})
}

func (s *Stage) Children() iter.Seq[ISprite] {
	return func(yield func(ISprite) bool) {
		for _, child := range s.children {
			if !yield(child) {
				return
			}
		}
	}
}

func (s *Stage) Render() {
	s.screen.FillTexture(ti.ColorWhite)

	for child := range s.Children() {

		bbox := child.BoundingBox(s.screen.Width(), s.screen.Height())

		mask := child.Mask()

		options := ti.RenderLayerOptions{
			X:        child.X(),
			Y:        child.Y(),
			Width:    child.Width(),
			Height:   child.Height(),
			Cx:       child.CenterX(),
			Cy:       child.CenterY(),
			Scale:    child.Scale(),
			Rotation: child.Rotation(),
			Alpha:    child.Alpha(),
			MinX:     int32(bbox.Min.X),
			MaxX:     int32(bbox.Max.X),
			MinY:     int32(bbox.Min.Y),
			MaxY:     int32(bbox.Max.Y),
		}

		if mask != nil {
			s.renderer.Module().RenderLayerWithMask(
				child.Texture(),
				mask.Texture(),
				s.screen.Texture(),
				options,
			)
		} else {
			s.renderer.Module().RenderLayerNoMask(
				child.Texture(),
				s.screen.Texture(),
				options,
			)
		}
	}
}

func (s *Stage) Texture() *ti.TiImage {
	return s.screen.Texture()
}

func (s *Stage) Release() {
	s.screen.Release()
}
